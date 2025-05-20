package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	productusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase   *productusecase.ProductUsecase
	validator *utils.ValidationHelper
}

func NewProductHandler(usc *productusecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Creates a new product with the provided information
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  body      product.CreateProductCommand  true  "Product information"
// @Success      201      {object}  utils.Result                   "Created product"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /product [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var params product.CreateProductCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateProductCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product with the provided information
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  body      product.UpdateProductCommand  true  "Updated product information"
// @Success      200      {object}  utils.Result                   "Updated product"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      404      {object}  utils.Result                   "Product not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /product [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var params product.UpdateProductCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateProductCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Deletes an existing product by its ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  body      product.DeleteProductCommand  true  "Product ID to delete"
// @Success      200      {object}  utils.Result                   "Deleted product confirmation"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      404      {object}  utils.Result                   "Product not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /product [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var params product.DeleteProductCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeleteProductCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Deleted(c, result)
}

// GetByIdProduct godoc
// @Summary      Get product by ID
// @Description  Retrieves a specific product by its ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  query     product.GetByIdProductQuery  true  "Product ID to retrieve"
// @Success      200      {object}  utils.Result                  "Product details"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "Unauthorized"
// @Failure      404      {object}  utils.Result                  "Product not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /product [get]
// @Security BearerAuth
func (h *ProductHandler) GetByIdProduct(c *gin.Context) {
	var params product.GetByIdProductQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdProductQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllProduct godoc
// @Summary      Get all products
// @Description  Retrieves all products with optional filtering
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  query     product.GetAllProductQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                 "List of products"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "Unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /product/all [get]
// @Security BearerAuth
func (h *ProductHandler) GetAllProduct(c *gin.Context) {
	var params product.GetAllProductQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllProductQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetByFiltersSortProduct godoc
// @Summary      Get products by filters and sorting
// @Description  Retrieves products based on specified filters and sorting criteria
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  body      product.GetByFiltersSortProductQuery  true  "Filter and sort parameters"
// @Success      200      {object}  utils.Result                           "Filtered and sorted products"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "Unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /product/filters-sort [post]
// @Security BearerAuth
func (h *ProductHandler) GetByFiltersSortProduct(c *gin.Context) {
	var params product.GetByFiltersSortProductQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.GetByFiltersSortProductQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// AdminGetAllProduct godoc
// @Summary      Admin: Get all products
// @Description  Admin endpoint to retrieve all products with additional information
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  query     product.AdminGetAllProductQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                      "List of all products"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "Unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /product/admin/all [get]
// @Security BearerAuth
func (h *ProductHandler) AdminGetAllProduct(c *gin.Context) {
	var params product.AdminGetAllProductQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllProductQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
