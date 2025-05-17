package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IGatewayRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Gateway, int64, error)
	GetBySiteID(siteID int64) (domain.Gateway, error)
	GetByID(id int64) (domain.Gateway, error)
	Create(gateway domain.Gateway) error
	Update(gateway domain.Gateway) error
	Delete(id int64) error
}
