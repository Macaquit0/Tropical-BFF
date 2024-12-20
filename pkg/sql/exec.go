package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Macaquit0/Tropical-BFF/pkg/errors"
	"github.com/Macaquit0/Tropical-BFF/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type Exec struct {
	log    *logger.Logger
	db     *pgxpool.Pool
	err    error
	row    *sql.Row
	method string
}

func NewExec(log *logger.Logger, db *pgxpool.Pool) *Exec {
	return &Exec{log, db, nil, nil, ""}
}

func (s *Exec) ExecOne(ctx context.Context, method, query string, args ...interface{}) error {
	_, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error on '%s': create statement. error: %s", method, err.Error()))
	}
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				return errors.NewDuplicatedEntryError(fmt.Sprintf("error on '%s': duplicated . error: %s", method, err.Error()))
			}
		}
		return errors.NewInternalServerError(fmt.Sprintf("error on '%s': exec. error: %s", method, err.Error()))
	}
	return nil
}
