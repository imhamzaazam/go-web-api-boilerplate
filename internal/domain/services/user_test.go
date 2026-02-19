package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/stretchr/testify/require"
)

type fakeQuerierForCreateProduct struct {
	pgsqlc.Querier
	tenantType         pgsqlc.TenantType
	createProductCalls int
}

func (fake *fakeQuerierForCreateProduct) GetTenantByID(ctx context.Context, id uuid.UUID) (pgsqlc.GetTenantByIDRow, error) {
	return pgsqlc.GetTenantByIDRow{
		ID:   id,
		Type: fake.tenantType,
	}, nil
}

func (fake *fakeQuerierForCreateProduct) CreateProduct(ctx context.Context, arg pgsqlc.CreateProductParams) (pgsqlc.CreateProductRow, error) {
	fake.createProductCalls++
	return pgsqlc.CreateProductRow{
		ID:       uuid.New(),
		TenantID: arg.TenantID,
		Name:     arg.Name,
		Sku:      arg.Sku,
		Price:    arg.Price,
	}, nil
}

func Test_ServiceManager_CreateProduct_MadeToOrderRule(t *testing.T) {
	t.Run("rejects made_to_order for non-bakery tenant", func(t *testing.T) {
		store := &fakeQuerierForCreateProduct{tenantType: pgsqlc.TenantTypeRestaurant}
		service := NewServiceManager(store)

		_, err := service.CreateProduct(context.Background(), ports.NewProduct{
			TenantID:    uuid.New(),
			Name:        "Chicken Karahi",
			Sku:         "RST-KARAHI-001",
			Price:       1599,
			VATPercent:  5,
			MadeToOrder: true,
		})

		require.NotNil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, err.HTTPCode)
		require.Equal(t, 0, store.createProductCalls)
	})

	t.Run("allows made_to_order for bakery tenant", func(t *testing.T) {
		store := &fakeQuerierForCreateProduct{tenantType: pgsqlc.TenantTypeBakery}
		service := NewServiceManager(store)

		product, err := service.CreateProduct(context.Background(), ports.NewProduct{
			TenantID:    uuid.New(),
			Name:        "Chocolate Fudge Cake",
			Sku:         "BK-CAKE-001",
			Price:       1299,
			VATPercent:  15,
			MadeToOrder: true,
		})

		require.Nil(t, err)
		require.Equal(t, "BK-CAKE-001", product.Sku)
		require.Equal(t, 1, store.createProductCalls)
	})
}
