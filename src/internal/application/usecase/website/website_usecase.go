package websiteusecase

import (
	"fmt"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/website"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type WebsiteUsecase struct {
	*usecase.BaseUsecase
	articleRepo         repository.IArticleRepository
	articleCategoryRepo repository.IArticleCategoryRepository
	productRepo         repository.IProductRepository
	productCategoryRepo repository.IProductCategoryRepository
	pageRepo            repository.IPageRepository
	headerFooterRepo    repository.IHeaderFooterRepository
	siteRepo            repository.ISiteRepository
}

func NewWebsiteUsecase(c contract.IContainer) *WebsiteUsecase {
	return &WebsiteUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		articleRepo:         c.GetArticleRepo(),
		articleCategoryRepo: c.GetArticleCategoryRepo(),
		productRepo:         c.GetProductRepo(),
		productCategoryRepo: c.GetProductCategoryRepo(),
		pageRepo:            c.GetPageRepo(),
		headerFooterRepo:    c.GetHeaderFooterRepo(),
		siteRepo:            c.GetSiteRepo(),
	}
}

func (u *WebsiteUsecase) GetByDomainPageQuery(params *website.GetByDomainPageQuery) (*resp.Response, error) {
	// Implementation to get page by domain and path
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find the page by slug for that site using the page repository

	// Return empty result for now
	return resp.NewResponse(resp.Success, "success"), nil
}

func (u *WebsiteUsecase) GetByDomainHeaderFooterQuery(params *website.GetByDomainHeaderFooterQuery) (*resp.Response, error) {
	// Implementation to get header/footer by domain
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find the header/footer for that site using the header/footer repository

	// Return empty result for now
	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) ProductSearchListQuery(params *website.ProductSearchListQuery) (*resp.Response, error) {
	// Implementation for searching products by domain
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find all products for that site using the product repository with pagination

	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) GetFiltersSortArticleQuery(params *website.GetFiltersSortArticleQuery) (*resp.Response, error) {
	// Implementation for getting articles with filters and sorting
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Get all articles for that site using the article repository with pagination
	// 3. Apply filters and sorting as specified in the params

	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) GetFiltersSortProductQuery(params *website.GetFiltersSortProductQuery) (*resp.Response, error) {
	// Implementation for getting products with filters and sorting
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Get all products for that site using the product repository with pagination
	// 3. Apply filters and sorting as specified in the params

	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) GetArticlesByCategorySlugQuery(params *website.GetArticlesByCategorySlugQuery) (*resp.Response, error) {
	// Implementation for getting articles by category slug
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find the category by slug using the category repository
	// 3. Find all articles for that category using the article repository with pagination

	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) GetProductsByCategorySlugQuery(params *website.GetProductsByCategorySlugQuery) (*resp.Response, error) {
	// Implementation for getting products by category slug
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find the category by slug using the category repository
	// 3. Find all products for that category using the product repository with pagination

	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) GetSingleArticleBySlugQuery(params *website.GetSingleArticleBySlugQuery) (*resp.Response, error) {
	// Implementation for getting a single article by slug
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find the article by slug using the article repository

	return resp.NewResponse(resp.Success, "success"), nil

}

func (u *WebsiteUsecase) GetSingleProductBySlugQuery(params *website.GetSingleProductBySlugQuery) (*resp.Response, error) {
	// Implementation for getting a single product by slug
	fmt.Println(params)

	// This is a placeholder implementation - in a real system, you would:
	// 1. Find the site by domain using the site repository
	// 2. Find the product by slug using the product repository

	return resp.NewResponse(resp.Success, "success"), nil

}
