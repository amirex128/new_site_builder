package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	"github.com/amirex128/new_site_builder/internal/application/dto/website"
	"github.com/amirex128/new_site_builder/internal/application/usecase/website"
	"github.com/gin-gonic/gin"
)

type WebsiteHandler struct {
	usecase   *websiteusecase.WebsiteUsecase
	validator *utils2.ValidationHelper
}

func NewWebsiteHandler(usc *websiteusecase.WebsiteUsecase) *WebsiteHandler {
	return &WebsiteHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// GetByDomainPage godoc
// @Summary      Get page by domain
// @Description  Retrieves page content for a specific domain
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetByDomainPageQuery  true  "Domain and page parameters"
// @success      200      {object}  utils.Result                  "Page content"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      404      {object}  utils.Result                  "Page or domain not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /website/page [get]
func (h *WebsiteHandler) GetByDomainPage(c *gin.Context) {
	var params website.GetByDomainPageQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByDomainPageQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByDomainHeaderFooter godoc
// @Summary      Get header and footer by domain
// @Description  Retrieves header and footer content for a specific domain
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetByDomainHeaderFooterQuery  true  "Domain parameters"
// @success      200      {object}  utils.Result                          "Header and footer content"
// @Failure      400      {object}  utils.Result                          "Validation error"
// @Failure      404      {object}  utils.Result                          "Domain not found"
// @Failure      500      {object}  utils.Result                          "Internal server error"
// @Router       /website/header-footer [get]
func (h *WebsiteHandler) GetByDomainHeaderFooter(c *gin.Context) {
	var params website.GetByDomainHeaderFooterQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByDomainHeaderFooterQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ProductSearchList godoc
// @Summary      Search products
// @Description  Searches for products based on provided criteria
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.ProductSearchListQuery  true  "Search parameters"
// @success      200      {object}  utils.Result                    "List of products matching search criteria"
// @Failure      400      {object}  utils.Result                    "Validation error"
// @Failure      500      {object}  utils.Result                    "Internal server error"
// @Router       /website/product/search [get]
func (h *WebsiteHandler) ProductSearchList(c *gin.Context) {
	var params website.ProductSearchListQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.ProductSearchListQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetFiltersSortArticle godoc
// @Summary      Get article filters and sorting options
// @Description  Retrieves available filters and sorting options for articles
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetFiltersSortArticleQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                        "Article filters and sorting options"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /website/article/filters-sort [get]
func (h *WebsiteHandler) GetFiltersSortArticle(c *gin.Context) {
	var params website.GetFiltersSortArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetFiltersSortArticleQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetFiltersSortProduct godoc
// @Summary      Get product filters and sorting options
// @Description  Retrieves available filters and sorting options for products
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetFiltersSortProductQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                        "Product filters and sorting options"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /website/product/filters-sort [get]
func (h *WebsiteHandler) GetFiltersSortProduct(c *gin.Context) {
	var params website.GetFiltersSortProductQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetFiltersSortProductQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetArticlesByCategorySlug godoc
// @Summary      Get articles by category slug
// @Description  Retrieves articles belonging to a specific category identified by slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetArticlesByCategorySlugQuery  true  "Category slug and query parameters"
// @success      200      {object}  utils.Result                            "List of articles in the category"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      404      {object}  utils.Result                            "Category not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /website/article/category [get]
func (h *WebsiteHandler) GetArticlesByCategorySlug(c *gin.Context) {
	var params website.GetArticlesByCategorySlugQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetArticlesByCategorySlugQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetProductsByCategorySlug godoc
// @Summary      Get products by category slug
// @Description  Retrieves products belonging to a specific category identified by slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetProductsByCategorySlugQuery  true  "Category slug and query parameters"
// @success      200      {object}  utils.Result                            "List of products in the category"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      404      {object}  utils.Result                            "Category not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /website/product/category [get]
func (h *WebsiteHandler) GetProductsByCategorySlug(c *gin.Context) {
	var params website.GetProductsByCategorySlugQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetProductsByCategorySlugQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetSingleArticleBySlug godoc
// @Summary      Get article by slug
// @Description  Retrieves a specific article identified by its slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  body      website.GetSingleArticleBySlugQuery  true  "Article slug"
// @success      200      {object}  utils.Result                         "Article details"
// @Failure      400      {object}  utils.Result                         "Validation error"
// @Failure      404      {object}  utils.Result                         "Article not found"
// @Failure      500      {object}  utils.Result                         "Internal server error"
// @Router       /website/article [get]
func (h *WebsiteHandler) GetSingleArticleBySlug(c *gin.Context) {
	var params website.GetSingleArticleBySlugQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetSingleArticleBySlugQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetSingleProductBySlug godoc
// @Summary      Get product by slug
// @Description  Retrieves a specific product identified by its slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  body      website.GetSingleProductBySlugQuery  true  "Product slug"
// @success      200      {object}  utils.Result                         "Product details"
// @Failure      400      {object}  utils.Result                         "Validation error"
// @Failure      404      {object}  utils.Result                         "Product not found"
// @Failure      500      {object}  utils.Result                         "Internal server error"
// @Router       /website/product [get]
func (h *WebsiteHandler) GetSingleProductBySlug(c *gin.Context) {
	var params website.GetSingleProductBySlugQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetSingleProductBySlugQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
