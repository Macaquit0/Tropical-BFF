package internalaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/sirupsen/logrus"
)

type AwsSesProvider struct {
	sesClient *sesv2.Client
}

func NewAwsSesProvider(awsEndpoint string) (*AwsSesProvider, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	if awsEndpoint != "" {
		cfg.BaseEndpoint = &awsEndpoint
	}

	return &AwsSesProvider{
		sesClient: sesv2.NewFromConfig(cfg),
	}, nil
}

type EmailParams struct {
	Destination string `binding:"required" json:"destination"`
	Subject     string `binding:"required" json:"subject"`
	Body        string `binding:"required" json:"body"`
	IsHTML      bool   `json:"is_html"`
}

func (a *AwsSesProvider) SendEmail(ctx context.Context, fromEmail string, params EmailParams) (string, error) {
	content := &sesTypes.EmailContent{
		Simple: &sesTypes.Message{
			Subject: &sesTypes.Content{
				Data: aws.String(params.Subject),
			},
			Body: &sesTypes.Body{},
		},
	}

	if params.IsHTML {
		content.Simple.Body.Html = &sesTypes.Content{Data: aws.String(params.Body)}
	} else {
		content.Simple.Body.Text = &sesTypes.Content{Data: aws.String(params.Body)}
	}

	sendEmailParams := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(fromEmail),
		Destination:      &sesTypes.Destination{ToAddresses: []string{params.Destination}},
		Content:          content,
	}

	message, err := a.sesClient.SendEmail(ctx, sendEmailParams)
	if err != nil {
		logrus.Errorf("error while sending email notification %v", err)
		return "", err
	}

	return *message.MessageId, nil
}
