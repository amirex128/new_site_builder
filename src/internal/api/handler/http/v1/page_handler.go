package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/page"
	pageusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/page"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	usecase   *pageusecase.PageUsecase
	validator *utils.ValidationHelper
}

func NewPageHandler(usc *pageusecase.PageUsecase) *PageHandler {
	return &PageHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *PageHandler) CreatePage(c *gin.Context) {
	var params page.CreatePageCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreatePageCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	var params page.UpdatePageCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdatePageCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	var params page.DeletePageCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeletePageCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *PageHandler) GetByIdPage(c *gin.Context) {
	var params page.GetByIdPageQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdPageQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *PageHandler) GetAllPage(c *gin.Context) {
	var params page.GetAllPageQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllPageQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *PageHandler) AdminGetAllPage(c *gin.Context) {
	var params page.AdminGetAllPageQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllPageQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
