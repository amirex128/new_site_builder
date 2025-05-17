package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	ticketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/ticket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func TicketInit(c contract.IContainer) *v1.TicketHandler {
	use := ticketusecase.NewTicketUsecase(c)
	handler := v1.NewTicketHandler(use)

	return handler
}
