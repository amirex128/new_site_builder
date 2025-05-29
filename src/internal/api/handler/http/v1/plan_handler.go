package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/plan"
	planusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/plan"
	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	usecase   *planusecase.PlanUsecase
	validator *utils.ValidationHelper
}

func NewPlanHandler(usc *planusecase.PlanUsecase) *PlanHandler {
	return &PlanHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// CreatePlan godoc
// @Summary      Create a new subscription plan
// @Description  Creates a new subscription plan with specified features and pricing
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  body      plan.CreatePlanCommand  true  "Plan information"
// @success      201      {object}  utils.Result            "Created plan"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /plan [post]
// @Security BearerAuth
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var params plan.CreatePlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreatePlanCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// UpdatePlan godoc
// @Summary      Update a subscription plan
// @Description  Updates an existing subscription plan with new features and pricing
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  body      plan.UpdatePlanCommand  true  "Updated plan information"
// @success      200      {object}  utils.Result            "Updated plan"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Plan not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /plan [put]
// @Security BearerAuth
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	var params plan.UpdatePlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdatePlanCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// DeletePlan godoc
// @Summary      Delete a subscription plan
// @Description  Deletes an existing subscription plan
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  body      plan.DeletePlanCommand  true  "Plan ID to delete"
// @success      200      {object}  utils.Result            "Deleted plan confirmation"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Plan not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /plan [delete]
// @Security BearerAuth
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	var params plan.DeletePlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeletePlanCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetByIdPlan godoc
// @Summary      Get plan by ID
// @Description  Retrieves a specific subscription plan by its ID
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  query     plan.GetByIDPlanQuery  true  "Plan ID to retrieve"
// @success      200      {object}  utils.Result           "Plan details"
// @Failure      400      {object}  utils.Result           "Validation error"
// @Failure      401      {object}  utils.Result           "unauthorized"
// @Failure      404      {object}  utils.Result           "Plan not found"
// @Failure      500      {object}  utils.Result           "Internal server error"
// @Router       /plan [get]
// @Security BearerAuth
func (h *PlanHandler) GetByIdPlan(c *gin.Context) {
	var params plan.GetByIDPlanQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIDPlanQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetAllPlan godoc
// @Summary      Get all subscription plans
// @Description  Retrieves all available subscription plans
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  query     plan.GetAllPlanQuery  true  "Query parameters"
// @success      200      {object}  utils.Result          "List of plans"
// @Failure      400      {object}  utils.Result          "Validation error"
// @Failure      401      {object}  utils.Result          "unauthorized"
// @Failure      500      {object}  utils.Result          "Internal server error"
// @Router       /plan/all [get]
// @Security BearerAuth
func (h *PlanHandler) GetAllPlan(c *gin.Context) {
	var params plan.GetAllPlanQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllPlanQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// CalculatePlanPrice godoc
// @Summary      Calculate plan price
// @Description  Calculates the price for a subscription plan based on selected options
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  query     plan.CalculatePlanPriceQuery  true  "Plan calculation parameters"
// @success      200      {object}  utils.Result                  "Calculated plan price"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "unauthorized"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /plan/calculate [get]
// @Security BearerAuth
func (h *PlanHandler) CalculatePlanPrice(c *gin.Context) {
	var params plan.CalculatePlanPriceQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CalculatePlanPriceQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}
