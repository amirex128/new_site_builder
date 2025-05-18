package v1

import (
	"net/http"

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

	result, err := h.usecase.CreatePlanCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	var params plan.UpdatePlanCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdatePlanCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	var params plan.DeletePlanCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeletePlanCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *PlanHandler) GetByIdPlan(c *gin.Context) {
	var params plan.GetByIDPlanQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIDPlanQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *PlanHandler) GetAllPlan(c *gin.Context) {
	var params plan.GetAllPlanQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllPlanQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *PlanHandler) CalculatePlanPrice(c *gin.Context) {
	var params plan.CalculatePlanPriceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CalculatePlanPriceQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
