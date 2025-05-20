package order

import "github.com/amirex128/new_site_builder/src/internal/domain/enums"

// CreateOrderRequestCommand represents a command to create an order request
type CreateOrderRequestCommand struct {
	Gateway             *enums.PaymentGatewaysEnum `json:"gateway" nameFa:"درگاه" validate:"required,enum"`
	FinalFrontReturnURL *string                    `json:"finalFrontReturnUrl" nameFa:"آدرس بازگشت پایانی" validate:"required_text=1,500"`
	Description         *string                    `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=0 1000"`
	SiteID              *int64                     `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	AddressID           *int64                     `json:"addressId" nameFa:"شناسه آدرس" validate:"required"`
	Courier             *enums.CourierEnum         `json:"courier,omitempty" nameFa:"ارسال کننده" validate:"enum_optional"`
}

// CreateOrderVerifyCommand represents a command to verify an order
type CreateOrderVerifyCommand struct {
	PaymentStatus   *string           `json:"paymentStatus" nameFa:"وضعیت پرداخت" validate:"required_text=1,50"`
	IsSuccess       *bool             `json:"isSuccess" nameFa:"آیا موفق بوده است" validate:"required_bool"`
	OrderData       map[string]string `json:"orderData" nameFa:"اطلاعات سفارش" validate:"required"`
	TransactionCode *string           `json:"transaction_code" nameFa:"کد تراکنش"`
}
