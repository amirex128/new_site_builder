package order

// CreateOrderRequestCommand represents a command to create an order request
type CreateOrderRequestCommand struct {
	Gateway             *PaymentGatewaysEnum `json:"gateway" validate:"required,enum"`
	FinalFrontReturnURL *string              `json:"finalFrontReturnUrl" validate:"required_text=1,500"`
	Description         *string              `json:"description,omitempty" validate:"optional_text=0 1000"`
	SiteID              *int64               `json:"siteId" validate:"required"`
	AddressID           *int64               `json:"addressId" validate:"required"`
	Courier             *CourierEnum         `json:"courier,omitempty" validate:"enum_optional"`
}

// CreateOrderVerifyCommand represents a command to verify an order
type CreateOrderVerifyCommand struct {
	PaymentStatus   *string           `json:"paymentStatus" validate:"required_text=1,50"`
	IsSuccess       *bool             `json:"isSuccess" validate:"required_bool"`
	OrderData       map[string]string `json:"orderData" validate:"required"`
	TransactionCode *string           `json:"transaction_code"`
}
