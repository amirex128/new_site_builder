package domain

import (
	"time"
)

// HeaderFooter represents Site.HeaderFooters table
type HeaderFooter struct {
	ID        int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID    int64      `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	Title     string     `json:"title" gorm:"column:title;type:longtext;not null"`
	IsMain    bool       `json:"is_main" gorm:"column:is_main;type:tinyint(1);not null"`
	Body      string     `json:"body" gorm:"column:body;type:longtext;null"`
	Type      string     `json:"type" gorm:"column:type;type:varchar(255);not null"`
	UserID    int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version   time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Site  *Site  `json:"site" gorm:"foreignKey:SiteID"`
	User  *User  `json:"user" gorm:"foreignKey:UserID"`
	Pages []Page `json:"pages" gorm:"many2many:page_header_footer_usages;"`
}

// TableName specifies the table name for HeaderFooter
func (HeaderFooter) TableName() string {
	return "header_footers"
}
