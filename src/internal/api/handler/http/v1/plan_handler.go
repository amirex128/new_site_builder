package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
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

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var params plan.CreatePlanCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).CreatePlanCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	var params plan.UpdatePlanCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpdatePlanCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	var params plan.DeletePlanCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).DeletePlanCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *PlanHandler) GetByIdPlan(c *gin.Context) {
	var params plan.GetByIDPlanQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetByIDPlanQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *PlanHandler) GetAllPlan(c *gin.Context) {
	var params plan.GetAllPlanQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllPlanQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *PlanHandler) CalculatePlanPrice(c *gin.Context) {
	var params plan.CalculatePlanPriceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).CalculatePlanPriceQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
