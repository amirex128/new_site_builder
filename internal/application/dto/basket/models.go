package basket

// BasketItemCommand represents an item in a basket
type BasketItemCommand struct {
	BasketItemID     *int64 `json:"basketItemId" nameFa:"شناسه آیتم سبد" validate:"gt=0"`
	Quantity         *int   `json:"quantity" nameFa:"تعداد" validate:"required,max=1000"`
	ProductID        *int64 `json:"productId" nameFa:"شناسه محصول" validate:"required"`
	ProductVariantID *int64 `json:"productVariantId,omitempty" nameFa:"شناسه ویژگی محصول" validate:"omitempty"`
}
