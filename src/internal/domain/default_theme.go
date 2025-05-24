package domain

import (
	"time"
)

// DefaultTheme represents Site.DefaultThemes table
type DefaultTheme struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name        string    `json:"name" gorm:"column:name;type:longtext;not null"`
	Description string    `json:"description" gorm:"column:description;type:longtext;null"`
	Demo        string    `json:"demo" gorm:"column:demo;type:longtext;null"`
	MediaID     int64     `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
	Pages       string    `json:"pages" gorm:"column:pages;type:longtext;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Media *Media `json:"media" gorm:"foreignKey:MediaID"`
}

// TableName specifies the table name for DefaultTheme
func (DefaultTheme) TableName() string {
	return "default_themes"
}
func (m *DefaultTheme) GetID() int64 {
	return m.ID
}
func (m *DefaultTheme) GetUserID() *int64 {
	return nil
}
func (m *DefaultTheme) GetCutomerID() *int64 {
	return nil
}
func (m *DefaultTheme) GetSiteID() *int64 {
	return nil
}
