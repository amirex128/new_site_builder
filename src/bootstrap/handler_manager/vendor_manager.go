package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	vendorusecase "github.com/amirex128/new_site_builder/src/internal/app/usecase/vend"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func VendorInit(c contract.IContainer) *v1.VendorHandler {
	use := vendorusecase.NewVendorUsecase(c)
	vendorList := v1.NewVendorHandler(use)

	return vendorList
}
