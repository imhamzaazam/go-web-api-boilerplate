package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"
)

func (service *ServiceManager) CreateUser(ctx context.Context, newUser ports.NewUser) (pgsqlc.CreateUserRow, *domainerr.DomainError) {
	hashedPassword, hashErr := utils.HashPassword(newUser.Password)
	if hashErr != nil {
		return pgsqlc.CreateUserRow{}, hashErr
	}

	args := pgsqlc.CreateUserParams{
		TenantID:  newUser.TenantID,
		Role:      pgsqlc.UserRole(newUser.Role),
		Email:     newUser.Email,
		Password:  hashedPassword,
		FullName:  newUser.FullName,
		IsStaff:   false,
		IsActive:  true,
		LastLogin: time.Now(),
	}

	user, err := service.store.CreateUser(ctx, args)
	if err != nil {
		return pgsqlc.CreateUserRow{}, domainerr.MatchPostgresError(err)
	}

	return user, nil
}

func (service *ServiceManager) LoginUser(ctx context.Context, loginUser ports.LoginUser) (pgsqlc.User, *domainerr.DomainError) {
	userRow, err := service.store.GetUserByTenantAndEmail(ctx, pgsqlc.GetUserByTenantAndEmailParams{
		TenantID: loginUser.TenantID,
		Email:    loginUser.Email,
	})
	if err != nil {
		return pgsqlc.User{}, domainerr.MatchPostgresError(err)
	}

	passErr := utils.CheckPassword(loginUser.Password, userRow.Password)
	if passErr != nil {
		return pgsqlc.User{}, passErr
	}

	return pgsqlc.User{
		ID:         userRow.ID,
		TenantID:   userRow.TenantID,
		Role:       userRow.Role,
		UID:        userRow.UID,
		Email:      userRow.Email,
		Password:   userRow.Password,
		FullName:   userRow.FullName,
		IsStaff:    userRow.IsStaff,
		IsActive:   userRow.IsActive,
		LastLogin:  userRow.LastLogin,
		CreatedAt:  userRow.CreatedAt,
		ModifiedAt: userRow.ModifiedAt,
	}, nil
}

func (service *ServiceManager) CreateUserSession(ctx context.Context, newUserSession ports.NewUserSession) (pgsqlc.Session, *domainerr.DomainError) {
	session, err := service.store.CreateSession(ctx, pgsqlc.CreateSessionParams{
		TenantID:     newUserSession.TenantID,
		UID:          newUserSession.RefreshTokenID,
		UserEmail:    newUserSession.Email,
		RefreshToken: newUserSession.RefreshToken,
		ExpiresAt:    newUserSession.RefreshTokenExpiresAt,
		UserAgent:    newUserSession.UserAgent,
		ClientIP:     newUserSession.ClientIP,
	})
	if err != nil {
		return pgsqlc.Session{}, domainerr.MatchPostgresError(err)
	}

	return session, nil
}

func (service *ServiceManager) GetUserSession(ctx context.Context, tenantID uuid.UUID, refreshTokenID uuid.UUID) (pgsqlc.Session, *domainerr.DomainError) {
	session, err := service.store.GetSessionByTenant(ctx, pgsqlc.GetSessionByTenantParams{
		TenantID: tenantID,
		UID:      refreshTokenID,
	})
	if err != nil {
		return pgsqlc.Session{}, domainerr.MatchPostgresError(err)
	}

	return session, nil
}

func (service *ServiceManager) GetUserByUID(ctx context.Context, tenantID uuid.UUID, userUID string) (pgsqlc.User, *domainerr.DomainError) {
	parsedUID, err := uuid.Parse(userUID)
	if err != nil {
		return pgsqlc.User{}, domainerr.NewDomainError(500, domainerr.UnexpectedError, err.Error(), err)
	}

	userRow, err := service.store.GetUserByTenantAndUID(ctx, pgsqlc.GetUserByTenantAndUIDParams{
		TenantID: tenantID,
		UID:      parsedUID,
	})
	if err != nil {
		return pgsqlc.User{}, domainerr.MatchPostgresError(err)
	}

	return pgsqlc.User{
		ID:         userRow.ID,
		TenantID:   userRow.TenantID,
		Role:       userRow.Role,
		UID:        userRow.UID,
		Email:      userRow.Email,
		Password:   userRow.Password,
		FullName:   userRow.FullName,
		IsStaff:    userRow.IsStaff,
		IsActive:   userRow.IsActive,
		LastLogin:  userRow.LastLogin,
		CreatedAt:  userRow.CreatedAt,
		ModifiedAt: userRow.ModifiedAt,
	}, nil
}
