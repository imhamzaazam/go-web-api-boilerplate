package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

type CreateTenantRequestDto struct {
	Name   string `json:"name" validate:"required"`
	Slug   string `json:"slug" validate:"required"`
	Domain string `json:"domain" validate:"required"`
	Type   string `json:"type" validate:"required,oneof=bakery pharmacy restaurant"`
}

type CreateTenantResponseDto struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Domain    string    `json:"domain"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (adapter *HTTPAdapter) createTenant(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	reqDto, err := httputils.Decode[CreateTenantRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	tenant, serviceErr := adapter.userService.CreateTenant(r.Context(), ports.NewTenant{
		Name:   reqDto.Name,
		Slug:   reqDto.Slug,
		Domain: reqDto.Domain,
		Type:   reqDto.Type,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusCreated, CreateTenantResponseDto{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Slug:      tenant.Slug,
		Domain:    tenant.Domain,
		Type:      string(tenant.Type),
		CreatedAt: tenant.CreatedAt,
	})
}

type CreateBranchRequestDto struct {
	TenantID string `json:"tenant_id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type CreateBranchResponseDto struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
}

func (adapter *HTTPAdapter) createBranch(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	reqDto, err := httputils.Decode[CreateBranchRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	tenantID, parseErr := uuid.Parse(reqDto.TenantID)
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"tenant_id": "The field must be a valid uuid"}, parseErr)
	}

	branch, serviceErr := adapter.userService.CreateBranch(r.Context(), ports.NewBranch{
		TenantID: tenantID,
		Name:     reqDto.Name,
		Code:     reqDto.Code,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusCreated, CreateBranchResponseDto{
		ID:       branch.ID,
		TenantID: branch.TenantID,
		Name:     branch.Name,
		Code:     branch.Code,
	})
}

type CreateSubscriptionRequestDto struct {
	TenantID string    `json:"tenant_id" validate:"required,uuid"`
	Plan     string    `json:"plan" validate:"required"`
	Status   string    `json:"status" validate:"required,oneof=trial active past_due suspended canceled"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type CreateSubscriptionResponseDto struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Plan     string    `json:"plan"`
	Status   string    `json:"status"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

func (adapter *HTTPAdapter) createSubscription(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	reqDto, err := httputils.Decode[CreateSubscriptionRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	tenantID, parseErr := uuid.Parse(reqDto.TenantID)
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"tenant_id": "The field must be a valid uuid"}, parseErr)
	}

	startsAt := reqDto.StartsAt
	endsAt := reqDto.EndsAt

	if startsAt.IsZero() {
		startsAt = time.Now().UTC()
	}
	if endsAt.IsZero() {
		endsAt = startsAt.Add(30 * 24 * time.Hour)
	}
	if !startsAt.Before(endsAt) {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"ends_at": "must be after starts_at"}, errors.New("invalid subscription dates"))
	}

	subscription, serviceErr := adapter.userService.CreateSubscription(r.Context(), ports.NewSubscription{
		TenantID: tenantID,
		Plan:     reqDto.Plan,
		Status:   reqDto.Status,
		StartsAt: startsAt,
		EndsAt:   endsAt,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusCreated, CreateSubscriptionResponseDto{
		ID:       subscription.ID,
		TenantID: subscription.TenantID,
		Plan:     subscription.Plan,
		Status:   string(subscription.Status),
		StartsAt: subscription.StartsAt,
		EndsAt:   subscription.EndsAt,
	})
}
