package domain

import (
	"time"
)

// ReturnItem represents Order.ReturnItem table
type ReturnItem struct {
	ID           int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ReturnReason string     `json:"return_reason" gorm:"column:return_reason;type:longtext;not null"`
	OrderStatus  int        `json:"order_status" gorm:"column:order_status;type:int;not null"`
	OrderItemID  int64      `json:"order_item_id" gorm:"column:order_item_id;type:bigint;not null;uniqueIndex"`
	ProductID    int64      `json:"product_id" gorm:"column:product_id;type:bigint;not null"`
	UserID       int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID   int64      `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version      time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted    bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	OrderItem *OrderItem `json:"order_item" gorm:"foreignKey:OrderItemID"`
	Product   *Product   `json:"article" gorm:"foreignKey:ProductID"`
	User      *User      `json:"user" gorm:"foreignKey:UserID"`
	Customer  *Customer  `json:"customer" gorm:"foreignKey:CustomerID"`
}

// TableName specifies the table name for ReturnItem
func (ReturnItem) TableName() string {
	return "return_item"
}
