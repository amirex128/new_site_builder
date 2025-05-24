package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IDefaultThemeRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.DefaultTheme], int64, error)
	GetByID(id int64) (domain.DefaultTheme, error)
	Create(theme domain.DefaultTheme) error
	Update(theme domain.DefaultTheme) error
	Delete(id int64) error
}
