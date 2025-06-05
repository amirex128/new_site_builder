package pageusageusecase

import (
	"errors"
	page_usage2 "github.com/amirex128/new_site_builder/internal/application/dto/page_usage"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	"github.com/amirex128/new_site_builder/internal/domain/enums"

	"gorm.io/gorm"
)

type PageUsageUsecase struct {
	*usecase.BaseUsecase
	pageRepo                  repository2.IPageRepository
	pageArticleUsageRepo      repository2.IPageArticleUsageRepository
	pageProductUsageRepo      repository2.IPageProductUsageRepository
	pageHeaderFooterUsageRepo repository2.IPageHeaderFooterUsageRepository
	articleRepo               repository2.IArticleRepository
	productRepo               repository2.IProductRepository
	headerFooterRepo          repository2.IHeaderFooterRepository
	siteRepo                  repository2.ISiteRepository
}

func NewPageUsageUsecase(c contract.IContainer) *PageUsageUsecase {
	return &PageUsageUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		pageRepo:                  c.GetPageRepo(),
		pageArticleUsageRepo:      c.GetPageArticleUsageRepo(),
		pageProductUsageRepo:      c.GetPageProductUsageRepo(),
		pageHeaderFooterUsageRepo: c.GetPageHeaderFooterUsageRepo(),
		articleRepo:               c.GetArticleRepo(),
		productRepo:               c.GetProductRepo(),
		headerFooterRepo:          c.GetHeaderFooterRepo(),
		siteRepo:                  c.GetSiteRepo(),
	}
}

func (u *PageUsageUsecase) SyncPageUsageCommand(params *page_usage2.SyncPageUsageCommand) (*resp.Response, error) {
	page, err := u.pageRepo.GetByIDAndSiteID(*params.PageID, *params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "صفحه مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	_, err = u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
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
	if page.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این صفحه دسترسی ندارید")
	}
	switch params.Type {
	case enums.PageArticleUsage:
		err = u.pageArticleUsageRepo.DeleteByPageIDAndSiteID(*params.PageID, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		if len(params.EntityIDs) > 0 {
			for _, articleID := range params.EntityIDs {
				_, err = u.articleRepo.GetByID(articleID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, resp.NewError(resp.Internal, err.Error())
				}
			}
			usages := make([]domain2.PageArticleUsage, 0, len(params.EntityIDs))
			for _, articleID := range params.EntityIDs {
				usages = append(usages, domain2.PageArticleUsage{
					PageID:    *params.PageID,
					ArticleID: articleID,
					SiteID:    *params.SiteID,
					UserID:    *userID,
				})
			}
			err = u.pageArticleUsageRepo.CreateBatch(usages)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	case enums.PageProductUsage:
		err = u.pageProductUsageRepo.DeleteByPageIDAndSiteID(*params.PageID, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		if len(params.EntityIDs) > 0 {
			for _, productID := range params.EntityIDs {
				_, err = u.productRepo.GetByID(productID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, resp.NewError(resp.Internal, err.Error())
				}
			}
			usages := make([]domain2.PageProductUsage, 0, len(params.EntityIDs))
			for _, productID := range params.EntityIDs {
				usages = append(usages, domain2.PageProductUsage{
					PageID:    *params.PageID,
					ProductID: productID,
					SiteID:    *params.SiteID,
					UserID:    *userID,
				})
			}
			err = u.pageProductUsageRepo.CreateBatch(usages)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	case enums.PageHeaderFooterUsage:
		err = u.pageHeaderFooterUsageRepo.DeleteByPageIDAndSiteID(*params.PageID, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		if len(params.EntityIDs) > 0 {
			for _, headerFooterID := range params.EntityIDs {
				_, err = u.headerFooterRepo.GetByID(headerFooterID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, resp.NewError(resp.Internal, err.Error())
				}
			}
			usages := make([]domain2.PageHeaderFooterUsage, 0, len(params.EntityIDs))
			for _, headerFooterID := range params.EntityIDs {
				usages = append(usages, domain2.PageHeaderFooterUsage{
					PageID:         *params.PageID,
					HeaderFooterID: headerFooterID,
					SiteID:         *params.SiteID,
					UserID:         *userID,
				})
			}
			err = u.pageHeaderFooterUsageRepo.CreateBatch(usages)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	default:
		return nil, resp.NewError(resp.BadRequest, "نوع استفاده نامعتبر است")
	}
	return resp.NewResponseData(resp.Success, resp.Data{"success": true}, "Page usage synchronized successfully"), nil
}

func (u *PageUsageUsecase) FindPageUsagesQuery(params *page_usage2.FindPageUsagesQuery) (*resp.Response, error) {
	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سایت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	switch params.Type {
	case enums.PageArticleUsage:
		usages, err := u.pageArticleUsageRepo.GetByArticleIDsAndSiteID(params.EntityIDs, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		pageIDMap := make(map[int64]bool)
		for _, usage := range usages {
			pageIDMap[usage.PageID] = true
		}
		pageIDs := make([]int64, 0, len(pageIDMap))
		for pageID := range pageIDMap {
			pageIDs = append(pageIDs, pageID)
		}
		var pages []domain2.Page
		if len(pageIDs) > 0 {
			pages, err = u.pageRepo.GetByIDs(pageIDs, *params.SiteID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
		return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"items": enhancePageUsageResponse(pages)}, "صفحات با موفقیت دریافت شدند"), nil
	case enums.PageProductUsage:
		usages, err := u.pageProductUsageRepo.GetByProductIDsAndSiteID(params.EntityIDs, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		pageIDMap := make(map[int64]bool)
		for _, usage := range usages {
			pageIDMap[usage.PageID] = true
		}
		pageIDs := make([]int64, 0, len(pageIDMap))
		for pageID := range pageIDMap {
			pageIDs = append(pageIDs, pageID)
		}
		var pages []domain2.Page
		if len(pageIDs) > 0 {
			pages, err = u.pageRepo.GetByIDs(pageIDs, *params.SiteID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
		return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"items": enhancePageUsageResponse(pages)}, "صفحات با موفقیت دریافت شدند"), nil
	case enums.PageHeaderFooterUsage:
		usages, err := u.pageHeaderFooterUsageRepo.GetByHeaderFooterIDsAndSiteID(params.EntityIDs, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		pageIDMap := make(map[int64]bool)
		for _, usage := range usages {
			pageIDMap[usage.PageID] = true
		}
		pageIDs := make([]int64, 0, len(pageIDMap))
		for pageID := range pageIDMap {
			pageIDs = append(pageIDs, pageID)
		}
		var pages []domain2.Page
		if len(pageIDs) > 0 {
			pages, err = u.pageRepo.GetByIDs(pageIDs, *params.SiteID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
		return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"items": enhancePageUsageResponse(pages)}, "صفحات با موفقیت دریافت شدند"), nil
	default:
		return nil, resp.NewError(resp.BadRequest, "نوع استفاده نامعتبر است")
	}
}

func enhancePageUsageResponse(pages []domain2.Page) []map[string]interface{} {
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
