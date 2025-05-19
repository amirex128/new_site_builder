package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/basket"
	basketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/basket"
	"github.com/gin-gonic/gin"
)

type BasketHandler struct {
	usecase   *basketusecase.BasketUsecase
	validator *utils.ValidationHelper
}

func NewBasketHandler(usc *basketusecase.BasketUsecase) *BasketHandler {
	return &BasketHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *BasketHandler) UpdateBasket(c *gin.Context) {
	var params basket.UpdateBasketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateBasketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *BasketHandler) GetBasket(c *gin.Context) {
	var params basket.GetBasketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetBasketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *BasketHandler) GetAllBasketUser(c *gin.Context) {
	var params basket.GetAllBasketUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllBasketUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *BasketHandler) AdminGetAllBasketUser(c *gin.Context) {
	var params basket.AdminGetAllBasketUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllBasketUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
