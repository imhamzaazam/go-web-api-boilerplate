package ports

import (
	"time"

	"github.com/google/uuid"
)

type NewUser struct {
	TenantID uuid.UUID
	Role     string
	FullName string
	Email    string
	Password string
}

type NewTenant struct {
	Name   string
	Slug   string
	Domain string
	Type   string
}

type NewBranch struct {
	TenantID uuid.UUID
	Name     string
	Code     string
}

type NewSubscription struct {
	TenantID uuid.UUID
	Plan     string
	Status   string
	StartsAt time.Time
	EndsAt   time.Time
}

type LoginUser struct {
	TenantID uuid.UUID
	Email    string
	Password string
}

type TenantContext struct {
	TenantID           uuid.UUID
	TenantSlug         string
	TenantType         string
	SubscriptionStatus string
}

type NewProduct struct {
	TenantID             uuid.UUID
	Name                 string
	Sku                  string
	Price                int64
	VATPercent           float64
	IsPreorder           bool
	MadeToOrder          bool
	RequiresPrescription bool
	AvailableForDelivery bool
	AvailableForPickup   bool
}

type NewDiscount struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	Code      string
	Name      string
	Type      string
	Value     float64
	StartsAt  time.Time
	EndsAt    time.Time
}

type NewProductAddon struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	Name      string
	Price     int64
}

type NewInventory struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	InStock   int32
}

type NewCartItem struct {
	TenantID        uuid.UUID
	UserUID         uuid.UUID
	ProductID       uuid.UUID
	Quantity        int32
	Note            string
	PrescriptionRef string
	HasNote         bool
	HasPrescription bool
}

type NewOrderLocation struct {
	AddressLine string
	City        string
	Lat         float64
	Lng         float64
}

type NewOrder struct {
	TenantID         uuid.UUID
	CartID           uuid.UUID
	PaymentMethodRef string
	FulfillmentType  string
	Location         *NewOrderLocation
}

type NewPaymentMethod struct {
	TenantID  uuid.UUID
	Type      string
	Label     string
	IsDefault bool
}

type PayOrder struct {
	TenantID         uuid.UUID
	OrderID          uuid.UUID
	PaymentMethodRef string
	Amount           int64
}

type NewUserSession struct {
	TenantID              uuid.UUID
	RefreshTokenID        uuid.UUID
	Email                 string
	RefreshToken          string
	UserAgent             string
	ClientIP              string
	RefreshTokenExpiresAt time.Time
}
