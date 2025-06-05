package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	header_footer2 "github.com/amirex128/new_site_builder/internal/application/dto/header_footer"
	"github.com/amirex128/new_site_builder/internal/application/usecase/header_footer"
	"github.com/gin-gonic/gin"
)

type HeaderFooterHandler struct {
	usecase   *headerfooterusecase.HeaderFooterUsecase
	validator *utils2.ValidationHelper
}

func NewHeaderFooterHandler(usc *headerfooterusecase.HeaderFooterUsecase) *HeaderFooterHandler {
	return &HeaderFooterHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateHeaderFooter godoc
// @Summary      Create a header/footer
// @Description  Creates a new header and footer template for websites
// @Tags         header-footer
// @Accept       json
// @Produce      json
// @Param        request  body      header_footer.CreateHeaderFooterCommand  true  "Header/footer information"
// @success      201      {object}  utils.Result                              "Created header/footer"
// @Failure      400      {object}  utils.Result                              "Validation error"
// @Failure      401      {object}  utils.Result                              "unauthorized"
// @Failure      500      {object}  utils.Result                              "Internal server error"
// @Router       /header-footer [post]
// @Security BearerAuth
func (h *HeaderFooterHandler) CreateHeaderFooter(c *gin.Context) {
	var params header_footer2.CreateHeaderFooterCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateHeaderFooterCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateHeaderFooter godoc
// @Summary      Update a header/footer
// @Description  Updates an existing header and footer template
// @Tags         header-footer
// @Accept       json
// @Produce      json
// @Param        request  body      header_footer.UpdateHeaderFooterCommand  true  "Updated header/footer information"
// @success      200      {object}  utils.Result                              "Updated header/footer"
// @Failure      400      {object}  utils.Result                              "Validation error"
// @Failure      401      {object}  utils.Result                              "unauthorized"
// @Failure      404      {object}  utils.Result                              "Header/footer not found"
// @Failure      500      {object}  utils.Result                              "Internal server error"
// @Router       /header-footer [put]
// @Security BearerAuth
func (h *HeaderFooterHandler) UpdateHeaderFooter(c *gin.Context) {
	var params header_footer2.UpdateHeaderFooterCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateHeaderFooterCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteHeaderFooter godoc
// @Summary      Delete a header/footer
// @Description  Deletes an existing header and footer template
// @Tags         header-footer
// @Accept       json
// @Produce      json
// @Param        request  body      header_footer.DeleteHeaderFooterCommand  true  "Header/footer ID to delete"
// @success      200      {object}  utils.Result                              "Deleted header/footer confirmation"
// @Failure      400      {object}  utils.Result                              "Validation error"
// @Failure      401      {object}  utils.Result                              "unauthorized"
// @Failure      404      {object}  utils.Result                              "Header/footer not found"
// @Failure      500      {object}  utils.Result                              "Internal server error"
// @Router       /header-footer [delete]
// @Security BearerAuth
func (h *HeaderFooterHandler) DeleteHeaderFooter(c *gin.Context) {
	var params header_footer2.DeleteHeaderFooterCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteHeaderFooterCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdHeaderFooter godoc
// @Summary      Get header/footer by ID
// @Description  Retrieves a specific header and footer template by its ID
// @Tags         header-footer
// @Accept       json
// @Produce      json
// @Param        request  query     header_footer.GetByIdHeaderFooterQuery  true  "Header/footer ID to retrieve"
// @success      200      {object}  utils.Result                             "Header/footer details"
// @Failure      400      {object}  utils.Result                             "Validation error"
// @Failure      401      {object}  utils.Result                             "unauthorized"
// @Failure      404      {object}  utils.Result                             "Header/footer not found"
// @Failure      500      {object}  utils.Result                             "Internal server error"
// @Router       /header-footer [get]
// @Security BearerAuth
func (h *HeaderFooterHandler) GetByIdHeaderFooter(c *gin.Context) {
	var params header_footer2.GetByIdHeaderFooterQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdHeaderFooterQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllHeaderFooter godoc
// @Summary      Get all header/footers
// @Description  Retrieves all header and footer templates
// @Tags         header-footer
// @Accept       json
// @Produce      json
// @Param        request  query     header_footer.GetAllHeaderFooterQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                            "List of header/footers"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /header-footer/all [get]
// @Security BearerAuth
func (h *HeaderFooterHandler) GetAllHeaderFooter(c *gin.Context) {
	var params header_footer2.GetAllHeaderFooterQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllHeaderFooterQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllHeaderFooter godoc
// @Summary      Admin: Get all header/footers
// @Description  Admin endpoint to retrieve all header and footer templates with additional information
// @Tags         header-footer
// @Accept       json
// @Produce      json
// @Param        request  query     header_footer.AdminGetAllHeaderFooterQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                                 "List of all header/footers"
// @Failure      400      {object}  utils.Result                                 "Validation error"
// @Failure      401      {object}  utils.Result                                 "unauthorized"
// @Failure      403      {object}  utils.Result                                 "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                                 "Internal server error"
// @Router       /header-footer/admin/all [get]
// @Security BearerAuth
func (h *HeaderFooterHandler) AdminGetAllHeaderFooter(c *gin.Context) {
	var params header_footer2.AdminGetAllHeaderFooterQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllHeaderFooterQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
