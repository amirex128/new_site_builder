package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	basket2 "github.com/amirex128/new_site_builder/internal/application/dto/basket"
	"github.com/amirex128/new_site_builder/internal/application/usecase/basket"
	"github.com/gin-gonic/gin"
)

type BasketHandler struct {
	usecase   *basketusecase.BasketUsecase
	validator *utils2.ValidationHelper
}

func NewBasketHandler(usc *basketusecase.BasketUsecase) *BasketHandler {
	return &BasketHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// UpdateBasket godoc
// @Summary      Update shopping basket
// @Description  Updates the user's shopping basket with new items or quantities
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        request  body      basket.UpdateBasketCommand  true  "Updated basket information"
// @success      200      {object}  utils.Result                 "Updated basket"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /basket [put]
// @Security BearerAuth
func (h *BasketHandler) UpdateBasket(c *gin.Context) {
	var params basket2.UpdateBasketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateBasketCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetBasket godoc
// @Summary      Get current basket
// @Description  Retrieves the current user's shopping basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        request  query     basket.GetBasketQuery  true  "Query parameters"
// @success      200      {object}  utils.Result            "Basket details"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "Basket not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /basket [get]
// @Security BearerAuth
func (h *BasketHandler) GetBasket(c *gin.Context) {
	var params basket2.GetBasketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetBasketQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllBasketUser godoc
// @Summary      Get all user baskets
// @Description  Retrieves all shopping baskets for the current user
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        request  query     basket.GetAllBasketUserQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                   "List of user baskets"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /basket/user/all [get]
// @Security BearerAuth
func (h *BasketHandler) GetAllBasketUser(c *gin.Context) {
	var params basket2.GetAllBasketUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllBasketUserQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllBasketUser godoc
// @Summary      Admin: Get all user baskets
// @Description  Admin endpoint to retrieve all shopping baskets across all users
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        request  query     basket.AdminGetAllBasketUserQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                        "List of all user baskets"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      401      {object}  utils.Result                        "unauthorized"
// @Failure      403      {object}  utils.Result                        "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /basket/admin/all [get]
// @Security BearerAuth
func (h *BasketHandler) AdminGetAllBasketUser(c *gin.Context) {
	var params basket2.AdminGetAllBasketUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllBasketUserQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
