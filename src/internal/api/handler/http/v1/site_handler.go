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

// CreateSite godoc
// @Summary      Create a new site
// @Description  Creates a new site for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.CreateSiteCommand  true  "Site information"
// @Success      201      {object}  resp.Result            "Created site"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /site [post]
// @Security     BearerAuth
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

// UpdateSite godoc
// @Summary      Update an existing site
// @Description  Updates an existing site for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.UpdateSiteCommand  true  "Updated site information"
// @Success      200      {object}  resp.Result            "Updated site"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      404      {object}  resp.Result            "Site not found"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /site [put]
// @Security     BearerAuth
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

// DeleteSite godoc
// @Summary      Delete a site
// @Description  Deletes a site for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.DeleteSiteCommand  true  "Site ID to delete"
// @Success      200      {object}  resp.Result            "Deleted site confirmation"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      404      {object}  resp.Result            "Site not found"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /site [delete]
// @Security     BearerAuth
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

	resp.Deleted(c, result)
}

// GetByIdSite godoc
// @Summary      Get site by ID
// @Description  Retrieves a specific site by its ID for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.GetByIdSiteQuery  true  "Site ID to retrieve"
// @Success      200      {object}  resp.Result           "Site details"
// @Failure      400      {object}  resp.Result           "Validation error"
// @Failure      401      {object}  resp.Result           "Unauthorized"
// @Failure      404      {object}  resp.Result           "Site not found"
// @Failure      500      {object}  resp.Result           "Internal server error"
// @Router       /site [get]
// @Security     BearerAuth
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

// GetAllSite godoc
// @Summary      Get all sites
// @Description  Retrieves all sites for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.GetAllSiteQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result          "List of sites"
// @Failure      400      {object}  resp.Result          "Validation error"
// @Failure      401      {object}  resp.Result          "Unauthorized"
// @Failure      500      {object}  resp.Result          "Internal server error"
// @Router       /site/all [get]
// @Security     BearerAuth
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

// AdminGetAllSite godoc
// @Summary      Admin: Get all sites
// @Description  Admin endpoint to retrieve all sites in the system
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.AdminGetAllSiteQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result               "List of all sites"
// @Failure      400      {object}  resp.Result               "Validation error"
// @Failure      401      {object}  resp.Result               "Unauthorized"
// @Failure      403      {object}  resp.Result               "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result               "Internal server error"
// @Router       /site/admin/all [get]
// @Security     BearerAuth
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
