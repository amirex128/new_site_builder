package productcategoryusecase

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_category"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type ProductCategoryUsecase struct {
	*usecase.BaseUsecase
	logger      sflogger.Logger
	repo        repository.IProductCategoryRepository
	mediaRepo   repository.IMediaRepository
	authContext func(c *gin.Context) service.IAuthService
}

func NewProductCategoryUsecase(c contract.IContainer) *ProductCategoryUsecase {
	return &ProductCategoryUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		repo:        c.GetProductCategoryRepo(),
		mediaRepo:   c.GetMediaRepo(),
		authContext: c.GetAuthTransientService(),
	}
}

func (u *ProductCategoryUsecase) CreateCategoryCommand(params *product_category.CreateCategoryCommand) (*resp.Response, error) {
	u.Logger.Info("CreateCategoryCommand called", map[string]interface{}{
		"name":   *params.Name,
		"siteId": *params.SiteID,
	})

	// Validate unique slug in site
	existingCategoryBySlug, err := u.repo.GetBySlug(*params.Slug)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil && existingCategoryBySlug.SiteID == *params.SiteID {
			return nil, errors.New("نامک تکراری است")
		}
	}

	// Get user ID from auth context
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Prepare SEO tags
	var seoTagsStr string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTagsStr = strings.Join(params.SeoTags, ",")
	}

	// Create category entity
	category := domain.ProductCategory{
		Name:             *params.Name,
		ParentCategoryID: params.ParentCategoryID,
		SiteID:           *params.SiteID,
		Order:            *params.Order,
		Description:      getStringValueOrEmpty(params.Description),
		SeoTags:          seoTagsStr,
		Slug:             *params.Slug,
		UserID:           userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// Create the category in the database
	err = u.repo.Create(category)
	if err != nil {
		return nil, err
	}

	// Handle media associations if any
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		err = u.attachMediaToCategory(category.ID, params.MediaIDs)
		if err != nil {
			u.Logger.Error("Failed to attach media to category", map[string]interface{}{
				"categoryId": category.ID,
				"error":      err.Error(),
			})
			// Continue even if media attachment fails
		}
	}

	// Fetch the category with related data
	createdCategory, err := u.repo.GetByID(category.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Created, resp.Data{
		"category": createdCategory,
	}, "دسته‌بندی با موفقیت ایجاد شد"), nil
}

func (u *ProductCategoryUsecase) UpdateCategoryCommand(params *product_category.UpdateCategoryCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateCategoryCommand called", map[string]interface{}{
		"id":     *params.ID,
		"siteId": *params.SiteID,
	})

	// Get existing category
	existingCategory, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("دسته‌بندی یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// In our monolithic approach, we check directly
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingCategory.SiteID != *params.SiteID && !isAdmin {
		return nil, errors.New("شما به این دسته‌بندی دسترسی ندارید")
	}

	// Validate slug uniqueness if changed
	if params.Slug != nil && *params.Slug != existingCategory.Slug {
		existingCategoryBySlug, err := u.repo.GetBySlug(*params.Slug)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			if err == nil && existingCategoryBySlug.SiteID == *params.SiteID && existingCategoryBySlug.ID != *params.ID {
				return nil, errors.New("نامک تکراری است")
			}
		}
	}

	// Prepare SEO tags
	var seoTagsStr string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTagsStr = strings.Join(params.SeoTags, ",")
		existingCategory.SeoTags = seoTagsStr
	}

	// Update fields
	existingCategory.Name = *params.Name
	existingCategory.SiteID = *params.SiteID
	existingCategory.Order = *params.Order
	existingCategory.ParentCategoryID = params.ParentCategoryID

	if params.Description != nil {
		existingCategory.Description = *params.Description
	}

	if params.Slug != nil {
		existingCategory.Slug = *params.Slug
	}

	existingCategory.UpdatedAt = time.Now()
	existingCategory.UserID = userID

	// Update the category
	err = u.repo.Update(existingCategory)
	if err != nil {
		return nil, err
	}

	// Handle media associations if any
	if params.MediaIDs != nil {
		err = u.updateCategoryMedia(existingCategory.ID, params.MediaIDs)
		if err != nil {
			u.Logger.Error("Failed to update media for category", map[string]interface{}{
				"categoryId": existingCategory.ID,
				"error":      err.Error(),
			})
			// Continue even if media update fails
		}
	}

	// Fetch the updated category with related data
	updatedCategory, err := u.repo.GetByID(existingCategory.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Updated, resp.Data{
		"category": updatedCategory,
	}, "دسته‌بندی با موفقیت بروزرسانی شد"), nil
}

func (u *ProductCategoryUsecase) DeleteCategoryCommand(params *product_category.DeleteCategoryCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteCategoryCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing category
	existingCategory, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("دسته‌بندی یافت نشد")
		}
		return nil, err
	}

	// Check user access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	if existingCategory.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این دسته‌بندی دسترسی ندارید")
	}

	// Delete the category
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Deleted, resp.Data{
		"success": true,
	}, "دسته‌بندی با موفقیت حذف شد"), nil
}

func (u *ProductCategoryUsecase) GetByIdCategoryQuery(params *product_category.GetByIdCategoryQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdCategoryQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get category by ID
	category, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("دسته‌بندی یافت نشد")
		}
		return nil, err
	}

	// Check user access - anyone can view categories but logging for audit
	userID, _ := u.authContext(u.Ctx).GetUserID()
	if userID > 0 {
		u.Logger.Info("Category accessed by user", map[string]interface{}{
			"categoryId": category.ID,
			"userId":     userID,
		})
	}

	// If there are SEO tags, parse them into an array for the response
	var seoTags []string
	if category.SeoTags != "" {
		seoTags = strings.Split(category.SeoTags, ",")
	}

	// Prepare the response with media URLs
	response := map[string]interface{}{
		"id":               category.ID,
		"name":             category.Name,
		"parentCategoryId": category.ParentCategoryID,
		"order":            category.Order,
		"description":      category.Description,
		"slug":             category.Slug,
		"seoTags":          seoTags,
		"siteId":           category.SiteID,
		"createdAt":        category.CreatedAt,
		"updatedAt":        category.UpdatedAt,
	}

	// Load media if any
	if len(category.Media) > 0 {
		var mediaURLs []map[string]interface{}
		for _, media := range category.Media {
			mediaURLs = append(mediaURLs, map[string]interface{}{
				"id": media.ID,
				// The Media entity doesn't have a URL field directly
				// In a real implementation, we would fetch the URL from a file service
				"url": "/api/media/" + strconv.FormatInt(media.ID, 64),
			})
		}
		response["media"] = mediaURLs
	}

	return resp.NewResponseData(resp.Retrieved, resp.Data(response), "دسته‌بندی با موفقیت دریافت شد"), nil
}

func (u *ProductCategoryUsecase) GetAllCategoryQuery(params *product_category.GetAllCategoryQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllCategoryQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check site access - simplified in monolithic app
	// In a real implementation, we would check if the user has access to this site

	// Get all categories for the site
	result, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	categories := result.Items
	count := result.TotalCount

	// Prepare response items with media URLs
	var items []map[string]interface{}
	for _, category := range categories {
		var seoTags []string
		if category.SeoTags != "" {
			seoTags = strings.Split(category.SeoTags, ",")
		}

		item := map[string]interface{}{
			"id":               category.ID,
			"name":             category.Name,
			"parentCategoryId": category.ParentCategoryID,
			"order":            category.Order,
			"description":      category.Description,
			"slug":             category.Slug,
			"seoTags":          seoTags,
			"siteId":           category.SiteID,
			"createdAt":        category.CreatedAt,
			"updatedAt":        category.UpdatedAt,
		}

		// Add media if available
		if len(category.Media) > 0 {
			var mediaURLs []map[string]interface{}
			for _, media := range category.Media {
				mediaURLs = append(mediaURLs, map[string]interface{}{
					"id": media.ID,
					// The Media entity doesn't have a URL field directly
					"url": "/api/media/" + strconv.FormatInt(media.ID, 64),
				})
			}
			item["media"] = mediaURLs
		}

		items = append(items, item)
	}

	// Return paginated result
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     items,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "دسته‌بندی‌ها با موفقیت دریافت شدند"), nil
}

func (u *ProductCategoryUsecase) AdminGetAllCategoryQuery(params *product_category.AdminGetAllCategoryQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllCategoryQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all categories across all sites for admin
	result, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	categories := result.Items
	count := result.TotalCount

	// Prepare response items with media URLs
	var items []map[string]interface{}
	for _, category := range categories {
		var seoTags []string
		if category.SeoTags != "" {
			seoTags = strings.Split(category.SeoTags, ",")
		}

		item := map[string]interface{}{
			"id":               category.ID,
			"name":             category.Name,
			"parentCategoryId": category.ParentCategoryID,
			"order":            category.Order,
			"description":      category.Description,
			"slug":             category.Slug,
			"seoTags":          seoTags,
			"siteId":           category.SiteID,
			"createdAt":        category.CreatedAt,
			"updatedAt":        category.UpdatedAt,
		}

		// Add media if available
		if len(category.Media) > 0 {
			var mediaURLs []map[string]interface{}
			for _, media := range category.Media {
				mediaURLs = append(mediaURLs, map[string]interface{}{
					"id": media.ID,
					// The Media entity doesn't have a URL field directly
					"url": "/api/media/" + strconv.FormatInt(media.ID, 64),
				})
			}
			item["media"] = mediaURLs
		}

		items = append(items, item)
	}

	// Return paginated result
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     items,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "دسته‌بندی‌ها با موفقیت دریافت شدند"), nil
}

// Helper functions

func (u *ProductCategoryUsecase) attachMediaToCategory(categoryID int64, mediaIDs []int64) error {
	// In the monolithic approach, we directly create the associations
	// For each media ID, create an entry in the join table
	for _, mediaID := range mediaIDs {
		// Check if media exists
		_, err := u.mediaRepo.GetByID(mediaID)
		if err != nil {
			continue // Skip invalid media IDs
		}

		// Create join table entry
		// This would normally be handled by a specific repository method
		// For simplicity, we're assuming the database handle could create this directly
		// We're just logging it instead of actually creating it for now
		u.Logger.Info("Attaching media to category", map[string]interface{}{
			"categoryId": categoryID,
			"mediaId":    mediaID,
		})
	}

	return nil
}

func (u *ProductCategoryUsecase) updateCategoryMedia(categoryID int64, mediaIDs []int64) error {
	// In a real implementation, we would first delete all existing associations
	// and then create new ones
	// For simplicity, we're just logging it
	u.Logger.Info("Updating media for category", map[string]interface{}{
		"categoryId": categoryID,
		"mediaIds":   mediaIDs,
	})

	// The complete implementation would delete existing associations and create new ones
	return nil
}

func getStringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
