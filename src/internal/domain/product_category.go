package domain

import (
	"time"
)

// ProductCategory represents Product.Categories table
type ProductCategory struct {
	ID               int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name             string     `json:"name" gorm:"column:name;type:longtext;not null"`
	ParentCategoryID *int64     `json:"parent_category_id" gorm:"column:parent_category_id;type:bigint;null"`
	Order            int        `json:"order" gorm:"column:order;type:int;not null"`
	Description      string     `json:"description" gorm:"column:description;type:longtext;null"`
	Slug             string     `json:"slug" gorm:"column:slug;type:longtext;not null"`
	SeoTags          string     `json:"seo_tags" gorm:"column:seo_tags;type:longtext;null"`
	SiteID           int64      `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	CategoryID       *int64     `json:"category_id" gorm:"column:category_id;type:bigint;null;index"`
	UserID           int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt        time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version          time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted        bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt        *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	ParentCategory  *ProductCategory   `json:"parent_category" gorm:"foreignKey:CategoryID"`
	ChildCategories []*ProductCategory `json:"child_categories" gorm:"foreignKey:CategoryID"`
	Products        []Product          `json:"products" gorm:"many2many:category_product;"`
	Media           []Media            `json:"media" gorm:"many2many:category_media;"`
	Site            *Site              `json:"site" gorm:"foreignKey:SiteID"`
	User            *User              `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for ProductCategory
func (ProductCategory) TableName() string {
	return "categories"
}

// CategoryProduct represents Product.CategoryProduct table - a join table
type CategoryProduct struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ProductID  int64 `json:"product_id" gorm:"column:product_id;type:bigint;not null;index"`
	CategoryID int64 `json:"category_id" gorm:"column:category_id;type:bigint;not null;index"`
}

// TableName specifies the table name for CategoryProduct
func (CategoryProduct) TableName() string {
	return "category_product"
}

// ProductCategoryMedia represents Product.CategoryMedia table - a join table
type ProductCategoryMedia struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	CategoryID int64 `json:"category_id" gorm:"column:category_id;type:bigint;not null;index"`
	MediaID    int64 `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
}

// TableName specifies the table name for ProductCategoryMedia
func (ProductCategoryMedia) TableName() string {
	return "category_media"
}
