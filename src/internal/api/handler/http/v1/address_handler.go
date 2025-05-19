package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
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
// @Summary      Update a user
// @Description  update an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "User ID"
// @Param        user  body      address.CreateAddressCommand    true  "User object"
// @Success      200   {object}  address.CreateAddressCommand
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/{id} [put]
// @Security     BearerAuth
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var params address.CreateAddressCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).CreateAddressCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	var params address.UpdateAddressCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpdateAddressCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	var params address.DeleteAddressCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).DeleteAddressCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *AddressHandler) GetByIdAddress(c *gin.Context) {
	var params address.GetByIdAddressQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetByIdAddressQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *AddressHandler) GetAllAddress(c *gin.Context) {
	var params address.GetAllAddressQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllAddressQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *AddressHandler) GetAllCity(c *gin.Context) {
	var params address.GetAllCityQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllCityQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *AddressHandler) GetAllProvince(c *gin.Context) {
	var params address.GetAllProvinceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllProvinceQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *AddressHandler) AdminGetAllAddress(c *gin.Context) {
	var params address.AdminGetAllAddressQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).AdminGetAllAddressQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
