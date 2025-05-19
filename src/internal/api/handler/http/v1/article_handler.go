package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
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

func (h *ArticleHandler) ArticleCreate(c *gin.Context) {
	var params article.CreateArticleCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateArticleCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *ArticleHandler) ArticleUpdate(c *gin.Context) {
	var params article.UpdateArticleCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateArticleCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *ArticleHandler) ArticleDelete(c *gin.Context) {
	var params article.DeleteArticleCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteArticleCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c, result)
}

func (h *ArticleHandler) ArticleGet(c *gin.Context) {
	var params article.GetByIdArticleQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdArticleQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *ArticleHandler) ArticleGetAll(c *gin.Context) {
	var params article.GetAllArticleQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllArticleQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *ArticleHandler) ArticleGetByFiltersSort(c *gin.Context) {
	var params article.GetByFiltersSortArticleQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByFiltersSortArticleQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *ArticleHandler) AdminArticleGetAll(c *gin.Context) {
	var params article.AdminGetAllArticleQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllArticleQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
