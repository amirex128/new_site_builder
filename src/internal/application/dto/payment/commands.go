package payment

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateOrUpdateGatewayCommand represents a command to create or update payment gateways
type CreateOrUpdateGatewayCommand struct {
	SiteID                *int64               `json:"siteId" validate:"required"`
	Saman                 *SamanGateway        `json:"saman,omitempty" validate:"required_if=IsActiveSaman 1"`
	IsActiveSaman         *enums.StatusEnum    `json:"isActiveSaman,omitempty" validate:"omitempty,enum"`
	Mellat                *MellatGateway       `json:"mellat,omitempty" validate:"required_if=IsActiveMellat 1"`
	IsActiveMellat        *enums.StatusEnum    `json:"isActiveMellat,omitempty" validate:"omitempty,enum"`
	Parsian               *ParsianGateway      `json:"parsian,omitempty" validate:"required_if=IsActiveParsian 1"`
	IsActiveParsian       *enums.StatusEnum    `json:"isActiveParsian,omitempty" validate:"omitempty,enum"`
	Pasargad              *PasargadGateway     `json:"pasargad,omitempty" validate:"required_if=IsActivePasargad 1"`
	IsActivePasargad      *enums.StatusEnum    `json:"isActivePasargad,omitempty" validate:"omitempty,enum"`
	IranKish              *IranKishGateway     `json:"iranKish,omitempty" validate:"required_if=IsActiveIranKish 1"`
	IsActiveIranKish      *enums.StatusEnum    `json:"isActiveIranKish,omitempty" validate:"omitempty,enum"`
	Melli                 *MelliGateway        `json:"melli,omitempty" validate:"required_if=IsActiveMelli 1"`
	IsActiveMelli         *enums.StatusEnum    `json:"isActiveMelli,omitempty" validate:"omitempty,enum"`
	AsanPardakht          *AsanPardakhtGateway `json:"asanPardakht,omitempty" validate:"required_if=IsActiveAsanPardakht 1"`
	IsActiveAsanPardakht  *enums.StatusEnum    `json:"isActiveAsanPardakht,omitempty" validate:"omitempty,enum"`
	Sepehr                *SepehrGateway       `json:"sepehr,omitempty" validate:"required_if=IsActiveSepehr 1"`
	IsActiveSepehr        *enums.StatusEnum    `json:"isActiveSepehr,omitempty" validate:"omitempty,enum"`
	ZarinPal              *ZarinPalGateway     `json:"zarinPal,omitempty" validate:"required_if=IsActiveZarinPal 1"`
	IsActiveZarinPal      *enums.StatusEnum    `json:"isActiveZarinPal,omitempty" validate:"omitempty,enum"`
	PayIr                 *PayIrGateway        `json:"payIr,omitempty" validate:"required_if=IsActivePayIr 1"`
	IsActivePayIr         *enums.StatusEnum    `json:"isActivePayIr,omitempty" validate:"omitempty,enum"`
	IdPay                 *IdPayGateway        `json:"idPay,omitempty" validate:"required_if=IsActiveIdPay 1"`
	IsActiveIdPay         *enums.StatusEnum    `json:"isActiveIdPay,omitempty" validate:"omitempty,enum"`
	YekPay                *YekPayGateway       `json:"yekPay,omitempty" validate:"required_if=IsActiveYekPay 1"`
	IsActiveYekPay        *enums.StatusEnum    `json:"isActiveYekPay,omitempty" validate:"omitempty,enum"`
	PayPing               *PayPingGateway      `json:"payPing,omitempty" validate:"required_if=IsActivePayPing 1"`
	IsActivePayPing       *enums.StatusEnum    `json:"isActivePayPing,omitempty" validate:"omitempty,enum"`
	IsActiveParbadVirtual *enums.StatusEnum    `json:"isActiveParbadVirtual,omitempty" validate:"omitempty,enum"`
}

// RequestGatewayCommand represents a command to request payment gateway
type RequestGatewayCommand struct {
	SiteID        *int64                           `json:"siteId" validate:"required"`
	Amount        *int64                           `json:"amount" validate:"required"`
	ServiceName   *string                          `json:"serviceName" validate:"required_text=1,100"`
	ServiceAction *string                          `json:"serviceAction" validate:"required_text=1,100"`
	OrderID       *int64                           `json:"orderId" validate:"required"`
	ReturnURL     *string                          `json:"returnUrl" validate:"required_text=1,500"`
	CallVerifyURL *enums.VerifyPaymentEndpointEnum `json:"callVerifyUrl" validate:"required,enum"`
	ClientIP      *string                          `json:"clientIp" validate:"required_text=1,50"`
	Gateway       *enums.PaymentGatewaysEnum       `json:"gateway" validate:"required,enum"`
	UserType      *enums.UserTypeEnum              `json:"userType" validate:"required,enum"`
	UserID        *int64                           `json:"userId" validate:"required"`
	OrderData     map[string]string                `json:"orderData" validate:"required,min=1"`
}

// VerifyPaymentCommand represents a command to verify payment
type VerifyPaymentCommand struct {
	TransactionCode *string `json:"transactionCode" validate:"required_text=1,100"`
	Result          *string `json:"result" validate:"required_text=1,100"`
	PaymentToken    *string `json:"paymentToken" validate:"required_text=1,100"`
}
