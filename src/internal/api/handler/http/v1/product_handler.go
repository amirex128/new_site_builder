package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	productusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase   *productusecase.ProductUsecase
	validator *utils.ValidationHelper
}

func NewProductHandler(usc *productusecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var params product.CreateProductCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).CreateProductCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var params product.UpdateProductCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpdateProductCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var params product.DeleteProductCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).DeleteProductCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *ProductHandler) GetByIdProduct(c *gin.Context) {
	var params product.GetByIdProductQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetByIdProductQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ProductHandler) GetAllProduct(c *gin.Context) {
	var params product.GetAllProductQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllProductQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ProductHandler) GetByFiltersSortProduct(c *gin.Context) {
	var params product.GetByFiltersSortProductQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetByFiltersSortProductQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ProductHandler) AdminGetAllProduct(c *gin.Context) {
	var params product.AdminGetAllProductQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).AdminGetAllProductQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
