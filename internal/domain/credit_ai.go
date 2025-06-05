package domain

import (
	"time"
)

// Credit represents Ai.Credits table
type Credit struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	UserID     int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID int64     `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

// TableName specifies the table name for Credit
func (Credit) TableName() string {
	return "credits"
}
func (m *Credit) GetID() int64 {
	return m.ID
}
func (m *Credit) GetUserID() *int64 {
	return &m.UserID
}
func (m *Credit) GetCustomerID() *int64 {
	return &m.CustomerID
}
func (m *Credit) GetSiteID() *int64 {
	return nil
}
