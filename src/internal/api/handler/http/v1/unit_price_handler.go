package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/unit_price"
	unitpriceusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/unit_price"
	"github.com/gin-gonic/gin"
)

type UnitPriceHandler struct {
	usecase   *unitpriceusecase.UnitPriceUsecase
	validator *utils.ValidationHelper
}

func NewUnitPriceHandler(usc *unitpriceusecase.UnitPriceUsecase) *UnitPriceHandler {
	return &UnitPriceHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *UnitPriceHandler) UpdateUnitPrice(c *gin.Context) {
	var params unit_price.UpdateUnitPriceCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpdateUnitPriceCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *UnitPriceHandler) CalculateUnitPrice(c *gin.Context) {
	var params unit_price.CalculateUnitPriceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).CalculateUnitPriceQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *UnitPriceHandler) GetAllUnitPrice(c *gin.Context) {
	var params unit_price.GetAllUnitPriceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllUnitPriceQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
