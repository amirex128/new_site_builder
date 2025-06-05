package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	article2 "github.com/amirex128/new_site_builder/internal/application/dto/article"
	"github.com/amirex128/new_site_builder/internal/application/usecase/article"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	usecase   *articleusecase.ArticleUsecase
	validator *utils2.ValidationHelper
}

func NewArticleHandler(usc *articleusecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// ArticleCreate godoc
// @Summary      Create a new article
// @Description  Creates a new article with the provided information
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.CreateArticleCommand  true  "Article information"
// @success      201      {object}  utils.Result                   "Created article"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /article [post]
// @Security BearerAuth
func (h *ArticleHandler) ArticleCreate(c *gin.Context) {
	var params article2.CreateArticleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateArticleCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ArticleUpdate godoc
// @Summary      Update an article
// @Description  Updates an existing article with the provided information
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.UpdateArticleCommand  true  "Updated article information"
// @success      200      {object}  utils.Result                   "Updated article"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      404      {object}  utils.Result                   "Article not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /article [put]
// @Security BearerAuth
func (h *ArticleHandler) ArticleUpdate(c *gin.Context) {
	var params article2.UpdateArticleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateArticleCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ArticleDelete godoc
// @Summary      Delete an article
// @Description  Deletes an existing article by its ID
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.DeleteArticleCommand  true  "Article ID to delete"
// @success      200      {object}  utils.Result                   "Deleted article confirmation"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      404      {object}  utils.Result                   "Article not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /article [delete]
// @Security BearerAuth
func (h *ArticleHandler) ArticleDelete(c *gin.Context) {
	var params article2.DeleteArticleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteArticleCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ArticleGet godoc
// @Summary      Get article by ID
// @Description  Retrieves a specific article by its ID
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  query     article.GetByIdArticleQuery  true  "Article ID to retrieve"
// @success      200      {object}  utils.Result                  "Article details"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "unauthorized"
// @Failure      404      {object}  utils.Result                  "Article not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /article [get]
// @Security BearerAuth
func (h *ArticleHandler) ArticleGet(c *gin.Context) {
	var params article2.GetByIdArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdArticleQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ArticleGetAll godoc
// @Summary      Get all articles
// @Description  Retrieves all articles with optional filtering
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  query     article.GetAllArticleQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                 "List of articles"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /article/all [get]
// @Security BearerAuth
func (h *ArticleHandler) ArticleGetAll(c *gin.Context) {
	var params article2.GetAllArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllArticleQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ArticleGetByFiltersSort godoc
// @Summary      Get articles by filters and sorting
// @Description  Retrieves articles based on specified filters and sorting criteria
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  body      article.GetByFiltersSortArticleQuery  true  "Filter and sort parameters"
// @success      200      {object}  utils.Result                           "Filtered and sorted articles"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /article/filters-sort [post]
// @Security BearerAuth
func (h *ArticleHandler) ArticleGetByFiltersSort(c *gin.Context) {
	var params article2.GetByFiltersSortArticleQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByFiltersSortArticleQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminArticleGetAll godoc
// @Summary      Admin: Get all articles
// @Description  Admin endpoint to retrieve all articles with additional information
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        request  query     article.AdminGetAllArticleQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                      "List of all articles"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /article/admin/all [get]
// @Security BearerAuth
func (h *ArticleHandler) AdminArticleGetAll(c *gin.Context) {
	var params article2.AdminGetAllArticleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllArticleQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
