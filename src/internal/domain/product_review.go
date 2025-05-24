package domain

import (
	"time"
)

// ProductReview represents Product.ProductReviews table
type ProductReview struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Rating     int       `json:"rating" gorm:"column:rating;type:int;not null"`
	Like       int       `json:"like" gorm:"column:like;type:int;not null"`
	Dislike    int       `json:"dislike" gorm:"column:dislike;type:int;not null"`
	Approved   bool      `json:"approved" gorm:"column:approved;type:tinyint(1);not null"`
	ReviewText string    `json:"review_text" gorm:"column:review_text;type:longtext;not null"`
	ProductID  int64     `json:"product_id" gorm:"column:product_id;type:bigint;not null"`
	SiteID     int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	UserID     int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID int64     `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Product  *Product  `json:"article" gorm:"foreignKey:ProductID"`
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	Site     *Site     `json:"site" gorm:"foreignKey:SiteID"`
}

// TableName specifies the table name for ProductReview
func (ProductReview) TableName() string {
	return "product_reviews"
}
func (m *ProductReview) GetID() int64 {
	return m.ID
}
func (m *ProductReview) GetUserID() *int64 {
	return &m.UserID
}
func (m *ProductReview) GetCutomerID() *int64 {
	return &m.CustomerID
}
func (m *ProductReview) GetSiteID() *int64 {
	return &m.SiteID
}
