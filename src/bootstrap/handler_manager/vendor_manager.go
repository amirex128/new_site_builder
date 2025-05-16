package handlermanager

import (
	v1 "go-boilerplate/src/internal/api/handler/http/v1"
	vendorusecase "go-boilerplate/src/internal/app/usecase/vend"
	"go-boilerplate/src/internal/contract"
)

func VendorInit(c contract.IContainer) *v1.VendorHandler {
	use := vendorusecase.NewVendorUsecase(c)
	vendorList := v1.NewVendorHandler(use)

	return vendorList
}
