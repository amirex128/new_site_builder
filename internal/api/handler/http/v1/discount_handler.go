package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	discount2 "github.com/amirex128/new_site_builder/internal/application/dto/discount"
	"github.com/amirex128/new_site_builder/internal/application/usecase/discount"
	"github.com/gin-gonic/gin"
)

type DiscountHandler struct {
	usecase   *discountusecase.DiscountUsecase
	validator *utils2.ValidationHelper
}

func NewDiscountHandler(usc *discountusecase.DiscountUsecase) *DiscountHandler {
	return &DiscountHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateDiscount godoc
// @Summary      Create a new discount
// @Description  Creates a new discount code or promotion
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  body      discount.CreateDiscountCommand  true  "Discount information"
// @success      201      {object}  utils.Result                     "Created discount"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /discount [post]
// @Security BearerAuth
func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var params discount2.CreateDiscountCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateDiscountCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateDiscount godoc
// @Summary      Update a discount
// @Description  Updates an existing discount code or promotion
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  body      discount.UpdateDiscountCommand  true  "Updated discount information"
// @success      200      {object}  utils.Result                     "Updated discount"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      404      {object}  utils.Result                     "Discount not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /discount [put]
// @Security BearerAuth
func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	var params discount2.UpdateDiscountCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateDiscountCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteDiscount godoc
// @Summary      Delete a discount
// @Description  Deletes an existing discount code or promotion
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  body      discount.DeleteDiscountCommand  true  "Discount ID to delete"
// @success      200      {object}  utils.Result                     "Deleted discount confirmation"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      404      {object}  utils.Result                     "Discount not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /discount [delete]
// @Security BearerAuth
func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	var params discount2.DeleteDiscountCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteDiscountCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdDiscount godoc
// @Summary      Get discount by ID
// @Description  Retrieves a specific discount by its ID
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  query     discount.GetByIdDiscountQuery  true  "Discount ID to retrieve"
// @success      200      {object}  utils.Result                    "Discount details"
// @Failure      400      {object}  utils.Result                    "Validation error"
// @Failure      401      {object}  utils.Result                    "unauthorized"
// @Failure      404      {object}  utils.Result                    "Discount not found"
// @Failure      500      {object}  utils.Result                    "Internal server error"
// @Router       /discount [get]
// @Security BearerAuth
func (h *DiscountHandler) GetByIdDiscount(c *gin.Context) {
	var params discount2.GetByIdDiscountQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdDiscountQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllDiscount godoc
// @Summary      Get all discounts
// @Description  Retrieves all discount codes and promotions
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  query     discount.GetAllDiscountQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                   "List of discounts"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /discount/all [get]
// @Security BearerAuth
func (h *DiscountHandler) GetAllDiscount(c *gin.Context) {
	var params discount2.GetAllDiscountQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllDiscountQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllDiscount godoc
// @Summary      Admin: Get all discounts
// @Description  Admin endpoint to retrieve all discount codes and promotions with additional information
// @Tags         discount
// @Accept       json
// @Produce      json
// @Param        request  query     discount.AdminGetAllDiscountQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                        "List of all discounts"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      401      {object}  utils.Result                        "unauthorized"
// @Failure      403      {object}  utils.Result                        "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /discount/admin/all [get]
// @Security BearerAuth
func (h *DiscountHandler) AdminGetAllDiscount(c *gin.Context) {
	var params discount2.AdminGetAllDiscountQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllDiscountQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
