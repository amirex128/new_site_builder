package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/website"
	websiteusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/website"
	"github.com/gin-gonic/gin"
)

type WebsiteHandler struct {
	usecase   *websiteusecase.WebsiteUsecase
	validator *utils.ValidationHelper
}

func NewWebsiteHandler(usc *websiteusecase.WebsiteUsecase) *WebsiteHandler {
	return &WebsiteHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// GetByDomainPage godoc
// @Summary      Get page by domain
// @Description  Retrieves page content for a specific domain
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetByDomainPageQuery  true  "Domain and page parameters"
// @Success      200      {object}  utils.Result                  "Page content"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      404      {object}  utils.Result                  "Page or domain not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /website/page [get]
func (h *WebsiteHandler) GetByDomainPage(c *gin.Context) {
	var params website.GetByDomainPageQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByDomainPageQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetByDomainHeaderFooter godoc
// @Summary      Get header and footer by domain
// @Description  Retrieves header and footer content for a specific domain
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetByDomainHeaderFooterQuery  true  "Domain parameters"
// @Success      200      {object}  utils.Result                          "Header and footer content"
// @Failure      400      {object}  utils.Result                          "Validation error"
// @Failure      404      {object}  utils.Result                          "Domain not found"
// @Failure      500      {object}  utils.Result                          "Internal server error"
// @Router       /website/header-footer [get]
func (h *WebsiteHandler) GetByDomainHeaderFooter(c *gin.Context) {
	var params website.GetByDomainHeaderFooterQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByDomainHeaderFooterQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// ProductSearchList godoc
// @Summary      Search products
// @Description  Searches for products based on provided criteria
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.ProductSearchListQuery  true  "Search parameters"
// @Success      200      {object}  utils.Result                    "List of products matching search criteria"
// @Failure      400      {object}  utils.Result                    "Validation error"
// @Failure      500      {object}  utils.Result                    "Internal server error"
// @Router       /website/product/search [get]
func (h *WebsiteHandler) ProductSearchList(c *gin.Context) {
	var params website.ProductSearchListQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.ProductSearchListQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetFiltersSortArticle godoc
// @Summary      Get article filters and sorting options
// @Description  Retrieves available filters and sorting options for articles
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetFiltersSortArticleQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                        "Article filters and sorting options"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /website/article/filters-sort [get]
func (h *WebsiteHandler) GetFiltersSortArticle(c *gin.Context) {
	var params website.GetFiltersSortArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetFiltersSortArticleQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetFiltersSortProduct godoc
// @Summary      Get product filters and sorting options
// @Description  Retrieves available filters and sorting options for products
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetFiltersSortProductQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                        "Product filters and sorting options"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /website/product/filters-sort [get]
func (h *WebsiteHandler) GetFiltersSortProduct(c *gin.Context) {
	var params website.GetFiltersSortProductQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetFiltersSortProductQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetArticlesByCategorySlug godoc
// @Summary      Get articles by category slug
// @Description  Retrieves articles belonging to a specific category identified by slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetArticlesByCategorySlugQuery  true  "Category slug and query parameters"
// @Success      200      {object}  utils.Result                            "List of articles in the category"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      404      {object}  utils.Result                            "Category not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /website/article/category [get]
func (h *WebsiteHandler) GetArticlesByCategorySlug(c *gin.Context) {
	var params website.GetArticlesByCategorySlugQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetArticlesByCategorySlugQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetProductsByCategorySlug godoc
// @Summary      Get products by category slug
// @Description  Retrieves products belonging to a specific category identified by slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  query     website.GetProductsByCategorySlugQuery  true  "Category slug and query parameters"
// @Success      200      {object}  utils.Result                            "List of products in the category"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      404      {object}  utils.Result                            "Category not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /website/product/category [get]
func (h *WebsiteHandler) GetProductsByCategorySlug(c *gin.Context) {
	var params website.GetProductsByCategorySlugQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetProductsByCategorySlugQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetSingleArticleBySlug godoc
// @Summary      Get article by slug
// @Description  Retrieves a specific article identified by its slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  body      website.GetSingleArticleBySlugQuery  true  "Article slug"
// @Success      200      {object}  utils.Result                         "Article details"
// @Failure      400      {object}  utils.Result                         "Validation error"
// @Failure      404      {object}  utils.Result                         "Article not found"
// @Failure      500      {object}  utils.Result                         "Internal server error"
// @Router       /website/article [get]
func (h *WebsiteHandler) GetSingleArticleBySlug(c *gin.Context) {
	var params website.GetSingleArticleBySlugQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.GetSingleArticleBySlugQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetSingleProductBySlug godoc
// @Summary      Get product by slug
// @Description  Retrieves a specific product identified by its slug
// @Tags         website
// @Accept       json
// @Produce      json
// @Param        request  body      website.GetSingleProductBySlugQuery  true  "Product slug"
// @Success      200      {object}  utils.Result                         "Product details"
// @Failure      400      {object}  utils.Result                         "Validation error"
// @Failure      404      {object}  utils.Result                         "Product not found"
// @Failure      500      {object}  utils.Result                         "Internal server error"
// @Router       /website/product [get]
func (h *WebsiteHandler) GetSingleProductBySlug(c *gin.Context) {
	var params website.GetSingleProductBySlugQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.GetSingleProductBySlugQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
