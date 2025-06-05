package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	address2 "github.com/amirex128/new_site_builder/internal/application/dto/address"
	"github.com/amirex128/new_site_builder/internal/application/usecase/address"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	usecase   *addressusecase.AddressUsecase
	validator *utils2.ValidationHelper
}

func NewAddressHandler(usc *addressusecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateAddress godoc
// @Summary      Create a new address
// @Description  Creates a new address for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  body      address.CreateAddressCommand  true  "Address information"
// @success      201      {object}  utils.Result                   "Created address"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /address [post]
// @Security BearerAuth
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var params address2.CreateAddressCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateAddressCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateAddress godoc
// @Summary      Update an existing address
// @Description  Updates an existing address for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  body      address.UpdateAddressCommand  true  "Updated address information"
// @success      200      {object}  utils.Result                   "Updated address"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      404      {object}  utils.Result                   "Address not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /address [put]
// @Security BearerAuth
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	var params address2.UpdateAddressCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateAddressCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteAddress godoc
// @Summary      Delete an address
// @Description  Deletes an address for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  body      address.DeleteAddressCommand  true  "Address ID to delete"
// @success      200      {object}  utils.Result                   "Deleted address confirmation"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      404      {object}  utils.Result                   "Address not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /address [delete]
// @Security BearerAuth
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	var params address2.DeleteAddressCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteAddressCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdAddress godoc
// @Summary      Get address by ID
// @Description  Retrieves a specific address by its ID for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetByIdAddressQuery  true  "Address ID to retrieve"
// @success      200      {object}  utils.Result                  "Address details"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "unauthorized"
// @Failure      404      {object}  utils.Result                  "Address not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /address [get]
// @Security BearerAuth
func (h *AddressHandler) GetByIdAddress(c *gin.Context) {
	var params address2.GetByIdAddressQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdAddressQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllAddress godoc
// @Summary      Get all addresses
// @Description  Retrieves all addresses for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetAllAddressQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                 "List of addresses"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /address/all [get]
// @Security BearerAuth
func (h *AddressHandler) GetAllAddress(c *gin.Context) {
	var params address2.GetAllAddressQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllAddressQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllCity godoc
// @Summary      Get all cities
// @Description  Retrieves all available cities
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetAllCityQuery  true  "Query parameters"
// @success      200      {object}  utils.Result              "List of cities"
// @Failure      400      {object}  utils.Result              "Validation error"
// @Failure      401      {object}  utils.Result              "unauthorized"
// @Failure      500      {object}  utils.Result              "Internal server error"
// @Router       /address/city/all [get]
// @Security BearerAuth
func (h *AddressHandler) GetAllCity(c *gin.Context) {
	var params address2.GetAllCityQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllCityQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllProvince godoc
// @Summary      Get all provinces
// @Description  Retrieves all available provinces
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetAllProvinceQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                  "List of provinces"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "unauthorized"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /address/province/all [get]
// @Security BearerAuth
func (h *AddressHandler) GetAllProvince(c *gin.Context) {
	var params address2.GetAllProvinceQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllProvinceQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllAddress godoc
// @Summary      Admin: Get all addresses
// @Description  Admin endpoint to retrieve all addresses in the system
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.AdminGetAllAddressQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                      "List of all addresses"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /address/admin/all [get]
// @Security BearerAuth
func (h *AddressHandler) AdminGetAllAddress(c *gin.Context) {
	var params address2.AdminGetAllAddressQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllAddressQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
