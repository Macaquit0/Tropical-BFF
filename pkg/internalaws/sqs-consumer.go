package internalaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/backend/bff-cognito/pkg/logger"
)

type SqsConsumer struct {
	log                 *logger.Logger
	client              *sqs.Client
	queueUrl            string
	maxNumberOfMessages int32
	waitTimeSeconds     int32
	visibilityTimeout   int32
}

type NewConsumerParams struct {
	QueueUrl            string
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	Cfg                 Config
}

func NewConsumer(log *logger.Logger, params NewConsumerParams) (*SqsConsumer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	if params.Cfg.AwsEndpoint != "" {
		cfg.BaseEndpoint = &params.Cfg.AwsEndpoint
	}

	cfg.Region = params.Cfg.Region

	client := sqs.NewFromConfig(cfg)
	return &SqsConsumer{
		log,
		client,
		params.QueueUrl,
		params.MaxNumberOfMessages,
		params.WaitTimeSeconds,
		params.VisibilityTimeout,
	}, nil
}

func (s *SqsConsumer) Run(ctx context.Context, process func(types.Message) error) {
	go func() {
		s.log.Info(ctx).Msg("Starting SQS consumer...")
		for {
			select {
			case <-ctx.Done():
				s.log.Info(ctx).Msg("Stopping SQS consumer...")
				return
			default:
				for {

					resp, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
						QueueUrl:            aws.String(s.queueUrl),
						MaxNumberOfMessages: s.maxNumberOfMessages,
						WaitTimeSeconds:     s.waitTimeSeconds,
						VisibilityTimeout:   s.visibilityTimeout,
					})

					if err != nil {
						s.log.Info(ctx).Msg("[sqs-consumer] - error receiving messages: %v", err)
						continue
					}

					if len(resp.Messages) == 0 {
						continue
					}
					for _, message := range resp.Messages {
						if err := process(message); err != nil {
							s.log.Error(ctx).Msg("[sqs-consumer] - error on process message. error: %s", err.Error())
							continue
						}

						_, err := s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
							QueueUrl:      aws.String(s.queueUrl),
							ReceiptHandle: message.ReceiptHandle,
						})
						if err != nil {
							s.log.Error(ctx).Msg("[sqs-consumer] - error on delete message. error: %s", err.Error())
						}
					}
				}
			}
		}
	}()
}
