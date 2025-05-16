package productusecase

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	productdto "go-boilerplate/src/internal/app/dto/product"
	"go-boilerplate/src/internal/app/initializer/product"
	"go-boilerplate/src/internal/contract"
	"go-boilerplate/src/internal/contract/repository"
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

	params = productinitializer.NewProductInitializer(u.repo, params).
		SetCPC().
		SetPickUp().
		SetDelivery().
		Build()

	fmt.Println(params)

	// work mysql elastic repository

	return nil, nil
}
