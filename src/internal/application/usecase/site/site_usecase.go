package siteusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/site"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type SiteUsecase struct {
	*usecase.BaseUsecase
	repo        repository.ISiteRepository
	settingRepo repository.ISettingRepository
	authContext func(c *gin.Context) service.IAuthService
}

func NewSiteUsecase(c contract.IContainer) *SiteUsecase {
	return &SiteUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		repo:        c.GetSiteRepo(),
		settingRepo: c.GetSettingRepo(),
		authContext: c.GetAuthTransientService(),
	}
}

func (u *SiteUsecase) CreateSiteCommand(params *site.CreateSiteCommand) (*resp.Response, error) {
	u.Logger.Info("CreateSiteCommand called", map[string]interface{}{
		"domain": *params.Domain,
		"name":   *params.Name,
	})

	// Check if domain already exists
	_, err := u.repo.GetByDomain(*params.Domain)
	if err == nil {
		return nil, resp.NewError(resp.BadRequest, "دامنه وارد شده قبلاً استفاده شده است")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		u.Logger.Error("Error checking domain existence", map[string]interface{}{
			"domain": *params.Domain,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بررسی دامنه")
	}

	// Get user ID from auth context
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		u.Logger.Error("Error getting user ID from auth context", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}

	// Convert enum values to strings
	var domainType enums.DomainTypeEnum
	var siteType enums.SiteTypeEnum
	var status enums.StatusEnum
	if params.DomainType != nil {
		domainType = *params.DomainType
	}
	if params.SiteType != nil {
		siteType = *params.SiteType
	}
	if params.Status != nil {
		status = *params.Status
	}

	// Create the Site entity
	site := domain.Site{
		Domain:     *params.Domain,
		DomainType: domainType,
		Name:       *params.Name,
		Status:     status,
		SiteType:   siteType,
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	// Create site in repository
	err = u.repo.Create(site)
	if err != nil {
		u.Logger.Error("Error creating site", map[string]interface{}{
			"domain": *params.Domain,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد سایت")
	}

	// Create default setting for site
	setting := domain.Setting{
		SiteID:     site.ID,
		UserID:     userID,
		CustomerID: 0, // Default customer ID or get from context if needed
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	err = u.settingRepo.Create(setting)
	if err != nil {
		u.Logger.Error("Failed to create settings for site", map[string]interface{}{
			"siteId": site.ID,
			"error":  err.Error(),
		})
		// Continue, as this is not a critical error
		// In a real scenario, you might want to rollback the site creation
	}

	// Reload site with setting relation
	createdSite, err := u.repo.GetByID(site.ID)
	if err != nil {
		u.Logger.Error("Error retrieving created site", map[string]interface{}{
			"siteId": site.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}

	return resp.NewResponseData(resp.Created, enhanceSiteResponse(createdSite), "سایت با موفقیت ایجاد شد"), nil
}

func (u *SiteUsecase) UpdateSiteCommand(params *site.UpdateSiteCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateSiteCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing site
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		u.Logger.Error("Error retrieving site", map[string]interface{}{
			"siteId": *params.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		u.Logger.Error("Error getting user ID from auth context", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		u.Logger.Error("Error checking admin status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در بررسی دسترسی مدیریت")
	}

	if existingSite.UserID != userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این سایت دسترسی ندارید")
	}

	// Update domain if provided
	if params.Domain != nil && *params.Domain != existingSite.Domain {
		// Check if new domain is available
		existingDomainSite, err := u.repo.GetByDomain(*params.Domain)
		if err == nil && existingDomainSite.ID != existingSite.ID {
			return nil, resp.NewError(resp.BadRequest, "دامنه وارد شده قبلاً استفاده شده است")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			u.Logger.Error("Error checking domain availability", map[string]interface{}{
				"domain": *params.Domain,
				"error":  err.Error(),
			})
			return nil, resp.NewError(resp.Internal, "خطا در بررسی دامنه")
		}

		existingSite.Domain = *params.Domain
	}

	// Update domain type if provided
	if params.DomainType != nil {
		existingSite.DomainType = *params.DomainType
	}

	// Update name if provided
	if params.Name != nil {
		existingSite.Name = *params.Name
	}

	// Update status if provided
	if params.Status != nil {
		existingSite.Status = *params.Status
	}

	// Update site type if provided
	if params.SiteType != nil {
		existingSite.SiteType = *params.SiteType
	}

	existingSite.UpdatedAt = time.Now()

	// Update in repository
	err = u.repo.Update(existingSite)
	if err != nil {
		u.Logger.Error("Error updating site", map[string]interface{}{
			"siteId": existingSite.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی سایت")
	}

	// Get updated site with relations
	updatedSite, err := u.repo.GetByID(existingSite.ID)
	if err != nil {
		u.Logger.Error("Error retrieving updated site", map[string]interface{}{
			"siteId": existingSite.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}

	return resp.NewResponseData(resp.Updated, enhanceSiteResponse(updatedSite), "سایت با موفقیت بروزرسانی شد"), nil
}

func (u *SiteUsecase) DeleteSiteCommand(params *site.DeleteSiteCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteSiteCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing site
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		u.Logger.Error("Error retrieving site for deletion", map[string]interface{}{
			"siteId": *params.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		u.Logger.Error("Error getting user ID from auth context", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		u.Logger.Error("Error checking admin status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در بررسی دسترسی مدیریت")
	}

	if existingSite.UserID != userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این سایت دسترسی ندارید")
	}

	// Delete site in repository
	err = u.repo.Delete(*params.ID)
	if err != nil {
		u.Logger.Error("Error deleting site", map[string]interface{}{
			"siteId": *params.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در حذف سایت")
	}

	return resp.NewResponseData(resp.Deleted, resp.Data{
		"id": *params.ID,
	}, "سایت با موفقیت حذف شد"), nil
}

func (u *SiteUsecase) GetByIdSiteQuery(params *site.GetByIdSiteQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdSiteQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get site by ID
	site, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		u.Logger.Error("Error retrieving site by ID", map[string]interface{}{
			"siteId": *params.ID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}

	return resp.NewResponseData(resp.Retrieved, enhanceSiteResponse(site), "اطلاعات سایت با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) GetByDomainSiteQuery(params *site.GetByDomainSiteQuery) (*resp.Response, error) {
	u.Logger.Info("GetByDomainSiteQuery called", map[string]interface{}{
		"domain": *params.Domain,
	})

	// Get site by domain
	site, err := u.repo.GetByDomain(*params.Domain)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		u.Logger.Error("Error retrieving site by domain", map[string]interface{}{
			"domain": *params.Domain,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}

	return resp.NewResponseData(resp.Retrieved, enhanceSiteResponse(site), "اطلاعات سایت با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) GetAllSiteQuery(params *site.GetAllSiteQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllSiteQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		u.Logger.Error("Error getting user ID from auth context", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}

	// Get all sites for the user with pagination
	sites, count, err := u.repo.GetAllByUserID(userID, params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error retrieving sites for user", map[string]interface{}{
			"userId": userID,
			"error":  err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی لیست سایت ها")
	}

	// Enhance site responses
	enhancedSites := make([]map[string]interface{}, 0, len(sites))
	for _, s := range sites {
		enhancedSites = append(enhancedSites, enhanceSiteResponse(s))
	}

	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     enhancedSites,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "لیست سایت ها با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) AdminGetAllSiteQuery(params *site.AdminGetAllSiteQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllSiteQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		u.Logger.Error("Error checking admin status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Unauthorized, "خطا در بررسی دسترسی مدیریت")
	}

	if !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all sites with pagination
	sites, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error retrieving all sites", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی لیست سایت ها")
	}

	// Enhance site responses
	enhancedSites := make([]map[string]interface{}, 0, len(sites))
	for _, s := range sites {
		enhancedSites = append(enhancedSites, enhanceSiteResponse(s))
	}

	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     enhancedSites,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "لیست سایت ها با موفقیت بازیابی شد"), nil
}

// Helper function to enhance site response with structured data
func enhanceSiteResponse(site domain.Site) resp.Data {
	response := resp.Data{
		"id":         site.ID,
		"domain":     site.Domain,
		"domainType": site.DomainType,
		"name":       site.Name,
		"status":     site.Status,
		"siteType":   site.SiteType,
		"userId":     site.UserID,
		"createdAt":  site.CreatedAt,
		"updatedAt":  site.UpdatedAt,
	}

	// Add settings if available
	if site.Setting != nil {
		response["setting"] = resp.Data{
			"id":         site.Setting.ID,
			"siteId":     site.Setting.SiteID,
			"userId":     site.Setting.UserID,
			"customerId": site.Setting.CustomerID,
		}
	}

	return response
}
