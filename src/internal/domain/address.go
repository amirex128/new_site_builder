package domain

import (
	"time"
)

// Address represents User.Addresses table
type Address struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Title       string    `json:"title" gorm:"column:title;type:longtext;null"`
	Latitude    *float32  `json:"latitude" gorm:"column:latitude;type:float;null"`
	Longitude   *float32  `json:"longitude" gorm:"column:longitude;type:float;null"`
	AddressLine string    `json:"address_line" gorm:"column:address_line;type:longtext;not null"`
	PostalCode  string    `json:"postal_code" gorm:"column:postal_code;type:longtext;not null"`
	CityID      int64     `json:"city_id" gorm:"column:city_id;type:bigint;not null;index"`
	ProvinceID  int64     `json:"province_id" gorm:"column:province_id;type:bigint;not null;index"`
	UserID      *int64    `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID  *int64    `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	City     *City     `json:"city" gorm:"foreignKey:CityID"`
	Province *Province `json:"province" gorm:"foreignKey:ProvinceID"`
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

// TableName specifies the table name for Address
func (Address) TableName() string {
	return "addresses"
}
func (m *Address) GetID() int64 {
	return m.ID
}
func (m *Address) GetUserID() *int64 {
	return m.UserID
}
func (m *Address) GetCutomerID() *int64 {
	return m.CustomerID
}
func (m *Address) GetSiteID() *int64 {
	return nil
}
