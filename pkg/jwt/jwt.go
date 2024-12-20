package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Macaquit0/Tropical-BFF/pkg/errors"
	"github.com/Macaquit0/Tropical-BFF/pkg/logger"

	"github.com/golang-jwt/jwt"
)

type Jwt struct {
	log    *logger.Logger
	secret string
}

func New(log *logger.Logger, secret string) *Jwt {
	return &Jwt{log, secret}
}

func (j *Jwt) GenerateToken(ctx context.Context, claims map[string]any, expire time.Time) (string, error) {
	claims["exp"] = expire.Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	tokenClaims := token.Claims.(jwt.MapClaims)
	for key, value := range claims {
		tokenClaims[key] = value
	}

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *Jwt) Decode(ctx context.Context, tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := make(map[string]any)
	if c, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for key, value := range c {
			claims[key] = value
		}
	} else {
		return nil, errors.NewUnauthorizedError()
	}

	return claims, nil
}

func (j *Jwt) DecodeOpen(ctx context.Context, token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %w", err)
	}

	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling payload: %w", err)
	}

	return claims, nil
}
