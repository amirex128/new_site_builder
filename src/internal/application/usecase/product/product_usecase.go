package productusecase

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	productdto "github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type ProductUsecase struct {
	logger sflogger.Logger
	repo   repository.ISampleRepository
}

func NewProductUsecase(c contract.IContainer) *ProductUsecase {
	return &ProductUsecase{
		logger: c.GetLogger(),
		repo:   c.GetSimpleRepo(),
	}
}

func (u *ProductUsecase) ProductList(params *productdto.ProductDto) ([]any, error) {

	fmt.Println(params)

	// work mysql elastic repository

	return nil, nil
}
