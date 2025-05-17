package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type DefaultThemeRepo struct {
	database *gorm.DB
}

func NewDefaultThemeRepository(db *gorm.DB) *DefaultThemeRepo {
	return &DefaultThemeRepo{
		database: db,
	}
}

func (r *DefaultThemeRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.DefaultTheme, int64, error) {
	var themes []domain.DefaultTheme
	var count int64

	query := r.database.Model(&domain.DefaultTheme{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&themes)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return themes, count, nil
}

func (r *DefaultThemeRepo) GetByID(id int64) (domain.DefaultTheme, error) {
	var theme domain.DefaultTheme
	result := r.database.First(&theme, id)
	if result.Error != nil {
		return theme, result.Error
	}
	return theme, nil
}

func (r *DefaultThemeRepo) Create(theme domain.DefaultTheme) error {
	result := r.database.Create(&theme)
	return result.Error
}

func (r *DefaultThemeRepo) Update(theme domain.DefaultTheme) error {
	result := r.database.Save(&theme)
	return result.Error
}

func (r *DefaultThemeRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.DefaultTheme{}, id)
	return result.Error
}
