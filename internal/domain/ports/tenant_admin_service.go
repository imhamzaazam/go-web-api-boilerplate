package ports

import (
	"context"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

type TenantAdminService interface {
	CreateTenant(ctx context.Context, tenant NewTenant) (pgsqlc.CreateTenantRow, *domainerr.DomainError)
	CreateBranch(ctx context.Context, branch NewBranch) (pgsqlc.CreateBranchRow, *domainerr.DomainError)
	CreateSubscription(ctx context.Context, subscription NewSubscription) (pgsqlc.CreateSubscriptionRow, *domainerr.DomainError)
	ResolveTenantContext(ctx context.Context, domain string) (TenantContext, *domainerr.DomainError)
}
