package bootstrap

import (
	handlermanager "github.com/amirex128/new_site_builder/src/bootstrap/handler_manager"
	"github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func HandlerBootstrap(c contract.IContainer) *HandlerManager {

	return &HandlerManager{
		ArticleHandlerV1: handlermanager.ArticleInit(c),
	}
}

type HandlerManager struct {
	ArticleHandlerV1 *v1.ArticleHandler
}
