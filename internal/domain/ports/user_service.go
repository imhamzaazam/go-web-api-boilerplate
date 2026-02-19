package ports

type UserService interface {
	TenantAdminService
	OwnerCatalogService
	OwnerCommerceService
	AuthUserService
}
