package v1

import (
	"github.com/gin-gonic/gin"
	"go-boilerplate/src/internal/api/helper"
	vendtransformer "go-boilerplate/src/internal/api/transformer/vend"
	vendorusecase "go-boilerplate/src/internal/app/usecase/vend"
	"net/http"
)

type VendorHandler struct {
	usecase *vendorusecase.VendorUsecase
}

func NewVendorHandler(usc *vendorusecase.VendorUsecase) *VendorHandler {
	return &VendorHandler{
		usecase: usc,
	}
}

func (h *VendorHandler) VendorList(c *gin.Context) {
	// initializer generate dto
	dto := vendtransformer.NewVendorTransformer(c).
		SetSuperType().
		Build()

	result, err := h.usecase.VendorList(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateBaseResponseWithError(nil, false, 1, err))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateBaseResponse(result, true, 0))
}
