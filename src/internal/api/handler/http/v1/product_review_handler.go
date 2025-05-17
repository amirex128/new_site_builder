package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_review"
	productreviewusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product_review"
	"github.com/gin-gonic/gin"
)

type ProductReviewHandler struct {
	usecase   *productreviewusecase.ProductReviewUsecase
	validator *utils.ValidationHelper
}

func NewProductReviewHandler(usc *productreviewusecase.ProductReviewUsecase) *ProductReviewHandler {
	return &ProductReviewHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *ProductReviewHandler) CreateProductReview(c *gin.Context) {
	var params product_review.CreateProductReviewCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateProductReviewCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *ProductReviewHandler) UpdateProductReview(c *gin.Context) {
	var params product_review.UpdateProductReviewCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateProductReviewCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *ProductReviewHandler) DeleteProductReview(c *gin.Context) {
	var params product_review.DeleteProductReviewCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteProductReviewCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *ProductReviewHandler) GetByIdProductReview(c *gin.Context) {
	var params product_review.GetByIdProductReviewQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdProductReviewQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ProductReviewHandler) GetAllProductReview(c *gin.Context) {
	var params product_review.GetAllProductReviewQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllProductReviewQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *ProductReviewHandler) AdminGetAllProductReview(c *gin.Context) {
	var params product_review.AdminGetAllProductReviewQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllProductReviewQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
