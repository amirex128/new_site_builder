package domain

import (
	"github.com/amirex128/new_site_builder/internal/domain/enums"
	"time"
)

// Customer represents User.Customers table
type Customer struct {
	ID                 int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID             int64            `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	FirstName          string           `json:"first_name" gorm:"column:first_name;type:longtext;null"`
	AvatarID           *int64           `json:"avatar_id" gorm:"column:avatar_id;type:bigint;null"`
	LastName           string           `json:"last_name" gorm:"column:last_name;type:longtext;null"`
	Email              string           `json:"email" gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	VerifyEmail        enums.StatusEnum `json:"verify_email" gorm:"column:verify_email;type:ENUM('inactive','active','pending');default:'inactive';null"`
	Password           string           `json:"password" gorm:"column:password;type:longtext;not null"`
	Salt               string           `json:"salt" gorm:"column:salt;type:longtext;not null"`
	NationalCode       string           `json:"national_code" gorm:"column:national_code;type:longtext;null"`
	Phone              string           `json:"phone" gorm:"column:phone;type:longtext;null"`
	VerifyPhone        enums.StatusEnum `json:"verify_phone" gorm:"column:verify_phone;type:ENUM('inactive','active','pending');default:'inactive';null"`
	IsActive           enums.StatusEnum `json:"is_active" gorm:"column:is_active;type:ENUM('inactive','active','pending');default:'inactive';not null"`
	VerifyCode         *int             `json:"verify_code" gorm:"column:verify_code;type:int;null"`
	ExpireVerifyCodeAt *time.Time       `json:"expire_verify_code_at" gorm:"column:expire_verify_code_at;type:datetime(6);null"`
	CreatedAt          time.Time        `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	IsDeleted          bool             `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt          *time.Time       `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Roles     []Role     `json:"roles" gorm:"many2many:customer_roles;"`
	Addresses []Address  `json:"addresses" gorm:"many2many:address_customer;"`
	Discounts []Discount `json:"discounts" gorm:"many2many:customer_discount;"`
}

// TableName specifies the table name for Customer
func (Customer) TableName() string {
	return "customers"
}

func (m *Customer) GetID() int64 {
	return m.ID
}
func (m *Customer) GetUserID() *int64 {
	return nil
}
func (m *Customer) GetCustomerID() *int64 {
	return nil
}
func (m *Customer) GetSiteID() *int64 {
	return &m.SiteID
}
