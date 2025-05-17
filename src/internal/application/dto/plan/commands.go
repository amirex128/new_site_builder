package plan

// CreatePlanCommand represents a command to create a new plan
type CreatePlanCommand struct {
	Name           *string           `json:"name" validate:"required,max=100" error:"required=نام پلن الزامی است|max=نام پلن نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Description    *string           `json:"description,omitempty" validate:"omitempty,max=1000" error:"max=توضیحات نمی‌تواند بیشتر از 1000 کاراکتر باشد"`
	Price          *int64            `json:"price" validate:"required,min=0" error:"required=قیمت الزامی است|min=قیمت نمی‌تواند کمتر از 0 باشد"`
	DiscountType   *DiscountTypeEnum `json:"discountType,omitempty" error:""`
	Discount       *int64            `json:"discount,omitempty" validate:"omitempty,min=0" error:"min=تخفیف نمی‌تواند کمتر از 0 باشد"`
	Duration       *int              `json:"duration" validate:"required,min=1,max=1000" error:"required=مدت زمان الزامی است|min=مدت زمان باید حداقل 1 باشد|max=مدت زمان نمی‌تواند بیشتر از 1000 باشد"`
	Feature        *string           `json:"feature,omitempty" validate:"omitempty,max=500" error:"max=توضیحات ویژگی‌ها نمی‌تواند بیشتر از 500 کاراکتر باشد"`
	SmsCredits     *int              `json:"smsCredits,omitempty" validate:"omitempty,min=0,max=1000" error:"min=اعتبار پیامک نمی‌تواند کمتر از 0 باشد|max=اعتبار پیامک نمی‌تواند بیشتر از 1000 باشد"`
	EmailCredits   *int              `json:"emailCredits,omitempty" validate:"omitempty,min=0,max=1000" error:"min=اعتبار ایمیل نمی‌تواند کمتر از 0 باشد|max=اعتبار ایمیل نمی‌تواند بیشتر از 1000 باشد"`
	StorageCredits *int              `json:"storageCredits,omitempty" validate:"omitempty,min=0,max=1000" error:"min=اعتبار فضای ذخیره‌سازی نمی‌تواند کمتر از 0 باشد|max=اعتبار فضای ذخیره‌سازی نمی‌تواند بیشتر از 1000 باشد"`
	AiCredits      *int              `json:"aiCredits,omitempty" validate:"omitempty,min=0,max=1000" error:"min=اعتبار هوش مصنوعی نمی‌تواند کمتر از 0 باشد|max=اعتبار هوش مصنوعی نمی‌تواند بیشتر از 1000 باشد"`
	AiImageCredits *int              `json:"aiImageCredits,omitempty" validate:"omitempty,min=0,max=1000" error:"min=اعتبار تصویر هوش مصنوعی نمی‌تواند کمتر از 0 باشد|max=اعتبار تصویر هوش مصنوعی نمی‌تواند بیشتر از 1000 باشد"`
}

// DeletePlanCommand represents a command to delete a plan
type DeletePlanCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=پلن الزامی است|gt=شناسه پلن باید بزرگتر از 0 باشد"`
}

// UpdatePlanCommand represents a command to update a plan
type UpdatePlanCommand struct {
	ID             *int64            `json:"id" validate:"required,gt=0" error:"required=پلن الزامی است|gt=شناسه پلن باید بزرگتر از 0 باشد"`
	Name           *string           `json:"name,omitempty" validate:"omitempty,max=100" error:"max=نام پلن نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Description    *string           `json:"description,omitempty" validate:"omitempty,max=1000" error:"max=توضیحات نمی‌تواند بیشتر از 1000 کاراکتر باشد"`
	Price          *int64            `json:"price,omitempty" validate:"omitempty,min=0" error:"min=قیمت نمی‌تواند کمتر از 0 باشد"`
	DiscountType   *DiscountTypeEnum `json:"discountType,omitempty" error:""`
	Discount       *int64            `json:"discount,omitempty" validate:"omitempty,min=0" error:"min=تخفیف نمی‌تواند کمتر از 0 باشد"`
	Duration       *int              `json:"duration,omitempty" validate:"omitempty,min=1" error:"min=مدت زمان باید حداقل 1 باشد"`
	Feature        *string           `json:"feature,omitempty" validate:"omitempty,max=500" error:"max=توضیحات ویژگی‌ها نمی‌تواند بیشتر از 500 کاراکتر باشد"`
	SmsCredits     *int              `json:"smsCredits,omitempty" validate:"omitempty,min=0" error:"min=اعتبار پیامک نمی‌تواند کمتر از 0 باشد"`
	EmailCredits   *int              `json:"emailCredits,omitempty" validate:"omitempty,min=0" error:"min=اعتبار ایمیل نمی‌تواند کمتر از 0 باشد"`
	StorageCredits *int              `json:"storageCredits,omitempty" validate:"omitempty,min=0" error:"min=اعتبار فضای ذخیره‌سازی نمی‌تواند کمتر از 0 باشد"`
	AiCredits      *int              `json:"aiCredits,omitempty" validate:"omitempty,min=0" error:"min=اعتبار هوش مصنوعی نمی‌تواند کمتر از 0 باشد"`
	AiImageCredits *int              `json:"aiImageCredits,omitempty" validate:"omitempty,min=0" error:"min=اعتبار تصویر هوش مصنوعی نمی‌تواند کمتر از 0 باشد"`
}
