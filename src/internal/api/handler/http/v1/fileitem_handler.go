package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/fileitem"
	fileitemusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/file_item"
	"github.com/gin-gonic/gin"
)

type FileItemHandler struct {
	usecase   *fileitemusecase.FileItemUsecase
	validator *utils.ValidationHelper
}

func NewFileItemHandler(usc *fileitemusecase.FileItemUsecase) *FileItemHandler {
	return &FileItemHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *FileItemHandler) CreateOrDirectoryItem(c *gin.Context) {
	var params fileitem.CreateOrDirectoryItemCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrDirectoryItemCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *FileItemHandler) DeleteFileItem(c *gin.Context) {
	var params fileitem.DeleteFileItemCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.DeleteFileItemCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *FileItemHandler) ForceDeleteFileItem(c *gin.Context) {
	var params fileitem.ForceDeleteFileItemCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.ForceDeleteFileItemCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Deleted(c)
}

func (h *FileItemHandler) UpdateFileItem(c *gin.Context) {
	var params fileitem.UpdateFileItemCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateFileItemCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *FileItemHandler) FileOperation(c *gin.Context) {
	var params fileitem.FileOperationCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.FileOperationCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *FileItemHandler) RestoreFileItem(c *gin.Context) {
	var params fileitem.RestoreFileItemCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.RestoreFileItemCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *FileItemHandler) GetByIds(c *gin.Context) {
	var params fileitem.GetByIdsQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdsQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *FileItemHandler) GetDeletedTreeDirectory(c *gin.Context) {
	var params fileitem.GetDeletedTreeDirectoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetDeletedTreeDirectoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *FileItemHandler) GetDownloadFileItemById(c *gin.Context) {
	var params fileitem.GetDownloadFileItemByIdQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetDownloadFileItemByIdQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *FileItemHandler) GetTreeDirectory(c *gin.Context) {
	var params fileitem.GetTreeDirectoryQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetTreeDirectoryQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
