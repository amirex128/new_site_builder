package domain

// UnitPrice represents User.UnitPrices table
type UnitPrice struct {
	ID           int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name         string `json:"name" gorm:"column:name;type:longtext;not null"`
	HasDay       bool   `json:"has_day" gorm:"column:has_day;type:tinyint(1);not null"`
	Price        int64  `json:"price" gorm:"column:price;type:bigint;not null"`
	DiscountType string `json:"discount_type" gorm:"column:discount_type;type:longtext;null"`
	Discount     *int64 `json:"discount" gorm:"column:discount;type:bigint;null"`
}

// TableName specifies the table name for UnitPrice
func (UnitPrice) TableName() string {
	return "unit_prices"
}
