package discount

import (
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
)

// CreateDiscountCommand represents a command to create a new discount
type CreateDiscountCommand struct {
	Code       *string                   `json:"code" validate:"required_text=1 100"`
	Quantity   *int                      `json:"quantity" validate:"required,max=1000"`
	Type       *product.DiscountTypeEnum `json:"type" validate:"required,enum"`
	Value      *int64                    `json:"value" validate:"required"`
	ExpiryDate *time.Time                `json:"expiryDate" validate:"required,gtfield=time.Now"`
	SiteID     *int64                    `json:"siteId" validate:"required"`
}

// UpdateDiscountCommand represents a command to update an existing discount
type UpdateDiscountCommand struct {
	ID         *int64                    `json:"id" validate:"required"`
	Code       *string                   `json:"code,omitempty" validate:"optional_text=1 100"`
	Quantity   *int                      `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Type       *product.DiscountTypeEnum `json:"type,omitempty" validate:"enum_optional"`
	Value      *int64                    `json:"value,omitempty" validate:"omitempty,gt=0"`
	ExpiryDate *time.Time                `json:"expiryDate,omitempty" validate:"omitempty,gtfield=time.Now"`
}

// DeleteDiscountCommand represents a command to delete a discount
type DeleteDiscountCommand struct {
	ID *int64 `json:"id" validate:"required"`
}
