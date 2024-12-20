package fbun

import (
	"context"
	"database/sql"
	"errors"
	errors "github.com/backend/bff-cognito/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func NewPostgres(url string) (*bun.DB, error) {
	config, err := pgxpool.ParseConfig(url)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(pool)
	bunClient := bun.NewDB(db, pgdialect.New())

	return bunClient, nil
}

func HandleError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ferrors.NewDuplicatedEntryError("duplicated entry")
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ferrors.NewNotFoundError("not found")
	}

	return err
}
