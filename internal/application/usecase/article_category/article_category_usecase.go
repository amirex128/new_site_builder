package articlecategoryusecase

import (
	"errors"
	article_category2 "github.com/amirex128/new_site_builder/internal/application/dto/article_category"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ArticleCategoryUsecase struct {
	*usecase.BaseUsecase
	categoryRepo repository2.IArticleCategoryRepository
	mediaRepo    repository2.IMediaRepository
}

func NewArticleCategoryUsecase(c contract.IContainer) *ArticleCategoryUsecase {
	return &ArticleCategoryUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		categoryRepo: c.GetArticleCategoryRepo(),
		mediaRepo:    c.GetMediaRepo(),
	}
}

func (u *ArticleCategoryUsecase) CreateCategoryCommand(params *article_category2.CreateCategoryCommand) (*resp.Response, error) {
	userID, _, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, err
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
	err = u.categoryRepo.Create(&newCategory)
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

func (u *ArticleCategoryUsecase) UpdateCategoryCommand(params *article_category2.UpdateCategoryCommand) (*resp.Response, error) {
	if params.ID == nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه دسته‌بندی اجباری است")
	}
	existingCategory, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "دسته بندی یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	err = u.CheckAccessUserModel(existingCategory)
	if err != nil {
		return nil, err
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

func (u *ArticleCategoryUsecase) DeleteCategoryCommand(params *article_category2.DeleteCategoryCommand) (*resp.Response, error) {
	if params.ID == nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه دسته‌بندی اجباری است")
	}
	existingCategory, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "دسته بندی یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	err = u.CheckAccessUserModel(existingCategory)
	if err != nil {
		return nil, err
	}
	err = u.categoryRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در حذف دسته‌بندی")
	}
	return resp.NewResponse(resp.Deleted, "دسته‌بندی با موفقیت حذف شد"), nil
}

func (u *ArticleCategoryUsecase) GetByIdCategoryQuery(params *article_category2.GetByIdCategoryQuery) (*resp.Response, error) {
	if params.ID == nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه دسته‌بندی اجباری است")
	}
	result, err := u.categoryRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "دسته بندی یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	mediaItems, err := u.categoryRepo.GetCategoryMedia(result.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"category": result,
		"media":    mediaItems,
	}, "دسته‌بندی با موفقیت دریافت شد"), nil
}

func (u *ArticleCategoryUsecase) GetAllCategoryQuery(params *article_category2.GetAllCategoryQuery) (*resp.Response, error) {
	result, err := u.categoryRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت دسته‌بندی‌ها")
	}
	categoriesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, category := range result.Items {
		media, err := u.categoryRepo.GetCategoryMedia(category.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		categoriesWithMedia[i] = map[string]interface{}{
			"category": category,
			"media":    media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, categoriesWithMedia, "دسته‌بندی‌ها با موفقیت دریافت شدند"), nil
}

func (u *ArticleCategoryUsecase) AdminGetAllCategoryQuery(params *article_category2.AdminGetAllCategoryQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}
	result, err := u.categoryRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت دسته‌بندی‌ها")
	}
	categoriesWithMedia := make([]map[string]interface{}, len(result.Items))
	for i, category := range result.Items {
		media, err := u.categoryRepo.GetCategoryMedia(category.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "رسانه یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		categoriesWithMedia[i] = map[string]interface{}{
			"category": category,
			"media":    media,
		}
	}
	return resp.NewResponseData(resp.Retrieved, categoriesWithMedia, "دسته‌بندی‌ها با موفقیت دریافت شدند (مدیر)"), nil
}
