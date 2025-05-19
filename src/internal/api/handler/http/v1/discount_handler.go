package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/discount"
	discountusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/discount"
	"github.com/gin-gonic/gin"
)

type DiscountHandler struct {
	usecase   *discountusecase.DiscountUsecase
	validator *utils.ValidationHelper
}

func NewDiscountHandler(usc *discountusecase.DiscountUsecase) *DiscountHandler {
	return &DiscountHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var params discount.CreateDiscountCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateDiscountCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	var params discount.UpdateDiscountCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateDiscountCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	var params discount.DeleteDiscountCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteDiscountCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *DiscountHandler) GetByIdDiscount(c *gin.Context) {
	var params discount.GetByIdDiscountQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdDiscountQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *DiscountHandler) GetAllDiscount(c *gin.Context) {
	var params discount.GetAllDiscountQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllDiscountQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *DiscountHandler) AdminGetAllDiscount(c *gin.Context) {
	var params discount.AdminGetAllDiscountQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllDiscountQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
