package sql

import (
	"context"
	"fmt"

	"github.com/Macaquit0/Tropical-BFF/pkg/errors"
	"github.com/Macaquit0/Tropical-BFF/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryRow struct {
	log    *logger.Logger
	db     *pgxpool.Pool
	err    error
	row    pgx.Row
	method string
}

func NewQueryRow(log *logger.Logger, db *pgxpool.Pool) *QueryRow {
	return &QueryRow{log, db, nil, nil, ""}
}

func (s *QueryRow) QueryOne(ctx context.Context, method, query string, args ...interface{}) *QueryRow {
	s.method = method
	row := s.db.QueryRow(ctx, query, args...)

	s.row = row
	return s
}

func (s *QueryRow) Parse(notFoundError string, args ...interface{}) error {
	if s.err != nil {
		return s.err
	}

	if err := s.row.Scan(args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.NewNotFoundError(notFoundError)
		}
		return errors.NewInternalServerError(fmt.Sprintf("error on '%s': scan. error: %s", s.method, err.Error()))
	}
	return nil
}
