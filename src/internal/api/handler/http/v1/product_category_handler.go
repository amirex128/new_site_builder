package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_category"
	productcategoryusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product_category"
	"github.com/gin-gonic/gin"
)

type ProductCategoryHandler struct {
	usecase   *productcategoryusecase.ProductCategoryUsecase
	validator *utils.ValidationHelper
}

func NewProductCategoryHandler(usc *productcategoryusecase.ProductCategoryUsecase) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// CreateCategory godoc
// @Summary      Create a new product category
// @Description  Creates a new category for products with the provided information
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  body      product_category.CreateCategoryCommand  true  "Category information"
// @Success      201      {object}  resp.Result                             "Created category"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      401      {object}  resp.Result                             "Unauthorized"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /product-category [post]
// @Security     BearerAuth
func (h *ProductCategoryHandler) CreateCategory(c *gin.Context) {
	var params product_category.CreateCategoryCommand
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

// UpdateCategory godoc
// @Summary      Update a product category
// @Description  Updates an existing product category with the provided information
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  body      product_category.UpdateCategoryCommand  true  "Updated category information"
// @Success      200      {object}  resp.Result                             "Updated category"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      401      {object}  resp.Result                             "Unauthorized"
// @Failure      404      {object}  resp.Result                             "Category not found"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /product-category [put]
// @Security     BearerAuth
func (h *ProductCategoryHandler) UpdateCategory(c *gin.Context) {
	var params product_category.UpdateCategoryCommand
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

// DeleteCategory godoc
// @Summary      Delete a product category
// @Description  Deletes an existing product category by its ID
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  body      product_category.DeleteCategoryCommand  true  "Category ID to delete"
// @Success      200      {object}  resp.Result                             "Deleted category confirmation"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      401      {object}  resp.Result                             "Unauthorized"
// @Failure      404      {object}  resp.Result                             "Category not found"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /product-category [delete]
// @Security     BearerAuth
func (h *ProductCategoryHandler) DeleteCategory(c *gin.Context) {
	var params product_category.DeleteCategoryCommand
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

// GetByIdCategory godoc
// @Summary      Get product category by ID
// @Description  Retrieves a specific product category by its ID
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  query     product_category.GetByIdCategoryQuery  true  "Category ID to retrieve"
// @Success      200      {object}  resp.Result                            "Category details"
// @Failure      400      {object}  resp.Result                            "Validation error"
// @Failure      401      {object}  resp.Result                            "Unauthorized"
// @Failure      404      {object}  resp.Result                            "Category not found"
// @Failure      500      {object}  resp.Result                            "Internal server error"
// @Router       /product-category [get]
// @Security     BearerAuth
func (h *ProductCategoryHandler) GetByIdCategory(c *gin.Context) {
	var params product_category.GetByIdCategoryQuery
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

// GetAllCategory godoc
// @Summary      Get all product categories
// @Description  Retrieves all product categories with optional filtering
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  query     product_category.GetAllCategoryQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                           "List of categories"
// @Failure      400      {object}  resp.Result                           "Validation error"
// @Failure      401      {object}  resp.Result                           "Unauthorized"
// @Failure      500      {object}  resp.Result                           "Internal server error"
// @Router       /product-category/all [get]
// @Security     BearerAuth
func (h *ProductCategoryHandler) GetAllCategory(c *gin.Context) {
	var params product_category.GetAllCategoryQuery
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

// AdminGetAllCategory godoc
// @Summary      Admin: Get all product categories
// @Description  Admin endpoint to retrieve all product categories with additional information
// @Tags         product-category
// @Accept       json
// @Produce      json
// @Param        request  query     product_category.AdminGetAllCategoryQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                                "List of all categories"
// @Failure      400      {object}  resp.Result                                "Validation error"
// @Failure      401      {object}  resp.Result                                "Unauthorized"
// @Failure      403      {object}  resp.Result                                "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result                                "Internal server error"
// @Router       /product-category/admin/all [get]
// @Security     BearerAuth
func (h *ProductCategoryHandler) AdminGetAllCategory(c *gin.Context) {
	var params product_category.AdminGetAllCategoryQuery
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
