package articlecategoryusecase

import (
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/article_category"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ArticleCategoryUsecase struct {
	*usecase.BaseUsecase
	categoryRepo repository.IArticleCategoryRepository
	mediaRepo    repository.IMediaRepository
}

func NewArticleCategoryUsecase(c contract.IContainer) *ArticleCategoryUsecase {
	return &ArticleCategoryUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		categoryRepo: c.GetArticleCategoryRepo(),
		mediaRepo:    c.GetMediaRepo(),
	}
}

func (u *ArticleCategoryUsecase) CreateCategoryCommand(params *article_category.CreateCategoryCommand) (*resp.Response, error) {
	// Implementation for creating a category based on .NET CreateCategoryCommand
	u.Logger.Info("Creating new category", map[string]interface{}{"name": *params.Name})

	// Convert SeoTags slice to string (comma-separated)
	var seoTags string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTags = strings.Join(params.SeoTags, ",")
	}

	var description string
	if params.Description != nil {
		description = *params.Description
	}

	// Create new category
	newCategory := domain.ArticleCategory{
		Name:             *params.Name,
		Slug:             *params.Slug,
		Description:      description,
		ParentCategoryID: params.ParentCategoryID,
		SiteID:           *params.SiteID,
		Order:            *params.Order,
		SeoTags:          seoTags,
		UserID:           1, // Should come from auth context in real impl
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// Save the category to get its ID
	err := u.categoryRepo.Create(newCategory)
	if err != nil {
		return nil, err
	}

	// Handle optional media relationships
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		for _, mediaID := range params.MediaIDs {
			err = u.categoryRepo.AddMediaToCategory(newCategory.ID, mediaID)
			if err != nil {
				u.Logger.Errorf("Failed to add media %d to category %d: %v", mediaID, newCategory.ID, err)
				// Continue with other media instead of failing completely
			}
		}
	}

	// Return the created category
	// Convert domain.ArticleCategory to map for response data
	categoryMap := map[string]interface{}{
		"id":               newCategory.ID,
		"name":             newCategory.Name,
		"slug":             newCategory.Slug,
		"description":      newCategory.Description,
		"parentCategoryId": newCategory.ParentCategoryID,
		"siteId":           newCategory.SiteID,
		"order":            newCategory.Order,
		"seoTags":          newCategory.SeoTags,
		"userId":           newCategory.UserID,
		"createdAt":        newCategory.CreatedAt,
		"updatedAt":        newCategory.UpdatedAt,
		"isDeleted":        newCategory.IsDeleted,
	}
	return resp.NewResponseData(resp.Created, categoryMap, "Article category created successfully"), nil
}

func (u *ArticleCategoryUsecase) UpdateCategoryCommand(params *article_category.UpdateCategoryCommand) (*resp.Response, error) {
	// Implementation for updating a category based on .NET UpdateCategoryCommand
	u.Logger.Info("Updating category", map[string]interface{}{"id": *params.ID})

	// Get existing category
	existingCategory, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this category
	// In a real implementation, check if the current user has rights to edit this category
	// This is equivalent to the gate.HasUserAccess(entity) call in .NET

	// Update fields if provided
	if params.Name != nil {
		existingCategory.Name = *params.Name
	}

	if params.Description != nil {
		existingCategory.Description = *params.Description
	}

	if params.ParentCategoryID != nil {
		existingCategory.ParentCategoryID = params.ParentCategoryID
	}

	if params.Slug != nil {
		existingCategory.Slug = *params.Slug
	}

	if params.Order != nil {
		existingCategory.Order = *params.Order
	}

	// Update SeoTags if provided
	if params.SeoTags != nil {
		existingCategory.SeoTags = strings.Join(params.SeoTags, ",")
	}

	existingCategory.UpdatedAt = time.Now()

	// Save changes
	err = u.categoryRepo.Update(existingCategory)
	if err != nil {
		return nil, err
	}

	// Handle media relationships if provided
	if params.MediaIDs != nil {
		// First remove all existing media associations
		err = u.categoryRepo.RemoveAllMediaFromCategory(existingCategory.ID)
		if err != nil {
			u.Logger.Errorf("Failed to remove media from category %d: %v", existingCategory.ID, err)
		}

		// Then add the new ones
		for _, mediaID := range params.MediaIDs {
			err = u.categoryRepo.AddMediaToCategory(existingCategory.ID, mediaID)
			if err != nil {
				u.Logger.Errorf("Failed to add media %d to category %d: %v", mediaID, existingCategory.ID, err)
			}
		}
	}

	// Convert domain.ArticleCategory to map for response data
	categoryMap := map[string]interface{}{
		"id":               existingCategory.ID,
		"name":             existingCategory.Name,
		"slug":             existingCategory.Slug,
		"description":      existingCategory.Description,
		"parentCategoryId": existingCategory.ParentCategoryID,
		"siteId":           existingCategory.SiteID,
		"order":            existingCategory.Order,
		"seoTags":          existingCategory.SeoTags,
		"userId":           existingCategory.UserID,
		"createdAt":        existingCategory.CreatedAt,
		"updatedAt":        existingCategory.UpdatedAt,
		"isDeleted":        existingCategory.IsDeleted,
	}
	return resp.NewResponseData(resp.Updated, categoryMap, "Article category updated successfully"), nil
}

func (u *ArticleCategoryUsecase) DeleteCategoryCommand(params *article_category.DeleteCategoryCommand) (*resp.Response, error) {
	// Implementation for deleting a category based on .NET DeleteCategoryCommand
	u.Logger.Info("Deleting category", map[string]interface{}{"id": *params.ID})

	// Get the category first to ensure it exists
	_, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this category
	// In a real implementation, check if the current user has rights to delete this category
	// This is equivalent to the gate.HasUserAccess(entity) call in .NET

	// Delete the category
	err = u.categoryRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponse(resp.Deleted, "Article category deleted successfully"), nil
}

func (u *ArticleCategoryUsecase) GetByIdCategoryQuery(params *article_category.GetByIdCategoryQuery) (*resp.Response, error) {
	// Implementation to get category by ID based on .NET GetByIdCategoryQuery
	u.Logger.Info("Getting category by ID", map[string]interface{}{"id": *params.ID})

	// Get the category
	result, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this category
	// In a real implementation, check if the current user has rights to view this category
	// This is equivalent to the gate.HasUserAccess(entity) call in .NET

	// Get media information
	mediaItems, err := u.categoryRepo.GetCategoryMedia(result.ID)
	if err != nil {
		u.Logger.Errorf("Failed to get media for category %d: %v", result.ID, err)
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"category": result,
		"media":    mediaItems,
	}, "Article category retrieved successfully"), nil
}

func (u *ArticleCategoryUsecase) GetAllCategoryQuery(params *article_category.GetAllCategoryQuery) (*resp.Response, error) {
	// Implementation to get all categories by site ID, based on .NET GetAllCategoryQuery
	u.Logger.Info("Getting all categories by site ID", map[string]interface{}{"siteID": *params.SiteID})

	// Check if user has access to this site
	// In a real implementation, check if the current user has rights to view categories for this site
	// This is equivalent to the gate.HasSiteAccess(request.SiteId) call in .NET

	result, count, err := u.categoryRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// For each category, get media information
	// In a more efficient implementation, this would be done in a single query
	categoriesWithMedia := make([]map[string]interface{}, len(result))
	for i, category := range result {
		media, err := u.categoryRepo.GetCategoryMedia(category.ID)
		if err != nil {
			u.Logger.Errorf("Failed to get media for category %d: %v", category.ID, err)
			media = []domain.Media{}
		}

		categoriesWithMedia[i] = map[string]interface{}{
			"category": category,
			"media":    media,
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items": categoriesWithMedia,
		"total": count,
	}, "Article categories retrieved successfully"), nil
}

func (u *ArticleCategoryUsecase) AdminGetAllCategoryQuery(params *article_category.AdminGetAllCategoryQuery) (*resp.Response, error) {
	// Implementation to get all categories for admin, based on .NET AdminGetAllCategoryQuery
	u.Logger.Info("Admin getting all categories", map[string]interface{}{})

	// Check if user has admin access
	// In a real implementation, check if the current user has admin rights
	// This is equivalent to the gate.HasSiteAccess() call in .NET

	result, count, err := u.categoryRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// For each category, get media information
	categoriesWithMedia := make([]map[string]interface{}, len(result))
	for i, category := range result {
		media, err := u.categoryRepo.GetCategoryMedia(category.ID)
		if err != nil {
			u.Logger.Errorf("Failed to get media for category %d: %v", category.ID, err)
			media = []domain.Media{}
		}

		categoriesWithMedia[i] = map[string]interface{}{
			"category": category,
			"media":    media,
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items": categoriesWithMedia,
		"total": count,
	}, "Article categories retrieved successfully (Admin)"), nil
}
