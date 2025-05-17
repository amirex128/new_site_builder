package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	customerticketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/customer_ticket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func CustomerTicketInit(c contract.IContainer) *v1.CustomerTicketHandler {
	use := customerticketusecase.NewCustomerTicketUsecase(c)
	handler := v1.NewCustomerTicketHandler(use)

	return handler
}
