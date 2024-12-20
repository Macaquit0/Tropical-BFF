package postgres

import (
	"context"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

func NewConnectionPool(postgresDns string, enableTracing string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(postgresDns)
	if err != nil {
		return nil, err
	}

	if enableTracing == "true" {
		config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithIncludeQueryParameters(), otelpgx.WithTracerProvider(otel.GetTracerProvider()))
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}
