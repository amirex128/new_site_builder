package user

import (
	"database/sql/driver"
	"errors"
)

// VerifyTypeEnum defines verification types
type VerifyTypeEnum string

const (
	VerifyEmailType         VerifyTypeEnum = "verify_email"
	VerifyPhoneType         VerifyTypeEnum = "verify_phone"
	ForgetPasswordEmailType VerifyTypeEnum = "forget_password_email"
	ForgetPasswordPhoneType VerifyTypeEnum = "forget_password_phone"
)

func (e *VerifyTypeEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !VerifyTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = VerifyTypeEnum(b)
	return nil
}

func (e VerifyTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid VerifyTypeEnum")
	}
	return string(e), nil
}

func (e VerifyTypeEnum) IsValid() bool {
	var types = []string{
		string(VerifyEmailType),
		string(VerifyPhoneType),
		string(ForgetPasswordEmailType),
		string(ForgetPasswordPhoneType),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// UnitPriceNameEnum defines unit price types
type UnitPriceNameEnum string

const (
	StorageMbCreditsName  UnitPriceNameEnum = "storage_mb_credits"
	PageViewCreditsName   UnitPriceNameEnum = "page_view_credits"
	FormSubmitCreditsName UnitPriceNameEnum = "form_submit_credits"
	SiteCreditsName       UnitPriceNameEnum = "site_credits"
	SmsCreditsName        UnitPriceNameEnum = "sms_credits"
	EmailCreditsName      UnitPriceNameEnum = "email_credits"
	AiCreditsName         UnitPriceNameEnum = "ai_credits"
	AiImageCreditsName    UnitPriceNameEnum = "ai_image_credits"
)

func (e *UnitPriceNameEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !UnitPriceNameEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = UnitPriceNameEnum(b)
	return nil
}

func (e UnitPriceNameEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid UnitPriceNameEnum")
	}
	return string(e), nil
}

func (e UnitPriceNameEnum) IsValid() bool {
	var types = []string{
		string(StorageMbCreditsName),
		string(PageViewCreditsName),
		string(FormSubmitCreditsName),
		string(SiteCreditsName),
		string(SmsCreditsName),
		string(EmailCreditsName),
		string(AiCreditsName),
		string(AiImageCreditsName),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// PaymentGatewaysEnum defines payment gateway types
type PaymentGatewaysEnum string

const (
	ZarinPalGatewayEnum PaymentGatewaysEnum = "zarinpal"
	IdPayGatewayEnum    PaymentGatewaysEnum = "idpay"
	NextPayGateway      PaymentGatewaysEnum = "nextpay"
)

func (e *PaymentGatewaysEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !PaymentGatewaysEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = PaymentGatewaysEnum(b)
	return nil
}

func (e PaymentGatewaysEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid PaymentGatewaysEnum")
	}
	return string(e), nil
}

func (e PaymentGatewaysEnum) IsValid() bool {
	var types = []string{
		string(ZarinPalGatewayEnum),
		string(IdPayGatewayEnum),
		string(NextPayGateway),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// ServiceNameEnum defines service names
type ServiceNameEnum string

const (
	UserService    ServiceNameEnum = "user"
	OrderService   ServiceNameEnum = "order"
	PaymentService ServiceNameEnum = "payment"
)

func (e *ServiceNameEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !ServiceNameEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ServiceNameEnum(b)
	return nil
}

func (e ServiceNameEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ServiceNameEnum")
	}
	return string(e), nil
}

func (e ServiceNameEnum) IsValid() bool {
	var types = []string{
		string(UserService),
		string(OrderService),
		string(PaymentService),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// VerifyPaymentEndpointEnum defines verify payment endpoint types
type VerifyPaymentEndpointEnum string

const (
	ChargeCreditVerifyEndpoint VerifyPaymentEndpointEnum = "charge_credit_verify"
	UpgradePlanVerifyEndpoint  VerifyPaymentEndpointEnum = "upgrade_plan_verify"
)

func (e *VerifyPaymentEndpointEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !VerifyPaymentEndpointEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = VerifyPaymentEndpointEnum(b)
	return nil
}

func (e VerifyPaymentEndpointEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid VerifyPaymentEndpointEnum")
	}
	return string(e), nil
}

func (e VerifyPaymentEndpointEnum) IsValid() bool {
	var types = []string{
		string(ChargeCreditVerifyEndpoint),
		string(UpgradePlanVerifyEndpoint),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// UserTypeEnum defines user types
type UserTypeEnum string

const (
	UserTypeValue     UserTypeEnum = "user"
	CustomerTypeValue UserTypeEnum = "customer"
)

func (e *UserTypeEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !UserTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = UserTypeEnum(b)
	return nil
}

func (e UserTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid UserTypeEnum")
	}
	return string(e), nil
}

func (e UserTypeEnum) IsValid() bool {
	var types = []string{
		string(UserTypeValue),
		string(CustomerTypeValue),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// StatusEnum defines status types
type StatusEnum string

const (
	DisabledStatus StatusEnum = "disabled"
	EnabledStatus  StatusEnum = "enabled"
)

func (e *StatusEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !StatusEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = StatusEnum(b)
	return nil
}

func (e StatusEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid StatusEnum")
	}
	return string(e), nil
}

func (e StatusEnum) IsValid() bool {
	var types = []string{
		string(DisabledStatus),
		string(EnabledStatus),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// AiTypeEnum defines AI types
type AiTypeEnum string

const (
	GPT35Type  AiTypeEnum = "gpt35"
	GPT4Type   AiTypeEnum = "gpt4"
	ClaudeType AiTypeEnum = "claude"
)

func (e *AiTypeEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !AiTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = AiTypeEnum(b)
	return nil
}

func (e AiTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid AiTypeEnum")
	}
	return string(e), nil
}

func (e AiTypeEnum) IsValid() bool {
	var types = []string{
		string(GPT35Type),
		string(GPT4Type),
		string(ClaudeType),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// DiscountType defines discount types
type DiscountType string

const (
	FixedDiscountType      DiscountType = "fixed"
	PercentageDiscountType DiscountType = "percentage"
)

func (e *DiscountType) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !DiscountType(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = DiscountType(b)
	return nil
}

func (e DiscountType) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid DiscountType")
	}
	return string(e), nil
}

func (e DiscountType) IsValid() bool {
	var types = []string{
		string(FixedDiscountType),
		string(PercentageDiscountType),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}
