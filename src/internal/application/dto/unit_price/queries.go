package unit_price

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// CalculateUnitPriceUnitPriceQuery represents a nested query for unit price calculation
type CalculateUnitPriceUnitPriceQuery struct {
	UnitPriceName  *UnitPriceNameEnum `json:"unitPriceName" validate:"required" error:"required=نام واحد قیمت الزامی است"`
	UnitPriceCount *int               `json:"unitPriceCount" validate:"required,gt=0" error:"required=تعداد واحد قیمت الزامی است|gt=تعداد واحد قیمت باید بزرگتر از 0 باشد"`
	UnitPriceDay   *int               `json:"unitPriceDay,omitempty" validate:"required_if=UnitPriceName StorageMbCredits" error:"required_if=تعداد روز الزامی است"`
}

// CalculateUnitPriceQuery represents a query to calculate unit prices
type CalculateUnitPriceQuery struct {
	UnitPrices []CalculateUnitPriceUnitPriceQuery `json:"unitPrices" validate:"required,min=1" error:"required=لیست واحدهای قیمت نمی‌تواند خالی باشد|min=لیست واحدهای قیمت نباید خالی باشد"`
}

// GetAllUnitPriceQuery represents a query to get all unit prices
type GetAllUnitPriceQuery struct {
	common.PaginationRequestDto
}
