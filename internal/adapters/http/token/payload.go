package token

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

var (
	ErrInvalidToken    = errors.New("token is invalid")
	ErrExpiredToken    = errors.New("token has expired")
	ErrInvalidInstance = errors.New("paseto maker is not initialized")
)

type Payload struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	TenantSlug         string    `json:"tenant_slug"`
	SubscriptionStatus string    `json:"subscription_status"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	IssuedAt           time.Time `json:"issued_at"`
	ExpiredAt          time.Time `json:"expired_at"`
}

func NewPayload(tenantID uuid.UUID, tenantSlug string, subscriptionStatus string, email string, role string, duration time.Duration) (*Payload, *domainerr.DomainError) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.UnexpectedError, err.Error(), err)
	}

	payload := &Payload{
		ID:                 tokenID,
		TenantID:           tenantID,
		TenantSlug:         tenantSlug,
		SubscriptionStatus: subscriptionStatus,
		Email:              email,
		Role:               role,
		IssuedAt:           time.Now(),
		ExpiredAt:          time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() *domainerr.DomainError {
	if time.Now().After(payload.ExpiredAt) {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.ExpiredToken, "Expired Token", ErrExpiredToken)
	}

	return nil
}
