package domain

import (
	"time"
)

// Role represents User.Roles table
type Role struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name      string    `json:"name" gorm:"column:name;type:longtext;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	IsDeleted bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`

	// Relations
	Users       []User       `json:"users" gorm:"many2many:role_user;"`
	Customers   []Customer   `json:"customers" gorm:"many2many:customer_roles;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:permission_roles;"`
	Plans       []Plan       `json:"plans" gorm:"many2many:role_plan;"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}
func (m *Role) GetID() int64 {
	return m.ID
}
func (m *Role) GetUserID() *int64 {
	return nil
}
func (m *Role) GetCutomerID() *int64 {
	return nil
}
func (m *Role) GetSiteID() *int64 {
	return nil
}
