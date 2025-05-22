package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/article"
	articleusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	usecase   *articleusecase.ArticleUsecase
	validator *utils.ValidationHelper
}

func NewArticleHandler(usc *articleusecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// ArticleCreate godoc
// @Summary      Create a new article
// @Description  Creates a new article with the provided information
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.CreateArticleCommand  true  "Article information"
// @Success      201      {object}  utils.Result                   "Created article"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /article [post]
// @Security BearerAuth
func (h *ArticleHandler) ArticleCreate(c *gin.Context) {
	var params article.CreateArticleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateArticleCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// ArticleUpdate godoc
// @Summary      Update an article
// @Description  Updates an existing article with the provided information
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.UpdateArticleCommand  true  "Updated article information"
// @Success      200      {object}  utils.Result                   "Updated article"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      404      {object}  utils.Result                   "Article not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /article [put]
// @Security BearerAuth
func (h *ArticleHandler) ArticleUpdate(c *gin.Context) {
	var params article.UpdateArticleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateArticleCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// ArticleDelete godoc
// @Summary      Delete an article
// @Description  Deletes an existing article by its ID
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.DeleteArticleCommand  true  "Article ID to delete"
// @Success      200      {object}  utils.Result                   "Deleted article confirmation"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      404      {object}  utils.Result                   "Article not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /article [delete]
// @Security BearerAuth
func (h *ArticleHandler) ArticleDelete(c *gin.Context) {
	var params article.DeleteArticleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeleteArticleCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// ArticleGet godoc
// @Summary      Get article by ID
// @Description  Retrieves a specific article by its ID
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  query     article.GetByIdArticleQuery  true  "Article ID to retrieve"
// @Success      200      {object}  utils.Result                  "Article details"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "Unauthorized"
// @Failure      404      {object}  utils.Result                  "Article not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /article [get]
// @Security BearerAuth
func (h *ArticleHandler) ArticleGet(c *gin.Context) {
	var params article.GetByIdArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdArticleQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// ArticleGetAll godoc
// @Summary      Get all articles
// @Description  Retrieves all articles with optional filtering
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  query     article.GetAllArticleQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                 "List of articles"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "Unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /article/all [get]
// @Security BearerAuth
func (h *ArticleHandler) ArticleGetAll(c *gin.Context) {
	var params article.GetAllArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllArticleQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// ArticleGetByFiltersSort godoc
// @Summary      Get articles by filters and sorting
// @Description  Retrieves articles based on specified filters and sorting criteria
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.GetByFiltersSortArticleQuery  true  "Filter and sort parameters"
// @Success      200      {object}  utils.Result                           "Filtered and sorted articles"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "Unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /article/filters-sort [post]
// @Security BearerAuth
func (h *ArticleHandler) ArticleGetByFiltersSort(c *gin.Context) {
	var params article.GetByFiltersSortArticleQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.GetByFiltersSortArticleQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// AdminArticleGetAll godoc
// @Summary      Admin: Get all articles
// @Description  Admin endpoint to retrieve all articles with additional information
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  query     article.AdminGetAllArticleQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                      "List of all articles"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "Unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /article/admin/all [get]
// @Security BearerAuth
func (h *ArticleHandler) AdminArticleGetAll(c *gin.Context) {
	var params article.AdminGetAllArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllArticleQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}
