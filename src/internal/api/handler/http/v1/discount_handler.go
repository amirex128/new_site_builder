package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
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

// CreateDiscount godoc
// @Summary      Create a new discount
// @Description  Creates a new discount code or promotion
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  body      discount.CreateDiscountCommand  true  "Discount information"
// @Success      201      {object}  resp.Result                     "Created discount"
// @Failure      400      {object}  resp.Result                     "Validation error"
// @Failure      401      {object}  resp.Result                     "Unauthorized"
// @Failure      500      {object}  resp.Result                     "Internal server error"
// @Router       /discount [post]
// @Security     BearerAuth
func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var params discount.CreateDiscountCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateDiscountCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// UpdateDiscount godoc
// @Summary      Update a discount
// @Description  Updates an existing discount code or promotion
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  body      discount.UpdateDiscountCommand  true  "Updated discount information"
// @Success      200      {object}  resp.Result                     "Updated discount"
// @Failure      400      {object}  resp.Result                     "Validation error"
// @Failure      401      {object}  resp.Result                     "Unauthorized"
// @Failure      404      {object}  resp.Result                     "Discount not found"
// @Failure      500      {object}  resp.Result                     "Internal server error"
// @Router       /discount [put]
// @Security     BearerAuth
func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	var params discount.UpdateDiscountCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateDiscountCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// DeleteDiscount godoc
// @Summary      Delete a discount
// @Description  Deletes an existing discount code or promotion
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  body      discount.DeleteDiscountCommand  true  "Discount ID to delete"
// @Success      200      {object}  resp.Result                     "Deleted discount confirmation"
// @Failure      400      {object}  resp.Result                     "Validation error"
// @Failure      401      {object}  resp.Result                     "Unauthorized"
// @Failure      404      {object}  resp.Result                     "Discount not found"
// @Failure      500      {object}  resp.Result                     "Internal server error"
// @Router       /discount [delete]
// @Security     BearerAuth
func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	var params discount.DeleteDiscountCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeleteDiscountCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Deleted(c, result)
}

// GetByIdDiscount godoc
// @Summary      Get discount by ID
// @Description  Retrieves a specific discount by its ID
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  query     discount.GetByIdDiscountQuery  true  "Discount ID to retrieve"
// @Success      200      {object}  resp.Result                    "Discount details"
// @Failure      400      {object}  resp.Result                    "Validation error"
// @Failure      401      {object}  resp.Result                    "Unauthorized"
// @Failure      404      {object}  resp.Result                    "Discount not found"
// @Failure      500      {object}  resp.Result                    "Internal server error"
// @Router       /discount [get]
// @Security     BearerAuth
func (h *DiscountHandler) GetByIdDiscount(c *gin.Context) {
	var params discount.GetByIdDiscountQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdDiscountQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllDiscount godoc
// @Summary      Get all discounts
// @Description  Retrieves all discount codes and promotions
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  query     discount.GetAllDiscountQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                   "List of discounts"
// @Failure      400      {object}  resp.Result                   "Validation error"
// @Failure      401      {object}  resp.Result                   "Unauthorized"
// @Failure      500      {object}  resp.Result                   "Internal server error"
// @Router       /discount/all [get]
// @Security     BearerAuth
func (h *DiscountHandler) GetAllDiscount(c *gin.Context) {
	var params discount.GetAllDiscountQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllDiscountQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// AdminGetAllDiscount godoc
// @Summary      Admin: Get all discounts
// @Description  Admin endpoint to retrieve all discount codes and promotions with additional information
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  query     discount.AdminGetAllDiscountQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                        "List of all discounts"
// @Failure      400      {object}  resp.Result                        "Validation error"
// @Failure      401      {object}  resp.Result                        "Unauthorized"
// @Failure      403      {object}  resp.Result                        "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result                        "Internal server error"
// @Router       /discount/admin/all [get]
// @Security     BearerAuth
func (h *DiscountHandler) AdminGetAllDiscount(c *gin.Context) {
	var params discount.AdminGetAllDiscountQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllDiscountQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
