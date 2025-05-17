package payment

// CreateOrUpdateGatewayCommand represents a command to create or update payment gateways
type CreateOrUpdateGatewayCommand struct {
	SiteID                *int64               `json:"siteId" validate:"required" error:"سایت الزامی است"`
	Saman                 *SamanGateway        `json:"saman,omitempty" validate:"required_if=IsActiveSaman 1" error:"اطلاعات درگاه سامان هنگام فعال بودن الزامی است"`
	IsActiveSaman         *StatusEnum          `json:"isActiveSaman,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه سامان نامعتبر است"`
	Mellat                *MellatGateway       `json:"mellat,omitempty" validate:"required_if=IsActiveMellat 1" error:"اطلاعات درگاه ملت هنگام فعال بودن الزامی است"`
	IsActiveMellat        *StatusEnum          `json:"isActiveMellat,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه ملت نامعتبر است"`
	Parsian               *ParsianGateway      `json:"parsian,omitempty" validate:"required_if=IsActiveParsian 1" error:"اطلاعات درگاه پارسیان هنگام فعال بودن الزامی است"`
	IsActiveParsian       *StatusEnum          `json:"isActiveParsian,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه پارسیان نامعتبر است"`
	Pasargad              *PasargadGateway     `json:"pasargad,omitempty" validate:"required_if=IsActivePasargad 1" error:"اطلاعات درگاه پاسارگاد هنگام فعال بودن الزامی است"`
	IsActivePasargad      *StatusEnum          `json:"isActivePasargad,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه پاسارگاد نامعتبر است"`
	IranKish              *IranKishGateway     `json:"iranKish,omitempty" validate:"required_if=IsActiveIranKish 1" error:"اطلاعات درگاه ایران کیش هنگام فعال بودن الزامی است"`
	IsActiveIranKish      *StatusEnum          `json:"isActiveIranKish,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه ایران کیش نامعتبر است"`
	Melli                 *MelliGateway        `json:"melli,omitempty" validate:"required_if=IsActiveMelli 1" error:"اطلاعات درگاه ملی هنگام فعال بودن الزامی است"`
	IsActiveMelli         *StatusEnum          `json:"isActiveMelli,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه ملی نامعتبر است"`
	AsanPardakht          *AsanPardakhtGateway `json:"asanPardakht,omitempty" validate:"required_if=IsActiveAsanPardakht 1" error:"اطلاعات درگاه آسان پرداخت هنگام فعال بودن الزامی است"`
	IsActiveAsanPardakht  *StatusEnum          `json:"isActiveAsanPardakht,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه آسان پرداخت نامعتبر است"`
	Sepehr                *SepehrGateway       `json:"sepehr,omitempty" validate:"required_if=IsActiveSepehr 1" error:"اطلاعات درگاه سپهر هنگام فعال بودن الزامی است"`
	IsActiveSepehr        *StatusEnum          `json:"isActiveSepehr,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه سپهر نامعتبر است"`
	ZarinPal              *ZarinPalGateway     `json:"zarinPal,omitempty" validate:"required_if=IsActiveZarinPal 1" error:"اطلاعات درگاه زرین‌پال هنگام فعال بودن الزامی است"`
	IsActiveZarinPal      *StatusEnum          `json:"isActiveZarinPal,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه زرین‌پال نامعتبر است"`
	PayIr                 *PayIrGateway        `json:"payIr,omitempty" validate:"required_if=IsActivePayIr 1" error:"اطلاعات درگاه پی‌آی‌آر هنگام فعال بودن الزامی است"`
	IsActivePayIr         *StatusEnum          `json:"isActivePayIr,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه پی‌آی‌آر نامعتبر است"`
	IdPay                 *IdPayGateway        `json:"idPay,omitempty" validate:"required_if=IsActiveIdPay 1" error:"اطلاعات درگاه آیدی‌پی هنگام فعال بودن الزامی است"`
	IsActiveIdPay         *StatusEnum          `json:"isActiveIdPay,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه آیدی‌پی نامعتبر است"`
	YekPay                *YekPayGateway       `json:"yekPay,omitempty" validate:"required_if=IsActiveYekPay 1" error:"اطلاعات درگاه یک‌پی هنگام فعال بودن الزامی است"`
	IsActiveYekPay        *StatusEnum          `json:"isActiveYekPay,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه یک‌پی نامعتبر است"`
	PayPing               *PayPingGateway      `json:"payPing,omitempty" validate:"required_if=IsActivePayPing 1" error:"اطلاعات درگاه پی‌پینگ هنگام فعال بودن الزامی است"`
	IsActivePayPing       *StatusEnum          `json:"isActivePayPing,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه پی‌پینگ نامعتبر است"`
	IsActiveParbadVirtual *StatusEnum          `json:"isActiveParbadVirtual,omitempty" validate:"omitempty,enum" error:"وضعیت درگاه مجازی پرباد نامعتبر است"`
}

// RequestGatewayCommand represents a command to request payment gateway
type RequestGatewayCommand struct {
	SiteID        *int64                     `json:"siteId" validate:"required" error:"سایت الزامی است"`
	Amount        *int64                     `json:"amount" validate:"required" error:"مبلغ الزامی است"`
	ServiceName   *string                    `json:"serviceName" validate:"required_text=1,100" error:"نام سرویس الزامی است"`
	ServiceAction *string                    `json:"serviceAction" validate:"required_text=1,100" error:"عملیات سرویس الزامی است"`
	OrderID       *int64                     `json:"orderId" validate:"required" error:"سفارش الزامی است"`
	ReturnURL     *string                    `json:"returnUrl" validate:"required_text=1,500" error:"آدرس بازگشت الزامی است"`
	CallVerifyURL *VerifyPaymentEndpointEnum `json:"callVerifyUrl" validate:"required,enum" error:"آدرس تأیید پرداخت الزامی است"`
	ClientIP      *string                    `json:"clientIp" validate:"required_text=1,50" error:"آی‌پی کاربر الزامی است"`
	Gateway       *PaymentGatewaysEnum       `json:"gateway" validate:"required,enum" error:"درگاه پرداخت الزامی است"`
	UserType      *UserTypeEnum              `json:"userType" validate:"required,enum" error:"نوع کاربر الزامی است"`
	UserID        *int64                     `json:"userId" validate:"required" error:"کاربر الزامی است"`
	OrderData     map[string]string          `json:"orderData" validate:"required,min=1" error:"اطلاعات سفارش الزامی است و نمی‌تواند خالی باشد"`
}

// VerifyPaymentCommand represents a command to verify payment
type VerifyPaymentCommand struct {
	TransactionCode *string `json:"transactionCode" validate:"required_text=1,100" error:"کد تراکنش الزامی است"`
	Result          *string `json:"result" validate:"required_text=1,100" error:"نتیجه پرداخت الزامی است"`
	PaymentToken    *string `json:"paymentToken" validate:"required_text=1,100" error:"توکن پرداخت الزامی است"`
}
