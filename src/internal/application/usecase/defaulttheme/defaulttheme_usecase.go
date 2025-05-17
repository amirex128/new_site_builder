package defaultthemeusecase

import (
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/defaulttheme"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type DefaultThemeUsecase struct {
	logger sflogger.Logger
	repo   repository.IDefaultThemeRepository
}

func NewDefaultThemeUsecase(c contract.IContainer) *DefaultThemeUsecase {
	return &DefaultThemeUsecase{
		logger: c.GetLogger(),
		repo:   c.GetDefaultThemeRepo(),
	}
}

func (u *DefaultThemeUsecase) CreateDefaultThemeCommand(params *defaulttheme.CreateDefaultThemeCommand) (any, error) {
	// Implementation for creating a default theme
	fmt.Println(params)

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

	err := u.repo.Create(theme)
	if err != nil {
		return nil, err
	}

	return theme, nil
}

func (u *DefaultThemeUsecase) UpdateDefaultThemeCommand(params *defaulttheme.UpdateDefaultThemeCommand) (any, error) {
	// Implementation for updating a default theme
	fmt.Println(params)

	existingTheme, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
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
		existingTheme.MediaID = int64(*params.MediaID)
	}

	if params.Pages != nil {
		existingTheme.Pages = *params.Pages
	}

	existingTheme.UpdatedAt = time.Now()

	err = u.repo.Update(existingTheme)
	if err != nil {
		return nil, err
	}

	return existingTheme, nil
}

func (u *DefaultThemeUsecase) DeleteDefaultThemeCommand(params *defaulttheme.DeleteDefaultThemeCommand) (any, error) {
	// Implementation for deleting a default theme
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *DefaultThemeUsecase) GetByIdDefaultThemeQuery(params *defaulttheme.GetByIdDefaultThemeQuery) (any, error) {
	// Implementation to get default theme by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *DefaultThemeUsecase) GetAllDefaultThemeQuery(params *defaulttheme.GetAllDefaultThemeQuery) (any, error) {
	// Implementation to get all default themes
	fmt.Println(params)

	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
