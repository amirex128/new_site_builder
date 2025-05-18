package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	articlecategoryusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article_category"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func ArticleCategoryInit(c contract.IContainer) *v1.ArticleCategoryHandler {
	use := articlecategoryusecase.NewArticleCategoryUsecase(c)
	handler := v1.NewBlogCategoryHandler(use)

	return handler
}
