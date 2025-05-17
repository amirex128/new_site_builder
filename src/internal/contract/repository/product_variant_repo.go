package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IProductVariantRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.ProductVariant, int64, error)
	GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductVariant, int64, error)
	GetByID(id int64) (domain.ProductVariant, error)
	Create(variant domain.ProductVariant) error
	Update(variant domain.ProductVariant) error
	Delete(id int64) error
}
