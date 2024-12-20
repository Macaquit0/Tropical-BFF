package internalaws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Macaquit0/Tropical-BFF/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var maxMessages = int32(10)

var all = "All"

type Message interface {
	Decode(out interface{}) error
	DecodeModified(out interface{}, changes interface{}) error
	Attribute(key string) string
}

type PublishMessage struct {
	Type    string         `json:"type"`
	Message PublisherInput `json:"message"`
}

type message struct {
	types.Message
	err chan error
}

func newMessage(m types.Message) *message {
	return &message{m, make(chan error, 1)}
}

func (m *message) body() []byte {
	return []byte(*m.Message.Body)
}

func (m *message) Decode(out interface{}) error {
	fmt.Println(string(m.body()))
	return json.Unmarshal(m.body(), &out)
}

func (m *message) ErrorResponse(ctx context.Context, err error) error {
	go func() {
		m.err <- err
	}()
	return err
}

func (m *message) Success(ctx context.Context) error {
	go func() {
		m.err <- nil
	}()

	return nil
}

var ErrMarshal = errors.New("invalid marshal")

const (
	DataTypeNumber = dataType("Number")
	DataTypeString = dataType("String")
)

type dataType string

type Config struct {
	Region            string
	Hostname          string
	AWSAccountID      string
	Env               string
	TopicPrefix       string
	TopicARN          string
	QueueURL          string
	VisibilityTimeout int
	RetryCount        int
	WorkerPool        int
	ExtensionLimit    *int
	ColdStart         bool
	AwsEndpoint       string
}

type Handler func(ctx context.Context, params PublisherInput) error

type AwsSqsConsumer interface {
	Consume()
}

type consumer struct {
	sqs               *sqs.Client
	handler           Handler
	env               string
	QueueURL          string
	Hostname          string
	VisibilityTimeout int
	workerPool        int
	extensionLimit    int
	logger            *logger.Logger
}

func NewConsumerDeprecated(logger *logger.Logger, c Config, handler Handler) (AwsSqsConsumer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	if c.AwsEndpoint != "" {
		cfg.BaseEndpoint = &c.AwsEndpoint
	}

	cfg.Region = c.Region

	cons := &consumer{
		sqs:               sqs.NewFromConfig(cfg),
		env:               c.Env,
		VisibilityTimeout: 30,
		workerPool:        30,
		extensionLimit:    2,
		handler:           handler,
		logger:            logger,
	}

	if c.VisibilityTimeout != 0 {
		cons.VisibilityTimeout = c.VisibilityTimeout
	}

	if c.WorkerPool != 0 {
		cons.workerPool = c.WorkerPool
	}

	if c.ExtensionLimit != nil {
		cons.extensionLimit = *c.ExtensionLimit
	}

	cons.QueueURL = c.QueueURL

	return cons, nil
}

func (c *consumer) Consume() {
	ctx := context.Background()

	jobs := make(chan *message)
	for w := 1; w <= c.workerPool; w++ {
		go c.worker(w, jobs)
	}

	for {
		output, err := c.sqs.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: &c.QueueURL, MaxNumberOfMessages: maxMessages, MessageAttributeNames: []string{all}})
		if err != nil {
			c.logger.Error(ctx).Msg("retrying in 10s %s", err.Error())
			time.Sleep(10 * time.Second)
			continue
		}

		for _, m := range output.Messages {
			jobs <- newMessage(m)
		}
	}
}

func (c *consumer) worker(id int, messages <-chan *message) {
	for m := range messages {
		if err := c.run(m); err != nil {
			c.logger.Error(context.Background()).Msg(err.Error())
		}
	}
}

func (c *consumer) run(m *message) error {
	var params *PublishMessage

	ctx := context.Background()
	c.logger.Info(ctx).Msg("worker starting to process messages")

	fmt.Println(params)
	if err := m.Decode(&params); err != nil {
		return m.ErrorResponse(ctx, err)
	}

	if params.Type == "" {
		return m.ErrorResponse(ctx, errors.New("missing payload type"))
	}

	c.logger.Info(ctx).Msg("processing the following event %s", params.Type)

	go c.extend(ctx, m)

	if err := c.handler(ctx, params.Message); err != nil {
		c.logger.Error(ctx).Msg("error processing processed %s", params.Type)
		return m.ErrorResponse(ctx, err)
	}

	c.logger.Info(ctx).Msg("message processed %s", params.Type)
	m.Success(ctx)

	return c.delete(ctx, m)
}

func (c *consumer) delete(ctx context.Context, m *message) error {
	_, err := c.sqs.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: &c.QueueURL, ReceiptHandle: m.ReceiptHandle})
	if err != nil {
		c.logger.Error(ctx).Msg("unable to delete %s", err.Error())
		return errors.New("unable to delete")
	}
	return nil
}

func (c *consumer) extend(ctx context.Context, m *message) {
	var count int
	extension := int32(c.VisibilityTimeout)

	for {
		if count >= c.extensionLimit {
			c.logger.Error(ctx).Msg("cannot extend a limit")
			return
		}

		count++
		time.Sleep(time.Duration(c.VisibilityTimeout-10) * time.Second)
		select {
		case <-m.err:
			return
		default:
			extension = extension + int32(c.VisibilityTimeout)
			_, err := c.sqs.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{QueueUrl: &c.QueueURL, ReceiptHandle: m.ReceiptHandle, VisibilityTimeout: extension})
			if err != nil {
				c.logger.Error(ctx).Msg("cannot extend a limit %s", err.Error())
				return
			}
		}
	}
}

type PublisherInput struct {
	Version uint64 `json:"version"`
	Data    any    `json:"data"`
}

type AwsSqsPublisher interface {
	Publish(ctx context.Context, message PublishMessage) error
}

type publisher struct {
	sqs               *sqs.Client
	env               string
	Hostname          string
	VisibilityTimeout int
	logger            *logger.Logger
	queueUrl          string
}

func NewPublisher(logger *logger.Logger, c Config) (AwsSqsPublisher, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	if c.AwsEndpoint != "" {
		cfg.BaseEndpoint = &c.AwsEndpoint
	}

	pub := publisher{
		sqs:               sqs.NewFromConfig(cfg),
		env:               c.Env,
		VisibilityTimeout: 30,
		queueUrl:          c.QueueURL,
		logger:            logger,
	}

	if c.VisibilityTimeout != 0 {
		pub.VisibilityTimeout = c.VisibilityTimeout
	}

	return pub, nil
}

func (c publisher) Publish(ctx context.Context, message PublishMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	params := sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(c.queueUrl),
	}

	out, err := c.sqs.SendMessage(ctx, &params)
	if err != nil {
		c.logger.Error(ctx).Msg("error sending message to queue %s", c.queueUrl)
		return err
	}

	c.logger.Debug(ctx).Msg("message sent to queue %s", *out.MessageId)

	return nil
}
