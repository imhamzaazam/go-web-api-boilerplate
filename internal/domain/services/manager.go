package services

import (
	"net/http"
	"strconv"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/jackc/pgx/v5/pgtype"
)

type ServiceManager struct {
	store pgsqlc.Querier
}

func NewServiceManager(store pgsqlc.Querier) *ServiceManager {
	return &ServiceManager{
		store: store,
	}
}

func numericFromFloat64(value float64) (pgtype.Numeric, *domainerr.DomainError) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(strconv.FormatFloat(value, 'f', -1, 64)); err != nil {
		return pgtype.Numeric{}, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.InternalError, "invalid numeric value", err)
	}

	return numeric, nil
}
