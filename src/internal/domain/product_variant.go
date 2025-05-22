package domain

import (
	"time"
)

// ProductVariant represents Product.ProductVariants table
type ProductVariant struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ProductID  int64     `json:"product_id" gorm:"column:product_id;type:bigint;not null;index"`
	Name       string    `json:"name" gorm:"column:name;type:longtext;not null"`
	Price      int64     `json:"price" gorm:"column:price;type:bigint;not null"`
	Stock      int       `json:"stock" gorm:"column:stock;type:int;not null"`
	UserID     int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID int64     `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Product  *Product  `json:"article" gorm:"foreignKey:ProductID"`
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

// TableName specifies the table name for ProductVariant
func (ProductVariant) TableName() string {
	return "product_variants"
}
