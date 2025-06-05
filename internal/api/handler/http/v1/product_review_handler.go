package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	product_review2 "github.com/amirex128/new_site_builder/internal/application/dto/product_review"
	"github.com/amirex128/new_site_builder/internal/application/usecase/product_review"
	"github.com/gin-gonic/gin"
)

type ProductReviewHandler struct {
	usecase   *productreviewusecase.ProductReviewUsecase
	validator *utils2.ValidationHelper
}

func NewProductReviewHandler(usc *productreviewusecase.ProductReviewUsecase) *ProductReviewHandler {
	return &ProductReviewHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateProductReview godoc
// @Summary      Create a product review
// @Description  Creates a new review for a product
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.CreateProductReviewCommand  true  "Review information"
// @success      201      {object}  utils.Result                                "Created review"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "unauthorized"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-review [post]
// @Security BearerAuth
func (h *ProductReviewHandler) CreateProductReview(c *gin.Context) {
	var params product_review2.CreateProductReviewCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateProductReviewCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateProductReview godoc
// @Summary      Update a product review
// @Description  Updates an existing product review
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.UpdateProductReviewCommand  true  "Updated review information"
// @success      200      {object}  utils.Result                                "Updated review"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "unauthorized"
// @Failure      404      {object}  utils.Result                                "Review not found"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-review [put]
// @Security BearerAuth
func (h *ProductReviewHandler) UpdateProductReview(c *gin.Context) {
	var params product_review2.UpdateProductReviewCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateProductReviewCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteProductReview godoc
// @Summary      Delete a product review
// @Description  Deletes an existing product review
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.DeleteProductReviewCommand  true  "Review ID to delete"
// @success      200      {object}  utils.Result                                "Deleted review confirmation"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "unauthorized"
// @Failure      404      {object}  utils.Result                                "Review not found"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-review [delete]
// @Security BearerAuth
func (h *ProductReviewHandler) DeleteProductReview(c *gin.Context) {
	var params product_review2.DeleteProductReviewCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteProductReviewCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdProductReview godoc
// @Summary      Get product review by ID
// @Description  Retrieves a specific product review by its ID
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.GetByIdProductReviewQuery  true  "Review ID to retrieve"
// @success      200      {object}  utils.Result                               "Review details"
// @Failure      400      {object}  utils.Result                               "Validation error"
// @Failure      401      {object}  utils.Result                               "unauthorized"
// @Failure      404      {object}  utils.Result                               "Review not found"
// @Failure      500      {object}  utils.Result                               "Internal server error"
// @Router       /product-review [get]
// @Security BearerAuth
func (h *ProductReviewHandler) GetByIdProductReview(c *gin.Context) {
	var params product_review2.GetByIdProductReviewQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdProductReviewQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllProductReview godoc
// @Summary      Get all product reviews
// @Description  Retrieves all product reviews with optional filtering
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.GetAllProductReviewQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                              "List of reviews"
// @Failure      400      {object}  utils.Result                              "Validation error"
// @Failure      401      {object}  utils.Result                              "unauthorized"
// @Failure      500      {object}  utils.Result                              "Internal server error"
// @Router       /product-review/all [get]
// @Security BearerAuth
func (h *ProductReviewHandler) GetAllProductReview(c *gin.Context) {
	var params product_review2.GetAllProductReviewQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllProductReviewQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllProductReview godoc
// @Summary      Admin: Get all product reviews
// @Description  Admin endpoint to retrieve all product reviews with additional information
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.AdminGetAllProductReviewQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                                   "List of all reviews"
// @Failure      400      {object}  utils.Result                                   "Validation error"
// @Failure      401      {object}  utils.Result                                   "unauthorized"
// @Failure      403      {object}  utils.Result                                   "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                                   "Internal server error"
// @Router       /product-review/admin/all [get]
// @Security BearerAuth
func (h *ProductReviewHandler) AdminGetAllProductReview(c *gin.Context) {
	var params product_review2.AdminGetAllProductReviewQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllProductReviewQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
