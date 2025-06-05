package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	product_category2 "github.com/amirex128/new_site_builder/internal/application/dto/product_category"
	"github.com/amirex128/new_site_builder/internal/application/usecase/product_category"
	"github.com/gin-gonic/gin"
)

type ProductCategoryHandler struct {
	usecase   *productcategoryusecase.ProductCategoryUsecase
	validator *utils2.ValidationHelper
}

func NewProductCategoryHandler(usc *productcategoryusecase.ProductCategoryUsecase) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateCategory godoc
// @Summary      Create a new product category
// @Description  Creates a new category for products with the provided information
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  body      product_category.CreateCategoryCommand  true  "Category information"
// @success      201      {object}  utils.Result                             "Created category"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /product-category [post]
// @Security BearerAuth
func (h *ProductCategoryHandler) CreateCategory(c *gin.Context) {
	var params product_category2.CreateCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateCategoryCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateCategory godoc
// @Summary      Update a product category
// @Description  Updates an existing product category with the provided information
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  body      product_category.UpdateCategoryCommand  true  "Updated category information"
// @success      200      {object}  utils.Result                             "Updated category"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Category not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /product-category [put]
// @Security BearerAuth
func (h *ProductCategoryHandler) UpdateCategory(c *gin.Context) {
	var params product_category2.UpdateCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateCategoryCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteCategory godoc
// @Summary      Delete a product category
// @Description  Deletes an existing product category by its ID
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  body      product_category.DeleteCategoryCommand  true  "Category ID to delete"
// @success      200      {object}  utils.Result                             "Deleted category confirmation"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Category not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /product-category [delete]
// @Security BearerAuth
func (h *ProductCategoryHandler) DeleteCategory(c *gin.Context) {
	var params product_category2.DeleteCategoryCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteCategoryCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdCategory godoc
// @Summary      Get product category by ID
// @Description  Retrieves a specific product category by its ID
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  query     product_category.GetByIdCategoryQuery  true  "Category ID to retrieve"
// @success      200      {object}  utils.Result                            "Category details"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      404      {object}  utils.Result                            "Category not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /product-category [get]
// @Security BearerAuth
func (h *ProductCategoryHandler) GetByIdCategory(c *gin.Context) {
	var params product_category2.GetByIdCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdCategoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllCategory godoc
// @Summary      Get all product categories
// @Description  Retrieves all product categories with optional filtering
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  query     product_category.GetAllCategoryQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                           "List of categories"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /product-category/all [get]
// @Security BearerAuth
func (h *ProductCategoryHandler) GetAllCategory(c *gin.Context) {
	var params product_category2.GetAllCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllCategoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllCategory godoc
// @Summary      Admin: Get all product categories
// @Description  Admin endpoint to retrieve all product categories with additional information
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  query     product_category.AdminGetAllCategoryQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                                "List of all categories"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "unauthorized"
// @Failure      403      {object}  utils.Result                                "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /product-category/admin/all [get]
// @Security BearerAuth
func (h *ProductCategoryHandler) AdminGetAllCategory(c *gin.Context) {
	var params product_category2.AdminGetAllCategoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllCategoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
