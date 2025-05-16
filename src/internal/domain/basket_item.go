package domain

import (
	"time"
)

// BasketItem represents Order.BasketItems table
type BasketItem struct {
	ID                           int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Quantity                     int        `json:"quantity" gorm:"column:quantity;type:int;not null"`
	RawPrice                     int64      `json:"raw_price" gorm:"column:raw_price;type:bigint;not null"`
	FinalRawPrice                int64      `json:"final_raw_price" gorm:"column:final_raw_price;type:bigint;not null"`
	FinalPriceWithCouponDiscount int64      `json:"final_price_with_coupon_discount" gorm:"column:final_price_with_coupon_discount;type:bigint;not null"`
	JustCouponPrice              int64      `json:"just_coupon_price" gorm:"column:just_coupon_price;type:bigint;not null"`
	JustDiscountPrice            int64      `json:"just_discount_price" gorm:"column:just_discount_price;type:bigint;not null"`
	BasketID                     int64      `json:"basket_id" gorm:"column:basket_id;type:bigint;not null;index"`
	ProductID                    int64      `json:"product_id" gorm:"column:product_id;type:bigint;not null"`
	ProductVariantID             int64      `json:"product_variant_id" gorm:"column:product_variant_id;type:bigint;not null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt                    time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version                      time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted                    bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt                    *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Basket         *Basket         `json:"basket" gorm:"foreignKey:BasketID"`
	Product        *Product        `json:"product" gorm:"foreignKey:ProductID"`
	ProductVariant *ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
}

// TableName specifies the table name for BasketItem
func (BasketItem) TableName() string {
	return "basket_items"
}
