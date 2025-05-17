package v1

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *SiteHandler) UpdateSite(c *gin.Context) {
	var params site.UpdateSiteCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateSiteCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *SiteHandler) DeleteSite(c *gin.Context) {
	var params site.DeleteSiteCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteSiteCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *SiteHandler) GetByIdSite(c *gin.Context) {
	var params site.GetByIdSiteQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdSiteQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *SiteHandler) GetAllSite(c *gin.Context) {
	var params site.GetAllSiteQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllSiteQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *SiteHandler) AdminGetAllSite(c *gin.Context) {
	var params site.AdminGetAllSiteQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllSiteQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
