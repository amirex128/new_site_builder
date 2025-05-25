package articlecategoryusecase

import (
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/article_category"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/gin-gonic/gin"
)

type ArticleCategoryUsecase struct {
	*usecase.BaseUsecase
	categoryRepo repository.IArticleCategoryRepository
	mediaRepo    repository.IMediaRepository
	authContext  func(c *gin.Context) service.IAuthService
}

func NewArticleCategoryUsecase(c contract.IContainer) *ArticleCategoryUsecase {
	return &ArticleCategoryUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		categoryRepo: c.GetArticleCategoryRepo(),
		mediaRepo:    c.GetMediaRepo(),
		authContext:  c.GetAuthTransientService(),
	}
}

func (u *ArticleCategoryUsecase) CreateCategoryCommand(params *article_category.CreateCategoryCommand) (*resp.Response, error) {
	userID, _, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	var seoTags string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTags = strings.Join(params.SeoTags, ",")
	}
	var description string
	if params.Description != nil {
		description = *params.Description
	}
	newCategory := domain.ArticleCategory{
		Name:             *params.Name,
		Slug:             *params.Slug,
		Description:      description,
		ParentCategoryID: params.ParentCategoryID,
		SiteID:           *params.SiteID,
		Order:            *params.Order,
		SeoTags:          seoTags,
		UserID:           *userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}
	err = u.categoryRepo.Create(newCategory)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد دسته‌بندی")
	}
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		for _, mediaID := range params.MediaIDs {
			err = u.categoryRepo.AddMediaToCategory(newCategory.ID, mediaID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن رسانه به دسته‌بندی")
			}
		}
	}
	return resp.NewResponseData(resp.Created, newCategory, "دسته‌بندی با موفقیت ایجاد شد"), nil
}

func (u *ArticleCategoryUsecase) UpdateCategoryCommand(params *article_category.UpdateCategoryCommand) (*resp.Response, error) {
	userID, _, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if params.ID == nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه دسته‌بندی اجباری است")
	}
	existingCategory, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "دسته‌بندی یافت نشد")
	}
	err = u.CheckAccessUserModel(&existingCategory, userID)
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if params.Name != nil {
		existingCategory.Name = *params.Name
	}
	if params.Description != nil {
		existingCategory.Description = *params.Description
	}
	if params.ParentCategoryID != nil {
		existingCategory.ParentCategoryID = params.ParentCategoryID
	}
	if params.Slug != nil {
		existingCategory.Slug = *params.Slug
	}
	if params.Order != nil {
		existingCategory.Order = *params.Order
	}
	if params.SeoTags != nil {
		existingCategory.SeoTags = strings.Join(params.SeoTags, ",")
	}
	existingCategory.UpdatedAt = time.Now()
	err = u.categoryRepo.Update(existingCategory)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ویرایش دسته‌بندی")
	}
	if params.MediaIDs != nil {
		_ = u.categoryRepo.RemoveAllMediaFromCategory(existingCategory.ID)
		for _, mediaID := range params.MediaIDs {
			err = u.categoryRepo.AddMediaToCategory(existingCategory.ID, mediaID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن رسانه به دسته‌بندی")
			}
		}
	}
	return resp.NewResponseData(resp.Updated, existingCategory, "دسته‌بندی با موفقیت ویرایش شد"), nil
}

func (u *ArticleCategoryUsecase) DeleteCategoryCommand(params *article_category.DeleteCategoryCommand) (*resp.Response, error) {
	userID, _, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if params.ID == nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه دسته‌بندی اجباری است")
	}
	existingCategory, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "دسته‌بندی یافت نشد")
	}
	err = u.CheckAccessUserModel(&existingCategory, userID)
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	err = u.categoryRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در حذف دسته‌بندی")
	}
	return resp.NewResponse(resp.Deleted, "دسته‌بندی با موفقیت حذف شد"), nil
}

func (u *ArticleCategoryUsecase) GetByIdCategoryQuery(params *article_category.GetByIdCategoryQuery) (*resp.Response, error) {
	if params.ID == nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه دسته‌بندی اجباری است")
	}
	result, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "دسته‌بندی یافت نشد")
	}
	mediaItems, err := u.categoryRepo.GetCategoryMedia(result.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"category": result,
		"media":    mediaItems,
	}, "دسته‌بندی با موفقیت دریافت شد"), nil
}

func (u *ArticleCategoryUsecase) GetAllCategoryQuery(params *article_category.GetAllCategoryQuery) (*resp.Response, error) {
	result, err := u.categoryRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت دسته‌بندی‌ها")
	}
	categoriesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, category := range result.Items {
		media, err := u.categoryRepo.GetCategoryMedia(category.ID)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
		}
		categoriesWithMedia[i] = map[string]interface{}{
			"category": category,
			"media":    media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, categoriesWithMedia, "دسته‌بندی‌ها با موفقیت دریافت شدند"), nil
}

func (u *ArticleCategoryUsecase) AdminGetAllCategoryQuery(params *article_category.AdminGetAllCategoryQuery) (*resp.Response, error) {
	result, err := u.categoryRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت دسته‌بندی‌ها")
	}
	categoriesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, category := range result.Items {
		media, err := u.categoryRepo.GetCategoryMedia(category.ID)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "رسانه‌ها یافت نشد")
		}
		categoriesWithMedia[i] = map[string]interface{}{
			"category": category,
			"media":    media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, categoriesWithMedia, "دسته‌بندی‌ها با موفقیت دریافت شدند (مدیر)"), nil
}
