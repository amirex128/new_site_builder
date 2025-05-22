package domain

import (
	"time"
)

// Setting represents Site.Settings table
type Setting struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID     int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null;uniqueIndex"`
	UserID     int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID int64     `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Site     *Site     `json:"site" gorm:"foreignKey:SiteID"`
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

// TableName specifies the table name for Setting
func (Setting) TableName() string {
	return "settings"
}
