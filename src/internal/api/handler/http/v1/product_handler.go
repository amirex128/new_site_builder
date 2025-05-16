package v1

import (
	"github.com/gin-gonic/gin"
	"go-boilerplate/src/internal/api/helper"
	producttransformer "go-boilerplate/src/internal/api/transformer/product"
	productusecase "go-boilerplate/src/internal/app/usecase/product"
	"net/http"
)

type ProductHandler struct {
	usecase *productusecase.ProductUsecase
}

func NewProductHandler(usc *productusecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		usecase: usc,
	}
}

func (h *ProductHandler) ProductList(c *gin.Context) {
	// initializer generate dto
	params := producttransformer.NewProductTransformer(c).
		SetSuperType().
		Build()

	result, err := h.usecase.ProductList(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateBaseResponseWithError(nil, false, 1, err))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateBaseResponse(result, true, 0))
}
