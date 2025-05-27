package siteusecase

import (
	"errors"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"regexp"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/site"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

var (
	allowedBaseDomains = []string{"squidweb.ir"}
)

type SiteUsecase struct {
	*usecase.BaseUsecase
	repo        repository.ISiteRepository
	settingRepo repository.ISettingRepository
}

func NewSiteUsecase(c contract.IContainer) *SiteUsecase {
	return &SiteUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		repo:        c.GetSiteRepo(),
		settingRepo: c.GetSettingRepo(),
	}
}

func (u *SiteUsecase) CreateSiteCommand(params *site.CreateSiteCommand) (*resp.Response, error) {
	_, err := u.repo.GetByDomain(*params.Domain)
	if err == nil {
		return nil, resp.NewError(resp.BadRequest, "دامنه وارد شده قبلاً استفاده شده است")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resp.NewError(resp.Internal, "خطا در بررسی دامنه")
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	domainInput := *params.Domain
	domainType := *params.DomainType

	if domainType == enums.DomainType {
		matched, _ := regexp.MatchString(`^([a-zA-Z0-9-]+\.)+(ir|com|net|org)$`, domainInput)
		if !matched {
			return nil, resp.NewError(resp.BadRequest, "دامنه اصلی باید با فرمت معتبر باشد (مثال: example.com)")
		}
	} else if domainType == enums.SubdomainType {
		validSubdomain := false

		for _, baseDomain := range allowedBaseDomains {
			pattern := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.` + regexp.QuoteMeta(baseDomain) + `$`
			matched, _ := regexp.MatchString(pattern, domainInput)
			if matched {
				validSubdomain = true
				break
			}
		}

		if !validSubdomain {
			return nil, resp.NewError(resp.BadRequest, "ساب دامنه باید با یکی از دامنه های مجاز باشد (مثال: shop.squidweb.ir)")
		}
	}

	site := domain.Site{
		Domain:     *params.Domain,
		DomainType: *params.DomainType,
		Name:       *params.Name,
		Status:     *params.Status,
		SiteType:   *params.SiteType,
		UserID:     *userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}
	err = u.repo.Create(&site)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد سایت")
	}
	setting := domain.Setting{
		SiteID:    site.ID,
		UserID:    *userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}
	err = u.settingRepo.Create(&setting)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد تنظیمات")
	}
	return resp.NewResponseData(resp.Created, site, "سایت با موفقیت ایجاد شد"), nil
}

func (u *SiteUsecase) UpdateSiteCommand(params *site.UpdateSiteCommand) (*resp.Response, error) {
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	err = u.CheckAccessUserModel(existingSite)
	if err != nil {
		return nil, err
	}

	if params.Domain != nil && *params.Domain != existingSite.Domain {
		existingDomainSite, err := u.repo.GetByDomain(*params.Domain)
		if err == nil && existingDomainSite.ID != existingSite.ID {
			return nil, resp.NewError(resp.BadRequest, "دامنه وارد شده قبلاً استفاده شده است")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.Internal, "خطا در بررسی دامنه")
		}

		domainInput := *params.Domain
		domainType := *params.DomainType

		if domainType == enums.DomainType {
			matched, _ := regexp.MatchString(`^([a-zA-Z0-9-]+\.)+(ir|com|net|org)$`, domainInput)
			if !matched {
				return nil, resp.NewError(resp.BadRequest, "دامنه اصلی باید با فرمت معتبر باشد (مثال: example.com)")
			}
		} else if domainType == enums.SubdomainType {
			validSubdomain := false

			for _, baseDomain := range allowedBaseDomains {
				pattern := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.` + regexp.QuoteMeta(baseDomain) + `$`
				matched, _ := regexp.MatchString(pattern, domainInput)
				if matched {
					validSubdomain = true
					break
				}
			}

			if !validSubdomain {
				return nil, resp.NewError(resp.BadRequest, "ساب دامنه باید با یکی از دامنه های مجاز باشد (مثال: shop.squidweb.ir)")
			}
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

	return resp.NewResponseData(resp.Updated, existingSite, "سایت با موفقیت بروزرسانی شد"), nil
}

func (u *SiteUsecase) DeleteSiteCommand(params *site.DeleteSiteCommand) (*resp.Response, error) {
	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
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
	return resp.NewResponseData(resp.Retrieved, site, "اطلاعات سایت با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) GetByDomainSiteQuery(params *site.GetByDomainSiteQuery) (*resp.Response, error) {
	site, err := u.repo.GetByDomain(*params.Domain)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات سایت")
	}
	return resp.NewResponseData(resp.Retrieved, site, "اطلاعات سایت با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) GetAllSiteQuery(params *site.GetAllSiteQuery) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	result, err := u.repo.GetAllByUserID(*userID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی لیست سایت ها")
	}

	return resp.NewResponseData(resp.Retrieved, result, "لیست سایت ها با موفقیت بازیابی شد"), nil
}

func (u *SiteUsecase) AdminGetAllSiteQuery(params *site.AdminGetAllSiteQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}
	result, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی لیست سایت ها")
	}

	return resp.NewResponseData(resp.Retrieved, result, "لیست سایت ها با موفقیت بازیابی شد"), nil
}
