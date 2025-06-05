package headerfooterusecase

import (
	"encoding/json"
	"errors"
	header_footer2 "github.com/amirex128/new_site_builder/internal/application/dto/header_footer"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	"github.com/amirex128/new_site_builder/internal/contract/common"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	"github.com/amirex128/new_site_builder/internal/domain"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

type HeaderFooterUsecase struct {
	*usecase.BaseUsecase
	logger   sflogger.Logger
	repo     repository2.IHeaderFooterRepository
	siteRepo repository2.ISiteRepository
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

func (u *HeaderFooterUsecase) CreateHeaderFooterCommand(params *header_footer2.CreateHeaderFooterCommand) (*resp.Response, error) {
	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت سایت ")
	}

	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, err
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

	return resp.NewResponseData(resp.Created, headerFooter, "هدر/فوتر با موفقیت ایجاد شد"), nil
}

func (u *HeaderFooterUsecase) UpdateHeaderFooterCommand(params *header_footer2.UpdateHeaderFooterCommand) (*resp.Response, error) {
	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "هدر/فوتر مورد نظر یافت نشد")
	}

	err = u.CheckAccessUserModel(existingHeaderFooter)
	if err != nil {
		return nil, err
	}

	_, err = u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "خطا در دریافت سایت")
	}

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

	err = u.repo.Update(existingHeaderFooter)
	if err != nil {
		return nil, resp.NewError("خطا در بروزرسانی هدر/فوتر")
	}

	return resp.NewResponseData(resp.Updated, existingHeaderFooter, "هدر/فوتر با موفقیت بروزرسانی شد"), nil
}

func (u *HeaderFooterUsecase) DeleteHeaderFooterCommand(params *header_footer2.DeleteHeaderFooterCommand) (*resp.Response, error) {

	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "هدر/فوتر مورد نظر یافت نشد")
	}

	err = u.CheckAccessUserModel(existingHeaderFooter)
	if err != nil {
		return nil, err
	}

	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponse(resp.Deleted, "هدر/فوتر با موفقیت حذف شد"), nil
}

func (u *HeaderFooterUsecase) GetByIdHeaderFooterQuery(params *header_footer2.GetByIdHeaderFooterQuery) (*resp.Response, error) {
	headerFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "هدر/فوتر مورد نظر یافت نشد")
	}

	err = u.CheckAccessUserModel(headerFooter)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Retrieved, headerFooter, "هدر/فوتر با موفقیت دریافت شد"), nil
}

func (u *HeaderFooterUsecase) GetAllHeaderFooterQuery(params *header_footer2.GetAllHeaderFooterQuery) (*resp.Response, error) {
	results, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	err = u.CheckAccessSiteModel(params.SiteID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Retrieved, results, "هدر/فوترها با موفقیت دریافت شدند"), nil
}

func (u *HeaderFooterUsecase) AdminGetAllHeaderFooterQuery(params *header_footer2.AdminGetAllHeaderFooterQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	results, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Retrieved, results, "هدر/فوترها با موفقیت دریافت شدند"), nil
}

func (u *HeaderFooterUsecase) GetHeaderFooterByDomainOrSiteIdQuery(params *header_footer2.GetHeaderFooterByDomainOrSiteIdQuery) (*resp.Response, error) {
	var siteID int64
	if params.SiteID != nil {
		siteID = *params.SiteID
	} else if params.Domain != nil {
		site, err := u.siteRepo.GetByDomain(*params.Domain)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "سایتی با این دامنه یافت نشد")
		}
		siteID = site.ID
	} else {
		return nil, errors.New("شناسه سایت یا دامنه الزامی است")
	}

	result, err := u.repo.GetAllBySiteID(siteID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Retrieved, result, "هدر/فوترها با موفقیت دریافت شدند"), nil
}
