package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	fileitemusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/file_item"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func FileItemInit(c contract.IContainer) *v1.FileItemHandler {
	use := fileitemusecase.NewFileItemUsecase(c)
	handler := v1.NewFileItemHandler(use)

	return handler
}
