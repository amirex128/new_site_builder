package payment

// Gateway models for different payment providers

// SamanGateway represents Saman payment gateway configuration
type SamanGateway struct {
	MerchantID *string `json:"merchantId" validate:"required_text=1 100" nameFa:"merchantId"`
	Password   *string `json:"password" validate:"required_text=1 100" nameFa:"password"`
}

// MellatGateway represents Mellat payment gateway configuration
type MellatGateway struct {
	TerminalID   *int64  `json:"terminalId" validate:"required" nameFa:"terminalId"`
	UserName     *string `json:"userName" validate:"required_text=1 100" nameFa:"userName"`
	UserPassword *string `json:"userPassword" validate:"required_text=1 100" nameFa:"userPassword"`
}

// ParsianGateway represents Parsian payment gateway configuration
type ParsianGateway struct {
	LoginAccount *string `json:"loginAccount" validate:"required_text=1 100" nameFa:"loginAccount"`
}

// PasargadGateway represents Pasargad payment gateway configuration
type PasargadGateway struct {
	MerchantCode *string `json:"merchantCode" validate:"required_text=1 100" nameFa:"merchantCode"`
	TerminalCode *string `json:"terminalCode" validate:"required_text=1 100" nameFa:"terminalCode"`
	PrivateKey   *string `json:"privateKey" validate:"required_text=1 1000" nameFa:"privateKey"`
}

// IranKishGateway represents IranKish payment gateway configuration
type IranKishGateway struct {
	TerminalID *string `json:"terminalId" validate:"required_text=1 100" nameFa:"terminalId"`
	AcceptorID *string `json:"acceptorId" validate:"required_text=1 100" nameFa:"acceptorId"`
	PassPhrase *string `json:"passPhrase" validate:"required_text=1 100" nameFa:"passPhrase"`
	PublicKey  *string `json:"publicKey" validate:"required_text=1 1000" nameFa:"publicKey"`
}

// MelliGateway represents Melli payment gateway configuration
type MelliGateway struct {
	TerminalID  *string `json:"terminalId" validate:"required_text=1 100" nameFa:"terminalId"`
	MerchantID  *string `json:"merchantId" validate:"required_text=1 100" nameFa:"merchantId"`
	TerminalKey *string `json:"terminalKey" validate:"required_text=1 100" nameFa:"terminalKey"`
}

// AsanPardakhtGateway represents AsanPardakht payment gateway configuration
type AsanPardakhtGateway struct {
	MerchantConfigurationID *string `json:"merchantConfigurationId" validate:"required_text=1 100" nameFa:"merchantConfigurationId"`
	UserName                *string `json:"userName" validate:"required_text=1 100" nameFa:"userName"`
	Password                *string `json:"password" validate:"required_text=1 100" nameFa:"password"`
	Key                     *string `json:"key" validate:"required_text=1 100" nameFa:"key"`
	IV                      *string `json:"iv" validate:"required_text=1 100" nameFa:"iv"`
}

// SepehrGateway represents Sepehr payment gateway configuration
type SepehrGateway struct {
	TerminalID *int64 `json:"terminalId" validate:"required" nameFa:"terminalId"`
}

// ZarinPalGateway represents ZarinPal payment gateway configuration
type ZarinPalGateway struct {
	MerchantID         *string `json:"merchantId" validate:"required_text=1 100" nameFa:"merchantId"`
	AuthorizationToken *string `json:"authorizationToken,omitempty" validate:"optional_text=0 100" nameFa:"authorizationToken"`
	IsSandbox          *bool   `json:"isSandbox,omitempty" validate:"optional" nameFa:"isSandbox"`
}

// PayIrGateway represents Pay.ir payment gateway configuration
type PayIrGateway struct {
	API           *string `json:"api" validate:"required_text=1 100" nameFa:"api"`
	IsTestAccount *bool   `json:"isTestAccount,omitempty" validate:"optional" nameFa:"isTestAccount"`
}

// IdPayGateway represents IdPay payment gateway configuration
type IdPayGateway struct {
	API           *string `json:"api" validate:"required_text=1 100" nameFa:"api"`
	IsTestAccount *bool   `json:"isTestAccount" validate:"required" nameFa:"isTestAccount"`
}

// YekPayGateway represents YekPay payment gateway configuration
type YekPayGateway struct {
	MerchantID *string `json:"merchantId" validate:"required_text=1 100" nameFa:"merchantId"`
}

// PayPingGateway represents PayPing payment gateway configuration
type PayPingGateway struct {
	AccessToken *string `json:"accessToken" validate:"required_text=1 100" nameFa:"accessToken"`
}
