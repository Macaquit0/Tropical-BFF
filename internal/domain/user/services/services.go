package services

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Service struct {
	CognitoClient *cognitoidentityprovider.Client
	AppClientID   string
}

func New(cognitoClient *cognitoidentityprovider.Client, appClientID string) *Service {
	return &Service{
		CognitoClient: cognitoClient,
		AppClientID:   appClientID,
	}
}
