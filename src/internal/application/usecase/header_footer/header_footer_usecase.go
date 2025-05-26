package headerfooterusecase

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/header_footer"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type HeaderFooterUsecase struct {
	*usecase.BaseUsecase
	logger   sflogger.Logger
	repo     repository.IHeaderFooterRepository
	siteRepo repository.ISiteRepository
}

func NewHeaderFooterUsecase(c contract.IContainer) *HeaderFooterUsecase {
	return &HeaderFooterUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		repo:     c.GetHeaderFooterRepo(),
		siteRepo: c.GetSiteRepo(),
	}
}

func (u *HeaderFooterUsecase) CreateHeaderFooterCommand(params *header_footer.CreateHeaderFooterCommand) (*resp.Response, error) {
	u.Logger.Info("CreateHeaderFooterCommand called", map[string]interface{}{
		"siteId": *params.SiteID,
		"title":  *params.Title,
		"type":   *params.Type,
	})

	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}

	bodyJSON, err := json.Marshal(params.Body)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	headerFooter := domain.HeaderFooter{
		SiteID:    *params.SiteID,
		Title:     *params.Title,
		IsMain:    *params.IsMain,
		Body:      string(bodyJSON),
		Type:      string(*params.Type),
		UserID:    *userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	err = u.repo.Create(&headerFooter)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	createdHeaderFooter, err := u.repo.GetByID(headerFooter.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Created, resp.Data{
		"headerFooter": createdHeaderFooter,
	}, "هدر/فوتر با موفقیت ایجاد شد"), nil
}

func (u *HeaderFooterUsecase) UpdateHeaderFooterCommand(params *header_footer.UpdateHeaderFooterCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateHeaderFooterCommand called", map[string]interface{}{
		"id":     *params.ID,
		"siteId": *params.SiteID,
		"type":   *params.Type,
	})

	// Get existing header/footer
	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("هدر/فوتر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingHeaderFooter.UserID != *userID && !isAdmin {
		return nil, errors.New("شما به این هدر/فوتر دسترسی ندارید")
	}

	// Check if site exists
	_, err = u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Update fields
	existingHeaderFooter.SiteID = *params.SiteID
	existingHeaderFooter.IsMain = *params.IsMain
	existingHeaderFooter.Type = string(*params.Type)

	if params.Title != nil {
		existingHeaderFooter.Title = *params.Title
	}

	if params.Body != nil {
		bodyJSON, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		existingHeaderFooter.Body = string(bodyJSON)
	}

	existingHeaderFooter.UpdatedAt = time.Now()

	// Update in repository
	err = u.repo.Update(existingHeaderFooter)
	if err != nil {
		return nil, err
	}

	// Get updated header/footer with relations
	updatedHeaderFooter, err := u.repo.GetByID(existingHeaderFooter.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Updated, resp.Data{
		"headerFooter": updatedHeaderFooter,
	}, "هدر/فوتر با موفقیت بروزرسانی شد"), nil
}

func (u *HeaderFooterUsecase) DeleteHeaderFooterCommand(params *header_footer.DeleteHeaderFooterCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteHeaderFooterCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing header/footer
	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("هدر/فوتر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingHeaderFooter.UserID != *userID && !isAdmin {
		return nil, errors.New("شما به این هدر/فوتر دسترسی ندارید")
	}

	// Delete header/footer
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Deleted, resp.Data{
		"id": existingHeaderFooter.ID,
	}, "هدر/فوتر با موفقیت حذف شد"), nil
}

func (u *HeaderFooterUsecase) GetByIdHeaderFooterQuery(params *header_footer.GetByIdHeaderFooterQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdHeaderFooterQuery called", map[string]interface{}{
		"id":     params.ID,
		"siteId": *params.SiteID,
	})

	if params.ID != nil {
		// Get by ID
		headerFooter, err := u.repo.GetByID(*params.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("هدر/فوتر مورد نظر یافت نشد")
			}
			return nil, err
		}

		// Check if header/footer belongs to the specified site
		if headerFooter.SiteID != *params.SiteID {
			return nil, errors.New("هدر/فوتر مورد نظر متعلق به سایت مشخص شده نیست")
		}

		// Check if type matches if specified
		if params.Type != nil && headerFooter.Type != string(*params.Type) {
			return nil, errors.New("نوع هدر/فوتر مطابقت ندارد")
		}

		return resp.NewResponseData(resp.Retrieved, resp.Data{
			"headerFooter": headerFooter,
		}, "هدر/فوتر با موفقیت دریافت شد"), nil
	} else if len(params.IDs) > 0 {
		// Get by multiple IDs
		// In a monolithic approach, we would need to implement a method to get by IDs in the repository
		// For now, we'll use a simple approach of getting each one individually
		var result []domain.HeaderFooter
		for _, id := range params.IDs {
			headerFooter, err := u.repo.GetByID(id)
			if err == nil && headerFooter.SiteID == *params.SiteID {
				result = append(result, *headerFooter)
			}
		}
		return resp.NewResponseData(resp.Retrieved, resp.Data{
			"headerFooters": result,
		}, "هدر/فوترها با موفقیت دریافت شدند"), nil
	}

	return nil, errors.New("شناسه هدر/فوتر الزامی است")
}

func (u *HeaderFooterUsecase) GetAllHeaderFooterQuery(params *header_footer.GetAllHeaderFooterQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllHeaderFooterQuery called", map[string]interface{}{
		"siteId": *params.SiteID,
		"type":   params.Type,
	})

	results, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	items := results.Items
	if params.Type != nil {
		filtered := make([]domain.HeaderFooter, 0, len(items))
		for _, hf := range items {
			if hf.Type == string(*params.Type) {
				filtered = append(filtered, hf)
			}
		}
		items = filtered
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     items,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "هدر/فوترها با موفقیت دریافت شدند"), nil
}

func (u *HeaderFooterUsecase) AdminGetAllHeaderFooterQuery(params *header_footer.AdminGetAllHeaderFooterQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllHeaderFooterQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	results, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"pagination": results,
	}, ""), nil
}

func (u *HeaderFooterUsecase) GetHeaderFooterByDomainOrSiteIdQuery(params *header_footer.GetHeaderFooterByDomainOrSiteIdQuery) (*resp.Response, error) {
	u.Logger.Info("GetHeaderFooterByDomainOrSiteIdQuery called", map[string]interface{}{
		"siteId": params.SiteID,
		"domain": params.Domain,
	})

	var siteID int64

	if params.SiteID != nil {
		// Use provided site ID
		siteID = *params.SiteID
	} else if params.Domain != nil {
		// Get site by domain
		site, err := u.siteRepo.GetByDomain(*params.Domain)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("سایت مورد نظر یافت نشد")
			}
			return nil, err
		}
		siteID = site.ID
	} else {
		return nil, errors.New("شناسه سایت یا دامنه الزامی است")
	}

	// Get main header/footer for site
	// In a real implementation, we would have a dedicated repository method for this
	// For now, we'll get all and filter for the main one
	result, err := u.repo.GetAllBySiteID(siteID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return nil, err
	}

	headerFooters := result.Items

	// Find main header/footer
	for _, hf := range headerFooters {
		if hf.IsMain {
			return resp.NewResponseData(resp.Retrieved, resp.Data{
				"headerFooter": hf,
			}, "هدر/فوتر اصلی با موفقیت دریافت شد"), nil
		}
	}

	return nil, errors.New("هدر/فوتر اصلی برای سایت مورد نظر یافت نشد")
}
