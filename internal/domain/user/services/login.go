package services

import (
	"context"

	errors "github.com/Macaquit0/Tropical-BFF/pkg/errors"
	"github.com/Macaquit0/Tropical-BFF/pkg/validate"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type LoginRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (s *Service) Login(ctx context.Context, request LoginRequestParams) (LoginResponse, error) {
	// Validate the request
	if err := validate.Validate(request, "error on login request validation"); err != nil {
		return LoginResponse{}, err
	}

	// Prepare Cognito input
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &s.AppClientID,
		AuthParameters: map[string]string{
			"USERNAME": request.Email,
			"PASSWORD": request.Password,
		},
	}

	// Call Cognito
	result, err := s.CognitoClient.InitiateAuth(ctx, input)
	if err != nil {
		return LoginResponse{}, errors.NewInternalServerError("failed to call Cognito")
	}

	// Check response
	if result.AuthenticationResult == nil || result.AuthenticationResult.IdToken == nil {
		return LoginResponse{}, errors.NewNotFoundError("authentication failed, missing authentication token")
	}

	return LoginResponse{
		Token: *result.AuthenticationResult.IdToken,
	}, nil
}
