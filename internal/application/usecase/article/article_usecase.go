package articleusecase

import (
	"errors"
	article2 "github.com/amirex128/new_site_builder/internal/application/dto/article"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ArticleUsecase struct {
	*usecase.BaseUsecase
	articleRepo  repository2.IArticleRepository
	categoryRepo repository2.IArticleCategoryRepository
	mediaRepo    repository2.IMediaRepository
}

func NewArticleUsecase(c contract.IContainer) *ArticleUsecase {
	return &ArticleUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		articleRepo:  c.GetArticleRepo(),
		categoryRepo: c.GetArticleCategoryRepo(),
		mediaRepo:    c.GetMediaRepo(),
	}
}

func (u *ArticleUsecase) CreateArticleCommand(params *article2.CreateArticleCommand) (*resp.Response, error) {
	userID, _, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, err
	}

	var seoTags string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTags = strings.Join(params.SeoTags, ",")
	}
	newArticle := domain.Article{
		Title:        *params.Title,
		Description:  *params.Description,
		Body:         *params.Body,
		Slug:         *params.Slug,
		SiteID:       *params.SiteID,
		SeoTags:      seoTags,
		UserID:       *userID,
		VisitedCount: 0,
		ReviewCount:  0,
		Rate:         0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}
	err = u.articleRepo.Create(&newArticle)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد مقاله")
	}
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		for _, mediaID := range params.MediaIDs {
			err = u.articleRepo.AddMediaToArticle(newArticle.ID, mediaID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن رسانه به مقاله")
			}
		}
	}
	if params.CategoryIDs != nil && len(params.CategoryIDs) > 0 {
		for _, categoryID := range params.CategoryIDs {
			err = u.articleRepo.AddCategoryToArticle(newArticle.ID, categoryID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن دسته بندی به مقاله")
			}
		}
	}
	return resp.NewResponseData(resp.Created, newArticle, "مقاله با موفقیت ایجاد شد"), nil
}

func (u *ArticleUsecase) UpdateArticleCommand(params *article2.UpdateArticleCommand) (*resp.Response, error) {
	existingArticle, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	err = u.CheckAccessUserModel(existingArticle)
	if err != nil {
		return nil, err
	}
	if params.Title != nil {
		existingArticle.Title = *params.Title
	}
	if params.Description != nil {
		existingArticle.Description = *params.Description
	}
	if params.Body != nil {
		existingArticle.Body = *params.Body
	}
	if params.Slug != nil {
		existingArticle.Slug = *params.Slug
	}
	if params.SeoTags != nil {
		existingArticle.SeoTags = strings.Join(params.SeoTags, ",")
	}
	existingArticle.UpdatedAt = time.Now()
	err = u.articleRepo.Update(existingArticle)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ویرایش مقاله")
	}
	if params.MediaIDs != nil {
		_ = u.articleRepo.RemoveAllMediaFromArticle(existingArticle.ID)
		for _, mediaID := range params.MediaIDs {
			err = u.articleRepo.AddMediaToArticle(existingArticle.ID, mediaID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن رسانه به مقاله")
			}
		}
	}
	if params.CategoryIDs != nil {
		_ = u.articleRepo.RemoveAllCategoriesFromArticle(existingArticle.ID)
		for _, categoryID := range params.CategoryIDs {
			err = u.articleRepo.AddCategoryToArticle(existingArticle.ID, categoryID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن دسته بندی به مقاله")
			}
		}
	}
	return resp.NewResponseData(resp.Updated, existingArticle, "مقاله با موفقیت ویرایش شد"), nil
}

func (u *ArticleUsecase) DeleteArticleCommand(params *article2.DeleteArticleCommand) (*resp.Response, error) {
	existingArticle, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	err = u.CheckAccessUserModel(existingArticle)
	if err != nil {
		return nil, err
	}
	err = u.articleRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در حذف مقاله")
	}
	return resp.NewResponse(resp.Deleted, "مقاله با موفقیت حذف شد"), nil
}

func (u *ArticleUsecase) GetByIdArticleQuery(params *article2.GetByIdArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	mediaItems, err := u.articleRepo.GetArticleMedia(result.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"article": result,
		"media":   mediaItems,
	}, "مقاله با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetSingleArticleQuery(params *article2.GetSingleArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetBySlugAndSiteID(*params.Slug, *params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	mediaItems, err := u.articleRepo.GetArticleMedia(result.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"article": result,
		"media":   mediaItems,
	}, "مقاله با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetAllArticleQuery(params *article2.GetAllArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, articlesWithMedia, "مقالات با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetArticleByCategoryQuery(params *article2.GetArticleByCategoryQuery) (*resp.Response, error) {
	category, err := u.categoryRepo.GetBySlugAndSiteID(*params.Slug, *params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "دسته بندی یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	result, err := u.articleRepo.GetAllByCategoryID(category.ID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":    articlesWithMedia,
		"category": category,
	}, "مقالات با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetByFiltersSortArticleQuery(params *article2.GetByFiltersSortArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetAllByFilterAndSort(
		*params.SiteID,
		params.SelectedFilters,
		params.SelectedSort,
		params.PaginationRequestDto,
	)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, articlesWithMedia, "مقالات با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) AdminGetAllArticleQuery(params *article2.AdminGetAllArticleQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	result, err := u.articleRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, articlesWithMedia, "مقالات با موفقیت دریافت شد (مدیر)"), nil
}
