package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	defaulttheme2 "github.com/amirex128/new_site_builder/internal/application/dto/defaulttheme"
	"github.com/amirex128/new_site_builder/internal/application/usecase/default_theme"
	"github.com/gin-gonic/gin"
)

type DefaultThemeHandler struct {
	usecase   *defaultthemeusecase.DefaultThemeUsecase
	validator *utils2.ValidationHelper
}

func NewDefaultThemeHandler(usc *defaultthemeusecase.DefaultThemeUsecase) *DefaultThemeHandler {
	return &DefaultThemeHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateDefaultTheme godoc
// @Summary      Create a new default theme
// @Description  Creates a new default theme template for websites
// @Tags         default-theme
// @Accept       json
// @Produce      json
// @Param        request  body      defaulttheme.CreateDefaultThemeCommand  true  "Theme information"
// @success      201      {object}  utils.Result                             "Created theme"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /default-theme [post]
// @Security BearerAuth
func (h *DefaultThemeHandler) CreateDefaultTheme(c *gin.Context) {
	var params defaulttheme2.CreateDefaultThemeCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateDefaultThemeCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateDefaultTheme godoc
// @Summary      Update a default theme
// @Description  Updates an existing default theme template
// @Tags         default-theme
// @Accept       json
// @Produce      json
// @Param        request  body      defaulttheme.UpdateDefaultThemeCommand  true  "Updated theme information"
// @success      200      {object}  utils.Result                             "Updated theme"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Theme not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /default-theme [put]
// @Security BearerAuth
func (h *DefaultThemeHandler) UpdateDefaultTheme(c *gin.Context) {
	var params defaulttheme2.UpdateDefaultThemeCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateDefaultThemeCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteDefaultTheme godoc
// @Summary      Delete a default theme
// @Description  Deletes an existing default theme template
// @Tags         default-theme
// @Accept       json
// @Produce      json
// @Param        request  body      defaulttheme.DeleteDefaultThemeCommand  true  "Theme ID to delete"
// @success      200      {object}  utils.Result                             "Deleted theme confirmation"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Theme not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /default-theme [delete]
// @Security BearerAuth
func (h *DefaultThemeHandler) DeleteDefaultTheme(c *gin.Context) {
	var params defaulttheme2.DeleteDefaultThemeCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteDefaultThemeCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdDefaultTheme godoc
// @Summary      Get default theme by ID
// @Description  Retrieves a specific default theme template by its ID
// @Tags         default-theme
// @Accept       json
// @Produce      json
// @Param        request  query     defaulttheme.GetByIdDefaultThemeQuery  true  "Theme ID to retrieve"
// @success      200      {object}  utils.Result                            "Theme details"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      404      {object}  utils.Result                            "Theme not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /default-theme [get]
// @Security BearerAuth
func (h *DefaultThemeHandler) GetByIdDefaultTheme(c *gin.Context) {
	var params defaulttheme2.GetByIdDefaultThemeQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdDefaultThemeQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllDefaultTheme godoc
// @Summary      Get all default themes
// @Description  Retrieves all default theme templates
// @Tags         default-theme
// @Accept       json
// @Produce      json
// @Param        request  query     defaulttheme.GetAllDefaultThemeQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                           "List of themes"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /default-theme/all [get]
// @Security BearerAuth
func (h *DefaultThemeHandler) GetAllDefaultTheme(c *gin.Context) {
	var params defaulttheme2.GetAllDefaultThemeQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllDefaultThemeQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
