package defaultthemeusecase

import (
	"errors"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/defaulttheme"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type DefaultThemeUsecase struct {
	logger         sflogger.Logger
	repo           repository.IDefaultThemeRepository
	mediaRepo      repository.IMediaRepository
	authContextSvc common.IAuthContextService
}

func NewDefaultThemeUsecase(c contract.IContainer) *DefaultThemeUsecase {
	return &DefaultThemeUsecase{
		logger:         c.GetLogger(),
		repo:           c.GetDefaultThemeRepo(),
		mediaRepo:      c.GetMediaRepo(),
		authContextSvc: c.GetAuthContextService(),
	}
}

func (u *DefaultThemeUsecase) CreateDefaultThemeCommand(params *defaulttheme.CreateDefaultThemeCommand) (any, error) {
	u.logger.Info("CreateDefaultThemeCommand called", map[string]interface{}{
		"name": *params.Name,
	})

	// Check if user is admin
	isAdmin, err := u.authContextSvc.IsAdmin()
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
	err = u.repo.Create(theme)
	if err != nil {
		return nil, err
	}

	// Get created theme with media relation
	createdTheme, err := u.repo.GetByID(theme.ID)
	if err != nil {
		return nil, err
	}

	return createdTheme, nil
}

func (u *DefaultThemeUsecase) UpdateDefaultThemeCommand(params *defaulttheme.UpdateDefaultThemeCommand) (any, error) {
	u.logger.Info("UpdateDefaultThemeCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Check if user is admin
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به ویرایش قالب پیش‌فرض هستند")
	}

	// Get existing theme
	existingTheme, err := u.repo.GetByID(*params.ID)
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
	err = u.repo.Update(existingTheme)
	if err != nil {
		return nil, err
	}

	// Get updated theme with media relation
	updatedTheme, err := u.repo.GetByID(existingTheme.ID)
	if err != nil {
		return nil, err
	}

	return updatedTheme, nil
}

func (u *DefaultThemeUsecase) DeleteDefaultThemeCommand(params *defaulttheme.DeleteDefaultThemeCommand) (any, error) {
	u.logger.Info("DeleteDefaultThemeCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Check if user is admin
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به حذف قالب پیش‌فرض هستند")
	}

	// Check if theme exists
	existingTheme, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("قالب پیش‌فرض مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Delete theme
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": existingTheme.ID,
	}, nil
}

func (u *DefaultThemeUsecase) GetByIdDefaultThemeQuery(params *defaulttheme.GetByIdDefaultThemeQuery) (any, error) {
	u.logger.Info("GetByIdDefaultThemeQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get theme by ID
	theme, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("قالب پیش‌فرض مورد نظر یافت نشد")
		}
		return nil, err
	}

	return theme, nil
}

func (u *DefaultThemeUsecase) GetAllDefaultThemeQuery(params *defaulttheme.GetAllDefaultThemeQuery) (any, error) {
	u.logger.Info("GetAllDefaultThemeQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get paginated list of themes
	themes, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items":     themes,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}
