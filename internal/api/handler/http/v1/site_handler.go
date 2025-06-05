package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	site2 "github.com/amirex128/new_site_builder/internal/application/dto/site"
	"github.com/amirex128/new_site_builder/internal/application/usecase/site"
	"github.com/gin-gonic/gin"
)

type SiteHandler struct {
	usecase   *siteusecase.SiteUsecase
	validator *utils2.ValidationHelper
}

func NewSiteHandler(usc *siteusecase.SiteUsecase) *SiteHandler {
	return &SiteHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateSite godoc
// @Summary      Create a new site
// @Description  Creates a new site for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.CreateSiteCommand  true  "Site information"
// @success      201      {object}  utils.Result            "Created site"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /site [post]
// @Security BearerAuth
func (h *SiteHandler) CreateSite(c *gin.Context) {
	var params site2.CreateSiteCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateSiteCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateSite godoc
// @Summary      Update an existing site
// @Description  Updates an existing site for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.UpdateSiteCommand  true  "Updated site information"
// @success      200      {object}  utils.Result            "Updated site"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Site not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /site [put]
// @Security BearerAuth
func (h *SiteHandler) UpdateSite(c *gin.Context) {
	var params site2.UpdateSiteCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateSiteCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteSite godoc
// @Summary      Delete a site
// @Description  Deletes a site for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  body      site.DeleteSiteCommand  true  "Site ID to delete"
// @success      200      {object}  utils.Result            "Deleted site confirmation"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Site not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /site [delete]
// @Security BearerAuth
func (h *SiteHandler) DeleteSite(c *gin.Context) {
	var params site2.DeleteSiteCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteSiteCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdSite godoc
// @Summary      Get site by ID
// @Description  Retrieves a specific site by its ID for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  query     site.GetByIdSiteQuery  true  "Site ID to retrieve"
// @success      200      {object}  utils.Result           "Site details"
// @Failure      400      {object}  utils.Result           "Validation error"
// @Failure      401      {object}  utils.Result           "unauthorized"
// @Failure      404      {object}  utils.Result           "Site not found"
// @Failure      500      {object}  utils.Result           "Internal server error"
// @Router       /site [get]
// @Security BearerAuth
func (h *SiteHandler) GetByIdSite(c *gin.Context) {
	var params site2.GetByIdSiteQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdSiteQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllSite godoc
// @Summary      Get all sites
// @Description  Retrieves all sites for the authenticated user
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  query     site.GetAllSiteQuery  true  "Query parameters"
// @success      200      {object}  utils.Result          "List of sites"
// @Failure      400      {object}  utils.Result          "Validation error"
// @Failure      401      {object}  utils.Result          "unauthorized"
// @Failure      500      {object}  utils.Result          "Internal server error"
// @Router       /site/all [get]
// @Security BearerAuth
func (h *SiteHandler) GetAllSite(c *gin.Context) {
	var params site2.GetAllSiteQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllSiteQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllSite godoc
// @Summary      Admin: Get all sites
// @Description  Admin endpoint to retrieve all sites in the system
// @Tags         site
// @Accept       json
// @Produce      json
// @Param        request  query     site.AdminGetAllSiteQuery  true  "Query parameters"
// @success      200      {object}  utils.Result               "List of all sites"
// @Failure      400      {object}  utils.Result               "Validation error"
// @Failure      401      {object}  utils.Result               "unauthorized"
// @Failure      403      {object}  utils.Result               "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result               "Internal server error"
// @Router       /site/admin/all [get]
// @Security BearerAuth
func (h *SiteHandler) AdminGetAllSite(c *gin.Context) {
	var params site2.AdminGetAllSiteQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllSiteQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
