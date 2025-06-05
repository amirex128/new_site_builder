package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IDefaultThemeRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.DefaultTheme], error)
	GetByID(id int64) (*domain.DefaultTheme, error)
	Create(theme *domain.DefaultTheme) error
	Update(theme *domain.DefaultTheme) error
	Delete(id int64) error
}
