package payment

// Gateway models for different payment providers

// SamanGateway represents Saman payment gateway configuration
type SamanGateway struct {
	MerchantID *string `json:"merchantId" validate:"required_text=1 100"`
	Password   *string `json:"password" validate:"required_text=1 100"`
}

// MellatGateway represents Mellat payment gateway configuration
type MellatGateway struct {
	TerminalID   *int64  `json:"terminalId" validate:"required"`
	UserName     *string `json:"userName" validate:"required_text=1 100"`
	UserPassword *string `json:"userPassword" validate:"required_text=1 100"`
}

// ParsianGateway represents Parsian payment gateway configuration
type ParsianGateway struct {
	LoginAccount *string `json:"loginAccount" validate:"required_text=1 100"`
}

// PasargadGateway represents Pasargad payment gateway configuration
type PasargadGateway struct {
	MerchantCode *string `json:"merchantCode" validate:"required_text=1 100"`
	TerminalCode *string `json:"terminalCode" validate:"required_text=1 100"`
	PrivateKey   *string `json:"privateKey" validate:"required_text=1 1000"`
}

// IranKishGateway represents IranKish payment gateway configuration
type IranKishGateway struct {
	TerminalID *string `json:"terminalId" validate:"required_text=1 100"`
	AcceptorID *string `json:"acceptorId" validate:"required_text=1 100"`
	PassPhrase *string `json:"passPhrase" validate:"required_text=1 100"`
	PublicKey  *string `json:"publicKey" validate:"required_text=1 1000"`
}

// MelliGateway represents Melli payment gateway configuration
type MelliGateway struct {
	TerminalID  *string `json:"terminalId" validate:"required_text=1 100"`
	MerchantID  *string `json:"merchantId" validate:"required_text=1 100"`
	TerminalKey *string `json:"terminalKey" validate:"required_text=1 100"`
}

// AsanPardakhtGateway represents AsanPardakht payment gateway configuration
type AsanPardakhtGateway struct {
	MerchantConfigurationID *string `json:"merchantConfigurationId" validate:"required_text=1 100"`
	UserName                *string `json:"userName" validate:"required_text=1 100"`
	Password                *string `json:"password" validate:"required_text=1 100"`
	Key                     *string `json:"key" validate:"required_text=1 100"`
	IV                      *string `json:"iv" validate:"required_text=1 100"`
}

// SepehrGateway represents Sepehr payment gateway configuration
type SepehrGateway struct {
	TerminalID *int64 `json:"terminalId" validate:"required"`
}

// ZarinPalGateway represents ZarinPal payment gateway configuration
type ZarinPalGateway struct {
	MerchantID         *string `json:"merchantId" validate:"required_text=1 100"`
	AuthorizationToken *string `json:"authorizationToken,omitempty" validate:"optional_text=0 100"`
	IsSandbox          *bool   `json:"isSandbox,omitempty" validate:"optional"`
}

// PayIrGateway represents Pay.ir payment gateway configuration
type PayIrGateway struct {
	API           *string `json:"api" validate:"required_text=1 100"`
	IsTestAccount *bool   `json:"isTestAccount,omitempty" validate:"optional"`
}

// IdPayGateway represents IdPay payment gateway configuration
type IdPayGateway struct {
	API           *string `json:"api" validate:"required_text=1 100"`
	IsTestAccount *bool   `json:"isTestAccount" validate:"required"`
}

// YekPayGateway represents YekPay payment gateway configuration
type YekPayGateway struct {
	MerchantID *string `json:"merchantId" validate:"required_text=1 100"`
}

// PayPingGateway represents PayPing payment gateway configuration
type PayPingGateway struct {
	AccessToken *string `json:"accessToken" validate:"required_text=1 100"`
}
