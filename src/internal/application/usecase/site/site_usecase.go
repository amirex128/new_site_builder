package siteusecase

import (
	"errors"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/gin-gonic/gin"
	"time"

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

func (u *SiteUsecase) CreateSiteCommand(params *site.CreateSiteCommand) (any, error) {
	u.Logger.Info("CreateSiteCommand called", map[string]interface{}{
		"domain": *params.Domain,
		"name":   *params.Name,
	})

	// Check if domain already exists
	_, err := u.repo.GetByDomain(*params.Domain)
	if err == nil {
		return nil, errors.New("دامنه وارد شده قبلاً استفاده شده است")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get user ID from auth context
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Convert enum values to strings
	domainType := getDomainTypeString(*params.DomainType)
	siteType := getSiteTypeString(*params.SiteType)
	status := getStatusString(*params.Status)

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
		return nil, err
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
		return nil, err
	}

	return enhanceSiteResponse(createdSite), nil
}

func (u *SiteUsecase) UpdateSiteCommand(params *site.UpdateSiteCommand) (any, error) {
	u.Logger.Info("UpdateSiteCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing site
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingSite.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این سایت دسترسی ندارید")
	}

	// Update domain if provided
	if params.Domain != nil && *params.Domain != existingSite.Domain {
		// Check if new domain is available
		existingDomainSite, err := u.repo.GetByDomain(*params.Domain)
		if err == nil && existingDomainSite.ID != existingSite.ID {
			return nil, errors.New("دامنه وارد شده قبلاً استفاده شده است")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		existingSite.Domain = *params.Domain
	}

	// Update domain type if provided
	if params.DomainType != nil {
		existingSite.DomainType = getDomainTypeString(*params.DomainType)
	}

	// Update name if provided
	if params.Name != nil {
		existingSite.Name = *params.Name
	}

	// Update status if provided
	if params.Status != nil {
		existingSite.Status = getStatusString(*params.Status)
	}

	// Update site type if provided
	if params.SiteType != nil {
		existingSite.SiteType = getSiteTypeString(*params.SiteType)
	}

	existingSite.UpdatedAt = time.Now()

	// Update in repository
	err = u.repo.Update(existingSite)
	if err != nil {
		return nil, err
	}

	// Get updated site with relations
	updatedSite, err := u.repo.GetByID(existingSite.ID)
	if err != nil {
		return nil, err
	}

	return enhanceSiteResponse(updatedSite), nil
}

func (u *SiteUsecase) DeleteSiteCommand(params *site.DeleteSiteCommand) (any, error) {
	u.Logger.Info("DeleteSiteCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing site
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingSite.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این سایت دسترسی ندارید")
	}

	// Delete site in repository
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": *params.ID,
	}, nil
}

func (u *SiteUsecase) GetByIdSiteQuery(params *site.GetByIdSiteQuery) (any, error) {
	u.Logger.Info("GetByIdSiteQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get site by ID
	site, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	return enhanceSiteResponse(site), nil
}

func (u *SiteUsecase) GetByDomainSiteQuery(params *site.GetByDomainSiteQuery) (any, error) {
	u.Logger.Info("GetByDomainSiteQuery called", map[string]interface{}{
		"domain": *params.Domain,
	})

	// Get site by domain
	site, err := u.repo.GetByDomain(*params.Domain)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	return enhanceSiteResponse(site), nil
}

func (u *SiteUsecase) GetAllSiteQuery(params *site.GetAllSiteQuery) (any, error) {
	u.Logger.Info("GetAllSiteQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Get all sites for the user with pagination
	sites, count, err := u.repo.GetAllByUserID(userID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Enhance site responses
	enhancedSites := make([]map[string]interface{}, 0, len(sites))
	for _, s := range sites {
		enhancedSites = append(enhancedSites, enhanceSiteResponse(s))
	}

	return map[string]interface{}{
		"items":     enhancedSites,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *SiteUsecase) AdminGetAllSiteQuery(params *site.AdminGetAllSiteQuery) (any, error) {
	u.Logger.Info("AdminGetAllSiteQuery called", map[string]interface{}{
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

	// Get all sites with pagination
	sites, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Enhance site responses
	enhancedSites := make([]map[string]interface{}, 0, len(sites))
	for _, s := range sites {
		enhancedSites = append(enhancedSites, enhanceSiteResponse(s))
	}

	return map[string]interface{}{
		"items":     enhancedSites,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// Helper function to enhance site response with structured data
func enhanceSiteResponse(site domain.Site) map[string]interface{} {
	response := map[string]interface{}{
		"id":         site.ID,
		"domain":     site.Domain,
		"domainType": getDomainTypeEnum(site.DomainType),
		"name":       site.Name,
		"status":     getStatusEnum(site.Status),
		"siteType":   getSiteTypeEnum(site.SiteType),
		"userId":     site.UserID,
		"createdAt":  site.CreatedAt,
		"updatedAt":  site.UpdatedAt,
	}

	// Add settings if available
	if site.Setting != nil {
		response["setting"] = map[string]interface{}{
			"id":         site.Setting.ID,
			"siteId":     site.Setting.SiteID,
			"userId":     site.Setting.UserID,
			"customerId": site.Setting.CustomerID,
		}
	}

	return response
}

// Helper functions to convert enum values between Go and .NET versions

func getDomainTypeString(domainType site.DomainTypeEnum) string {
	switch domainType {
	case site.Domain:
		return "Domain"
	case site.Subdomain:
		return "Subdomain"
	default:
		return "Domain"
	}
}

func getDomainTypeEnum(domainType string) site.DomainTypeEnum {
	switch domainType {
	case "Domain":
		return site.Domain
	case "Subdomain":
		return site.Subdomain
	default:
		return site.Domain
	}
}

func getSiteTypeString(siteType site.SiteTypeEnum) string {
	switch siteType {
	case site.Shop:
		return "Shop"
	case site.Blog:
		return "Blog"
	case site.Business:
		return "Business"
	default:
		return "Shop"
	}
}

func getSiteTypeEnum(siteType string) site.SiteTypeEnum {
	switch siteType {
	case "Shop":
		return site.Shop
	case "Blog":
		return site.Blog
	case "Business":
		return site.Business
	default:
		return site.Shop
	}
}

func getStatusString(status site.StatusEnum) string {
	switch status {
	case site.Active:
		return "Active"
	case site.Inactive:
		return "Inactive"
	case site.Pending:
		return "Pending"
	case site.Deleted:
		return "Deleted"
	default:
		return "Active"
	}
}

func getStatusEnum(status string) site.StatusEnum {
	switch status {
	case "Active":
		return site.Active
	case "Inactive":
		return site.Inactive
	case "Pending":
		return site.Pending
	case "Deleted":
		return site.Deleted
	default:
		return site.Active
	}
}
