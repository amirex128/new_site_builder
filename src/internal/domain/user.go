package domain

import (
	"time"
)

// User represents User.Users table
type User struct {
	ID                       int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	FirstName                string     `json:"first_name" gorm:"column:first_name;type:longtext;null"`
	LastName                 string     `json:"last_name" gorm:"column:last_name;type:longtext;null"`
	Email                    string     `json:"email" gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	AvatarID                 *int64     `json:"avatar_id" gorm:"column:avatar_id;type:bigint;null"`
	VerifyEmail              string     `json:"verify_email" gorm:"column:verify_email;type:longtext;null"`
	Password                 string     `json:"password" gorm:"column:password;type:longtext;not null"`
	Salt                     string     `json:"salt" gorm:"column:salt;type:longtext;not null"`
	NationalCode             string     `json:"national_code" gorm:"column:national_code;type:longtext;null"`
	Phone                    string     `json:"phone" gorm:"column:phone;type:longtext;null"`
	VerifyPhone              string     `json:"verify_phone" gorm:"column:verify_phone;type:longtext;null"`
	IsActive                 string     `json:"is_active" gorm:"column:is_active;type:longtext;not null"`
	AiTypeEnum               string     `json:"ai_type_enum" gorm:"column:ai_type_enum;type:longtext;not null"`
	UserTypeEnum             string     `json:"user_type_enum" gorm:"column:user_type_enum;type:longtext;not null"`
	PlanID                   *int64     `json:"plan_id" gorm:"column:plan_id;type:bigint;null"`
	PlanStartedAt            *time.Time `json:"plan_started_at" gorm:"column:plan_started_at;type:datetime(6);null"`
	PlanExpiredAt            *time.Time `json:"plan_expired_at" gorm:"column:plan_expired_at;type:datetime(6);null"`
	VerifyCode               *int       `json:"verify_code" gorm:"column:verify_code;type:int;null"`
	ExpireVerifyCodeAt       *time.Time `json:"expire_verify_code_at" gorm:"column:expire_verify_code_at;type:datetime(6);null"`
	AiCredits                int        `json:"ai_credits" gorm:"column:ai_credits;type:int;not null"`
	AiImageCredits           int        `json:"ai_image_credits" gorm:"column:ai_image_credits;type:int;not null"`
	StorageMbCredits         int        `json:"storage_mb_credits" gorm:"column:storage_mb_credits;type:int;not null"`
	StorageMbCreditsExpireAt *time.Time `json:"storage_mb_credits_expire_at" gorm:"column:storage_mb_credits_expire_at;type:datetime(6);null"`
	EmailCredits             int        `json:"email_credits" gorm:"column:email_credits;type:int;not null"`
	SmsCredits               int        `json:"sms_credits" gorm:"column:sms_credits;type:int;not null"`
	UseCustomEmailSmtp       string     `json:"use_custom_email_smtp" gorm:"column:use_custom_email_smtp;type:longtext;not null"`
	SmtpHost                 string     `json:"smtp_host" gorm:"column:smtp_host;type:longtext;null"`
	SmtpPort                 *int       `json:"smtp_port" gorm:"column:smtp_port;type:int;null"`
	SmtpUsername             string     `json:"smtp_username" gorm:"column:smtp_username;type:longtext;null"`
	SmtpPassword             string     `json:"smtp_password" gorm:"column:smtp_password;type:longtext;null"`
	SmtpEnableSsl            *bool      `json:"smtp_enable_ssl" gorm:"column:smtp_enable_ssl;type:tinyint(1);null"`
	SmtpSenderEmail          string     `json:"smtp_sender_email" gorm:"column:smtp_sender_email;type:longtext;null"`
	IsAdmin                  bool       `json:"is_admin" gorm:"column:is_admin;type:tinyint(1);not null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt                time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version                  time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted                bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt                *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Roles     []Role    `json:"roles" gorm:"many2many:role_user;"`
	Addresses []Address `json:"addresses" gorm:"many2many:address_user;"`
	Plan      *Plan     `json:"plan" gorm:"foreignKey:PlanID"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// Customer represents User.Customers table
type Customer struct {
	ID                 int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID             int64      `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	FirstName          string     `json:"first_name" gorm:"column:first_name;type:longtext;null"`
	AvatarID           *int64     `json:"avatar_id" gorm:"column:avatar_id;type:bigint;null"`
	LastName           string     `json:"last_name" gorm:"column:last_name;type:longtext;null"`
	Email              string     `json:"email" gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	VerifyEmail        string     `json:"verify_email" gorm:"column:verify_email;type:longtext;null"`
	Password           string     `json:"password" gorm:"column:password;type:longtext;not null"`
	Salt               string     `json:"salt" gorm:"column:salt;type:longtext;not null"`
	NationalCode       string     `json:"national_code" gorm:"column:national_code;type:longtext;null"`
	Phone              string     `json:"phone" gorm:"column:phone;type:longtext;null"`
	VerifyPhone        string     `json:"verify_phone" gorm:"column:verify_phone;type:longtext;null"`
	IsActive           string     `json:"is_active" gorm:"column:is_active;type:longtext;not null"`
	VerifyCode         *int       `json:"verify_code" gorm:"column:verify_code;type:int;null"`
	ExpireVerifyCodeAt *time.Time `json:"expire_verify_code_at" gorm:"column:expire_verify_code_at;type:datetime(6);null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version            time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted          bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Roles     []Role     `json:"roles" gorm:"many2many:customer_roles;"`
	Addresses []Address  `json:"addresses" gorm:"many2many:address_customer;"`
	Discounts []Discount `json:"discounts" gorm:"many2many:customer_discount;"`
}

// TableName specifies the table name for Customer
func (Customer) TableName() string {
	return "customers"
}

// RoleUser represents User.RoleUser table - a join table
type RoleUser struct {
	ID     int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	RoleID int64 `json:"role_id" gorm:"column:role_id;type:bigint;not null;index"`
	UserID int64 `json:"user_id" gorm:"column:user_id;type:bigint;not null;index"`
}

// TableName specifies the table name for RoleUser
func (RoleUser) TableName() string {
	return "role_user"
}

// CustomerRole represents User.CustomerRoles table - a join table
type CustomerRole struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	RoleID     int64 `json:"role_id" gorm:"column:role_id;type:bigint;not null;index"`
	CustomerID int64 `json:"customer_id" gorm:"column:customer_id;type:bigint;not null;index"`
}

// TableName specifies the table name for CustomerRole
func (CustomerRole) TableName() string {
	return "customer_roles"
}

// PermissionRole represents User.PermissionRoles table - a join table
type PermissionRole struct {
	ID           int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	RoleID       int64 `json:"role_id" gorm:"column:role_id;type:bigint;not null;index"`
	PermissionID int64 `json:"permission_id" gorm:"column:permission_id;type:bigint;not null;index"`
}

// TableName specifies the table name for PermissionRole
func (PermissionRole) TableName() string {
	return "permission_roles"
}

// RolePlan represents User.RolePlan table - a join table
type RolePlan struct {
	ID     int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	RoleID int64 `json:"role_id" gorm:"column:role_id;type:bigint;not null;index"`
	PlanID int64 `json:"plan_id" gorm:"column:plan_id;type:bigint;not null;index"`
}

// TableName specifies the table name for RolePlan
func (RolePlan) TableName() string {
	return "role_plan"
}

// AddressUser represents User.AddressUser table - a join table
type AddressUser struct {
	ID        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	AddressID int64 `json:"address_id" gorm:"column:address_id;type:bigint;not null;index"`
	UserID    int64 `json:"user_id" gorm:"column:user_id;type:bigint;not null;index"`
}

// TableName specifies the table name for AddressUser
func (AddressUser) TableName() string {
	return "address_user"
}

// AddressCustomer represents User.AddressCustomer table - a join table
type AddressCustomer struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	AddressID  int64 `json:"address_id" gorm:"column:address_id;type:bigint;not null;index"`
	CustomerID int64 `json:"customer_id" gorm:"column:customer_id;type:bigint;not null;index"`
}

// TableName specifies the table name for AddressCustomer
func (AddressCustomer) TableName() string {
	return "address_customer"
}
