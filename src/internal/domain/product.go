package domain

import (
	"time"
)

// Product represents Product.Products table
type Product struct {
	ID              int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name            string     `json:"name" gorm:"column:name;type:longtext;not null"`
	Description     string     `json:"description" gorm:"column:description;type:longtext;null"`
	Status          string     `json:"status" gorm:"column:status;type:longtext;not null"`
	Weight          int        `json:"weight" gorm:"column:weight;type:int;not null"`
	SellingCount    int        `json:"selling_count" gorm:"column:selling_count;type:int;not null"`
	VisitedCount    int        `json:"visited_count" gorm:"column:visited_count;type:int;not null"`
	ReviewCount     int        `json:"review_count" gorm:"column:review_count;type:int;not null"`
	Rate            int        `json:"rate" gorm:"column:rate;type:int;not null"`
	Badges          string     `json:"badges" gorm:"column:badges;type:longtext;null"`
	FreeSend        bool       `json:"free_send" gorm:"column:free_send;type:tinyint(1);not null"`
	LongDescription string     `json:"long_description" gorm:"column:long_description;type:longtext;null"`
	Slug            string     `json:"slug" gorm:"column:slug;type:longtext;not null"`
	SeoTags         string     `json:"seo_tags" gorm:"column:seo_tags;type:longtext;null"`
	SiteID          int64      `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	UserID          int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version         time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted       bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Categories []ProductCategory `json:"categories" gorm:"many2many:category_product;"`
	Media      []Media           `json:"media" gorm:"many2many:product_media;"`
	Pages      []Page            `json:"pages" gorm:"many2many:page_product_usages;"`
	Site       *Site             `json:"site" gorm:"foreignKey:SiteID"`
	User       *User             `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// ProductMedia represents Product.ProductMedia table - a join table
type ProductMedia struct {
	ID        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ProductID int64 `json:"product_id" gorm:"column:product_id;type:bigint;not null;index"`
	MediaID   int64 `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
}

// TableName specifies the table name for ProductMedia
func (ProductMedia) TableName() string {
	return "product_media"
}
