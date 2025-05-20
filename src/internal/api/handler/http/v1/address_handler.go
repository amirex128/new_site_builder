package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/address"
	addressusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/address"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	usecase   *addressusecase.AddressUsecase
	validator *utils.ValidationHelper
}

func NewAddressHandler(usc *addressusecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// CreateAddress godoc
// @Summary      Create a new address
// @Description  Creates a new address for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  body      address.CreateAddressCommand  true  "Address information"
// @Success      201      {object}  utils.Result                   "Created address"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /address [post]
// @Security BearerAuth
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var params address.CreateAddressCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	panic("s")
	result, err := h.usecase.CreateAddressCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// UpdateAddress godoc
// @Summary      Update an existing address
// @Description  Updates an existing address for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  body      address.UpdateAddressCommand  true  "Updated address information"
// @Success      200      {object}  utils.Result                   "Updated address"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      404      {object}  utils.Result                   "Address not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /address [put]
// @Security BearerAuth
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	var params address.UpdateAddressCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateAddressCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// DeleteAddress godoc
// @Summary      Delete an address
// @Description  Deletes an address for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  body      address.DeleteAddressCommand  true  "Address ID to delete"
// @Success      200      {object}  utils.Result                   "Deleted address confirmation"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      404      {object}  utils.Result                   "Address not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /address [delete]
// @Security BearerAuth
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	var params address.DeleteAddressCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.DeleteAddressCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Deleted(c, result)
}

// GetByIdAddress godoc
// @Summary      Get address by ID
// @Description  Retrieves a specific address by its ID for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetByIdAddressQuery  true  "Address ID to retrieve"
// @Success      200      {object}  utils.Result                  "Address details"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "Unauthorized"
// @Failure      404      {object}  utils.Result                  "Address not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /address [get]
// @Security BearerAuth
func (h *AddressHandler) GetByIdAddress(c *gin.Context) {
	var params address.GetByIdAddressQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdAddressQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllAddress godoc
// @Summary      Get all addresses
// @Description  Retrieves all addresses for the authenticated user
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetAllAddressQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                 "List of addresses"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "Unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /address/all [get]
// @Security BearerAuth
func (h *AddressHandler) GetAllAddress(c *gin.Context) {
	var params address.GetAllAddressQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllAddressQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllCity godoc
// @Summary      Get all cities
// @Description  Retrieves all available cities
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetAllCityQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result              "List of cities"
// @Failure      400      {object}  utils.Result              "Validation error"
// @Failure      401      {object}  utils.Result              "Unauthorized"
// @Failure      500      {object}  utils.Result              "Internal server error"
// @Router       /address/city/all [get]
// @Security BearerAuth
func (h *AddressHandler) GetAllCity(c *gin.Context) {
	var params address.GetAllCityQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllCityQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllProvince godoc
// @Summary      Get all provinces
// @Description  Retrieves all available provinces
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.GetAllProvinceQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                  "List of provinces"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "Unauthorized"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /address/province/all [get]
// @Security BearerAuth
func (h *AddressHandler) GetAllProvince(c *gin.Context) {
	var params address.GetAllProvinceQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllProvinceQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// AdminGetAllAddress godoc
// @Summary      Admin: Get all addresses
// @Description  Admin endpoint to retrieve all addresses in the system
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        request  query     address.AdminGetAllAddressQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                      "List of all addresses"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "Unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /address/admin/all [get]
// @Security BearerAuth
func (h *AddressHandler) AdminGetAllAddress(c *gin.Context) {
	var params address.AdminGetAllAddressQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllAddressQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
