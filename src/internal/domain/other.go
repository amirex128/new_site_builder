package domain

import (
	"time"
)

// Media is a reference model for all tables with media relationships
type Media struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`

	// Relations
	Articles       []Article        `json:"articles" gorm:"many2many:article_media;"`
	BlogCategories []BlogCategory   `json:"blog_categories" gorm:"many2many:category_media;"`
	ProductCategories []ProductCategory `json:"product_categories" gorm:"many2many:category_media;"`
	Products       []Product        `json:"products" gorm:"many2many:product_media;"`
	Pages          []Page           `json:"pages" gorm:"many2many:page_media;"`
	Tickets        []Ticket         `json:"tickets" gorm:"many2many:ticket_media;"`
	CustomerTickets []CustomerTicket `json:"customer_tickets" gorm:"many2many:customer_ticket_media;"`
	DefaultThemes  []DefaultTheme   `json:"default_themes" gorm:"foreignKey:MediaID"`
}

// TableName specifies the table name for Media
func (Media) TableName() string {
	return "media"
}

// FileItem represents Drive.FileItems table
type FileItem struct {
	ID          int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name        string     `json:"name" gorm:"column:name;type:longtext;not null"`
	BucketName  string     `json:"bucket_name" gorm:"column:bucket_name;type:longtext;not null"`
	ServerKey   string     `json:"server_key" gorm:"column:server_key;type:longtext;not null"`
	FilePath    string     `json:"file_path" gorm:"column:file_path;type:longtext;not null"`
	IsDirectory bool       `json:"is_directory" gorm:"column:is_directory;type:tinyint(1);not null"`
	Size        int64      `json:"size" gorm:"column:size;type:bigint;not null"`
	MimeType    string     `json:"mime_type" gorm:"column:mime_type;type:longtext;not null"`
	ParentID    *int64     `json:"parent_id" gorm:"column:parent_id;type:bigint;null;index"`
	Permission  string     `json:"permission" gorm:"column:permission;type:longtext;not null"`
	UserID      int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version     time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted   bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Parent   *FileItem   `json:"parent" gorm:"foreignKey:ParentID"`
	Children []*FileItem `json:"children" gorm:"foreignKey:ParentID"`
	User     *User       `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for FileItem
func (FileItem) TableName() string {
	return "file_items"
}

// Storage represents Drive.Storages table
type Storage struct {
	ID          int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	UsedSpaceKb int64      `json:"used_space_kb" gorm:"column:used_space_kb;type:bigint;not null"`
	QuotaKb     int64      `json:"quota_kb" gorm:"column:quota_kb;type:bigint;not null"`
	ChargedAt   time.Time  `json:"charged_at" gorm:"column:charged_at;type:datetime(6);not null"`
	ExpireAt    time.Time  `json:"expire_at" gorm:"column:expire_at;type:datetime(6);not null"`
	UserID      int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version     time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted   bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	User *User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Storage
func (Storage) TableName() string {
	return "storages"
}

// Credit represents Ai.Credits table
type Credit struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	UserID     int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID int64      `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version    time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted  bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

// TableName specifies the table name for Credit
func (Credit) TableName() string {
	return "credits"
} 