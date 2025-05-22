package domain

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"time"
)

// FileItem represents Drive.FileItems table
type FileItem struct {
	ID          int64                        `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name        string                       `json:"name" gorm:"column:name;type:longtext;not null"`
	BucketName  string                       `json:"bucket_name" gorm:"column:bucket_name;type:longtext;not null"`
	ServerKey   string                       `json:"server_key" gorm:"column:server_key;type:longtext;not null"`
	FilePath    string                       `json:"file_path" gorm:"column:file_path;type:longtext;not null"`
	IsDirectory bool                         `json:"is_directory" gorm:"column:is_directory;type:tinyint(1);not null"`
	Size        int64                        `json:"size" gorm:"column:size;type:bigint;not null"`
	MimeType    string                       `json:"mime_type" gorm:"column:mime_type;type:longtext;not null"`
	ParentID    *int64                       `json:"parent_id" gorm:"column:parent_id;type:bigint;null;index"`
	Permission  enums.FileItemPermissionEnum `json:"permission" gorm:"column:permission;type:longtext;not null"`
	UserID      int64                        `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt   time.Time                    `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time                    `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version     time.Time                    `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted   bool                         `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt   *time.Time                   `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Parent   *FileItem   `json:"parent" gorm:"foreignKey:ParentID"`
	Children []*FileItem `json:"children" gorm:"foreignKey:ParentID"`
	User     *User       `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for FileItem
func (FileItem) TableName() string {
	return "file_items"
}
