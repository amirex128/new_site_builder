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
	_, err := u.repo.GetByDomain(*params.Domain)
	if err == nil {
		return nil, resp.NewError(resp.BadRequest, "دامنه وارد شده قبلاً استفاده شده است")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resp.NewError(resp.Internal, "خطا در بررسی دامنه")
	}
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
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
	site := domain.Site{
		Domain:     *params.Domain,
		DomainType: domainType,
		Name:       *params.Name,
		Status:     status,
		SiteType:   siteType,
		UserID:     *userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}
	err = u.repo.Create(site)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد سایت")
	}
	setting := domain.Setting{
		SiteID:     site.ID,
		UserID:     *userID,
		CustomerID: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}
	err = u.settingRepo.Create(setting)
	if err != nil {
		// continue
	}
	createdSite, err := u.repo.GetByID(site.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	return resp.NewResponseData(resp.Created, enhanceSiteResponse(createdSite), "سایت با موفقیت ایجاد شد"), nil
}

func (u *SiteUsecase) UpdateSiteCommand(params *site.UpdateSiteCommand) (*resp.Response, error) {
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در بررسی دسترسی مدیریت")
	}
	if userID != nil && existingSite.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این سایت دسترسی ندارید")
	}
	if params.Domain != nil && *params.Domain != existingSite.Domain {
		existingDomainSite, err := u.repo.GetByDomain(*params.Domain)
		if err == nil && existingDomainSite.ID != existingSite.ID {
			return nil, resp.NewError(resp.BadRequest, "دامنه وارد شده قبلاً استفاده شده است")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.Internal, "خطا در بررسی دامنه")
		}
		existingSite.Domain = *params.Domain
	}
	if params.DomainType != nil {
		existingSite.DomainType = *params.DomainType
	}
	if params.Name != nil {
		existingSite.Name = *params.Name
	}
	if params.Status != nil {
		existingSite.Status = *params.Status
	}
	if params.SiteType != nil {
		existingSite.SiteType = *params.SiteType
	}
	existingSite.UpdatedAt = time.Now()
	err = u.repo.Update(existingSite)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی سایت")
	}
	updatedSite, err := u.repo.GetByID(existingSite.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	return resp.NewResponseData(resp.Updated, enhanceSiteResponse(updatedSite), "سایت با موفقیت بروزرسانی شد"), nil
}

func (u *SiteUsecase) DeleteSiteCommand(params *site.DeleteSiteCommand) (*resp.Response, error) {
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در بررسی دسترسی مدیریت")
	}
	if userID != nil && existingSite.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این سایت دسترسی ندارید")
	}
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در حذف سایت")
	}
	return resp.NewResponseData(resp.Deleted, resp.Data{"id": *params.ID}, "سایت با موفقیت حذف شد"), nil
}

func (u *SiteUsecase) GetByIdSiteQuery(params *site.GetByIdSiteQuery) (*resp.Response, error) {
	site, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	return resp.NewResponseData(resp.Retrieved, enhanceSiteResponse(site), "اطلاعات سایت با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) GetByDomainSiteQuery(params *site.GetByDomainSiteQuery) (*resp.Response, error) {
	site, err := u.repo.GetByDomain(*params.Domain)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	return resp.NewResponseData(resp.Retrieved, enhanceSiteResponse(site), "اطلاعات سایت با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) GetAllSiteQuery(params *site.GetAllSiteQuery) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	result, err := u.repo.GetAllByUserID(*userID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی لیست سایت ها")
	}
	sites := result.Items
	count := result.TotalCount
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
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}
	result, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی لیست سایت ها")
	}
	sites := result.Items
	count := result.TotalCount
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
