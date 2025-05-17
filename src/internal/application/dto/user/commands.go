package user

// LoginUserCommand represents a command to log in a user
type LoginUserCommand struct {
	Email    *string `json:"email" validate:"required,email" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد"`
	Password *string `json:"password" validate:"required,min=8,max=100" error:"required=رمز عبور الزامی است|min=رمز عبور باید حداقل 8 کاراکتر باشد|max=رمز عبور نمی‌تواند بیشتر از 100 کاراکتر باشد"`
}

// RegisterUserCommand represents a command to register a new user
type RegisterUserCommand struct {
	Email    *string `json:"email" validate:"required,email,max=100" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد|max=ایمیل نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Password *string `json:"password" validate:"required,min=8,max=100" error:"required=رمز عبور الزامی است|min=رمز عبور باید حداقل 8 کاراکتر باشد|max=رمز عبور نمی‌تواند بیشتر از 100 کاراکتر باشد"`
}

// RequestVerifyAndForgetUserCommand represents a command to request verification or password reset
type RequestVerifyAndForgetUserCommand struct {
	Email *string         `json:"email,omitempty" validate:"required_if=Type 0 Type 2,omitempty,email" error:"required_if=ایمیل الزامی است|email=ایمیل باید معتبر باشد"`
	Phone *string         `json:"phone,omitempty" validate:"required_if=Type 1 Type 3,omitempty,pattern=^09\\d{9}$" error:"required_if=شماره تلفن الزامی است|pattern=شماره تلفن نامعتبر است"`
	Type  *VerifyTypeEnum `json:"type" validate:"required" error:"required=نوع تأیید الزامی است"`
}

// SmptSettings represents SMTP settings for user profile
type SmptSettings struct {
	Host     string `json:"host" validate:"required,max=100" error:"required=هاست SMTP الزامی است|max=هاست SMTP نباید بیشتر از 100 کاراکتر باشد"`
	Port     int    `json:"port" validate:"required,min=1,max=65535" error:"required=پورت SMTP الزامی است|min=پورت SMTP باید حداقل 1 باشد|max=پورت SMTP باید حداکثر 65535 باشد"`
	Username string `json:"username" validate:"required,max=100" error:"required=نام کاربری SMTP الزامی است|max=نام کاربری SMTP نباید بیشتر از 100 کاراکتر باشد"`
	Password string `json:"password" validate:"required,max=100" error:"required=رمز عبور SMTP الزامی است|max=رمز عبور SMTP نباید بیشتر از 100 کاراکتر باشد"`
}

// UpdateProfileUserCommand represents a command to update a user's profile
type UpdateProfileUserCommand struct {
	FirstName          *string       `json:"firstName,omitempty" validate:"omitempty,max=100" error:"max=نام نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	LastName           *string       `json:"lastName,omitempty" validate:"omitempty,max=100" error:"max=نام خانوادگی نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Email              *string       `json:"email,omitempty" validate:"omitempty,email" error:"email=فرمت ایمیل نامعتبر است"`
	Password           *string       `json:"password,omitempty" validate:"omitempty,min=6,max=100" error:"min=رمز عبور باید حداقل 6 کاراکتر باشد|max=رمز عبور نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	NationalCode       *string       `json:"nationalCode,omitempty" validate:"omitempty,max=100" error:"max=کد ملی نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Phone              *string       `json:"phone" validate:"required,pattern=^09\\d{9}$" error:"required=شماره تلفن الزامی است|pattern=شماره تلفن باید یک شماره موبایل ایرانی معتبر باشد"`
	AddressIDs         []int64       `json:"addressIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های آدرس باید بزرگتر از 0 باشند"`
	AiTypeEnum         *AiTypeEnum   `json:"aiTypeEnum,omitempty" error:""`
	UseCustomEmailSmtp *StatusEnum   `json:"useCustomEmailSmtp,omitempty" error:""`
	Smtp               *SmptSettings `json:"smtp,omitempty" validate:"omitempty" error:""`
}

// UnitPriceQuery represents a nested query for unit price in charge credit request
type UnitPriceQuery struct {
	UnitPriceName  *UnitPriceNameEnum `json:"unitPriceName" validate:"required" error:"required=نام واحد قیمت الزامی است"`
	UnitPriceCount *int               `json:"unitPriceCount" validate:"required,min=1,max=1000" error:"required=تعداد واحد قیمت الزامی است|min=تعداد واحد قیمت باید حداقل 1 باشد|max=تعداد واحد قیمت نمی‌تواند بیشتر از 1000 باشد"`
	UnitPriceDay   *int               `json:"unitPriceDay,omitempty" validate:"omitempty" error:""`
}

// ChargeCreditRequestUserCommand represents a command to request charging credit
type ChargeCreditRequestUserCommand struct {
	Gateway             *PaymentGatewaysEnum `json:"gateway" validate:"required" error:"required=درگاه پرداخت الزامی است"`
	FinalFrontReturnUrl *string              `json:"finalFrontReturnUrl" validate:"required,max=500,url" error:"required=آدرس بازگشت الزامی است|max=آدرس بازگشت نمی‌تواند بیشتر از 500 کاراکتر باشد|url=آدرس بازگشت باید یک URL معتبر باشد"`
	UnitPrices          []UnitPriceQuery     `json:"unitPrices" validate:"required,min=1" error:"required=واحدهای قیمت الزامی هستند|min=حداقل یک واحد قیمت باید وجود داشته باشد"`
}

// ChargeCreditVerifyUserCommand represents a command to verify charge credit
type ChargeCreditVerifyUserCommand struct {
	PaymentStatus *string           `json:"paymentStatus" validate:"required,max=50" error:"required=وضعیت پرداخت الزامی است|max=وضعیت پرداخت نمی‌تواند بیشتر از 50 کاراکتر باشد"`
	IsSuccess     *bool             `json:"isSuccess" validate:"required" error:"required=وضعیت موفقیت پرداخت الزامی است"`
	OrderData     map[string]string `json:"orderData" validate:"required" error:"required=اطلاعات سفارش الزامی است"`
}

// UpgradePlanRequestUserCommand represents a command to request plan upgrade
type UpgradePlanRequestUserCommand struct {
	Gateway             *PaymentGatewaysEnum `json:"gateway" validate:"required" error:"required=درگاه پرداخت الزامی است"`
	FinalFrontReturnUrl *string              `json:"finalFrontReturnUrl" validate:"required,url" error:"required=آدرس بازگشت الزامی است|url=آدرس بازگشت باید یک URL معتبر باشد"`
	PlanID              *int64               `json:"planId" validate:"required,gt=0" error:"required=پلن الزامی است|gt=شناسه پلن باید بزرگتر از 0 باشد"`
}

// UpgradePlanVerifyUserCommand represents a command to verify plan upgrade
type UpgradePlanVerifyUserCommand struct {
	PaymentStatus *string           `json:"paymentStatus" validate:"required" error:"required=وضعیت پرداخت الزامی است"`
	IsSuccess     *bool             `json:"isSuccess" validate:"required" error:"required=نتیجه پرداخت الزامی است"`
	OrderData     map[string]string `json:"orderData" validate:"required" error:"required=اطلاعات سفارش الزامی است"`
}
