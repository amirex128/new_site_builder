package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	article_category2 "github.com/amirex128/new_site_builder/internal/application/dto/article_category"
	blogcategoryusecase "github.com/amirex128/new_site_builder/internal/application/usecase/article_category"
	"github.com/gin-gonic/gin"
)

type ArticleCategoryHandler struct {
	usecase   *blogcategoryusecase.ArticleCategoryUsecase
	validator *utils2.ValidationHelper
}

func NewBlogCategoryHandler(usc *blogcategoryusecase.ArticleCategoryUsecase) *ArticleCategoryHandler {
	return &ArticleCategoryHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CategoryCreate godoc
// @Summary      Create a new article category
// @Description  Creates a new category for articles with the provided information
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  body      article_category.CreateCategoryCommand  true  "Category information"
// @success      201      {object}  utils.Result                             "Created category"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /article-category [post]
// @Security BearerAuth
func (h *ArticleCategoryHandler) CategoryCreate(c *gin.Context) {
	var params article_category2.CreateCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateCategoryCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// CategoryUpdate godoc
// @Summary      Update an article category
// @Description  Updates an existing article category with the provided information
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  body      article_category.UpdateCategoryCommand  true  "Updated category information"
// @success      200      {object}  utils.Result                             "Updated category"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Category not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /article-category [put]
// @Security BearerAuth
func (h *ArticleCategoryHandler) CategoryUpdate(c *gin.Context) {
	var params article_category2.UpdateCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateCategoryCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// CategoryDelete godoc
// @Summary      Delete an article category
// @Description  Deletes an existing article category by its ID
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  body      article_category.DeleteCategoryCommand  true  "Category ID to delete"
// @success      200      {object}  utils.Result                             "Deleted category confirmation"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Category not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /article-category [delete]
// @Security BearerAuth
func (h *ArticleCategoryHandler) CategoryDelete(c *gin.Context) {
	var params article_category2.DeleteCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteCategoryCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// CategoryGet godoc
// @Summary      Get article category by ID
// @Description  Retrieves a specific article category by its ID
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  query     article_category.GetByIdCategoryQuery  true  "Category ID to retrieve"
// @success      200      {object}  utils.Result                            "Category details"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      404      {object}  utils.Result                            "Category not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /article-category [get]
// @Security BearerAuth
func (h *ArticleCategoryHandler) CategoryGet(c *gin.Context) {
	var params article_category2.GetByIdCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdCategoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// CategoryGetAll godoc
// @Summary      Get all article categories
// @Description  Retrieves all article categories with optional filtering
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  query     article_category.GetAllCategoryQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                           "List of categories"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /article-category/all [get]
// @Security BearerAuth
func (h *ArticleCategoryHandler) CategoryGetAll(c *gin.Context) {
	var params article_category2.GetAllCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllCategoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminCategoryGetAll godoc
// @Summary      Admin: Get all article categories
// @Description  Admin endpoint to retrieve all article categories with additional information
// @Tags         article-category
// @Accept       json
// @Produce      json
// @Param        request  query     article_category.AdminGetAllCategoryQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                                "List of all categories"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "unauthorized"
// @Failure      403      {object}  utils.Result                                "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /article-category/admin/all [get]
// @Security BearerAuth
func (h *ArticleCategoryHandler) AdminCategoryGetAll(c *gin.Context) {
	var params article_category2.AdminGetAllCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllCategoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
