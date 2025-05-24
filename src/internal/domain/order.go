package domain

import (
	"time"
)

// Order represents Order.Orders table
type Order struct {
	ID                           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID                       int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	TotalRawPrice                int64     `json:"total_raw_price" gorm:"column:total_raw_price;type:bigint;not null"`
	TotalCouponDiscount          int64     `json:"total_coupon_discount" gorm:"column:total_coupon_discount;type:bigint;not null"`
	TotalPriceWithCouponDiscount int64     `json:"total_price_with_coupon_discount" gorm:"column:total_price_with_coupon_discount;type:bigint;not null"`
	CourierPrice                 int64     `json:"courier_price" gorm:"column:courier_price;type:bigint;not null"`
	Courier                      string    `json:"courier" gorm:"column:courier;type:longtext;not null"`
	OrderStatus                  string    `json:"order_status" gorm:"column:order_status;type:longtext;not null"`
	TotalFinalPrice              int64     `json:"total_final_price" gorm:"column:total_final_price;type:bigint;not null"`
	Description                  string    `json:"description" gorm:"column:description;type:longtext;null"`
	TotalWeight                  int       `json:"total_weight" gorm:"column:total_weight;type:int;not null"`
	TrackingCode                 string    `json:"tracking_code" gorm:"column:tracking_code;type:longtext;null"`
	BasketID                     int64     `json:"basket_id" gorm:"column:basket_id;type:bigint;not null"`
	DiscountID                   *int64    `json:"discount_id" gorm:"column:discount_id;type:bigint;null"`
	AddressID                    int64     `json:"address_id" gorm:"column:address_id;type:bigint;not null"`
	CustomerID                   int64     `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt                    time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Address  *Address  `json:"address" gorm:"foreignKey:AddressID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	Site     *Site     `json:"site" gorm:"foreignKey:SiteID"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}
func (m *Order) GetID() int64 {
	return m.ID
}
func (m *Order) GetUserID() *int64 {
	return nil
}
func (m *Order) GetCustomerID() *int64 {
	return &m.CustomerID
}
func (m *Order) GetSiteID() *int64 {
	return &m.SiteID
}
