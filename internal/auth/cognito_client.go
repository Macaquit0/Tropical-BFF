package auth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoClient struct {
	client      *cognitoidentityprovider.Client
	appClientID string
}

func NewCognitoClient(region, appClientID string) *CognitoClient {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("Erro ao carregar a configuração AWS: " + err.Error())
	}

	return &CognitoClient{
		client:      cognitoidentityprovider.NewFromConfig(cfg),
		appClientID: appClientID,
	}
}

func (c *CognitoClient) LoginWithFacebook(ctx context.Context, facebookToken string) (string, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserSrpAuth,
		ClientId: aws.String(c.appClientID),
		AuthParameters: map[string]string{
			"ACCESS_TOKEN": facebookToken,
		},
	}

	result, err := c.client.InitiateAuth(ctx, input)
	if err != nil {
		return "", err
	}

	return *result.AuthenticationResult.IdToken, nil
}