package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
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

// UpdateUnitPrice godoc
// @Summary      Update unit price
// @Description  Updates an existing unit price configuration
// @Tags         unit-price
// @Accept       json
// @Produce      json
// @Param        request  body      unit_price.UpdateUnitPriceCommand  true  "Unit price update information"
// @Success      200      {object}  resp.Result                        "Updated unit price"
// @Failure      400      {object}  resp.Result                        "Validation error"
// @Failure      401      {object}  resp.Result                        "Unauthorized"
// @Failure      404      {object}  resp.Result                        "Unit price not found"
// @Failure      500      {object}  resp.Result                        "Internal server error"
// @Router       /unit-price [put]
// @Security     BearerAuth
func (h *UnitPriceHandler) UpdateUnitPrice(c *gin.Context) {
	var params unit_price.UpdateUnitPriceCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateUnitPriceCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// CalculateUnitPrice godoc
// @Summary      Calculate unit price
// @Description  Calculates the unit price based on provided parameters
// @Tags         unit-price
// @Accept       json
// @Produce      json
// @Param        request  query     unit_price.CalculateUnitPriceQuery  true  "Parameters for price calculation"
// @Success      200      {object}  resp.Result                         "Calculated unit price"
// @Failure      400      {object}  resp.Result                         "Validation error"
// @Failure      401      {object}  resp.Result                         "Unauthorized"
// @Failure      500      {object}  resp.Result                         "Internal server error"
// @Router       /unit-price/calculate [get]
// @Security     BearerAuth
func (h *UnitPriceHandler) CalculateUnitPrice(c *gin.Context) {
	var params unit_price.CalculateUnitPriceQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.CalculateUnitPriceQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllUnitPrice godoc
// @Summary      Get all unit prices
// @Description  Retrieves all unit price configurations
// @Tags         unit-price
// @Accept       json
// @Produce      json
// @Param        request  query     unit_price.GetAllUnitPriceQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                      "List of unit prices"
// @Failure      400      {object}  resp.Result                      "Validation error"
// @Failure      401      {object}  resp.Result                      "Unauthorized"
// @Failure      500      {object}  resp.Result                      "Internal server error"
// @Router       /unit-price/all [get]
// @Security     BearerAuth
func (h *UnitPriceHandler) GetAllUnitPrice(c *gin.Context) {
	var params unit_price.GetAllUnitPriceQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllUnitPriceQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
