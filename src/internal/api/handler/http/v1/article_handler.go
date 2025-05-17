package v1

import (
	producttransformer "github.com/amirex128/new_site_builder/src/internal/api/transformer/product"
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	productusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleHandler struct {
	usecase *productusecase.ArticleUsecase
}

func NewArticleHandler(usc *productusecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		usecase: usc,
	}
}

func (h *ArticleHandler) ArticleList(c *gin.Context) {
	// initializer generate dto
	params := producttransformer.NewArticleTransformer(c).
		SetSuperType().
		Build()

	result, err := h.usecase.ArticleList(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.Created().Succeeded)
		return
	}

	c.JSON(http.StatusOK, utils.GenerateBaseResponse(result, true, 0))
}
