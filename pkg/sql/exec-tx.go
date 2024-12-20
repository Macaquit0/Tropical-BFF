package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/backend/bff-cognito/pkg/errors"
	"github.com/backend/bff-cognito/pkg/logger"
	"github.com/lib/pq"
)

type ExecTx struct {
	log    *logger.Logger
	tx     *sql.Tx
	err    error
	row    *sql.Row
	method string
}

func NewExecTx(log *logger.Logger, db *sql.DB) (*ExecTx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.NewInternalServerErrorWithError("error on begin db transaction", err)
	}

	return &ExecTx{log, tx, nil, nil, ""}, nil
}

func (s *ExecTx) ExecOne(ctx context.Context, method, query string, args ...interface{}) error {
	stmt, err := s.tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error on '%s': create statement. error: %s", method, err.Error()))
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.log.Error(ctx).Msg(fmt.Sprintf("error on '%s': close statement. error: %s", method, err.Error()))
		}
	}()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				return errors.NewDuplicatedEntryError(fmt.Sprintf("error on '%s': duplicated . error: %s", method, err.Error()))
			}
		}
		return errors.NewInternalServerError(fmt.Sprintf("error on '%s': exec statement. error: %s", method, err.Error()))
	}
	return nil
}

func (s *ExecTx) Commit() error {
	return s.tx.Commit()
}

func (s *ExecTx) Rollback() error {
	return s.tx.Rollback()
}
