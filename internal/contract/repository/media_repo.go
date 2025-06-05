package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IMediaRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Media], error)
	GetByID(id int64) (*domain.Media, error)
	Create(media *domain.Media) error
	Update(media *domain.Media) error
	Delete(id int64) error
}
