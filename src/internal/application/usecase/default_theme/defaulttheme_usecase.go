package defaultthemeusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/defaulttheme"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type DefaultThemeUsecase struct {
	*usecase.BaseUsecase
	defaultThemeRepo repository.IDefaultThemeRepository
	mediaRepo        repository.IMediaRepository
	authContext      func(c *gin.Context) service.IAuthService
}

func NewDefaultThemeUsecase(c contract.IContainer) *DefaultThemeUsecase {
	return &DefaultThemeUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		defaultThemeRepo: c.GetDefaultThemeRepo(),
		mediaRepo:        c.GetMediaRepo(),
		authContext:      c.GetAuthTransientService(),
	}
}

func (u *DefaultThemeUsecase) CreateDefaultThemeCommand(params *defaulttheme.CreateDefaultThemeCommand) (*resp.Response, error) {
	u.Logger.Info("CreateDefaultThemeCommand called", map[string]interface{}{
		"name": *params.Name,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به ایجاد قالب پیش‌فرض هستند")
	}

	// Verify media exists
	_, err = u.mediaRepo.GetByID(int64(*params.MediaID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تصویر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Create the DefaultTheme entity
	var description string
	if params.Description != nil {
		description = *params.Description
	}

	var demo string
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

	// Create in repository
	err = u.defaultThemeRepo.Create(theme)
	if err != nil {
		return nil, err
	}

	// Get created theme with media relation
	createdTheme, err := u.defaultThemeRepo.GetByID(theme.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Success,
		resp.Data{
			"default_theme": createdTheme,
		},
		"success"), nil
}

func (u *DefaultThemeUsecase) UpdateDefaultThemeCommand(params *defaulttheme.UpdateDefaultThemeCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateDefaultThemeCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به ویرایش قالب پیش‌فرض هستند")
	}

	// Get existing theme
	existingTheme, err := u.defaultThemeRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("قالب پیش‌فرض مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Update fields if provided
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
		// Verify media exists
		_, err = u.mediaRepo.GetByID(int64(*params.MediaID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("تصویر مورد نظر یافت نشد")
			}
			return nil, err
		}
		existingTheme.MediaID = int64(*params.MediaID)
	}

	if params.Pages != nil {
		existingTheme.Pages = *params.Pages
	}

	existingTheme.UpdatedAt = time.Now()

	// Update in repository
	err = u.defaultThemeRepo.Update(existingTheme)
	if err != nil {
		return nil, err
	}

	// Get updated theme with media relation
	updatedTheme, err := u.defaultThemeRepo.GetByID(existingTheme.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Success,
		resp.Data{
			"default_theme": updatedTheme,
		},
		"success"), nil

}

func (u *DefaultThemeUsecase) DeleteDefaultThemeCommand(params *defaulttheme.DeleteDefaultThemeCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteDefaultThemeCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به حذف قالب پیش‌فرض هستند")
	}

	// Check if theme exists
	existingTheme, err := u.defaultThemeRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("قالب پیش‌فرض مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Delete theme
	err = u.defaultThemeRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Success,
		resp.Data{
			"default_theme": existingTheme,
		},
		"success"), nil
}

func (u *DefaultThemeUsecase) GetByIdDefaultThemeQuery(params *defaulttheme.GetByIdDefaultThemeQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdDefaultThemeQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get theme by ID
	theme, err := u.defaultThemeRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("قالب پیش‌فرض مورد نظر یافت نشد")
		}
		return nil, err
	}

	return resp.NewResponseData(resp.Success,
		resp.Data{
			"default_theme": theme,
		},
		"success"), nil
}

func (u *DefaultThemeUsecase) GetAllDefaultThemeQuery(params *defaulttheme.GetAllDefaultThemeQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllDefaultThemeQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	themesResult, err := u.defaultThemeRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return resp.NewResponseData(resp.Success,
		resp.Data{
			"default_theme": themesResult.Items,
			"total":         themesResult.TotalCount,
			"page":          themesResult.PageNumber,
			"pageSize":      params.PageSize,
			"totalPage":     themesResult.TotalPages,
		},
		"success"), nil
}
