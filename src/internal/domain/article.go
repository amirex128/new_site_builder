package domain

import (
	"time"
)

// Article represents Blog.Articles table
type Article struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Title        string    `json:"title" gorm:"column:title;type:longtext;null"`
	Description  string    `json:"description" gorm:"column:description;type:longtext;null"`
	Body         string    `json:"body" gorm:"column:body;type:longtext;null"`
	Slug         string    `json:"slug" gorm:"column:slug;type:longtext;not null"`
	SiteID       int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	VisitedCount int       `json:"visited_count" gorm:"column:visited_count;type:int;not null"`
	ReviewCount  int       `json:"review_count" gorm:"column:review_count;type:int;not null"`
	Rate         int       `json:"rate" gorm:"column:rate;type:int;not null"`
	Badges       string    `json:"badges" gorm:"column:badges;type:longtext;null"`
	SeoTags      string    `json:"seo_tags" gorm:"column:seo_tags;type:longtext;null"`
	UserID       int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Categories []ArticleCategory `json:"categories" gorm:"many2many:article_category;"`
	Media      []Media           `json:"media" gorm:"many2many:article_media;"`
	Pages      []Page            `json:"pages" gorm:"many2many:page_article_usages;"`
	Site       *Site             `json:"site" gorm:"foreignKey:SiteID"`
	User       *User             `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Article
func (Article) TableName() string {
	return "articles"
}
func (m *Article) GetID() int64 {
	return m.ID
}
func (m *Article) GetUserID() *int64 {
	return &m.UserID
}
func (m *Article) GetCustomerID() *int64 {
	return nil
}
func (m *Article) GetSiteID() *int64 {
	return &m.SiteID
}

// ArticleMedia represents Blog.ArticleMedia table - a join table
type ArticleMedia struct {
	ID        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ArticleID int64 `json:"article_id" gorm:"column:article_id;type:bigint;not null;index"`
	MediaID   int64 `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
}

// TableName specifies the table name for ArticleMedia
func (ArticleMedia) TableName() string {
	return "article_media"
}

// ArticleArticleCategory represents Blog.ArticleArticleCategory table - a join table
type ArticleArticleCategory struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ArticleID  int64 `json:"article_id" gorm:"column:article_id;type:bigint;not null;index"`
	CategoryID int64 `json:"category_id" gorm:"column:category_id;type:bigint;not null;index"`
}

// TableName specifies the table name for ArticleArticleCategory
func (ArticleArticleCategory) TableName() string {
	return "article_category"
}
