package unit_price

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CalculateUnitPriceUnitPriceQuery represents a nested query for unit price calculation
type CalculateUnitPriceUnitPriceQuery struct {
	UnitPriceName  *enums.UnitPriceNameEnum `json:"unitPriceName" form:"unitPriceName" validate:"required"`
	UnitPriceCount *int                     `json:"unitPriceCount" form:"unitPriceCount" validate:"required,gt=0"`
	UnitPriceDay   *int                     `json:"unitPriceDay,omitempty" form:"unitPriceDay" validate:"required_if=UnitPriceName StorageMbCredits"`
}

// CalculateUnitPriceQuery represents a query to calculate unit prices
type CalculateUnitPriceQuery struct {
	UnitPrices []CalculateUnitPriceUnitPriceQuery `json:"unitPrices" form:"unitPrices" validate:"required,min=1"`
}

// GetAllUnitPriceQuery represents a query to get all unit prices
type GetAllUnitPriceQuery struct {
	common.PaginationRequestDto
}
