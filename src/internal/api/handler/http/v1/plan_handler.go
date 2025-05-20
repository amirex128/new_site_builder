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
// @Success      201      {object}  resp.Result            "Created plan"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /plan [post]
// @Security     BearerAuth
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var params plan.CreatePlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreatePlanCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// UpdatePlan godoc
// @Summary      Update a subscription plan
// @Description  Updates an existing subscription plan with new features and pricing
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  body      plan.UpdatePlanCommand  true  "Updated plan information"
// @Success      200      {object}  resp.Result            "Updated plan"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      404      {object}  resp.Result            "Plan not found"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /plan [put]
// @Security     BearerAuth
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	var params plan.UpdatePlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdatePlanCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// DeletePlan godoc
// @Summary      Delete a subscription plan
// @Description  Deletes an existing subscription plan
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  body      plan.DeletePlanCommand  true  "Plan ID to delete"
// @Success      200      {object}  resp.Result            "Deleted plan confirmation"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      404      {object}  resp.Result            "Plan not found"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /plan [delete]
// @Security     BearerAuth
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	var params plan.DeletePlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeletePlanCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Deleted(c, result)
}

// GetByIdPlan godoc
// @Summary      Get plan by ID
// @Description  Retrieves a specific subscription plan by its ID
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  query     plan.GetByIDPlanQuery  true  "Plan ID to retrieve"
// @Success      200      {object}  resp.Result           "Plan details"
// @Failure      400      {object}  resp.Result           "Validation error"
// @Failure      401      {object}  resp.Result           "Unauthorized"
// @Failure      404      {object}  resp.Result           "Plan not found"
// @Failure      500      {object}  resp.Result           "Internal server error"
// @Router       /plan [get]
// @Security     BearerAuth
func (h *PlanHandler) GetByIdPlan(c *gin.Context) {
	var params plan.GetByIDPlanQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIDPlanQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllPlan godoc
// @Summary      Get all subscription plans
// @Description  Retrieves all available subscription plans
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  query     plan.GetAllPlanQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result          "List of plans"
// @Failure      400      {object}  resp.Result          "Validation error"
// @Failure      401      {object}  resp.Result          "Unauthorized"
// @Failure      500      {object}  resp.Result          "Internal server error"
// @Router       /plan/all [get]
// @Security     BearerAuth
func (h *PlanHandler) GetAllPlan(c *gin.Context) {
	var params plan.GetAllPlanQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllPlanQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// CalculatePlanPrice godoc
// @Summary      Calculate plan price
// @Description  Calculates the price for a subscription plan based on selected options
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        request  query     plan.CalculatePlanPriceQuery  true  "Plan calculation parameters"
// @Success      200      {object}  resp.Result                  "Calculated plan price"
// @Failure      400      {object}  resp.Result                  "Validation error"
// @Failure      401      {object}  resp.Result                  "Unauthorized"
// @Failure      500      {object}  resp.Result                  "Internal server error"
// @Router       /plan/calculate [get]
// @Security     BearerAuth
func (h *PlanHandler) CalculatePlanPrice(c *gin.Context) {
	var params plan.CalculatePlanPriceQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.CalculatePlanPriceQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
