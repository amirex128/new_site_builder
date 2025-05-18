package v1

import (
	"net/http"

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

	result, err := h.usecase.UpdateUnitPriceCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *UnitPriceHandler) CalculateUnitPrice(c *gin.Context) {
	var params unit_price.CalculateUnitPriceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CalculateUnitPriceQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *UnitPriceHandler) GetAllUnitPrice(c *gin.Context) {
	var params unit_price.GetAllUnitPriceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllUnitPriceQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
