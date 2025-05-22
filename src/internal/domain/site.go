package domain

import (
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// Site represents Site.Sites table
type Site struct {
	ID         int64                `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Domain     string               `json:"domain" gorm:"column:domain;type:longtext;not null"`
	DomainType enums.DomainTypeEnum `json:"domain_type" gorm:"column:domain_type;type:longtext;not null"`
	Name       string               `json:"name" gorm:"column:name;type:longtext;not null"`
	Status     enums.StatusEnum     `json:"status" gorm:"column:status;type:longtext;not null"`
	SiteType   enums.SiteTypeEnum   `json:"site_type" gorm:"column:site_type;type:longtext;not null"`
	UserID     int64                `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt  time.Time            `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time            `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version    time.Time            `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted  bool                 `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt  *time.Time           `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Setting       *Setting       `json:"setting" gorm:"foreignKey:SiteID"`
	Pages         []Page         `json:"pages" gorm:"foreignKey:SiteID"`
	HeaderFooters []HeaderFooter `json:"header_footers" gorm:"foreignKey:SiteID"`
	User          *User          `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Site
func (Site) TableName() string {
	return "sites"
}
