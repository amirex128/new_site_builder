package pageusageusecase

import (
	"errors"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/page_usage"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"gorm.io/gorm"
)

type PageUsageUsecase struct {
	*usecase.BaseUsecase
	pageRepo                  repository.IPageRepository
	pageArticleUsageRepo      repository.IPageArticleUsageRepository
	pageProductUsageRepo      repository.IPageProductUsageRepository
	pageHeaderFooterUsageRepo repository.IPageHeaderFooterUsageRepository
	articleRepo               repository.IArticleRepository
	productRepo               repository.IProductRepository
	headerFooterRepo          repository.IHeaderFooterRepository
	siteRepo                  repository.ISiteRepository
	authContext               func(c *gin.Context) service.IAuthService
}

func NewPageUsageUsecase(c contract.IContainer) *PageUsageUsecase {
	return &PageUsageUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		pageRepo:                  c.GetPageRepo(),
		pageArticleUsageRepo:      c.GetPageArticleUsageRepo(),
		pageProductUsageRepo:      c.GetPageProductUsageRepo(),
		pageHeaderFooterUsageRepo: c.GetPageHeaderFooterUsageRepo(),
		articleRepo:               c.GetArticleRepo(),
		productRepo:               c.GetProductRepo(),
		headerFooterRepo:          c.GetHeaderFooterRepo(),
		siteRepo:                  c.GetSiteRepo(),
		authContext:               c.GetAuthTransientService(),
	}
}

func (u *PageUsageUsecase) SyncPageUsageCommand(params *page_usage.SyncPageUsageCommand) (any, error) {
	u.Logger.Info("SyncPageUsageCommand called", map[string]interface{}{
		"pageId":    *params.PageID,
		"siteId":    *params.SiteID,
		"type":      params.Type,
		"entityIds": params.EntityIDs,
	})

	// Check if page exists
	page, err := u.pageRepo.GetByIDAndSiteID(*params.PageID, *params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("صفحه مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check if site exists
	_, err = u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Get user ID from auth context
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Check user access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if page.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این صفحه دسترسی ندارید")
	}

	// Handle different types of usages
	switch params.Type {
	case enums.ArticleUsage:
		// Delete existing usages
		err = u.pageArticleUsageRepo.DeleteByPageIDAndSiteID(*params.PageID, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Create new usages
		if len(params.EntityIDs) > 0 {
			// Verify articles exist
			for _, articleID := range params.EntityIDs {
				_, err = u.articleRepo.GetByID(articleID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
			}

			// Create usages
			usages := make([]domain.PageArticleUsage, 0, len(params.EntityIDs))
			for _, articleID := range params.EntityIDs {
				usages = append(usages, domain.PageArticleUsage{
					PageID:    *params.PageID,
					ArticleID: articleID,
					SiteID:    *params.SiteID,
					UserID:    userID,
				})
			}
			err = u.pageArticleUsageRepo.CreateBatch(usages)
			if err != nil {
				return nil, err
			}
		}

	case enums.ProductUsage:
		// Delete existing usages
		err = u.pageProductUsageRepo.DeleteByPageIDAndSiteID(*params.PageID, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Create new usages
		if len(params.EntityIDs) > 0 {
			// Verify products exist
			for _, productID := range params.EntityIDs {
				_, err = u.productRepo.GetByID(productID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
			}

			// Create usages
			usages := make([]domain.PageProductUsage, 0, len(params.EntityIDs))
			for _, productID := range params.EntityIDs {
				usages = append(usages, domain.PageProductUsage{
					PageID:    *params.PageID,
					ProductID: productID,
					SiteID:    *params.SiteID,
					UserID:    userID,
				})
			}
			err = u.pageProductUsageRepo.CreateBatch(usages)
			if err != nil {
				return nil, err
			}
		}

	case enums.HeaderFooterUsage:
		// Delete existing usages
		err = u.pageHeaderFooterUsageRepo.DeleteByPageIDAndSiteID(*params.PageID, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Create new usages
		if len(params.EntityIDs) > 0 {
			// Verify header/footers exist
			for _, headerFooterID := range params.EntityIDs {
				_, err = u.headerFooterRepo.GetByID(headerFooterID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
			}

			// Create usages
			usages := make([]domain.PageHeaderFooterUsage, 0, len(params.EntityIDs))
			for _, headerFooterID := range params.EntityIDs {
				usages = append(usages, domain.PageHeaderFooterUsage{
					PageID:         *params.PageID,
					HeaderFooterID: headerFooterID,
					SiteID:         *params.SiteID,
					UserID:         userID,
				})
			}
			err = u.pageHeaderFooterUsageRepo.CreateBatch(usages)
			if err != nil {
				return nil, err
			}
		}

	default:
		return nil, errors.New("نوع استفاده نامعتبر است")
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *PageUsageUsecase) FindPageUsagesQuery(params *page_usage.FindPageUsagesQuery) (any, error) {
	u.Logger.Info("FindPageUsagesQuery called", map[string]interface{}{
		"entityIds": params.EntityIDs,
		"siteId":    *params.SiteID,
		"type":      params.Type,
	})

	// Check if site exists
	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Find usages based on type
	switch params.Type {
	case enums.ArticleUsage:
		usages, err := u.pageArticleUsageRepo.GetByArticleIDsAndSiteID(params.EntityIDs, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Extract page IDs from usages
		pageIDMap := make(map[int64]bool)
		for _, usage := range usages {
			pageIDMap[usage.PageID] = true
		}

		pageIDs := make([]int64, 0, len(pageIDMap))
		for pageID := range pageIDMap {
			pageIDs = append(pageIDs, pageID)
		}

		// Get page details
		var pages []domain.Page
		if len(pageIDs) > 0 {
			pages, err = u.pageRepo.GetByIDs(pageIDs, *params.SiteID)
			if err != nil {
				return nil, err
			}
		}

		return enhancePageUsageResponse(pages), nil

	case enums.ProductUsage:
		usages, err := u.pageProductUsageRepo.GetByProductIDsAndSiteID(params.EntityIDs, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Extract page IDs from usages
		pageIDMap := make(map[int64]bool)
		for _, usage := range usages {
			pageIDMap[usage.PageID] = true
		}

		pageIDs := make([]int64, 0, len(pageIDMap))
		for pageID := range pageIDMap {
			pageIDs = append(pageIDs, pageID)
		}

		// Get page details
		var pages []domain.Page
		if len(pageIDs) > 0 {
			pages, err = u.pageRepo.GetByIDs(pageIDs, *params.SiteID)
			if err != nil {
				return nil, err
			}
		}

		return enhancePageUsageResponse(pages), nil

	case enums.HeaderFooterUsage:
		usages, err := u.pageHeaderFooterUsageRepo.GetByHeaderFooterIDsAndSiteID(params.EntityIDs, *params.SiteID)
		if err != nil {
			return nil, err
		}

		// Extract page IDs from usages
		pageIDMap := make(map[int64]bool)
		for _, usage := range usages {
			pageIDMap[usage.PageID] = true
		}

		pageIDs := make([]int64, 0, len(pageIDMap))
		for pageID := range pageIDMap {
			pageIDs = append(pageIDs, pageID)
		}

		// Get page details
		var pages []domain.Page
		if len(pageIDs) > 0 {
			pages, err = u.pageRepo.GetByIDs(pageIDs, *params.SiteID)
			if err != nil {
				return nil, err
			}
		}

		return enhancePageUsageResponse(pages), nil

	default:
		return nil, errors.New("نوع استفاده نامعتبر است")
	}
}

// Helper function to create a standardized response for page usages
func enhancePageUsageResponse(pages []domain.Page) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(pages))

	for _, page := range pages {
		pageData := map[string]interface{}{
			"id":     page.ID,
			"title":  page.Title,
			"slug":   page.Slug,
			"siteId": page.SiteID,
		}
		result = append(result, pageData)
	}

	return result
}
