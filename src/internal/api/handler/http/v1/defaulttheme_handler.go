package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/defaulttheme"
	defaultthemeusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/defaulttheme"
	"github.com/gin-gonic/gin"
)

type DefaultThemeHandler struct {
	usecase   *defaultthemeusecase.DefaultThemeUsecase
	validator *utils.ValidationHelper
}

func NewDefaultThemeHandler(usc *defaultthemeusecase.DefaultThemeUsecase) *DefaultThemeHandler {
	return &DefaultThemeHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *DefaultThemeHandler) CreateDefaultTheme(c *gin.Context) {
	var params defaulttheme.CreateDefaultThemeCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateDefaultThemeCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *DefaultThemeHandler) UpdateDefaultTheme(c *gin.Context) {
	var params defaulttheme.UpdateDefaultThemeCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateDefaultThemeCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *DefaultThemeHandler) DeleteDefaultTheme(c *gin.Context) {
	var params defaulttheme.DeleteDefaultThemeCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteDefaultThemeCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *DefaultThemeHandler) GetByIdDefaultTheme(c *gin.Context) {
	var params defaulttheme.GetByIdDefaultThemeQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdDefaultThemeQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *DefaultThemeHandler) GetAllDefaultTheme(c *gin.Context) {
	var params defaulttheme.GetAllDefaultThemeQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllDefaultThemeQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
