package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	page2 "github.com/amirex128/new_site_builder/internal/application/dto/page"
	"github.com/amirex128/new_site_builder/internal/application/usecase/page"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	usecase   *pageusecase.PageUsecase
	validator *utils2.ValidationHelper
}

func NewPageHandler(usc *pageusecase.PageUsecase) *PageHandler {
	return &PageHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreatePage godoc
// @Summary      Create a new page
// @Description  Creates a new web page for a website
// @Tags         page
// @Accept       json
// @Produce      json
// @Param        request  body      page.CreatePageCommand  true  "Page information"
// @success      201      {object}  utils.Result            "Created page"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /page [post]
// @Security BearerAuth
func (h *PageHandler) CreatePage(c *gin.Context) {
	var params page2.CreatePageCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreatePageCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdatePage godoc
// @Summary      Update a page
// @Description  Updates an existing web page with new content and settings
// @Tags         page
// @Accept       json
// @Produce      json
// @Param        request  body      page.UpdatePageCommand  true  "Updated page information"
// @success      200      {object}  utils.Result            "Updated page"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Page not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /page [put]
// @Security BearerAuth
func (h *PageHandler) UpdatePage(c *gin.Context) {
	var params page2.UpdatePageCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdatePageCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeletePage godoc
// @Summary      Delete a page
// @Description  Deletes an existing web page
// @Tags         page
// @Accept       json
// @Produce      json
// @Param        request  body      page.DeletePageCommand  true  "Page ID to delete"
// @success      200      {object}  utils.Result            "Deleted page confirmation"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Page not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /page [delete]
// @Security BearerAuth
func (h *PageHandler) DeletePage(c *gin.Context) {
	var params page2.DeletePageCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeletePageCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdPage godoc
// @Summary      Get page by ID
// @Description  Retrieves a specific web page by its ID
// @Tags         page
// @Accept       json
// @Produce      json
// @Param        request  query     page.GetByIdPageQuery  true  "Page ID to retrieve"
// @success      200      {object}  utils.Result           "Page details"
// @Failure      400      {object}  utils.Result           "Validation error"
// @Failure      401      {object}  utils.Result           "unauthorized"
// @Failure      404      {object}  utils.Result           "Page not found"
// @Failure      500      {object}  utils.Result           "Internal server error"
// @Router       /page [get]
// @Security BearerAuth
func (h *PageHandler) GetByIdPage(c *gin.Context) {
	var params page2.GetByIdPageQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdPageQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllPage godoc
// @Summary      Get all pages
// @Description  Retrieves all web pages for the authenticated user
// @Tags         page
// @Accept       json
// @Produce      json
// @Param        request  query     page.GetAllPageQuery  true  "Query parameters"
// @success      200      {object}  utils.Result          "List of pages"
// @Failure      400      {object}  utils.Result          "Validation error"
// @Failure      401      {object}  utils.Result          "unauthorized"
// @Failure      500      {object}  utils.Result          "Internal server error"
// @Router       /page/all [get]
// @Security BearerAuth
func (h *PageHandler) GetAllPage(c *gin.Context) {
	var params page2.GetAllPageQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllPageQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllPage godoc
// @Summary      Admin: Get all pages
// @Description  Admin endpoint to retrieve all web pages across all websites
// @Tags         page
// @Accept       json
// @Produce      json
// @Param        request  query     page.AdminGetAllPageQuery  true  "Query parameters"
// @success      200      {object}  utils.Result               "List of all pages"
// @Failure      400      {object}  utils.Result               "Validation error"
// @Failure      401      {object}  utils.Result               "unauthorized"
// @Failure      403      {object}  utils.Result               "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result               "Internal server error"
// @Router       /page/admin/all [get]
// @Security BearerAuth
func (h *PageHandler) AdminGetAllPage(c *gin.Context) {
	var params page2.AdminGetAllPageQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllPageQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
