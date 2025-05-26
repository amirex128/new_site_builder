package defaultthemeusecase

import (
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/defaulttheme"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type DefaultThemeUsecase struct {
	*usecase.BaseUsecase
	defaultThemeRepo repository.IDefaultThemeRepository
	mediaRepo        repository.IMediaRepository
}

func NewDefaultThemeUsecase(c contract.IContainer) *DefaultThemeUsecase {
	return &DefaultThemeUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		defaultThemeRepo: c.GetDefaultThemeRepo(),
		mediaRepo:        c.GetMediaRepo(),
	}
}

func (u *DefaultThemeUsecase) CreateDefaultThemeCommand(params *defaulttheme.CreateDefaultThemeCommand) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	_, err = u.mediaRepo.GetByID(int64(*params.MediaID))
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "تصویر مورد نظر یافت نشد")
	}

	description := ""
	if params.Description != nil {
		description = *params.Description
	}
	demo := ""
	if params.Demo != nil {
		demo = *params.Demo
	}

	theme := domain.DefaultTheme{
		Name:        *params.Name,
		Description: description,
		Demo:        demo,
		MediaID:     int64(*params.MediaID),
		Pages:       *params.Pages,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	err = u.defaultThemeRepo.Create(&theme)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطایی در ایجاد قالب پیش‌فرض رخ داده است")
	}

	return resp.NewResponseData(resp.Created, theme, "قالب پیش‌فرض با موفقیت ایجاد شد"), nil
}

func (u *DefaultThemeUsecase) UpdateDefaultThemeCommand(params *defaulttheme.UpdateDefaultThemeCommand) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	existingTheme, err := u.defaultThemeRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "قالب پیش‌فرض مورد نظر یافت نشد")
	}

	if params.Name != nil {
		existingTheme.Name = *params.Name
	}
	if params.Description != nil {
		existingTheme.Description = *params.Description
	}
	if params.Demo != nil {
		existingTheme.Demo = *params.Demo
	}
	if params.MediaID != nil {
		_, err = u.mediaRepo.GetByID(int64(*params.MediaID))
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "تصویر مورد نظر یافت نشد")
		}
		existingTheme.MediaID = int64(*params.MediaID)
	}
	if params.Pages != nil {
		existingTheme.Pages = *params.Pages
	}
	existingTheme.UpdatedAt = time.Now()

	err = u.defaultThemeRepo.Update(existingTheme)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطایی در بروزرسانی قالب پیش‌فرض رخ داده است")
	}

	return resp.NewResponseData(resp.Updated, existingTheme, "قالب پیش‌فرض با موفقیت بروزرسانی شد"), nil
}

func (u *DefaultThemeUsecase) DeleteDefaultThemeCommand(params *defaulttheme.DeleteDefaultThemeCommand) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	existingTheme, err := u.defaultThemeRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "قالب پیش‌فرض مورد نظر یافت نشد")
	}

	err = u.defaultThemeRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطایی در حذف قالب پیش‌فرض رخ داده است")
	}

	return resp.NewResponseData(resp.Deleted, existingTheme, "قالب پیش‌فرض با موفقیت حذف شد"), nil
}

func (u *DefaultThemeUsecase) GetByIdDefaultThemeQuery(params *defaulttheme.GetByIdDefaultThemeQuery) (*resp.Response, error) {
	theme, err := u.defaultThemeRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "قالب پیش‌فرض مورد نظر یافت نشد")
	}

	return resp.NewResponseData(resp.Retrieved, theme, "قالب پیش‌فرض با موفقیت دریافت شد"), nil
}

func (u *DefaultThemeUsecase) GetAllDefaultThemeQuery(params *defaulttheme.GetAllDefaultThemeQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	themesResult, err := u.defaultThemeRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطایی در دریافت قالب‌های پیش‌فرض رخ داده است")
	}

	return resp.NewResponseData(resp.Retrieved, themesResult, "لیست قالب‌های پیش‌فرض با موفقیت دریافت شد"), nil
}
