package authsession

import (
	"context"
	"time"

	"github.com/Macaquit0/Tropical-BFF/pkg/fbun"
	"github.com/uptrace/bun"
)

type Repository struct {
	bun *bun.DB
}

func NewRepository(bun *bun.DB) *Repository {
	return &Repository{bun}
}

func (r *Repository) Insert(ctx context.Context, authSession AuthSession) error {
	_, err := r.bun.NewInsert().Model(&authSession).Exec(ctx)
	if err != nil {
		return fbun.HandleError(err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, authSession AuthSession) error {
	_, err := r.bun.NewUpdate().Model(&authSession).Where("id = ?", authSession.Id).Exec(ctx)
	if err != nil {
		return fbun.HandleError(err)
	}

	return nil
}

func (r *Repository) GetById(ctx context.Context, id string) (AuthSession, error) {
	var p AuthSession
	err := r.bun.NewSelect().Model(&p).Where("id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return p, fbun.HandleError(err)
	}

	return p, nil
}

func (r *Repository) RevokeByPersonId(ctx context.Context, partnerId string) error {
	_, err := r.bun.NewUpdate().Model(&AuthSession{}).Set("expires_at = ?", time.Now().UTC()).Set("updated_at = ?", time.Now().UTC()).Where("partner_id = ?", partnerId).Exec(ctx)

	if err != nil {
		return fbun.HandleError(err)
	}

	return nil
}
