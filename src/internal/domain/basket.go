package domain

import (
	"time"
)

// Basket represents Order.Baskets table
type Basket struct {
	ID                           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID                       int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	TotalRawPrice                int64     `json:"total_raw_price" gorm:"column:total_raw_price;type:bigint;not null"`
	TotalCouponDiscount          int64     `json:"total_coupon_discount" gorm:"column:total_coupon_discount;type:bigint;not null"`
	TotalPriceWithCouponDiscount int64     `json:"total_price_with_coupon_discount" gorm:"column:total_price_with_coupon_discount;type:bigint;not null"`
	DiscountID                   *int64    `json:"discount_id" gorm:"column:discount_id;type:bigint;null"`
	CustomerID                   int64     `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt                    time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Items    []BasketItem `json:"items" gorm:"foreignKey:BasketID"`
	Customer *Customer    `json:"customer" gorm:"foreignKey:CustomerID"`
	Discount *Discount    `json:"discount" gorm:"foreignKey:DiscountID"`
	Site     *Site        `json:"site" gorm:"foreignKey:SiteID"`
}

// TableName specifies the table name for Basket
func (Basket) TableName() string {
	return "baskets"
}
func (m *Basket) GetID() int64 {
	return m.ID
}
func (m *Basket) GetUserID() *int64 {
	return nil
}
func (m *Basket) GetCustomerID() *int64 {
	return &m.CustomerID
}
func (m *Basket) GetSiteID() *int64 {
	return &m.SiteID
}
