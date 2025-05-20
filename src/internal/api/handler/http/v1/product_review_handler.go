package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
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

// CreateProductReview godoc
// @Summary      Create a product review
// @Description  Creates a new review for a product
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.CreateProductReviewCommand  true  "Review information"
// @Success      201      {object}  utils.Result                                "Created review"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "Unauthorized"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-review [post]
// @Security BearerAuth
func (h *ProductReviewHandler) CreateProductReview(c *gin.Context) {
	var params product_review.CreateProductReviewCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateProductReviewCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// UpdateProductReview godoc
// @Summary      Update a product review
// @Description  Updates an existing product review
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.UpdateProductReviewCommand  true  "Updated review information"
// @Success      200      {object}  utils.Result                                "Updated review"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "Unauthorized"
// @Failure      404      {object}  utils.Result                                "Review not found"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-review [put]
// @Security BearerAuth
func (h *ProductReviewHandler) UpdateProductReview(c *gin.Context) {
	var params product_review.UpdateProductReviewCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateProductReviewCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// DeleteProductReview godoc
// @Summary      Delete a product review
// @Description  Deletes an existing product review
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.DeleteProductReviewCommand  true  "Review ID to delete"
// @Success      200      {object}  utils.Result                                "Deleted review confirmation"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "Unauthorized"
// @Failure      404      {object}  utils.Result                                "Review not found"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-review [delete]
// @Security BearerAuth
func (h *ProductReviewHandler) DeleteProductReview(c *gin.Context) {
	var params product_review.DeleteProductReviewCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeleteProductReviewCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Deleted(c, result)
}

// GetByIdProductReview godoc
// @Summary      Get product review by ID
// @Description  Retrieves a specific product review by its ID
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.GetByIdProductReviewQuery  true  "Review ID to retrieve"
// @Success      200      {object}  utils.Result                               "Review details"
// @Failure      400      {object}  utils.Result                               "Validation error"
// @Failure      401      {object}  utils.Result                               "Unauthorized"
// @Failure      404      {object}  utils.Result                               "Review not found"
// @Failure      500      {object}  utils.Result                               "Internal server error"
// @Router       /product-review [get]
// @Security BearerAuth
func (h *ProductReviewHandler) GetByIdProductReview(c *gin.Context) {
	var params product_review.GetByIdProductReviewQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdProductReviewQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllProductReview godoc
// @Summary      Get all product reviews
// @Description  Retrieves all product reviews with optional filtering
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.GetAllProductReviewQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                              "List of reviews"
// @Failure      400      {object}  utils.Result                              "Validation error"
// @Failure      401      {object}  utils.Result                              "Unauthorized"
// @Failure      500      {object}  utils.Result                              "Internal server error"
// @Router       /product-review/all [get]
// @Security BearerAuth
func (h *ProductReviewHandler) GetAllProductReview(c *gin.Context) {
	var params product_review.GetAllProductReviewQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.GetAllProductReviewQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// AdminGetAllProductReview godoc
// @Summary      Admin: Get all product reviews
// @Description  Admin endpoint to retrieve all product reviews with additional information
// @Tags         product-review
// @Accept       json
// @Produce      json
// @Param        request  body      product_review.AdminGetAllProductReviewQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                                   "List of all reviews"
// @Failure      400      {object}  utils.Result                                   "Validation error"
// @Failure      401      {object}  utils.Result                                   "Unauthorized"
// @Failure      403      {object}  utils.Result                                   "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                                   "Internal server error"
// @Router       /product-review/admin/all [get]
// @Security BearerAuth
func (h *ProductReviewHandler) AdminGetAllProductReview(c *gin.Context) {
	var params product_review.AdminGetAllProductReviewQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllProductReviewQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
