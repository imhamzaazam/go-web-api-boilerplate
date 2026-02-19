package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

type AuthUserService interface {
	CreateUser(ctx context.Context, newUser NewUser) (pgsqlc.CreateUserRow, *domainerr.DomainError)
	LoginUser(ctx context.Context, loginUser LoginUser) (pgsqlc.User, *domainerr.DomainError)
	CreateUserSession(ctx context.Context, newUserSession NewUserSession) (pgsqlc.Session, *domainerr.DomainError)
	GetUserSession(ctx context.Context, tenantID uuid.UUID, refreshTokenID uuid.UUID) (pgsqlc.Session, *domainerr.DomainError)
	GetUserByUID(ctx context.Context, tenantID uuid.UUID, userUID string) (pgsqlc.User, *domainerr.DomainError)
}
