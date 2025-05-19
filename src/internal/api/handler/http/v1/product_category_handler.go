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

func (h *ProductCategoryHandler) CreateCategory(c *gin.Context) {
	var params product_category.CreateCategoryCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateCategoryCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *ProductCategoryHandler) UpdateCategory(c *gin.Context) {
	var params product_category.UpdateCategoryCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateCategoryCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *ProductCategoryHandler) DeleteCategory(c *gin.Context) {
	var params product_category.DeleteCategoryCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteCategoryCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c, result)
}

func (h *ProductCategoryHandler) GetByIdCategory(c *gin.Context) {
	var params product_category.GetByIdCategoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdCategoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *ProductCategoryHandler) GetAllCategory(c *gin.Context) {
	var params product_category.GetAllCategoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllCategoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *ProductCategoryHandler) AdminGetAllCategory(c *gin.Context) {
	var params product_category.AdminGetAllCategoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllCategoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
