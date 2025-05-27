package pageusecase

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/page"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type PageUsecase struct {
	*usecase.BaseUsecase
	repo             repository.IPageRepository
	siteRepo         repository.ISiteRepository
	headerFooterRepo repository.IHeaderFooterRepository
	mediaRepo        repository.IMediaRepository
	authContext      func(c *gin.Context) service.IAuthService
}

func NewPageUsecase(c contract.IContainer) *PageUsecase {
	return &PageUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		repo:             c.GetPageRepo(),
		siteRepo:         c.GetSiteRepo(),
		headerFooterRepo: c.GetHeaderFooterRepo(),
		mediaRepo:        c.GetMediaRepo(),
	}
}

func (u *PageUsecase) CreatePageCommand(params *page.CreatePageCommand) (*resp.Response, error) {
	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	_, err = u.headerFooterRepo.GetByID(*params.HeaderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "هدر مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	_, err = u.headerFooterRepo.GetByID(*params.FooterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "فوتر مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	_, err = u.repo.GetBySlug(*params.Slug, *params.SiteID)
	if err == nil {
		return nil, resp.NewError(resp.BadRequest, "نامک (slug) تکراری است")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, err
	}
	bodyJSON, err := json.Marshal(params.Body)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	var seoTags string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTags = strings.Join(params.SeoTags, ",")
	}
	pageEntity := domain.Page{
		SiteID:      *params.SiteID,
		HeaderID:    *params.HeaderID,
		FooterID:    *params.FooterID,
		Slug:        *params.Slug,
		Title:       *params.Title,
		Description: getStringValueOrEmpty(params.Description),
		Body:        string(bodyJSON),
		SeoTags:     seoTags,
		UserID:      *userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}
	err = u.repo.Create(&pageEntity)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		err = u.repo.AddMediaToPage(pageEntity.ID, params.MediaIDs)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
	}
	createdPage, err := u.repo.GetByID(pageEntity.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Created, createdPage, "صفحه با موفقیت ایجاد شد"), nil
}

func (u *PageUsecase) UpdatePageCommand(params *page.UpdatePageCommand) (*resp.Response, error) {
	existingPage, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "صفحه مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, err
	}
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}
	if existingPage.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این صفحه دسترسی ندارید")
	}
	if params.SiteID != nil {
		_, err = u.siteRepo.GetByID(*params.SiteID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		existingPage.SiteID = *params.SiteID
	}
	if params.HeaderID != nil {
		_, err = u.headerFooterRepo.GetByID(*params.HeaderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "هدر مورد نظر یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		existingPage.HeaderID = *params.HeaderID
	}
	if params.FooterID != nil {
		_, err = u.headerFooterRepo.GetByID(*params.FooterID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "فوتر مورد نظر یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		existingPage.FooterID = *params.FooterID
	}
	if params.Slug != nil && *params.Slug != existingPage.Slug {
		siteID := existingPage.SiteID
		if params.SiteID != nil {
			siteID = *params.SiteID
		}
		conflictPage, err := u.repo.GetBySlug(*params.Slug, siteID)
		if err == nil && conflictPage.ID != existingPage.ID {
			return nil, resp.NewError(resp.BadRequest, "نامک (slug) تکراری است")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		existingPage.Slug = *params.Slug
	}
	if params.Title != nil {
		existingPage.Title = *params.Title
	}
	if params.Description != nil {
		existingPage.Description = *params.Description
	}
	if params.Body != nil {
		bodyJSON, err := json.Marshal(params.Body)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		existingPage.Body = string(bodyJSON)
	}
	if params.SeoTags != nil {
		if len(params.SeoTags) > 0 {
			existingPage.SeoTags = strings.Join(params.SeoTags, ",")
		} else {
			existingPage.SeoTags = ""
		}
	}
	existingPage.UpdatedAt = time.Now()
	err = u.repo.Update(existingPage)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.MediaIDs != nil {
		err = u.repo.RemoveAllMediaFromPage(existingPage.ID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		if len(params.MediaIDs) > 0 {
			err = u.repo.AddMediaToPage(existingPage.ID, params.MediaIDs)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	}
	updatedPage, err := u.repo.GetByID(existingPage.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Updated, updatedPage, "صفحه با موفقیت بروزرسانی شد"), nil
}

func (u *PageUsecase) DeletePageCommand(params *page.DeletePageCommand) (*resp.Response, error) {
	existingPage, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "صفحه مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	err = u.CheckAccessUserModel(existingPage)
	if err != nil {
		return nil, err
	}
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Deleted, map[string]interface{}{"id": *params.ID}, "صفحه با موفقیت حذف شد"), nil
}

func (u *PageUsecase) GetByIdPageQuery(params *page.GetByIdPageQuery) (*resp.Response, error) {
	if params.ID != nil {
		pageItem, err := u.repo.GetByIDAndSiteID(*params.ID, *params.SiteID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "صفحه مورد نظر یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		return resp.NewResponseData(resp.Retrieved, pageItem, "صفحه با موفقیت دریافت شد"), nil
	} else if params.IDs != nil && len(params.IDs) > 0 {
		pages, err := u.repo.GetByIDs(params.IDs, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		return resp.NewResponseData(resp.Retrieved, pages, "صفحات با موفقیت دریافت شدند"), nil
	}
	return nil, resp.NewError(resp.BadRequest, "شناسه صفحه یا صفحات الزامی است")
}

func (u *PageUsecase) GetAllPageQuery(params *page.GetAllPageQuery) (*resp.Response, error) {
	err := u.CheckAccessSiteID(params.SiteID)
	if err != nil {
		return nil, err
	}
	pagesResult, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Retrieved, pagesResult, "لیست صفحات با موفقیت دریافت شد"), nil
}

func (u *PageUsecase) AdminGetAllPageQuery(params *page.AdminGetAllPageQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}
	pagesResult, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, pagesResult, "لیست صفحات ادمین با موفقیت دریافت شد"), nil
}

func (u *PageUsecase) GetByPathPageQuery(params *page.GetByPathPageQuery) (*resp.Response, error) {
	pages, err := u.repo.GetByPaths(params.Paths, *params.SiteID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, pages, "صفحات با موفقیت دریافت شدند"), nil
}

func getStringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
