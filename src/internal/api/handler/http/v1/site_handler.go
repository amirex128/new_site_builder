package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/site"
	siteusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/site"
	"github.com/gin-gonic/gin"
)

type SiteHandler struct {
	usecase   *siteusecase.SiteUsecase
	validator *utils.ValidationHelper
}

func NewSiteHandler(usc *siteusecase.SiteUsecase) *SiteHandler {
	return &SiteHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *SiteHandler) CreateSite(c *gin.Context) {
	var params site.CreateSiteCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateSiteCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *SiteHandler) UpdateSite(c *gin.Context) {
	var params site.UpdateSiteCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateSiteCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	var params site.DeleteSiteCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteSiteCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *SiteHandler) GetByIdSite(c *gin.Context) {
	var params site.GetByIdSiteQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdSiteQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *SiteHandler) GetAllSite(c *gin.Context) {
	var params site.GetAllSiteQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllSiteQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *SiteHandler) AdminGetAllSite(c *gin.Context) {
	var params site.AdminGetAllSiteQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllSiteQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
