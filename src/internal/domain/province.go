package domain

import "github.com/amirex128/new_site_builder/src/internal/domain/enums"

// Province represents User.Provinces table
type Province struct {
	ID     int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name   string           `json:"name" gorm:"column:name;type:longtext;not null"`
	Slug   string           `json:"slug" gorm:"column:slug;type:longtext;not null"`
	Status enums.StatusEnum `json:"status" gorm:"column:status;type:ENUM('inactive','active','pending');default:'inactive';not null"`

	// Relations
	Cities    []City    `json:"cities" gorm:"foreignKey:ProvinceID"`
	Addresses []Address `json:"addresses" gorm:"foreignKey:ProvinceID"`
}

// TableName specifies the table name for Province
func (Province) TableName() string {
	return "provinces"
}
func (m *Province) GetID() int64 {
	return m.ID
}
func (m *Province) GetUserID() *int64 {
	return nil
}
func (m *Province) GetCustomerID() *int64 {
	return nil
}
func (m *Province) GetSiteID() *int64 {
	return nil
}
