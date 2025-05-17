package payment

// Gateway models for different payment providers

// SamanGateway represents Saman payment gateway configuration
type SamanGateway struct {
	MerchantID *string `json:"merchantId" validate:"required_text=1,100" error:"شناسه پذیرنده سامان الزامی است"`
	Password   *string `json:"password" validate:"required_text=1,100" error:"رمز عبور سامان الزامی است"`
}

// MellatGateway represents Mellat payment gateway configuration
type MellatGateway struct {
	TerminalID   *int64  `json:"terminalId" validate:"required" error:"شناسه ترمینال ملت الزامی است"`
	UserName     *string `json:"userName" validate:"required_text=1,100" error:"نام کاربری ملت الزامی است"`
	UserPassword *string `json:"userPassword" validate:"required_text=1,100" error:"رمز عبور ملت الزامی است"`
}

// ParsianGateway represents Parsian payment gateway configuration
type ParsianGateway struct {
	LoginAccount *string `json:"loginAccount" validate:"required_text=1,100" error:"حساب کاربری پارسیان الزامی است"`
}

// PasargadGateway represents Pasargad payment gateway configuration
type PasargadGateway struct {
	MerchantCode *string `json:"merchantCode" validate:"required_text=1,100" error:"کد پذیرنده پاسارگاد الزامی است"`
	TerminalCode *string `json:"terminalCode" validate:"required_text=1,100" error:"کد ترمینال پاسارگاد الزامی است"`
	PrivateKey   *string `json:"privateKey" validate:"required_text=1,1000" error:"کلید خصوصی پاسارگاد الزامی است"`
}

// IranKishGateway represents IranKish payment gateway configuration
type IranKishGateway struct {
	TerminalID *string `json:"terminalId" validate:"required_text=1,100" error:"شناسه ترمینال ایران کیش الزامی است"`
	AcceptorID *string `json:"acceptorId" validate:"required_text=1,100" error:"شناسه پذیرنده ایران کیش الزامی است"`
	PassPhrase *string `json:"passPhrase" validate:"required_text=1,100" error:"عبارت عبور ایران کیش الزامی است"`
	PublicKey  *string `json:"publicKey" validate:"required_text=1,1000" error:"کلید عمومی ایران کیش الزامی است"`
}

// MelliGateway represents Melli payment gateway configuration
type MelliGateway struct {
	TerminalID  *string `json:"terminalId" validate:"required_text=1,100" error:"شناسه ترمینال ملی الزامی است"`
	MerchantID  *string `json:"merchantId" validate:"required_text=1,100" error:"شناسه پذیرنده ملی الزامی است"`
	TerminalKey *string `json:"terminalKey" validate:"required_text=1,100" error:"کلید ترمینال ملی الزامی است"`
}

// AsanPardakhtGateway represents AsanPardakht payment gateway configuration
type AsanPardakhtGateway struct {
	MerchantConfigurationID *string `json:"merchantConfigurationId" validate:"required_text=1,100" error:"شناسه تنظیمات پذیرنده آسان پرداخت الزامی است"`
	UserName                *string `json:"userName" validate:"required_text=1,100" error:"نام کاربری آسان پرداخت الزامی است"`
	Password                *string `json:"password" validate:"required_text=1,100" error:"رمز عبور آسان پرداخت الزامی است"`
	Key                     *string `json:"key" validate:"required_text=1,100" error:"کلید آسان پرداخت الزامی است"`
	IV                      *string `json:"iv" validate:"required_text=1,100" error:"بردار اولیه آسان پرداخت الزامی است"`
}

// SepehrGateway represents Sepehr payment gateway configuration
type SepehrGateway struct {
	TerminalID *int64 `json:"terminalId" validate:"required" error:"شناسه ترمینال سپهر الزامی است"`
}

// ZarinPalGateway represents ZarinPal payment gateway configuration
type ZarinPalGateway struct {
	MerchantID         *string `json:"merchantId" validate:"required_text=1,100" error:"شناسه پذیرنده زرین‌پال الزامی است"`
	AuthorizationToken *string `json:"authorizationToken,omitempty" validate:"optional_text=0,100"`
	IsSandbox          *bool   `json:"isSandbox,omitempty" validate:"optional"`
}

// PayIrGateway represents Pay.ir payment gateway configuration
type PayIrGateway struct {
	API           *string `json:"api" validate:"required_text=1,100" error:"کلید API پی‌آی‌آر الزامی است"`
	IsTestAccount *bool   `json:"isTestAccount,omitempty" validate:"optional"`
}

// IdPayGateway represents IdPay payment gateway configuration
type IdPayGateway struct {
	API           *string `json:"api" validate:"required_text=1,100" error:"کلید API آیدی‌پی الزامی است"`
	IsTestAccount *bool   `json:"isTestAccount" validate:"required" error:"محیط آزمایشی آیدی‌پی باید مشخص شود"`
}

// YekPayGateway represents YekPay payment gateway configuration
type YekPayGateway struct {
	MerchantID *string `json:"merchantId" validate:"required_text=1,100" error:"شناسه پذیرنده یک‌پی الزامی است"`
}

// PayPingGateway represents PayPing payment gateway configuration
type PayPingGateway struct {
	AccessToken *string `json:"accessToken" validate:"required_text=1,100" error:"کلید API پی‌پینگ الزامی است"`
}
