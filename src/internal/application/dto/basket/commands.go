package basket

// UpdateBasketCommand represents a command to update a basket
type UpdateBasketCommand struct {
	SimpleAdd   *bool               `json:"simpleAdd" nameFa:"افزودن ساده" validate:"required_bool"`
	SiteID      *int64              `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Code        *string             `json:"code,omitempty" nameFa:"کد" validate:"optional_text=0 50"`
	BasketItems []BasketItemCommand `json:"basketItems" nameFa:"آیتم های سبد" validate:"required,min=1,dive"`
}
