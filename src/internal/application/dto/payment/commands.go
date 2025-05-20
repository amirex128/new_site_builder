package payment

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateOrUpdateGatewayCommand represents a command to create or update payment gateways
type CreateOrUpdateGatewayCommand struct {
	SiteID                *int64               `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Saman                 *SamanGateway        `json:"saman,omitempty" nameFa:"سامان" validate:"required_if=IsActiveSaman 1"`
	IsActiveSaman         *enums.StatusEnum    `json:"isActiveSaman,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	Mellat                *MellatGateway       `json:"mellat,omitempty" nameFa:"ملت" validate:"required_if=IsActiveMellat 1"`
	IsActiveMellat        *enums.StatusEnum    `json:"isActiveMellat,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	Parsian               *ParsianGateway      `json:"parsian,omitempty" nameFa:"پارسیان" validate:"required_if=IsActiveParsian 1"`
	IsActiveParsian       *enums.StatusEnum    `json:"isActiveParsian,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	Pasargad              *PasargadGateway     `json:"pasargad,omitempty" nameFa:"پسارگاد" validate:"required_if=IsActivePasargad 1"`
	IsActivePasargad      *enums.StatusEnum    `json:"isActivePasargad,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	IranKish              *IranKishGateway     `json:"iranKish,omitempty" nameFa:"ایرانکیش" validate:"required_if=IsActiveIranKish 1"`
	IsActiveIranKish      *enums.StatusEnum    `json:"isActiveIranKish,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	Melli                 *MelliGateway        `json:"melli,omitempty" nameFa:"ملی" validate:"required_if=IsActiveMelli 1"`
	IsActiveMelli         *enums.StatusEnum    `json:"isActiveMelli,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	AsanPardakht          *AsanPardakhtGateway `json:"asanPardakht,omitempty" nameFa:"آسان پرداخت" validate:"required_if=IsActiveAsanPardakht 1"`
	IsActiveAsanPardakht  *enums.StatusEnum    `json:"isActiveAsanPardakht,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	Sepehr                *SepehrGateway       `json:"sepehr,omitempty" nameFa:"سپهر" validate:"required_if=IsActiveSepehr 1"`
	IsActiveSepehr        *enums.StatusEnum    `json:"isActiveSepehr,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	ZarinPal              *ZarinPalGateway     `json:"zarinPal,omitempty" nameFa:"زرین پال" validate:"required_if=IsActiveZarinPal 1"`
	IsActiveZarinPal      *enums.StatusEnum    `json:"isActiveZarinPal,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	PayIr                 *PayIrGateway        `json:"payIr,omitempty" nameFa:"پرداخت ایر" validate:"required_if=IsActivePayIr 1"`
	IsActivePayIr         *enums.StatusEnum    `json:"isActivePayIr,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	IdPay                 *IdPayGateway        `json:"idPay,omitempty" nameFa:"ایدی پی" validate:"required_if=IsActiveIdPay 1"`
	IsActiveIdPay         *enums.StatusEnum    `json:"isActiveIdPay,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	YekPay                *YekPayGateway       `json:"yekPay,omitempty" nameFa:"یکپی" validate:"required_if=IsActiveYekPay 1"`
	IsActiveYekPay        *enums.StatusEnum    `json:"isActiveYekPay,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	PayPing               *PayPingGateway      `json:"payPing,omitempty" nameFa:"پینگ" validate:"required_if=IsActivePayPing 1"`
	IsActivePayPing       *enums.StatusEnum    `json:"isActivePayPing,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
	IsActiveParbadVirtual *enums.StatusEnum    `json:"isActiveParbadVirtual,omitempty" nameFa:"فعال" validate:"omitempty,enum"`
}

// RequestGatewayCommand represents a command to request payment gateway
type RequestGatewayCommand struct {
	SiteID        *int64                           `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Amount        *int64                           `json:"amount" nameFa:"مبلغ" validate:"required"`
	ServiceName   *string                          `json:"serviceName" nameFa:"نام سرویس" validate:"required_text=1,100"`
	ServiceAction *string                          `json:"serviceAction" nameFa:"عملیات سرویس" validate:"required_text=1,100"`
	OrderID       *int64                           `json:"orderId" nameFa:"شناسه سفارش" validate:"required"`
	ReturnURL     *string                          `json:"returnUrl" nameFa:"آدرس بازگشت" validate:"required_text=1,500"`
	CallVerifyURL *enums.VerifyPaymentEndpointEnum `json:"callVerifyUrl" nameFa:"آدرس فراخوانی اعتبارسنجی" validate:"required,enum"`
	ClientIP      *string                          `json:"clientIp" nameFa:"آدرس کلاینت" validate:"required_text=1,50"`
	Gateway       *enums.PaymentGatewaysEnum       `json:"gateway" nameFa:"درگاه" validate:"required,enum"`
	UserType      *enums.UserTypeEnum              `json:"userType" nameFa:"نوع کاربر" validate:"required,enum"`
	UserID        *int64                           `json:"userId" nameFa:"شناسه کاربر" validate:"required"`
	OrderData     map[string]string                `json:"orderData" nameFa:"داده سفارش" validate:"required,min=1"`
}

// VerifyPaymentCommand represents a command to verify payment
type VerifyPaymentCommand struct {
	TransactionCode *string `json:"transactionCode" nameFa:"کد تراکنش" validate:"required_text=1,100"`
	Result          *string `json:"result" nameFa:"نتیجه" validate:"required_text=1,100"`
	PaymentToken    *string `json:"paymentToken" nameFa:"توکن پرداخت" validate:"required_text=1,100"`
}
