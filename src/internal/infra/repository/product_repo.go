package repository

import (
	"strings"

	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

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

func (r *ProductRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error) {
	var products []domain.Product
	var count int64

	query := r.database.Model(&domain.Product{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(products, paginationRequestDto, count)
}

func (r *ProductRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error) {
	var products []domain.Product
	var count int64

	query := r.database.Model(&domain.Product{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(products, paginationRequestDto, count)
}

func (r *ProductRepo) GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Product], error) {
	var products []domain.Product
	var count int64

	query := r.database.Model(&domain.Product{}).
		Joins("JOIN product_category ON product_category.product_id = products.id").
		Where("product_category.category_id = ?", categoryID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(products, paginationRequestDto, count)
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

func (r *ProductRepo) GetAllByFilterAndSort(
	siteID int64,
	filters map[enums.ProductFilterEnum][]string,
	sort *string,
	paginationRequestDto common.PaginationRequestDto,
) (*common.PaginationResponseDto[domain.Product], error) {
	var products []domain.Product
	var count int64

	query := r.database.Model(&domain.Product{}).
		Where("site_id = ?", siteID)

	// Apply filters if they exist
	if filters != nil {
		for filterType, values := range filters {
			switch filterType {
			case "price_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("price BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "rating_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("rate BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "selling_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("selling_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "visited_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("visited_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "review_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("review_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "weight_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("weight BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "added_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("created_at BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "updated_range":
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("updated_at BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "category_ids":
				if len(values) > 0 {
					categoryIds := strings.Join(values, ",")
					query = query.Joins("JOIN product_category pc ON pc.product_id = products.id").
						Where("pc.category_id IN (?)", categoryIds)
				}
			case "product_ids":
				if len(values) > 0 {
					productIds := strings.Join(values, ",")
					query = query.Where("id IN (?)", productIds)
				}
			case "free_send":
				if len(values) > 0 {
					query = query.Where("free_send = ?", values[0])
				}
			}
		}
	}

	// Apply sorting if specified
	if sort != nil {
		switch *sort {
		case "name_az":
			query = query.Order("name ASC")
		case "name_za":
			query = query.Order("name DESC")
		case "recently_added":
			query = query.Order("created_at DESC")
		case "recently_updated":
			query = query.Order("updated_at DESC")
		case "most_visited":
			query = query.Order("visited_count DESC")
		case "least_visited":
			query = query.Order("visited_count ASC")
		case "most_rated":
			query = query.Order("rate DESC")
		case "least_rated":
			query = query.Order("rate ASC")
		case "most_reviewed":
			query = query.Order("review_count DESC")
		case "least_reviewed":
			query = query.Order("review_count ASC")
		}
	} else {
		query = query.Order("updated_at DESC")
	}

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(products, paginationRequestDto, count)
}
