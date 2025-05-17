package payment

// CreateOrUpdateGatewayCommand represents a command to create or update payment gateways
type CreateOrUpdateGatewayCommand struct {
	SiteID                *int64               `json:"siteId" validate:"required"`
	Saman                 *SamanGateway        `json:"saman,omitempty" validate:"required_if=IsActiveSaman 1"`
	IsActiveSaman         *StatusEnum          `json:"isActiveSaman,omitempty" validate:"omitempty,enum"`
	Mellat                *MellatGateway       `json:"mellat,omitempty" validate:"required_if=IsActiveMellat 1"`
	IsActiveMellat        *StatusEnum          `json:"isActiveMellat,omitempty" validate:"omitempty,enum"`
	Parsian               *ParsianGateway      `json:"parsian,omitempty" validate:"required_if=IsActiveParsian 1"`
	IsActiveParsian       *StatusEnum          `json:"isActiveParsian,omitempty" validate:"omitempty,enum"`
	Pasargad              *PasargadGateway     `json:"pasargad,omitempty" validate:"required_if=IsActivePasargad 1"`
	IsActivePasargad      *StatusEnum          `json:"isActivePasargad,omitempty" validate:"omitempty,enum"`
	IranKish              *IranKishGateway     `json:"iranKish,omitempty" validate:"required_if=IsActiveIranKish 1"`
	IsActiveIranKish      *StatusEnum          `json:"isActiveIranKish,omitempty" validate:"omitempty,enum"`
	Melli                 *MelliGateway        `json:"melli,omitempty" validate:"required_if=IsActiveMelli 1"`
	IsActiveMelli         *StatusEnum          `json:"isActiveMelli,omitempty" validate:"omitempty,enum"`
	AsanPardakht          *AsanPardakhtGateway `json:"asanPardakht,omitempty" validate:"required_if=IsActiveAsanPardakht 1"`
	IsActiveAsanPardakht  *StatusEnum          `json:"isActiveAsanPardakht,omitempty" validate:"omitempty,enum"`
	Sepehr                *SepehrGateway       `json:"sepehr,omitempty" validate:"required_if=IsActiveSepehr 1"`
	IsActiveSepehr        *StatusEnum          `json:"isActiveSepehr,omitempty" validate:"omitempty,enum"`
	ZarinPal              *ZarinPalGateway     `json:"zarinPal,omitempty" validate:"required_if=IsActiveZarinPal 1"`
	IsActiveZarinPal      *StatusEnum          `json:"isActiveZarinPal,omitempty" validate:"omitempty,enum"`
	PayIr                 *PayIrGateway        `json:"payIr,omitempty" validate:"required_if=IsActivePayIr 1"`
	IsActivePayIr         *StatusEnum          `json:"isActivePayIr,omitempty" validate:"omitempty,enum"`
	IdPay                 *IdPayGateway        `json:"idPay,omitempty" validate:"required_if=IsActiveIdPay 1"`
	IsActiveIdPay         *StatusEnum          `json:"isActiveIdPay,omitempty" validate:"omitempty,enum"`
	YekPay                *YekPayGateway       `json:"yekPay,omitempty" validate:"required_if=IsActiveYekPay 1"`
	IsActiveYekPay        *StatusEnum          `json:"isActiveYekPay,omitempty" validate:"omitempty,enum"`
	PayPing               *PayPingGateway      `json:"payPing,omitempty" validate:"required_if=IsActivePayPing 1"`
	IsActivePayPing       *StatusEnum          `json:"isActivePayPing,omitempty" validate:"omitempty,enum"`
	IsActiveParbadVirtual *StatusEnum          `json:"isActiveParbadVirtual,omitempty" validate:"omitempty,enum"`
}

// RequestGatewayCommand represents a command to request payment gateway
type RequestGatewayCommand struct {
	SiteID        *int64                     `json:"siteId" validate:"required"`
	Amount        *int64                     `json:"amount" validate:"required"`
	ServiceName   *string                    `json:"serviceName" validate:"required_text=1,100"`
	ServiceAction *string                    `json:"serviceAction" validate:"required_text=1,100"`
	OrderID       *int64                     `json:"orderId" validate:"required"`
	ReturnURL     *string                    `json:"returnUrl" validate:"required_text=1,500"`
	CallVerifyURL *VerifyPaymentEndpointEnum `json:"callVerifyUrl" validate:"required,enum"`
	ClientIP      *string                    `json:"clientIp" validate:"required_text=1,50"`
	Gateway       *PaymentGatewaysEnum       `json:"gateway" validate:"required,enum"`
	UserType      *UserTypeEnum              `json:"userType" validate:"required,enum"`
	UserID        *int64                     `json:"userId" validate:"required"`
	OrderData     map[string]string          `json:"orderData" validate:"required,min=1"`
}

// VerifyPaymentCommand represents a command to verify payment
type VerifyPaymentCommand struct {
	TransactionCode *string `json:"transactionCode" validate:"required_text=1,100"`
	Result          *string `json:"result" validate:"required_text=1,100"`
	PaymentToken    *string `json:"paymentToken" validate:"required_text=1,100"`
}
