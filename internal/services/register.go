package services

import (
	"context"

	errors "github.com/Macaquit0/Tropical-BFF/pkg/errors"
	"github.com/Macaquit0/Tropical-BFF/pkg/validate"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CreateUserRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
}

func (s *Service) CreateUser(ctx context.Context, request CreateUserRequestParams) (CreateUserResponse, error) {
	// Validate the request
	if err := validate.Validate(request, "error on user registration validation"); err != nil {
		return CreateUserResponse{}, err
	}

	// Prepare Cognito input
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: &s.AppClientID,
		Username: &request.Email,
		Password: &request.Password,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(request.Name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(request.Email),
			},
		},
	}

	// Call Cognito
	_, err := s.CognitoClient.SignUp(ctx, input)
	if err != nil {
		return CreateUserResponse{}, errors.NewInternalServerError("failed to register user")
	}

	return CreateUserResponse{
		Message: "User successfully registered",
	}, nil
}
