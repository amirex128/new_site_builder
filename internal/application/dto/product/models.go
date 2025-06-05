package product

import (
	enums2 "github.com/amirex128/new_site_builder/internal/domain/enums"
	"time"
)

// CouponCommand represents a coupon for a article
type CouponCommand struct {
	Quantity   *int                     `json:"quantity" validate:"required"`
	Type       *enums2.DiscountTypeEnum `json:"type" validate:"required,enum"`
	Value      *int64                   `json:"value" validate:"required"`
	ExpiryDate *time.Time               `json:"expiryDate" validate:"required,gtfield=time.Now"`
}

// ProductAttributeCommand represents a article attribute
type ProductAttributeCommand struct {
	ID    *int64                           `json:"id,omitempty" validate:"omitempty"`
	Type  *enums2.ProductAttributeTypeEnum `json:"type" validate:"required,enum"`
	Name  *string                          `json:"name" validate:"required_text=1 100"`
	Value *string                          `json:"value" validate:"required_text=1 500"`
}

// ProductVariantCommand represents a article variant
type ProductVariantCommand struct {
	ID    *int64  `json:"id,omitempty" validate:"omitempty"`
	Name  *string `json:"name" validate:"required_text=1 200"`
	Price *int64  `json:"price" validate:"required"`
	Stock *int    `json:"stock" validate:"required"`
}

// OrderBasketItem represents an item in an order basket for price calculation
type OrderBasketItem struct {
	BasketItemID       *int64 `json:"basketItemId" validate:"required"`
	Quantity           *int   `json:"quantity" validate:"required,max=1000"`
	ProductID          *int64 `json:"productId" validate:"required"`
	ProductVariationID *int64 `json:"productVariationId" validate:"required"`
}
