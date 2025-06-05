package domain

import (
	"time"
)

// Page represents Site.Pages table
type Page struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID      int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	HeaderID    int64     `json:"header_id" gorm:"column:header_id;type:bigint;not null"`
	FooterID    int64     `json:"footer_id" gorm:"column:footer_id;type:bigint;not null"`
	Slug        string    `json:"slug" gorm:"column:slug;type:longtext;not null"`
	Title       string    `json:"title" gorm:"column:title;type:longtext;not null"`
	Description string    `json:"description" gorm:"column:description;type:longtext;null"`
	Body        string    `json:"body" gorm:"column:body;type:longtext;null"`
	SeoTags     string    `json:"seo_tags" gorm:"column:seo_tags;type:longtext;null"`
	UserID      int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Site          *Site          `json:"site" gorm:"foreignKey:SiteID"`
	Header        *HeaderFooter  `json:"header" gorm:"foreignKey:HeaderID"`
	Footer        *HeaderFooter  `json:"footer" gorm:"foreignKey:FooterID"`
	Media         []Media        `json:"media" gorm:"many2many:page_media;"`
	Articles      []Article      `json:"articles" gorm:"many2many:page_article_usages;"`
	Products      []Product      `json:"products" gorm:"many2many:page_product_usages;"`
	HeaderFooters []HeaderFooter `json:"header_footers" gorm:"many2many:page_header_footer_usages;"`
}

// TableName specifies the table name for Page
func (Page) TableName() string {
	return "pages"
}
func (m *Page) GetID() int64 {
	return m.ID
}
func (m *Page) GetUserID() *int64 {
	return &m.UserID
}
func (m *Page) GetCustomerID() *int64 {
	return nil
}
func (m *Page) GetSiteID() *int64 {
	return &m.SiteID
}

// PageMedia represents Site.PageMedia table - a join table
type PageMedia struct {
	ID      int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	PageID  int64 `json:"page_id" gorm:"column:page_id;type:bigint;not null;index"`
	MediaID int64 `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
}

// TableName specifies the table name for PageMedia
func (PageMedia) TableName() string {
	return "page_media"
}
