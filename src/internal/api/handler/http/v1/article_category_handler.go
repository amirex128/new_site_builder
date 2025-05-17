package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/article_category"
	blogcategoryusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article_category"
	"github.com/gin-gonic/gin"
)

type ArticleCategoryHandler struct {
	usecase   *blogcategoryusecase.ArticleCategoryUsecase
	validator *utils.ValidationHelper
}

func NewBlogCategoryHandler(usc *blogcategoryusecase.ArticleCategoryUsecase) *ArticleCategoryHandler {
	return &ArticleCategoryHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *ArticleCategoryHandler) CategoryCreate(c *gin.Context) {
	var params article_category.CreateCategoryCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateCategoryCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *ArticleCategoryHandler) CategoryUpdate(c *gin.Context) {
	var params article_category.UpdateCategoryCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateCategoryCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *ArticleCategoryHandler) CategoryDelete(c *gin.Context) {
	var params article_category.DeleteCategoryCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteCategoryCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *ArticleCategoryHandler) CategoryGet(c *gin.Context) {
	var params article_category.GetByIdCategoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdCategoryQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ArticleCategoryHandler) CategoryGetAll(c *gin.Context) {
	var params article_category.GetAllCategoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllCategoryQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ArticleCategoryHandler) AdminCategoryGetAll(c *gin.Context) {
	var params article_category.AdminGetAllCategoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllCategoryQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
