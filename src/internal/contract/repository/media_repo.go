package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IMediaRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Media, int64, error)
	GetByID(id int64) (domain.Media, error)
	Create(media domain.Media) error
	Update(media domain.Media) error
	Delete(id int64) error
}
