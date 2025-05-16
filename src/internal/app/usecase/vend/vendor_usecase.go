package vendorusecase

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	vendordto "go-boilerplate/src/internal/app/dto/vend"
	vendorinitializer "go-boilerplate/src/internal/app/initializer/vend"
	"go-boilerplate/src/internal/contract"
	"go-boilerplate/src/internal/contract/repository"
)

type VendorUsecase struct {
	logger sflogger.Logger
	repo   repository.ISampleRepository
}

func NewVendorUsecase(c contract.IContainer) *VendorUsecase {
	return &VendorUsecase{
		logger: c.GetLogger(),
		repo:   c.GetSimpleRepo(),
	}
}

func (u *VendorUsecase) VendorList(dto *vendordto.VendorDto) ([]any, error) {

	params := vendorinitializer.NewVendorInitializer(u.repo, dto).
		SetCPC().
		SetPickUp().
		SetDelivery().
		Build()

	fmt.Println(params)

	// work mysql elastic repository

	return nil, nil
}
