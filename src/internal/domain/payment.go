package domain

import (
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// Payment represents Payment.Payments table
type Payment struct {
	ID                  int64              `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID              int64              `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	PaymentStatusEnum   enums.StatusEnum   `json:"payment_status_enum" gorm:"column:payment_status_enum;type:ENUM('inactive','active','pending');default:'inactive';not null"`
	UserType            enums.UserTypeEnum `json:"user_type" gorm:"column:user_type;type:ENUM('user','customer','guest');default:'user';null"`
	TrackingNumber      int64              `json:"tracking_number" gorm:"column:tracking_number;type:bigint;not null"`
	Gateway             string             `json:"gateway" gorm:"column:gateway;type:longtext;not null"`
	GatewayAccountName  string             `json:"gateway_account_name" gorm:"column:gateway_account_name;type:longtext;not null"`
	Amount              int64              `json:"amount" gorm:"column:amount;type:bigint;not null"`
	ServiceName         string             `json:"service_name" gorm:"column:service_name;type:longtext;not null"`
	ServiceAction       string             `json:"service_action" gorm:"column:service_action;type:longtext;not null"`
	OrderID             int64              `json:"order_id" gorm:"column:order_id;type:bigint;not null"`
	ReturnUrl           string             `json:"return_url" gorm:"column:return_url;type:longtext;not null"`
	CallVerifyUrl       string             `json:"call_verify_url" gorm:"column:call_verify_url;type:longtext;not null"`
	ClientIp            string             `json:"client_ip" gorm:"column:client_ip;type:longtext;not null"`
	Message             string             `json:"message" gorm:"column:message;type:longtext;null"`
	GatewayResponseCode string             `json:"gateway_response_code" gorm:"column:gateway_response_code;type:longtext;null"`
	TransactionCode     string             `json:"transaction_code" gorm:"column:transaction_code;type:longtext;null"`
	AdditionalData      string             `json:"additional_data" gorm:"column:additional_data;type:longtext;null"`
	OrderData           string             `json:"order_data" gorm:"column:order_data;type:longtext;null"`
	UserID              int64              `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID          int64              `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt           time.Time          `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt           time.Time          `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version             time.Time          `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted           bool               `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt           *time.Time         `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	Order    *Order    `json:"order" gorm:"foreignKey:OrderID"`
	Site     *Site     `json:"site" gorm:"foreignKey:SiteID"`
}

// TableName specifies the table name for Payment
func (Payment) TableName() string {
	return "payments"
}
func (m *Payment) GetID() int64 {
	return m.ID
}
func (m *Payment) GetUserID() *int64 {
	return &m.UserID
}
func (m *Payment) GetCustomerID() *int64 {
	return &m.CustomerID
}
func (m *Payment) GetSiteID() *int64 {
	return &m.SiteID
}
