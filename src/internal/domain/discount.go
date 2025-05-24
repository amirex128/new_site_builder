package domain

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"time"
)

// Discount represents Product.Discounts table
type Discount struct {
	ID         int64                  `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Code       string                 `json:"code" gorm:"column:code;type:longtext;not null"`
	Quantity   int                    `json:"quantity" gorm:"column:quantity;type:int;not null"`
	Type       enums.DiscountTypeEnum `json:"type" gorm:"column:type;type:longtext;not null"`
	Value      int64                  `json:"value" gorm:"column:value;type:bigint;not null"`
	ExpiryDate time.Time              `json:"expiry_date" gorm:"column:expiry_date;type:datetime(6);not null"`
	SiteID     int64                  `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	UserID     int64                  `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt  time.Time              `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time              `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version    time.Time              `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted  bool                   `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt  *time.Time             `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Products  []Product  `json:"products" gorm:"many2many:discount_product;"`
	Customers []Customer `json:"customers" gorm:"many2many:customer_discount;"`
	Site      *Site      `json:"site" gorm:"foreignKey:SiteID"`
	User      *User      `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Discount
func (Discount) TableName() string {
	return "discounts"
}
func (m *Discount) GetID() int64 {
	return m.ID
}
func (m *Discount) GetUserID() *int64 {
	return &m.UserID
}
func (m *Discount) GetCutomerID() *int64 {
	return nil
}
func (m *Discount) GetSiteID() *int64 {
	return &m.SiteID
}

// DiscountProduct represents Product.DiscountProduct table - a join table
type DiscountProduct struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ProductID  int64 `json:"product_id" gorm:"column:product_id;type:bigint;not null;index"`
	DiscountID int64 `json:"discount_id" gorm:"column:discount_id;type:bigint;not null"`
}

// TableName specifies the table name for DiscountProduct
func (DiscountProduct) TableName() string {
	return "discount_product"
}

// CustomerDiscount represents Product.CustomerDiscount table - a join table
type CustomerDiscount struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	DiscountID int64 `json:"discount_id" gorm:"column:discount_id;type:bigint;not null;index"`
	CustomerID int64 `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
}

// TableName specifies the table name for CustomerDiscount
func (CustomerDiscount) TableName() string {
	return "customer_discount"
}
