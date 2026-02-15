package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

type CreateUserRequestDto struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponseDto struct {
	UID      uuid.UUID `json:"uid"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

func (adapter *HTTPAdapter) createUser(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	reqUser, err := httputils.Decode[CreateUserRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqUser)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	createdUser, err := adapter.userService.CreateUser(r.Context(), ports.NewUser{
		TenantID: tenantCtx.TenantID,
		Role:     "employee",
		FullName: reqUser.FullName,
		Email:    reqUser.Email,
		Password: reqUser.Password,
	})
	if err != nil {
		return err
	}

	err = httputils.Encode(w, r, http.StatusCreated, CreateUserResponseDto{
		UID:      createdUser.UID,
		FullName: createdUser.FullName,
		Email:    createdUser.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adapter *HTTPAdapter) getUserByUID(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	payload := r.Context().Value(middleware.KeyAuthUser).(*token.Payload)

	userID := chi.URLParam(r, "uid")

	user, serviceErr := adapter.userService.GetUserByUID(r.Context(), payload.TenantID, userID)
	if serviceErr != nil {
		return serviceErr
	}

	err := httputils.Encode(w, r, http.StatusOK, CreateUserResponseDto{
		UID:      user.UID,
		Email:    user.Email,
		FullName: user.FullName,
	})
	if err != nil {
		return err
	}

	return nil
}
