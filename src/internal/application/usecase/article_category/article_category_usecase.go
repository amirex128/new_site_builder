package articlecategoryusecase

import (
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/article_category"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ArticleCategoryUsecase struct {
	logger sflogger.Logger
	repo   repository.IArticleCategoryRepository
}

func NewArticleCategoryUsecase(c contract.IContainer) *ArticleCategoryUsecase {
	return &ArticleCategoryUsecase{
		logger: c.GetLogger(),
		repo:   c.GetArticleCategoryRepo(),
	}
}

func (u *ArticleCategoryUsecase) CreateCategoryCommand(params *article_category.CreateCategoryCommand) (any, error) {
	// Implementation for creating a category
	fmt.Println(params)

	var description string
	if params.Description != nil {
		description = *params.Description
	}

	category := domain.BlogCategory{
		Name:             *params.Name,
		Slug:             *params.Slug,
		Description:      description,
		ParentCategoryID: params.ParentCategoryID,
		SiteID:           *params.SiteID,
		Order:            *params.Order,
		UserID:           1, // Should come from auth context
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

func (u *ArticleCategoryUsecase) UpdateCategoryCommand(params *article_category.UpdateCategoryCommand) (any, error) {
	// Implementation for updating a category
	fmt.Println(params)

	existingCategory, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	if params.Name != nil {
		existingCategory.Name = *params.Name
	}

	if params.Slug != nil {
		existingCategory.Slug = *params.Slug
	}

	if params.Description != nil {
		existingCategory.Description = *params.Description
	}

	if params.ParentCategoryID != nil {
		existingCategory.ParentCategoryID = params.ParentCategoryID
	}

	if params.SiteID != nil {
		existingCategory.SiteID = *params.SiteID
	}

	if params.Order != nil {
		existingCategory.Order = *params.Order
	}

	existingCategory.UpdatedAt = time.Now()

	err = u.repo.Update(existingCategory)
	if err != nil {
		return nil, err
	}

	return existingCategory, nil
}

func (u *ArticleCategoryUsecase) DeleteCategoryCommand(params *article_category.DeleteCategoryCommand) (any, error) {
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

func (u *ArticleCategoryUsecase) GetByIdCategoryQuery(params *article_category.GetByIdCategoryQuery) (any, error) {
	// Implementation to get category by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *ArticleCategoryUsecase) GetAllCategoryQuery(params *article_category.GetAllCategoryQuery) (any, error) {
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

func (u *ArticleCategoryUsecase) AdminGetAllCategoryQuery(params *article_category.AdminGetAllCategoryQuery) (any, error) {
	// Implementation for admin to get all categories
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
