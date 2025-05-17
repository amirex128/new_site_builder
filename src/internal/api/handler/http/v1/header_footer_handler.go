package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/header_footer"
	headerfooterusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/header_footer"
	"github.com/gin-gonic/gin"
)

type HeaderFooterHandler struct {
	usecase   *headerfooterusecase.HeaderFooterUsecase
	validator *utils.ValidationHelper
}

func NewHeaderFooterHandler(usc *headerfooterusecase.HeaderFooterUsecase) *HeaderFooterHandler {
	return &HeaderFooterHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *HeaderFooterHandler) CreateHeaderFooter(c *gin.Context) {
	var params header_footer.CreateHeaderFooterCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateHeaderFooterCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *HeaderFooterHandler) UpdateHeaderFooter(c *gin.Context) {
	var params header_footer.UpdateHeaderFooterCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateHeaderFooterCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *HeaderFooterHandler) DeleteHeaderFooter(c *gin.Context) {
	var params header_footer.DeleteHeaderFooterCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteHeaderFooterCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *HeaderFooterHandler) GetByIdHeaderFooter(c *gin.Context) {
	var params header_footer.GetByIdHeaderFooterQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdHeaderFooterQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *HeaderFooterHandler) GetAllHeaderFooter(c *gin.Context) {
	var params header_footer.GetAllHeaderFooterQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllHeaderFooterQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *HeaderFooterHandler) AdminGetAllHeaderFooter(c *gin.Context) {
	var params header_footer.AdminGetAllHeaderFooterQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllHeaderFooterQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
