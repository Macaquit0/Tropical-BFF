package authsession

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/Macaquit0/Tropical-BFF/pkg/errors"
)

type Uuid interface {
	NewString() string
}

type Jwt interface {
	GenerateToken(ctx context.Context, claims map[string]any, expire time.Time) (string, error)
	Decode(ctx context.Context, tokenString string) (map[string]any, error)
}

type Repository interface {
	RevokeByPersonId(ctx context.Context, partnerId string) error
	Insert(ctx context.Context, authSession AuthSession) error
	GetById(ctx context.Context, id string) (AuthSession, error)
	Update(ctx context.Context, authSession AuthSession) error
}

type Service struct {
	repository        Repository
	uuid              Uuid
	jwt               Jwt
}
func NewService(repository *Repository, uuid Uuid, jwt Jwt) *Service {
	return &Service{
		repository:        repository,
		uuid:              uuid,
		jwt:               jwt,
	}
}

func (a *Service) GenerateToken(ctx context.Context, partnerId string) (string, error) {

	now := time.Now().UTC()
	authSession := AuthSession{
		Id:        a.uuid.NewString(),
		PartnerId:  partnerId,
		ExpiresAt: now.AddDate(0, 0, 7),
		CreatedAt: now,
		UpdatedAt: now,
	}

	tokenClaims := map[string]any{
		"id": authSession.Id,
	}

	if err := a.repository.RevokeByPersonId(ctx, partnerId); err != nil {
		return "", err
	}

	token, err := a.jwt.GenerateToken(ctx, tokenClaims, authSession.ExpiresAt)
	if err != nil {
		return "", err
	}

	if err := a.repository.Insert(ctx, authSession); err != nil {
		return "", err
	}

	return token, nil
}

func (a *Service) Authorize(ctx context.Context, token string) (AuthSession, error) {
	id, err := a.decode(ctx, token)
	if err != nil {
		return AuthSession{}, errors.NewUnauthorizedError()
	}

	authSession, err := a.repository.GetById(ctx, id)
	if err != nil {
		return AuthSession{}, errors.NewUnauthorizedError()
	}

	now := time.Now().UTC()
	if now.After(authSession.ExpiresAt) {
		return AuthSession{}, errors.NewUnauthorizedError()
	}
	return authSession, nil
}

func (a *Service) RevokeByToken(ctx context.Context, token string) error {
	id, err := a.decode(ctx, token)
	if err != nil {
		return err
	}

	authSession, err := a.repository.GetById(ctx, id)
	if err != nil {
		return errors.NewUnauthorizedError()
	}

	now := time.Now().UTC()
	authSession.UpdatedAt = now

	if err := a.repository.Update(ctx, authSession); err != nil {
		return err
	}
	return nil
}

func (a *Service) decode(ctx context.Context, token string) (string, error) {
	claims, err := a.jwt.Decode(ctx, token)
	if err != nil {
		return "", err
	}

	id, ok := claims["id"]
	if !ok {
		return "", errors.NewInternalServerError("error on retrieve auth  session id from claims")
	}

	return id.(string), nil
}

func (s *Service) PartnerLogin(ctx context.Context, email, password string) (string, error) {
	partner, err := s.partnerRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.NewValidationError("invalid credentials")
	}

	if !partner.IsActive() {
		return "", errors.NewValidationError("partner is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(partner.Password), []byte(password)); err != nil {
		return "", errors.NewValidationError("invalid credentials")
	}

	now := time.Now().UTC()
	session := AuthSession{
		Id:        s.uuid.NewString(),
		PartnerId: partner.ID,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Insert(ctx, session); err != nil {
		return "", errors.NewInternalServerError("failed to create session")
	}

	claims := map[string]interface{}{
		"id":         session.Id,
		"partner_id": partner.ID,
	}

	token, err := s.jwt.GenerateToken(ctx, claims, session.ExpiresAt)
	if err != nil {
		return "", errors.NewInternalServerError("failed to generate access token")
	}

	return token, nil
}

