package domain

import (
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// ProductAttribute represents Product.ProductAttributes table
type ProductAttribute struct {
	ID        int64                          `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	ProductID int64                          `json:"product_id" gorm:"column:product_id;type:bigint;not null;index"`
	Type      enums.ProductAttributeTypeEnum `json:"type" gorm:"column:type;type:longtext;not null"`
	Name      string                         `json:"name" gorm:"column:name;type:longtext;not null"`
	Value     string                         `json:"value" gorm:"column:value;type:longtext;not null"`
	CreatedAt time.Time                      `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt time.Time                      `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version   time.Time                      `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted bool                           `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time                     `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Product *Product `json:"article" gorm:"foreignKey:ProductID"`
}

// TableName specifies the table name for ProductAttribute
func (ProductAttribute) TableName() string {
	return "product_attributes"
}
func (m *ProductAttribute) GetID() int64 {
	return m.ID
}
func (m *ProductAttribute) GetUserID() *int64 {
	return nil
}
func (m *ProductAttribute) GetCustomerID() *int64 {
	return nil
}
func (m *ProductAttribute) GetSiteID() *int64 {
	return nil
}
