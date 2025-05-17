package basket

// UpdateBasketCommand represents a command to update a basket
type UpdateBasketCommand struct {
	SimpleAdd   *bool               `json:"simpleAdd" validate:"required"`
	SiteID      *int64              `json:"siteId" validate:"required"`
	Code        *string             `json:"code,omitempty" validate:"optional_text=0,50"`
	BasketItems []BasketItemCommand `json:"basketItems" validate:"required,min=1,dive"`
}
