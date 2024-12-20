package internalaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_SEND_TYPE = "Transactional"
)

type AwsSnsProvider struct {
	snsClient *sns.Client
}
type SMSParams struct {
	PhoneNumber string `binding:"required" json:"phone_number"`
	Message     string `binding:"required" json:"message"`
}

func NewAWSNotificationProvider() (*AwsSnsProvider, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &AwsSnsProvider{
		snsClient: sns.NewFromConfig(cfg),
	}, nil
}

func (a *AwsSnsProvider) SendSMS(ctx context.Context, params SMSParams) (string, error) {
	publishParams := &sns.PublishInput{
		Message:     aws.String(params.Message),
		PhoneNumber: aws.String(params.PhoneNumber),
	}

	publishParams.MessageAttributes = make(map[string]types.MessageAttributeValue)
	publishParams.MessageAttributes["DefaultSMSType"] = types.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(DEFAULT_SEND_TYPE),
	}

	result, err := a.snsClient.Publish(ctx, publishParams)
	if err != nil {
		logrus.Errorf("error while sending SMS notification %v", err)
		return "", err
	}

	return *result.MessageId, nil
}
