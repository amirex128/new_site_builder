package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	paymentusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/payment"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func PaymentInit(c contract.IContainer) *v1.PaymentHandler {
	use := paymentusecase.NewPaymentUsecase(c)
	handler := v1.NewPaymentHandler(use)

	return handler
}
