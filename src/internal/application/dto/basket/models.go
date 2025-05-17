package basket

// BasketItemCommand represents an item in a basket
type BasketItemCommand struct {
	BasketItemID     *int64 `json:"basketItemId" validate:"gt=0"`
	Quantity         *int   `json:"quantity" validate:"required,max=1000"`
	ProductID        *int64 `json:"productId" validate:"required"`
	ProductVariantID *int64 `json:"productVariantId,omitempty" validate:"omitempty"`
}
