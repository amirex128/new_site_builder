package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	productusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func ArticleInit(c contract.IContainer) *v1.ArticleHandler {
	use := productusecase.NewArticleUsecase(c)
	productList := v1.NewArticleHandler(use)

	return productList
}
