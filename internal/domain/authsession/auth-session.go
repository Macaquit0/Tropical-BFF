package authsession

import (
	"time"

	"github.com/uptrace/bun"
)

type AuthSession struct {
	bun.BaseModel `bun:"table:auth_session"`
	Id            string `bun:",pk"`
	UserId     string
	ExpiresAt     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
