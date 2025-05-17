package productcategoryusecase

import (
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_category"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ProductCategoryUsecase struct {
	logger sflogger.Logger
	repo   repository.IProductCategoryRepository
}

func NewProductCategoryUsecase(c contract.IContainer) *ProductCategoryUsecase {
	return &ProductCategoryUsecase{
		logger: c.GetLogger(),
		repo:   c.GetProductCategoryRepo(),
	}
}

func (u *ProductCategoryUsecase) CreateCategoryCommand(params *product_category.CreateCategoryCommand) (any, error) {
	// Implementation for creating a category
	fmt.Println(params)

	var description string
	if params.Description != nil {
		description = *params.Description
	}

	category := domain.ProductCategory{
		Name:             *params.Name,
		ParentCategoryID: params.ParentCategoryID,
		SiteID:           *params.SiteID,
		Order:            *params.Order,
		Description:      description,
		Slug:             *params.Slug,
		UserID:           1, // Should be set from authenticated user context
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	err := u.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *ProductCategoryUsecase) UpdateCategoryCommand(params *product_category.UpdateCategoryCommand) (any, error) {
	// Implementation for updating a category
	fmt.Println(params)

	existingCategory, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	existingCategory.Name = *params.Name
	existingCategory.SiteID = *params.SiteID

	if params.Order != nil {
		existingCategory.Order = *params.Order
	}

	existingCategory.ParentCategoryID = params.ParentCategoryID

	if params.Description != nil {
		existingCategory.Description = *params.Description
	}

	if params.Slug != nil {
		existingCategory.Slug = *params.Slug
	}

	existingCategory.UpdatedAt = time.Now()

	err = u.repo.Update(existingCategory)
	if err != nil {
		return nil, err
	}

	return existingCategory, nil
}

func (u *ProductCategoryUsecase) DeleteCategoryCommand(params *product_category.DeleteCategoryCommand) (any, error) {
	// Implementation for deleting a category
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *ProductCategoryUsecase) GetByIdCategoryQuery(params *product_category.GetByIdCategoryQuery) (any, error) {
	// Implementation to get category by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *ProductCategoryUsecase) GetAllCategoryQuery(params *product_category.GetAllCategoryQuery) (any, error) {
	// Implementation to get all categories
	fmt.Println(params)

	result, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *ProductCategoryUsecase) AdminGetAllCategoryQuery(params *product_category.AdminGetAllCategoryQuery) (any, error) {
	// Implementation to get all categories for admin
	fmt.Println(params)

	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
