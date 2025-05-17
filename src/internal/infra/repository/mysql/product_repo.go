package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type ProductRepo struct {
	database *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		database: db,
	}
}

func (r *ProductRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Product, int64, error) {
	var products []domain.Product
	var count int64

	query := r.database.Model(&domain.Product{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, count, nil
}

func (r *ProductRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Product, int64, error) {
	var products []domain.Product
	var count int64

	query := r.database.Model(&domain.Product{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, count, nil
}

func (r *ProductRepo) GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Product, int64, error) {
	var products []domain.Product
	var count int64

	// For many-to-many relationship using the join table
	query := r.database.Model(&domain.Product{}).
		Joins("JOIN product_category ON product_category.product_id = products.id").
		Where("product_category.category_id = ?", categoryID)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, count, nil
}

func (r *ProductRepo) GetByID(id int64) (domain.Product, error) {
	var product domain.Product
	result := r.database.First(&product, id)
	if result.Error != nil {
		return product, result.Error
	}
	return product, nil
}

func (r *ProductRepo) GetBySlug(slug string) (domain.Product, error) {
	var product domain.Product
	result := r.database.Where("slug = ?", slug).First(&product)
	if result.Error != nil {
		return product, result.Error
	}
	return product, nil
}

func (r *ProductRepo) Create(product domain.Product) error {
	result := r.database.Create(&product)
	return result.Error
}

func (r *ProductRepo) Update(product domain.Product) error {
	result := r.database.Save(&product)
	return result.Error
}

func (r *ProductRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Product{}, id)
	return result.Error
}
