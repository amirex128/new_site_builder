package v1

import (
	"net/http"

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

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var params address.CreateAddressCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).CreateAddressCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	var params address.UpdateAddressCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpdateAddressCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	var params address.DeleteAddressCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).DeleteAddressCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Deleted().WithData(result))
}

func (h *AddressHandler) GetByIdAddress(c *gin.Context) {
	var params address.GetByIdAddressQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetByIdAddressQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *AddressHandler) GetAllAddress(c *gin.Context) {
	var params address.GetAllAddressQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllAddressQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *AddressHandler) GetAllCity(c *gin.Context) {
	var params address.GetAllCityQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllCityQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *AddressHandler) GetAllProvince(c *gin.Context) {
	var params address.GetAllProvinceQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetAllProvinceQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *AddressHandler) AdminGetAllAddress(c *gin.Context) {
	var params address.AdminGetAllAddressQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).AdminGetAllAddressQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
