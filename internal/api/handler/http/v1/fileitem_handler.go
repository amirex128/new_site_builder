package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	fileitem2 "github.com/amirex128/new_site_builder/internal/application/dto/fileitem"
	"github.com/amirex128/new_site_builder/internal/application/usecase/file_item"
	"github.com/gin-gonic/gin"
)

type FileItemHandler struct {
	usecase   *fileitemusecase.FileItemUsecase
	validator *utils2.ValidationHelper
}

func NewFileItemHandler(usc *fileitemusecase.FileItemUsecase) *FileItemHandler {
	return &FileItemHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateOrDirectoryItem godoc
// @Summary      Create file or directory
// @Description  Creates a new file or directory in the file system
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.CreateOrDirectoryItemCommand  true  "File or directory information"
// @success      201      {object}  utils.Result                            "Created file or directory"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /file-item [post]
// @Security BearerAuth
func (h *FileItemHandler) CreateOrDirectoryItem(c *gin.Context) {
	var params fileitem2.CreateOrDirectoryItemCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateOrDirectoryItemCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// DeleteFileItem godoc
// @Summary      Delete a file or directory
// @Description  Soft deletes a file or directory (moves to trash)
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.DeleteFileItemCommand  true  "File or directory ID to delete"
// @success      200      {object}  utils.Result                     "Deleted file or directory confirmation"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      404      {object}  utils.Result                     "File or directory not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /file-item [delete]
// @Security BearerAuth
func (h *FileItemHandler) DeleteFileItem(c *gin.Context) {
	var params fileitem2.DeleteFileItemCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.DeleteFileItemCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// ForceDeleteFileItem godoc
// @Summary      Permanently delete a file or directory
// @Description  Permanently deletes a file or directory from the system
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.ForceDeleteFileItemCommand  true  "File or directory ID to permanently delete"
// @success      200      {object}  utils.Result                          "Permanently deleted file or directory confirmation"
// @Failure      400      {object}  utils.Result                          "Validation error"
// @Failure      401      {object}  utils.Result                          "unauthorized"
// @Failure      404      {object}  utils.Result                          "File or directory not found"
// @Failure      500      {object}  utils.Result                          "Internal server error"
// @Router       /file-item/force [delete]
// @Security BearerAuth
func (h *FileItemHandler) ForceDeleteFileItem(c *gin.Context) {
	var params fileitem2.ForceDeleteFileItemCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.ForceDeleteFileItemCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// UpdateFileItem godoc
// @Summary      Update a file or directory
// @Description  Updates an existing file or directory with new information
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.UpdateFileItemCommand  true  "Updated file or directory information"
// @success      200      {object}  utils.Result                     "Updated file or directory"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      404      {object}  utils.Result                     "File or directory not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /file-item [put]
// @Security BearerAuth
func (h *FileItemHandler) UpdateFileItem(c *gin.Context) {
	var params fileitem2.UpdateFileItemCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateFileItemCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// FileOperation godoc
// @Summary      Perform file operations
// @Description  Performs operations like copy, move, or rename on files and directories
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.FileOperationCommand  true  "Operation details"
// @success      200      {object}  utils.Result                    "Operation result"
// @Failure      400      {object}  utils.Result                    "Validation error"
// @Failure      401      {object}  utils.Result                    "unauthorized"
// @Failure      404      {object}  utils.Result                    "File or directory not found"
// @Failure      500      {object}  utils.Result                    "Internal server error"
// @Router       /file-item/operation [post]
// @Security BearerAuth
func (h *FileItemHandler) FileOperation(c *gin.Context) {
	var params fileitem2.FileOperationCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.FileOperationCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// RestoreFileItem godoc
// @Summary      Restore deleted file or directory
// @Description  Restores a previously deleted file or directory from trash
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.RestoreFileItemCommand  true  "File or directory ID to restore"
// @success      200      {object}  utils.Result                      "Restored file or directory"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      404      {object}  utils.Result                      "File or directory not found in trash"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /file-item/restore [put]
// @Security BearerAuth
func (h *FileItemHandler) RestoreFileItem(c *gin.Context) {
	var params fileitem2.RestoreFileItemCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.RestoreFileItemCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIds godoc
// @Summary      Get files or directories by IDs
// @Description  Retrieves multiple files or directories by their IDs
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  query     fileitem.GetByIdsQuery  true  "File or directory IDs to retrieve"
// @success      200      {object}  utils.Result            "Files or directories details"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      404      {object}  utils.Result            "One or more files or directories not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /file-item/ids [get]
// @Security BearerAuth
func (h *FileItemHandler) GetByIds(c *gin.Context) {
	var params fileitem2.GetByIdsQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdsQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetDeletedTreeDirectory godoc
// @Summary      Get deleted directory tree
// @Description  Retrieves the tree structure of deleted files and directories
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  query     fileitem.GetDeletedTreeDirectoryQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                            "Deleted directory tree"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /file-item/tree/deleted [get]
// @Security BearerAuth
func (h *FileItemHandler) GetDeletedTreeDirectory(c *gin.Context) {
	var params fileitem2.GetDeletedTreeDirectoryQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetDeletedTreeDirectoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetDownloadFileItemById godoc
// @Summary      Download a file
// @Description  Retrieves a file for download by its ID
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.GetDownloadFileItemByIdQuery  true  "File ID to download"
// @success      200      {object}  utils.Result                            "File download information"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      404      {object}  utils.Result                            "File not found"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /file-item/download [get]
// @Security BearerAuth
func (h *FileItemHandler) GetDownloadFileItemById(c *gin.Context) {
	var params fileitem2.GetDownloadFileItemByIdQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetDownloadFileItemByIdQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetTreeDirectory godoc
// @Summary      Get directory tree
// @Description  Retrieves the tree structure of files and directories
// @Tags         file-item
// @Accept       json
// @Produce      json
// @Param        request  body      fileitem.GetTreeDirectoryQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                     "Directory tree"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /file-item/tree [get]
// @Security BearerAuth
func (h *FileItemHandler) GetTreeDirectory(c *gin.Context) {
	var params fileitem2.GetTreeDirectoryQuery
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetTreeDirectoryQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
