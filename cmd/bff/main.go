package main

import (
	"context"
	"os"

	"github.com/Macaquit0/Tropical-BFF/internal/domain/authsession"
	"github.com/Macaquit0/Tropical-BFF/pkg/config"
	sharedhttp "github.com/Macaquit0/Tropical-BFF/pkg/http"
	"github.com/Macaquit0/Tropical-BFF/pkg/jwt"
	"github.com/Macaquit0/Tropical-BFF/pkg/logger"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type ApiConfig struct {
	Env             string `env:"API_ENV" envDefault:"development"`
	HttpPort        string `env:"API_HTTP_PORT" envDefault:"8080"`
	JwtSecret       string `env:"JWT_SECRET"`
	CognitoRegion   string `env:"COGNITO_REGION"`
	CognitoClientID string `env:"COGNITO_APP_CLIENT_ID"`
}

var build = "1.0.0"

const appName = "Tropical-BFF"

func main() {
	ctx := context.Background()
	log := logger.New(build, appName)

	if err := run(ctx, log); err != nil {
		log.Error(ctx).Msg("error on startup %s: %v", appName, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------- Config --------------
	cfg := ApiConfig{}
	if err := config.LoadConfigFromEnv(&cfg); err != nil {
		log.Error(ctx).Msg("error loading config: %v", err)
		return err
	}

	// -------------- AWS Cognito --------------
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.CognitoRegion))
	if err != nil {
		log.Error(ctx).Msg("error loading AWS config: %v", err)
		return err
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(awsCfg)

	// -------------- JWT --------------
	jwtService := jwt.New(log, cfg.JwtSecret)

	// -------------- HTTP Server --------------
	serverOpts := sharedhttp.ServerOpts{
		ServerPort:    cfg.HttpPort,
		ServerVersion: build,
		ServerEnv:     cfg.Env,
		ServerName:    appName,
	}
	server := sharedhttp.NewServer(log, serverOpts)

	// -------------- Services and Handlers --------------
	authSessionRepository := authsession.NewRepository(bunDb)
	authSessionService := authsession.NewService(authSessionRepository, partnerRepository, uid, jwtService)

	authHandler.Routes() // Registrar rotas relacionadas à autenticação

	// -------------- Start Server --------------
	if err := server.Init(); err != nil {
		log.Error(ctx).Msg("error initializing server: %v", err)
		return err
	}

	return nil
}