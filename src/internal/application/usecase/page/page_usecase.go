package pageusecase

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/page"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type PageUsecase struct {
	logger           sflogger.Logger
	repo             repository.IPageRepository
	siteRepo         repository.ISiteRepository
	headerFooterRepo repository.IHeaderFooterRepository
	mediaRepo        repository.IMediaRepository
	authContextSvc   common.IAuthContextService
}

func NewPageUsecase(c contract.IContainer) *PageUsecase {
	return &PageUsecase{
		logger:           c.GetLogger(),
		repo:             c.GetPageRepo(),
		siteRepo:         c.GetSiteRepo(),
		headerFooterRepo: c.GetHeaderFooterRepo(),
		mediaRepo:        c.GetMediaRepo(),
		authContextSvc:   c.GetAuthContextTransientService(),
	}
}

func (u *PageUsecase) CreatePageCommand(params *page.CreatePageCommand) (any, error) {
	u.logger.Info("CreatePageCommand called", map[string]interface{}{
		"siteId": *params.SiteID,
		"title":  *params.Title,
		"slug":   *params.Slug,
	})

	// Check if site exists
	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check if header exists
	_, err = u.headerFooterRepo.GetByID(*params.HeaderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("هدر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check if footer exists
	_, err = u.headerFooterRepo.GetByID(*params.FooterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("فوتر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check if slug is unique for this site
	_, err = u.repo.GetBySlug(*params.Slug, *params.SiteID)
	if err == nil {
		return nil, errors.New("نامک (slug) تکراری است")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get user ID from auth context
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Convert body to JSON string
	bodyJSON, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}

	// Prepare SEO tags
	var seoTags string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTags = strings.Join(params.SeoTags, ",")
	}

	// Create the Page entity
	page := domain.Page{
		SiteID:      *params.SiteID,
		HeaderID:    *params.HeaderID,
		FooterID:    *params.FooterID,
		Slug:        *params.Slug,
		Title:       *params.Title,
		Description: getStringValueOrEmpty(params.Description),
		Body:        string(bodyJSON),
		SeoTags:     seoTags,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	// Create in repository
	err = u.repo.Create(page)
	if err != nil {
		return nil, err
	}

	// Add media if any
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		err = u.repo.AddMediaToPage(page.ID, params.MediaIDs)
		if err != nil {
			u.logger.Error("Failed to add media to page", map[string]interface{}{
				"pageId": page.ID,
				"error":  err.Error(),
			})
			// Continue, as this is not a critical error
		}
	}

	// Get created page with relations
	createdPage, err := u.repo.GetByID(page.ID)
	if err != nil {
		return nil, err
	}

	return enhancePageResponse(createdPage), nil
}

func (u *PageUsecase) UpdatePageCommand(params *page.UpdatePageCommand) (any, error) {
	u.logger.Info("UpdatePageCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing page
	existingPage, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("صفحه مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingPage.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این صفحه دسترسی ندارید")
	}

	// Verify site ID if provided
	if params.SiteID != nil {
		_, err = u.siteRepo.GetByID(*params.SiteID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("سایت مورد نظر یافت نشد")
			}
			return nil, err
		}
		existingPage.SiteID = *params.SiteID
	}

	// Verify header ID if provided
	if params.HeaderID != nil {
		_, err = u.headerFooterRepo.GetByID(*params.HeaderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("هدر مورد نظر یافت نشد")
			}
			return nil, err
		}
		existingPage.HeaderID = *params.HeaderID
	}

	// Verify footer ID if provided
	if params.FooterID != nil {
		_, err = u.headerFooterRepo.GetByID(*params.FooterID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("فوتر مورد نظر یافت نشد")
			}
			return nil, err
		}
		existingPage.FooterID = *params.FooterID
	}

	// Check if slug is unique for this site if it's being changed
	if params.Slug != nil && *params.Slug != existingPage.Slug {
		siteID := existingPage.SiteID
		if params.SiteID != nil {
			siteID = *params.SiteID
		}

		conflictPage, err := u.repo.GetBySlug(*params.Slug, siteID)
		if err == nil && conflictPage.ID != existingPage.ID {
			return nil, errors.New("نامک (slug) تکراری است")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		existingPage.Slug = *params.Slug
	}

	// Update other fields if provided
	if params.Title != nil {
		existingPage.Title = *params.Title
	}

	if params.Description != nil {
		existingPage.Description = *params.Description
	}

	if params.Body != nil {
		bodyJSON, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		existingPage.Body = string(bodyJSON)
	}

	// Update SEO tags if provided
	if params.SeoTags != nil {
		if len(params.SeoTags) > 0 {
			existingPage.SeoTags = strings.Join(params.SeoTags, ",")
		} else {
			existingPage.SeoTags = ""
		}
	}

	existingPage.UpdatedAt = time.Now()

	// Update in repository
	err = u.repo.Update(existingPage)
	if err != nil {
		return nil, err
	}

	// Update media if provided
	if params.MediaIDs != nil {
		// Remove existing media
		err = u.repo.RemoveAllMediaFromPage(existingPage.ID)
		if err != nil {
			u.logger.Error("Failed to remove existing media from page", map[string]interface{}{
				"pageId": existingPage.ID,
				"error":  err.Error(),
			})
			// Continue as this is not a critical error
		}

		// Add new media if any
		if len(params.MediaIDs) > 0 {
			err = u.repo.AddMediaToPage(existingPage.ID, params.MediaIDs)
			if err != nil {
				u.logger.Error("Failed to add media to page", map[string]interface{}{
					"pageId": existingPage.ID,
					"error":  err.Error(),
				})
				// Continue as this is not a critical error
			}
		}
	}

	// Get updated page with relations
	updatedPage, err := u.repo.GetByID(existingPage.ID)
	if err != nil {
		return nil, err
	}

	return enhancePageResponse(updatedPage), nil
}

func (u *PageUsecase) DeletePageCommand(params *page.DeletePageCommand) (any, error) {
	u.logger.Info("DeletePageCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing page
	existingPage, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("صفحه مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingPage.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این صفحه دسترسی ندارید")
	}

	// Delete the page
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": *params.ID,
	}, nil
}

func (u *PageUsecase) GetByIdPageQuery(params *page.GetByIdPageQuery) (any, error) {
	u.logger.Info("GetByIdPageQuery called", map[string]interface{}{
		"id":     params.ID,
		"ids":    params.IDs,
		"siteId": *params.SiteID,
	})

	// Check if we're looking for a single page or multiple pages
	if params.ID != nil {
		// Get single page by ID
		pageItem, err := u.repo.GetByIDAndSiteID(*params.ID, *params.SiteID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("صفحه مورد نظر یافت نشد")
			}
			return nil, err
		}

		return enhancePageResponse(pageItem), nil
	} else if params.IDs != nil && len(params.IDs) > 0 {
		// Get multiple pages by IDs
		pages, err := u.repo.GetByIDs(params.IDs, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Enhance each page response
		enhancedPages := make([]map[string]interface{}, 0, len(pages))
		for _, p := range pages {
			enhancedPages = append(enhancedPages, enhancePageResponse(p))
		}

		return enhancedPages, nil
	}

	return nil, errors.New("شناسه صفحه یا صفحات الزامی است")
}

func (u *PageUsecase) GetAllPageQuery(params *page.GetAllPageQuery) (any, error) {
	u.logger.Info("GetAllPageQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get pages for site with pagination
	pages, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Enhance each page response
	enhancedPages := make([]map[string]interface{}, 0, len(pages))
	for _, p := range pages {
		enhancedPages = append(enhancedPages, enhancePageResponse(p))
	}

	return map[string]interface{}{
		"items":     enhancedPages,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *PageUsecase) AdminGetAllPageQuery(params *page.AdminGetAllPageQuery) (any, error) {
	u.logger.Info("AdminGetAllPageQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all pages with pagination
	pages, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Enhance each page response
	enhancedPages := make([]map[string]interface{}, 0, len(pages))
	for _, p := range pages {
		enhancedPages = append(enhancedPages, enhancePageResponse(p))
	}

	return map[string]interface{}{
		"items":     enhancedPages,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *PageUsecase) GetByPathPageQuery(params *page.GetByPathPageQuery) (any, error) {
	u.logger.Info("GetByPathPageQuery called", map[string]interface{}{
		"paths":  params.Paths,
		"siteId": *params.SiteID,
	})

	// Get pages by paths (slugs)
	pages, err := u.repo.GetByPaths(params.Paths, *params.SiteID)
	if err != nil {
		return nil, err
	}

	// Enhance each page response
	enhancedPages := make([]map[string]interface{}, 0, len(pages))
	for _, p := range pages {
		enhancedPages = append(enhancedPages, enhancePageResponse(p))
	}

	return enhancedPages, nil
}

// Helper function to enhance page response with structured data
func enhancePageResponse(page domain.Page) map[string]interface{} {
	response := map[string]interface{}{
		"id":          page.ID,
		"siteId":      page.SiteID,
		"headerId":    page.HeaderID,
		"footerId":    page.FooterID,
		"slug":        page.Slug,
		"title":       page.Title,
		"description": page.Description,
		"createdAt":   page.CreatedAt,
		"updatedAt":   page.UpdatedAt,
	}

	// Parse body JSON if available
	if page.Body != "" {
		var bodyObj interface{}
		if err := json.Unmarshal([]byte(page.Body), &bodyObj); err == nil {
			response["body"] = bodyObj
		} else {
			// If parsing fails, return raw body
			response["body"] = page.Body
		}
	}

	// Parse SEO tags if available
	if page.SeoTags != "" {
		response["seoTags"] = strings.Split(page.SeoTags, ",")
	} else {
		response["seoTags"] = []string{}
	}

	// Add media information if available
	if page.Media != nil && len(page.Media) > 0 {
		mediaItems := make([]map[string]interface{}, 0, len(page.Media))
		for _, media := range page.Media {
			mediaItems = append(mediaItems, map[string]interface{}{
				"id":  media.ID,
				"url": "/api/media/" + strconv.FormatInt(media.ID, 10),
			})
		}
		response["media"] = mediaItems
	} else {
		response["media"] = []map[string]interface{}{}
	}

	return response
}

// Helper function to handle nil string pointers
func getStringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
