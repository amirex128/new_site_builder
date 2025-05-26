package articleusecase

import (
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/article"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ArticleUsecase struct {
	*usecase.BaseUsecase
	articleRepo  repository.IArticleRepository
	categoryRepo repository.IArticleCategoryRepository
	mediaRepo    repository.IMediaRepository
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

func (u *ArticleUsecase) CreateArticleCommand(params *article.CreateArticleCommand) (*resp.Response, error) {
	userID, _, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
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

func (u *ArticleUsecase) UpdateArticleCommand(params *article.UpdateArticleCommand) (*resp.Response, error) {
	userID, _, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	existingArticle, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
	}
	err = u.CheckAccessUserModel(existingArticle, userID)
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
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

func (u *ArticleUsecase) DeleteArticleCommand(params *article.DeleteArticleCommand) (*resp.Response, error) {
	userID, _, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	existingArticle, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
	}
	err = u.CheckAccessUserModel(existingArticle, userID)
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	err = u.articleRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در حذف مقاله")
	}
	return resp.NewResponse(resp.Deleted, "مقاله با موفقیت حذف شد"), nil
}

func (u *ArticleUsecase) GetByIdArticleQuery(params *article.GetByIdArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
	}
	mediaItems, err := u.articleRepo.GetArticleMedia(result.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"article": result,
		"media":   mediaItems,
	}, "مقاله با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetSingleArticleQuery(params *article.GetSingleArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetBySlugAndSiteID(*params.Slug, *params.SiteID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "مقاله یافت نشد")
	}
	mediaItems, err := u.articleRepo.GetArticleMedia(result.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"article": result,
		"media":   mediaItems,
	}, "مقاله با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetAllArticleQuery(params *article.GetAllArticleQuery) (*resp.Response, error) {
	result, err := u.articleRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, articlesWithMedia, "مقالات با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) GetArticleByCategoryQuery(params *article.GetArticleByCategoryQuery) (*resp.Response, error) {
	category, err := u.categoryRepo.GetBySlugAndSiteID(*params.Slug, *params.SiteID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "دسته‌بندی یافت نشد")
	}
	result, err := u.articleRepo.GetAllByCategoryID(category.ID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
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

func (u *ArticleUsecase) GetByFiltersSortArticleQuery(params *article.GetByFiltersSortArticleQuery) (*resp.Response, error) {
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
			return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, articlesWithMedia, "مقالات با موفقیت دریافت شد"), nil
}

func (u *ArticleUsecase) AdminGetAllArticleQuery(params *article.AdminGetAllArticleQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط ادمین ها میتوانند به این منور دسترسی داشته باشند")
	}

	result, err := u.articleRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت مقالات")
	}
	articlesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, article := range result.Items {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
		}
		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, articlesWithMedia, "مقالات با موفقیت دریافت شد (مدیر)"), nil
}
