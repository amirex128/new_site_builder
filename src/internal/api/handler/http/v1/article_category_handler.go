package v1

import (
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

// CategoryCreate godoc
// @Summary      Create a new article category
// @Description  Creates a new category for articles with the provided information
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  body      article_category.CreateCategoryCommand  true  "Category information"
// @Success      201      {object}  resp.Result                             "Created category"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      401      {object}  resp.Result                             "Unauthorized"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /article-category [post]
// @Security     BearerAuth
func (h *ArticleCategoryHandler) CategoryCreate(c *gin.Context) {
	var params article_category.CreateCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateCategoryCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

// CategoryUpdate godoc
// @Summary      Update an article category
// @Description  Updates an existing article category with the provided information
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  body      article_category.UpdateCategoryCommand  true  "Updated category information"
// @Success      200      {object}  resp.Result                             "Updated category"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      401      {object}  resp.Result                             "Unauthorized"
// @Failure      404      {object}  resp.Result                             "Category not found"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /article-category [put]
// @Security     BearerAuth
func (h *ArticleCategoryHandler) CategoryUpdate(c *gin.Context) {
	var params article_category.UpdateCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateCategoryCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

// CategoryDelete godoc
// @Summary      Delete an article category
// @Description  Deletes an existing article category by its ID
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  body      article_category.DeleteCategoryCommand  true  "Category ID to delete"
// @Success      200      {object}  resp.Result                             "Deleted category confirmation"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      401      {object}  resp.Result                             "Unauthorized"
// @Failure      404      {object}  resp.Result                             "Category not found"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /article-category [delete]
// @Security     BearerAuth
func (h *ArticleCategoryHandler) CategoryDelete(c *gin.Context) {
	var params article_category.DeleteCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeleteCategoryCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c, result)
}

// CategoryGet godoc
// @Summary      Get article category by ID
// @Description  Retrieves a specific article category by its ID
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  query     article_category.GetByIdCategoryQuery  true  "Category ID to retrieve"
// @Success      200      {object}  resp.Result                            "Category details"
// @Failure      400      {object}  resp.Result                            "Validation error"
// @Failure      401      {object}  resp.Result                            "Unauthorized"
// @Failure      404      {object}  resp.Result                            "Category not found"
// @Failure      500      {object}  resp.Result                            "Internal server error"
// @Router       /article-category [get]
// @Security     BearerAuth
func (h *ArticleCategoryHandler) CategoryGet(c *gin.Context) {
	var params article_category.GetByIdCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdCategoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

// CategoryGetAll godoc
// @Summary      Get all article categories
// @Description  Retrieves all article categories with optional filtering
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  query     article_category.GetAllCategoryQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                           "List of categories"
// @Failure      400      {object}  resp.Result                           "Validation error"
// @Failure      401      {object}  resp.Result                           "Unauthorized"
// @Failure      500      {object}  resp.Result                           "Internal server error"
// @Router       /article-category/all [get]
// @Security     BearerAuth
func (h *ArticleCategoryHandler) CategoryGetAll(c *gin.Context) {
	var params article_category.GetAllCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllCategoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

// AdminCategoryGetAll godoc
// @Summary      Admin: Get all article categories
// @Description  Admin endpoint to retrieve all article categories with additional information
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  query     article_category.AdminGetAllCategoryQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                                "List of all categories"
// @Failure      400      {object}  resp.Result                                "Validation error"
// @Failure      401      {object}  resp.Result                                "Unauthorized"
// @Failure      403      {object}  resp.Result                                "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result                                "Internal server error"
// @Router       /article-category/admin/all [get]
// @Security     BearerAuth
func (h *ArticleCategoryHandler) AdminCategoryGetAll(c *gin.Context) {
	var params article_category.AdminGetAllCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllCategoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
