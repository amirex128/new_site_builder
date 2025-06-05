package domain

// Plan represents User.Plans table
type Plan struct {
	ID               int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name             string `json:"name" gorm:"column:name;type:longtext;not null"`
	ShowStatus       string `json:"show_status" gorm:"column:show_status;type:longtext;not null"`
	Description      string `json:"description" gorm:"column:description;type:longtext;null"`
	Price            int64  `json:"price" gorm:"column:price;type:bigint;not null"`
	DiscountType     string `json:"discount_type" gorm:"column:discount_type;type:longtext;null"`
	Discount         *int64 `json:"discount" gorm:"column:discount;type:bigint;null"`
	Duration         int    `json:"duration" gorm:"column:duration;type:int;not null"`
	Feature          string `json:"feature" gorm:"column:feature;type:longtext;null"`
	SmsCredits       int    `json:"sms_credits" gorm:"column:sms_credits;type:int;not null"`
	EmailCredits     int    `json:"email_credits" gorm:"column:email_credits;type:int;not null"`
	StorageMbCredits int    `json:"storage_mb_credits" gorm:"column:storage_mb_credits;type:int;not null"`
	AiCredits        int    `json:"ai_credits" gorm:"column:ai_credits;type:int;not null"`
	AiImageCredits   int    `json:"ai_image_credits" gorm:"column:ai_image_credits;type:int;not null"`

	// Relations
	Roles []Role `json:"roles" gorm:"many2many:role_plan;"`
	Users []User `json:"users" gorm:"foreignKey:PlanID"`
}

// TableName specifies the table name for Plan
func (Plan) TableName() string {
	return "plans"
}
func (m *Plan) GetID() int64 {
	return m.ID
}
func (m *Plan) GetUserID() *int64 {
	return nil
}
func (m *Plan) GetCustomerID() *int64 {
	return nil
}
func (m *Plan) GetSiteID() *int64 {
	return nil
}
