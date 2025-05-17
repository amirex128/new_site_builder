package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
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

func (h *WebsiteHandler) GetByDomainPage(c *gin.Context) {
	var params website.GetByDomainPageQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByDomainPageQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetByDomainHeaderFooter(c *gin.Context) {
	var params website.GetByDomainHeaderFooterQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByDomainHeaderFooterQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) ProductSearchList(c *gin.Context) {
	var params website.ProductSearchListQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.ProductSearchListQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetFiltersSortArticle(c *gin.Context) {
	var params website.GetFiltersSortArticleQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetFiltersSortArticleQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetFiltersSortProduct(c *gin.Context) {
	var params website.GetFiltersSortProductQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetFiltersSortProductQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetArticlesByCategorySlug(c *gin.Context) {
	var params website.GetArticlesByCategorySlugQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetArticlesByCategorySlugQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetProductsByCategorySlug(c *gin.Context) {
	var params website.GetProductsByCategorySlugQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetProductsByCategorySlugQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetSingleArticleBySlug(c *gin.Context) {
	var params website.GetSingleArticleBySlugQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetSingleArticleBySlugQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *WebsiteHandler) GetSingleProductBySlug(c *gin.Context) {
	var params website.GetSingleProductBySlugQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetSingleProductBySlugQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
