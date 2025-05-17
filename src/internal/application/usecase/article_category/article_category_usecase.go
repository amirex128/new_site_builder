package blogcategoryusecase

import (
	"encoding/json"
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

func NewBlogCategoryUsecase(c contract.IContainer) *ArticleCategoryUsecase {
	return &ArticleCategoryUsecase{
		logger: c.GetLogger(),
		repo:   c.GetArticleCategoryRepo(),
	}
}

func (u *ArticleCategoryUsecase) CreateCategoryCommand(params *article_category.CreateCategoryCommand) (any, error) {
	// Implementation for creating a new category
	var description string
	if params.Description != nil {
		description = *params.Description
	}

	// Convert slice of strings to a single string for SeoTags
	seoTags := ""
	if len(params.SeoTags) > 0 {
		seoTagsBytes, err := json.Marshal(params.SeoTags)
		if err == nil {
			seoTags = string(seoTagsBytes)
		}
	}

	category := domain.BlogCategory{
		Name:             *params.Name,
		SiteID:           *params.SiteID,
		Order:            *params.Order,
		ParentCategoryID: params.ParentCategoryID,
		Description:      description,
		Slug:             *params.Slug,
		SeoTags:          seoTags,
		UserID:           0, // This would need to come from the authenticated user
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	err := u.repo.Create(category)
	if err != nil {
		return nil, err
	}

	// TODO: Handle MediaIDs through a separate join table operation if needed

	return category, nil
}

func (u *ArticleCategoryUsecase) UpdateCategoryCommand(params *article_category.UpdateCategoryCommand) (any, error) {
	// Implementation for updating an existing category
	// First get the existing category
	existingCategory, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update the fields
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

	// Convert slice of strings to a single string for SeoTags
	if len(params.SeoTags) > 0 {
		seoTagsBytes, err := json.Marshal(params.SeoTags)
		if err == nil {
			existingCategory.SeoTags = string(seoTagsBytes)
		}
	}

	existingCategory.UpdatedAt = time.Now()

	err = u.repo.Update(existingCategory)
	if err != nil {
		return nil, err
	}

	// TODO: Handle MediaIDs through a separate join table operation if needed

	return existingCategory, nil
}

func (u *ArticleCategoryUsecase) DeleteCategoryCommand(params *article_category.DeleteCategoryCommand) (any, error) {
	// Implementation for deleting a category
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
	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *ArticleCategoryUsecase) GetAllCategoryQuery(params *article_category.GetAllCategoryQuery) (any, error) {
	// Implementation to get all categories by site ID
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
	// Implementation to get all categories for admin
	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
