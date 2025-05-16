package domain

import (
	"time"
)

// Coupon represents Product.Coupons table
type Coupon struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ProductID  int64      `json:"product_id" gorm:"column:product_id;type:bigint;not null;uniqueIndex"`
	Quantity   int        `json:"quantity" gorm:"column:quantity;type:int;not null"`
	Type       string     `json:"type" gorm:"column:type;type:longtext;not null"`
	Value      int64      `json:"value" gorm:"column:value;type:bigint;not null"`
	ExpiryDate time.Time  `json:"expiry_date" gorm:"column:expiry_date;type:datetime(6);not null"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version    time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted  bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Product *Product `json:"product" gorm:"foreignKey:ProductID"`
}

// TableName specifies the table name for Coupon
func (Coupon) TableName() string {
	return "coupons"
}
