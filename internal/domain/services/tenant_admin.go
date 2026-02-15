package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/jackc/pgx/v5"
)

func (service *ServiceManager) CreateTenant(ctx context.Context, tenant ports.NewTenant) (pgsqlc.CreateTenantRow, *domainerr.DomainError) {
	createdTenant, err := service.store.CreateTenant(ctx, pgsqlc.CreateTenantParams{
		Name:   tenant.Name,
		Slug:   tenant.Slug,
		Domain: tenant.Domain,
		Type:   pgsqlc.TenantType(tenant.Type),
	})
	if err != nil {
		return pgsqlc.CreateTenantRow{}, domainerr.MatchPostgresError(err)
	}

	return createdTenant, nil
}

func (service *ServiceManager) CreateBranch(ctx context.Context, branch ports.NewBranch) (pgsqlc.CreateBranchRow, *domainerr.DomainError) {
	createdBranch, err := service.store.CreateBranch(ctx, pgsqlc.CreateBranchParams{
		TenantID: branch.TenantID,
		Name:     branch.Name,
		Code:     branch.Code,
	})
	if err != nil {
		return pgsqlc.CreateBranchRow{}, domainerr.MatchPostgresError(err)
	}

	return createdBranch, nil
}

func (service *ServiceManager) CreateSubscription(ctx context.Context, subscription ports.NewSubscription) (pgsqlc.CreateSubscriptionRow, *domainerr.DomainError) {
	createdSubscription, err := service.store.CreateSubscription(ctx, pgsqlc.CreateSubscriptionParams{
		TenantID: subscription.TenantID,
		Plan:     subscription.Plan,
		Status:   pgsqlc.SubscriptionStatus(subscription.Status),
		StartsAt: subscription.StartsAt,
		EndsAt:   subscription.EndsAt,
	})
	if err != nil {
		return pgsqlc.CreateSubscriptionRow{}, domainerr.MatchPostgresError(err)
	}

	return createdSubscription, nil
}

func (service *ServiceManager) ResolveTenantContext(ctx context.Context, domain string) (ports.TenantContext, *domainerr.DomainError) {
	tenantByDomain, err := service.store.GetTenantByDomain(ctx, domain)
	if err != nil {
		return ports.TenantContext{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant host", err)
	}

	subscription, err := service.store.GetLatestSubscriptionByTenant(ctx, tenantByDomain.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.TenantContext{
				TenantID:           tenantByDomain.ID,
				TenantSlug:         tenantByDomain.Slug,
				SubscriptionStatus: "trial",
			}, nil
		}
		return ports.TenantContext{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "subscription not found for tenant", err)
	}

	status := string(subscription.Status)
	if status != "trial" && status != "active" {
		return ports.TenantContext{}, domainerr.NewDomainError(http.StatusPaymentRequired, "subscription/inactive", "Subscription is not active", errors.New("subscription inactive"))
	}

	return ports.TenantContext{
		TenantID:           tenantByDomain.ID,
		TenantSlug:         tenantByDomain.Slug,
		TenantType:         string(tenantByDomain.Type),
		SubscriptionStatus: status,
	}, nil
}
