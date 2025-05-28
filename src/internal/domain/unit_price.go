package domain

import "github.com/amirex128/new_site_builder/src/internal/domain/enums"

// UnitPrice represents User.UnitPrices table
type UnitPrice struct {
	ID           int64                  `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name         string                 `json:"name" gorm:"column:name;type:longtext;not null"`
	HasDay       bool                   `json:"has_day" gorm:"column:has_day;type:tinyint(1);not null"`
	Price        int64                  `json:"price" gorm:"column:price;type:bigint;not null"`
	DiscountType enums.DiscountTypeEnum `json:"discount_type" gorm:"column:discount_type;type:ENUM('fixed','percentage');default:'fixed';null"`
	Discount     *int64                 `json:"discount" gorm:"column:discount;type:bigint;null"`
}

// TableName specifies the table name for UnitPrice
func (UnitPrice) TableName() string {
	return "unit_prices"
}
func (m *UnitPrice) GetID() int64 {
	return m.ID
}
func (m *UnitPrice) GetUserID() *int64 {
	return nil
}
func (m *UnitPrice) GetCustomerID() *int64 {
	return nil
}
func (m *UnitPrice) GetSiteID() *int64 {
	return nil
}
