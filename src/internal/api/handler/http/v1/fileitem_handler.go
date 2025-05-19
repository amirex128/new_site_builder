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

// CreateOrDirectoryItem godoc
// @Summary      Create file or directory
// @Description  Creates a new file or directory in the file system
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.CreateOrDirectoryItemCommand  true  "File or directory information"
// @Success      201      {object}  resp.Result                            "Created file or directory"
// @Failure      400      {object}  resp.Result                            "Validation error"
// @Failure      401      {object}  resp.Result                            "Unauthorized"
// @Failure      500      {object}  resp.Result                            "Internal server error"
// @Router       /file-item [post]
// @Security     BearerAuth
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

// DeleteFileItem godoc
// @Summary      Delete a file or directory
// @Description  Soft deletes a file or directory (moves to trash)
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.DeleteFileItemCommand  true  "File or directory ID to delete"
// @Success      200      {object}  resp.Result                     "Deleted file or directory confirmation"
// @Failure      400      {object}  resp.Result                     "Validation error"
// @Failure      401      {object}  resp.Result                     "Unauthorized"
// @Failure      404      {object}  resp.Result                     "File or directory not found"
// @Failure      500      {object}  resp.Result                     "Internal server error"
// @Router       /file-item [delete]
// @Security     BearerAuth
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

	resp.Deleted(c, result)
}

// ForceDeleteFileItem godoc
// @Summary      Permanently delete a file or directory
// @Description  Permanently deletes a file or directory from the system
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.ForceDeleteFileItemCommand  true  "File or directory ID to permanently delete"
// @Success      200      {object}  resp.Result                          "Permanently deleted file or directory confirmation"
// @Failure      400      {object}  resp.Result                          "Validation error"
// @Failure      401      {object}  resp.Result                          "Unauthorized"
// @Failure      404      {object}  resp.Result                          "File or directory not found"
// @Failure      500      {object}  resp.Result                          "Internal server error"
// @Router       /file-item/force [delete]
// @Security     BearerAuth
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

	resp.Deleted(c, result)
}

// UpdateFileItem godoc
// @Summary      Update a file or directory
// @Description  Updates an existing file or directory with new information
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.UpdateFileItemCommand  true  "Updated file or directory information"
// @Success      200      {object}  resp.Result                     "Updated file or directory"
// @Failure      400      {object}  resp.Result                     "Validation error"
// @Failure      401      {object}  resp.Result                     "Unauthorized"
// @Failure      404      {object}  resp.Result                     "File or directory not found"
// @Failure      500      {object}  resp.Result                     "Internal server error"
// @Router       /file-item [put]
// @Security     BearerAuth
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

// FileOperation godoc
// @Summary      Perform file operations
// @Description  Performs operations like copy, move, or rename on files and directories
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.FileOperationCommand  true  "Operation details"
// @Success      200      {object}  resp.Result                    "Operation result"
// @Failure      400      {object}  resp.Result                    "Validation error"
// @Failure      401      {object}  resp.Result                    "Unauthorized"
// @Failure      404      {object}  resp.Result                    "File or directory not found"
// @Failure      500      {object}  resp.Result                    "Internal server error"
// @Router       /file-item/operation [post]
// @Security     BearerAuth
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

// RestoreFileItem godoc
// @Summary      Restore deleted file or directory
// @Description  Restores a previously deleted file or directory from trash
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.RestoreFileItemCommand  true  "File or directory ID to restore"
// @Success      200      {object}  resp.Result                      "Restored file or directory"
// @Failure      400      {object}  resp.Result                      "Validation error"
// @Failure      401      {object}  resp.Result                      "Unauthorized"
// @Failure      404      {object}  resp.Result                      "File or directory not found in trash"
// @Failure      500      {object}  resp.Result                      "Internal server error"
// @Router       /file-item/restore [put]
// @Security     BearerAuth
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

// GetByIds godoc
// @Summary      Get files or directories by IDs
// @Description  Retrieves multiple files or directories by their IDs
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.GetByIdsQuery  true  "File or directory IDs to retrieve"
// @Success      200      {object}  resp.Result            "Files or directories details"
// @Failure      400      {object}  resp.Result            "Validation error"
// @Failure      401      {object}  resp.Result            "Unauthorized"
// @Failure      404      {object}  resp.Result            "One or more files or directories not found"
// @Failure      500      {object}  resp.Result            "Internal server error"
// @Router       /file-item/ids [get]
// @Security     BearerAuth
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

// GetDeletedTreeDirectory godoc
// @Summary      Get deleted directory tree
// @Description  Retrieves the tree structure of deleted files and directories
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.GetDeletedTreeDirectoryQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                            "Deleted directory tree"
// @Failure      400      {object}  resp.Result                            "Validation error"
// @Failure      401      {object}  resp.Result                            "Unauthorized"
// @Failure      500      {object}  resp.Result                            "Internal server error"
// @Router       /file-item/tree/deleted [get]
// @Security     BearerAuth
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

// GetDownloadFileItemById godoc
// @Summary      Download a file
// @Description  Retrieves a file for download by its ID
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.GetDownloadFileItemByIdQuery  true  "File ID to download"
// @Success      200      {object}  resp.Result                            "File download information"
// @Failure      400      {object}  resp.Result                            "Validation error"
// @Failure      401      {object}  resp.Result                            "Unauthorized"
// @Failure      404      {object}  resp.Result                            "File not found"
// @Failure      500      {object}  resp.Result                            "Internal server error"
// @Router       /file-item/download [get]
// @Security     BearerAuth
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

// GetTreeDirectory godoc
// @Summary      Get directory tree
// @Description  Retrieves the tree structure of files and directories
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.GetTreeDirectoryQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                     "Directory tree"
// @Failure      400      {object}  resp.Result                     "Validation error"
// @Failure      401      {object}  resp.Result                     "Unauthorized"
// @Failure      500      {object}  resp.Result                     "Internal server error"
// @Router       /file-item/tree [get]
// @Security     BearerAuth
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
