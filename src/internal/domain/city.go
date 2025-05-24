package domain

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"time"
)

// City represents User.Cities table
type City struct {
	ID         int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name       string           `json:"name" gorm:"column:name;type:longtext;not null"`
	Slug       string           `json:"slug" gorm:"column:slug;type:longtext;not null"`
	Status     enums.StatusEnum `json:"status" gorm:"column:status;type:longtext;not null"`
	ProvinceID int64            `json:"province_id" gorm:"column:province_id;type:bigint;not null;index"`
	Version    time.Time        `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`

	// Relations
	Province  *Province `json:"province" gorm:"foreignKey:ProvinceID"`
	Addresses []Address `json:"addresses" gorm:"foreignKey:CityID"`
}

// TableName specifies the table name for City
func (City) TableName() string {
	return "cities"
}
func (m *City) GetID() int64 {
	return m.ID
}
func (m *City) GetUserID() *int64 {
	return nil
}
func (m *City) GetCustomerID() *int64 {
	return nil
}
func (m *City) GetSiteID() *int64 {
	return nil
}
