package domain

import (
	"time"
)

// Storage represents Drive.Storages table
type Storage struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	UsedSpaceKb int64     `json:"used_space_kb" gorm:"column:used_space_kb;type:bigint;not null"`
	QuotaKb     int64     `json:"quota_kb" gorm:"column:quota_kb;type:bigint;not null"`
	ChargedAt   time.Time `json:"charged_at" gorm:"column:charged_at;type:datetime(6);not null"`
	ExpireAt    time.Time `json:"expire_at" gorm:"column:expire_at;type:datetime(6);not null"`
	UserID      int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	User *User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Storage
func (Storage) TableName() string {
	return "storages"
}
