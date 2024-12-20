package sql

import (
	"context"
	"fmt"

	"github.com/backend/bff-cognito/pkg/errors"
	"github.com/backend/bff-cognito/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Query struct {
	log    *logger.Logger
	db     *pgxpool.Pool
	err    error
	method string
}

func NewQuery(log *logger.Logger, db *pgxpool.Pool) *Query {
	return &Query{log, db, nil, ""}
}

func (s *Query) Query(ctx context.Context, method, query string, args ...interface{}) (pgx.Rows, error) {
	s.method = method
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		s.err = errors.NewInternalServerError(fmt.Sprintf("error on '%s': create statement. error: %s", s.method, err.Error()))
		return nil, err
	}

	return rows, nil
}
