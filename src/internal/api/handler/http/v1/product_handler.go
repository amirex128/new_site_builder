package v1

import (
	producttransformer "github.com/amirex128/new_site_builder/src/internal/api/transformer/product"
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	productusecase "github.com/amirex128/new_site_builder/src/internal/app/usecase/product"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, resp.Created().Succeeded)
		return
	}

	c.JSON(http.StatusOK, utils.GenerateBaseResponse(result, true, 0))
}
