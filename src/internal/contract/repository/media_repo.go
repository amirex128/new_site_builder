package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IMediaRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Media, int64, error)
	GetByID(id int64) (domain.Media, error)
	Create(media domain.Media) error
	Update(media domain.Media) error
	Delete(id int64) error
}
