package fileitemusecase

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	contractStorage "github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/fileitem"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"github.com/gin-gonic/gin"
)

type FileItemUsecase struct {
	*usecase.BaseUsecase
	FileItemRepo   repository.IFileItemRepository
	storageRepo    repository.IStorageRepository
	storageService contractStorage.IStorageService
	authContext    func(c *gin.Context) contractStorage.IAuthService
}

func NewFileItemUsecase(c contract.IContainer) *FileItemUsecase {
	return &FileItemUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		FileItemRepo:   c.GetFileItemRepo(),
		storageRepo:    c.GetStorageRepo(),
		storageService: c.GetStorageService(),
		authContext:    c.GetAuthTransientService(),
	}
}

// Add helper function to convert fileitem.FileItemPermissionEnum to service.FileItemPermissionEnum
func toServicePermissionEnum(p enums.FileItemPermissionEnum) contractStorage.FileItemPermissionEnum {
	switch p {
	case enums.FileItemPrivatePermission:
		return contractStorage.Private
	case enums.FileItemPublicPermission:
		return contractStorage.Public
	default:
		return contractStorage.Private // default fallback
	}
}

// CreateOrDirectoryItemCommand handles the creation of a new file or directory
func (u *FileItemUsecase) CreateOrDirectoryItemCommand(params *fileitem.CreateOrDirectoryItemCommand) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	storage, err := u.storageRepo.GetByUserID(*userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving storage: %v", err)
	}

	// Check if storage is expired
	expired, err := u.storageRepo.CheckHasExpired(storage.ID)
	if err != nil {
		return nil, fmt.Errorf("error checking storage expiry: %v", err)
	}
	if expired {
		return nil, fmt.Errorf("storage has expired")
	}

	// For files, check quota
	if !*params.IsDirectory && params.File != nil {
		exceeded, err := u.storageRepo.CheckQuotaExceeded(storage.ID, int64(params.File.Size))
		if err != nil {
			return nil, fmt.Errorf("error checking quota: %v", err)
		}
		if exceeded {
			return nil, fmt.Errorf("storage quota exceeded")
		}
	}

	// Get parent directory and its path
	var parentPath string
	if params.ParentID != nil && *params.ParentID > 0 {
		parent, err := u.FileItemRepo.GetByID(*params.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent directory not found: %v", err)
		}
		parentPath = getFullPath(parent)
	}

	// Generate server info
	bucketName := fmt.Sprintf("user-%d", *userID)
	serverKey := fmt.Sprintf("S%d", *userID%2+1)

	// Create bucket if not exists
	if err := u.storageService.CreateBucketIfNotExists(serverKey, bucketName); err != nil {
		return nil, fmt.Errorf("error creating bucket: %v", err)
	}

	var fileItem domain.FileItem
	fileItem.UserID = *userID
	fileItem.ServerKey = serverKey
	fileItem.BucketName = bucketName
	fileItem.CreatedAt = time.Now()
	fileItem.UpdatedAt = time.Now()
	fileItem.IsDeleted = false
	fileItem.IsDirectory = *params.IsDirectory
	fileItem.Permission = *params.Permission

	if params.ParentID != nil {
		fileItem.ParentID = params.ParentID
	}

	// Set file path using parent's path
	fileItem.FilePath = parentPath

	if *params.IsDirectory {
		// Process directory creation
		name := *params.Name
		if !strings.HasSuffix(name, "/") {
			name += "/"
		}
		name = normalizeFileName(name)

		// Check if directory already exists and create a unique name if needed
		fullPath, finalName := u.resolveNameConflict(parentPath, name, "", serverKey, bucketName)

		fileItem.Name = finalName
		fileItem.MimeType = "directory"
		fileItem.Size = 0

		// Create directory in storage
		_, err := u.storageService.CreateFileOrDirectoryIfNotExists(
			serverKey,
			bucketName,
			fullPath,
			toServicePermissionEnum(*params.Permission),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating directory in storage: %v", err)
		}
	} else {
		// Process file upload
		file := params.File
		if file == nil {
			return nil, fmt.Errorf("file is required for non-directory items")
		}

		// Normalize file name
		rawFileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		normalizedName := normalizeFileName(rawFileName)
		extension := filepath.Ext(file.Filename)

		// Check if file already exists and create a unique name if needed
		fullPath, finalName := u.resolveNameConflict(parentPath, normalizedName, extension, serverKey, bucketName)

		fileItem.Name = finalName
		fileItem.MimeType = file.Header.Get("Content-Type")
		if fileItem.MimeType == "" {
			fileItem.MimeType = "application/octet-stream"
		}
		fileItem.Size = int64(file.Size) / 1024 // Convert bytes to KB

		// Open the file for reading
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening uploaded file: %v", err)
		}
		defer src.Close()

		// Create file in storage
		_, err = u.storageService.CreateFileOrDirectoryIfNotExists(
			serverKey,
			bucketName,
			fullPath,
			toServicePermissionEnum(*params.Permission),
			src,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating file in storage: %v", err)
		}

		// Update storage usage
		if err := u.storageRepo.SetIncreaseUsedSpaceKb(storage.ID, fileItem.Size); err != nil {
			return nil, fmt.Errorf("error updating storage usage: %v", err)
		}

		// Update parent directory size if applicable
		if params.ParentID != nil && *params.ParentID > 0 {
			if err := u.FileItemRepo.UpdateSize(*params.ParentID, fileItem.Size); err != nil {
				return nil, fmt.Errorf("error updating parent directory size: %v", err)
			}
		}
	}

	// Save file item to database
	if err := u.FileItemRepo.Create(&fileItem); err != nil {
		return nil, fmt.Errorf("error saving file item: %v", err)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"fileItem": fileItem,
	}, ""), nil
}

// DeleteFileItemCommand marks a file item as deleted
func (u *FileItemUsecase) DeleteFileItemCommand(params *fileitem.DeleteFileItemCommand) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Check if file exists
	fileItem, err := u.FileItemRepo.GetByID(*params.ID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Check if user has access
	if fileItem.UserID != *userID {
		return nil, fmt.Errorf("access denied")
	}

	// Mark as deleted
	if err := u.FileItemRepo.SetDelete(*params.ID); err != nil {
		return nil, fmt.Errorf("error deleting file: %v", err)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"id": params.ID,
	}, ""), nil
}

// ForceDeleteFileItemCommand permanently deletes a file item
func (u *FileItemUsecase) ForceDeleteFileItemCommand(params *fileitem.ForceDeleteFileItemCommand) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Check if file exists
	fileItem, err := u.FileItemRepo.GetByID(*params.ID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Check if user has access
	if fileItem.UserID != *userID {
		return nil, fmt.Errorf("access denied")
	}

	// Delete from storage
	fullPath := getFullPath(fileItem)
	success, err := u.storageService.RemoveFileOrDirectoryIfExists(fileItem.ServerKey, fileItem.BucketName, fullPath)
	if err != nil {
		return nil, fmt.Errorf("error removing from storage: %v", err)
	}

	if !success {
		u.Logger.Warn("File not found in storage but will be removed from database", map[string]interface{}{
			"fileId": fileItem.ID,
			"path":   fullPath,
		})
	}

	// Delete from database (this will also recursively delete children if it's a directory)
	if err := u.FileItemRepo.ForceDelete(*params.ID); err != nil {
		return nil, fmt.Errorf("error permanently deleting file: %v", err)
	}

	// Decrease storage usage
	if fileItem.Size > 0 {
		storage, err := u.storageRepo.GetByUserID(*userID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving storage: %v", err)
		}

		if err := u.storageRepo.SetIncreaseUsedSpaceKb(storage.ID, -fileItem.Size); err != nil {
			return nil, fmt.Errorf("error updating storage usage: %v", err)
		}
	}

	// Update parent directory size if applicable
	if fileItem.ParentID != nil && *fileItem.ParentID > 0 {
		if err := u.FileItemRepo.UpdateSize(*fileItem.ParentID, -fileItem.Size); err != nil {
			return nil, fmt.Errorf("error updating parent directory size: %v", err)
		}
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"id": params.ID,
	}, ""), nil
}

// RestoreFileItemCommand restores a deleted file item
func (u *FileItemUsecase) RestoreFileItemCommand(params *fileitem.RestoreFileItemCommand) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Restore the file
	result := u.FileItemRepo.SetRestore(*params.ID)
	if result != nil {
		return nil, fmt.Errorf("error restoring file: %v", result)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"id": params.ID,
	}, ""), nil
}

// UpdateFileItemCommand updates file item properties
func (u *FileItemUsecase) UpdateFileItemCommand(params *fileitem.UpdateFileItemCommand) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Check if file exists
	fileItem, err := u.FileItemRepo.GetByID(*params.ID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Check if user has access
	if fileItem.UserID != *userID {
		return nil, fmt.Errorf("access denied")
	}

	// Update permission if requested
	if *params.IsChangePermission && params.Permission != nil {
		newPermission := *params.Permission
		fileItem.Permission = newPermission

		// Update permission in storage
		fullPath := getFullPath(fileItem)
		if err := u.storageService.AddOrRemoveOrChangePermission(
			fileItem.ServerKey,
			fileItem.BucketName,
			fullPath,
			toServicePermissionEnum(*params.Permission),
			false); err != nil {
			return nil, fmt.Errorf("error updating permission in storage: %v", err)
		}
	}

	// Save changes
	if err := u.FileItemRepo.Update(fileItem); err != nil {
		return nil, fmt.Errorf("error updating file: %v", err)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"fileItem": fileItem,
	}, ""), nil
}

// FileOperationCommand handles file operations like copy, move, rename
func (u *FileItemUsecase) FileOperationCommand(params *fileitem.FileOperationCommand) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Check if file exists
	fileItem, err := u.FileItemRepo.GetByID(*params.ID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Check if user has access
	if fileItem.UserID != *userID {
		return nil, fmt.Errorf("access denied")
	}

	// Get source path
	oldFullPath := getFullPath(fileItem)
	var newFullPath string

	switch params.OperationType {
	case enums.FileItemRenameOperation:
		if params.NewName == nil {
			return nil, fmt.Errorf("new name is required for rename operation")
		}

		// Get new name with proper formatting
		newName := *params.NewName
		if fileItem.IsDirectory && !strings.HasSuffix(newName, "/") {
			newName += "/"
		}

		// Rename in storage
		newPath := filepath.Join(filepath.Dir(oldFullPath), newName)
		newFullPath = newPath

		permission := string(fileItem.Permission)
		_, err := u.storageService.RenameOrMoveFileOrDirectory(
			fileItem.ServerKey,
			fileItem.BucketName,
			oldFullPath,
			newFullPath,
			toServicePermissionEnum(enums.FileItemPermissionEnum(permission)))
		if err != nil {
			return nil, fmt.Errorf("error renaming in storage: %v", err)
		}

		// Update in database
		fileItem.Name = newName
		if err := u.FileItemRepo.Update(fileItem); err != nil {
			return nil, fmt.Errorf("error updating file name: %v", err)
		}

		// If it's a directory, update all children paths
		if fileItem.IsDirectory {
			// This would require recursive updates to all children
			// Implementation depends on additional methods to update file paths
		}

	case enums.FileItemMoveOperation:
		if params.NewParentID == nil {
			return nil, fmt.Errorf("new parent ID is required for move operation")
		}

		// Get the parent directory
		var newParent *domain.FileItem
		var newParentPath string

		if *params.NewParentID == 0 {
			// Moving to root
			newParentPath = ""
		} else {
			newParent, err = u.FileItemRepo.GetByID(*params.NewParentID)
			if err != nil {
				return nil, fmt.Errorf("new parent directory not found: %v", err)
			}

			// Check if new parent is a directory
			if !newParent.IsDirectory {
				return nil, fmt.Errorf("destination must be a directory")
			}

			newParentPath = getFullPath(newParent)
		}

		// Set new path
		newFullPath = filepath.Join(newParentPath, fileItem.Name)

		// Move in storage
		permission := string(fileItem.Permission)
		_, err := u.storageService.RenameOrMoveFileOrDirectory(
			fileItem.ServerKey,
			fileItem.BucketName,
			oldFullPath,
			newFullPath,
			toServicePermissionEnum(enums.FileItemPermissionEnum(permission)))
		if err != nil {
			return nil, fmt.Errorf("error moving in storage: %v", err)
		}

		// Update file path and parent ID in database
		fileItem.FilePath = newParentPath
		if *params.NewParentID == 0 {
			fileItem.ParentID = nil
		} else {
			fileItem.ParentID = params.NewParentID
		}

		if err := u.FileItemRepo.Update(fileItem); err != nil {
			return nil, fmt.Errorf("error updating file: %v", err)
		}

		// Update size for old and new parent
		if fileItem.ParentID != nil && *fileItem.ParentID > 0 {
			if err := u.FileItemRepo.UpdateSize(*fileItem.ParentID, -fileItem.Size); err != nil {
				return nil, fmt.Errorf("error updating old parent size: %v", err)
			}
		}

		if *params.NewParentID > 0 {
			if err := u.FileItemRepo.UpdateSize(*params.NewParentID, fileItem.Size); err != nil {
				return nil, fmt.Errorf("error updating new parent size: %v", err)
			}
		}

		// If it's a directory, update all children paths
		if fileItem.IsDirectory {
			// This would require recursive updates to all children
			// Implementation depends on additional methods to update file paths
		}

	case enums.FileItemCopyOperation:
		if params.NewParentID == nil {
			return nil, fmt.Errorf("new parent ID is required for copy operation")
		}

		// Get storage to check quota
		storage, err := u.storageRepo.GetByUserID(*userID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving storage: %v", err)
		}

		// Check if storage has enough space
		exceeded, err := u.storageRepo.CheckQuotaExceeded(storage.ID, fileItem.Size*1024) // Convert KB to bytes
		if err != nil {
			return nil, fmt.Errorf("error checking quota: %v", err)
		}
		if exceeded {
			return nil, fmt.Errorf("storage quota exceeded")
		}

		// Get the parent directory
		var newParent *domain.FileItem
		var newParentPath string

		if *params.NewParentID == 0 {
			// Copying to root
			newParentPath = ""
		} else {
			newParent, err = u.FileItemRepo.GetByID(*params.NewParentID)
			if err != nil {
				return nil, fmt.Errorf("new parent directory not found: %v", err)
			}

			// Check if new parent is a directory
			if !newParent.IsDirectory {
				return nil, fmt.Errorf("destination must be a directory")
			}

			newParentPath = getFullPath(newParent)
		}

		// Set new path
		newFullPath = filepath.Join(newParentPath, fileItem.Name)

		// Copy in storage
		permission := string(fileItem.Permission)
		success, err := u.storageService.CopyFileOrDirectory(
			fileItem.ServerKey,
			fileItem.BucketName,
			oldFullPath,
			newFullPath,
			toServicePermissionEnum(enums.FileItemPermissionEnum(permission)))
		if err != nil {
			return nil, fmt.Errorf("error copying in storage: %v", err)
		}

		if !success {
			return nil, fmt.Errorf("failed to copy file in storage")
		}

		// Create new file item
		newFileItem := fileItem
		newFileItem.ID = 0 // Clear ID for new record
		newFileItem.FilePath = newParentPath
		if *params.NewParentID == 0 {
			newFileItem.ParentID = nil
		} else {
			newFileItem.ParentID = params.NewParentID
		}
		newFileItem.CreatedAt = time.Now()
		newFileItem.UpdatedAt = time.Now()

		if err := u.FileItemRepo.Create(newFileItem); err != nil {
			return nil, fmt.Errorf("error creating new file item: %v", err)
		}

		// Update new parent size
		if *params.NewParentID > 0 {
			if err := u.FileItemRepo.UpdateSize(*params.NewParentID, fileItem.Size); err != nil {
				return nil, fmt.Errorf("error updating new parent size: %v", err)
			}
		}

		// Update storage usage
		if err := u.storageRepo.SetIncreaseUsedSpaceKb(storage.ID, fileItem.Size); err != nil {
			return nil, fmt.Errorf("error updating storage usage: %v", err)
		}

		return resp.NewResponseData(resp.Success, resp.Data{
			"fileItem": newFileItem,
		}, ""), nil
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"fileItem": fileItem,
	}, ""), nil
}

// GetByIdsQuery retrieves file items by IDs
func (u *FileItemUsecase) GetByIdsQuery(params *fileitem.GetByIdsQuery) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Extract IDs from the request
	var ids []int64
	for _, item := range params.IdsOrder {
		if item.ID != nil {
			ids = append(ids, *item.ID)
		}
	}

	// Get file items
	fileItems, err := u.FileItemRepo.GetByIDs(ids)
	if err != nil {
		return nil, fmt.Errorf("error retrieving file items: %v", err)
	}

	// Create result DTOs with download links
	var result []map[string]interface{}
	for _, fileItem := range fileItems {
		// Check if user has access
		if fileItem.UserID != *userID {
			continue
		}

		// Generate link based on isTemporary flag
		var link string
		var err error

		if *params.IsTemporary {
			if params.ExpireMinutes == nil {
				params.ExpireMinutes = new(int)
				*params.ExpireMinutes = 5 // Default expiry
			}
			link, err = u.storageService.GeneratePreSignedUrl(
				fileItem.ServerKey,
				fileItem.BucketName,
				getFullPath(&fileItem),
				time.Duration(*params.ExpireMinutes)*time.Minute)
			if err != nil {
				return nil, fmt.Errorf("error generating presigned URL: %v", err)
			}
		} else {
			link = u.storageService.GenerateUrl(
				fileItem.ServerKey,
				fileItem.BucketName,
				getFullPath(&fileItem))
		}

		// Find the requested order
		var order int
		for _, item := range params.IdsOrder {
			if item.ID != nil && *item.ID == fileItem.ID && item.Order != nil {
				order = *item.Order
				break
			}
		}

		// Create DTO
		fileItemDTO := map[string]interface{}{
			"id":          fileItem.ID,
			"name":        fileItem.Name,
			"filePath":    fileItem.FilePath,
			"isDirectory": fileItem.IsDirectory,
			"size":        fileItem.Size,
			"mimeType":    fileItem.MimeType,
			"parentId":    fileItem.ParentID,
			"permission":  fileItem.Permission,
			"userId":      fileItem.UserID,
			"createdAt":   fileItem.CreatedAt,
			"updatedAt":   fileItem.UpdatedAt,
			"link":        link,
			"order":       order,
		}

		result = append(result, fileItemDTO)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"items": result,
	}, ""), nil
}

// GetDeletedTreeDirectoryQuery retrieves deleted file items
func (u *FileItemUsecase) GetDeletedTreeDirectoryQuery(params *fileitem.GetDeletedTreeDirectoryQuery) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Get deleted items
	items, err := u.FileItemRepo.GetDeletedItems(*userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving deleted items: %v", err)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"items": items,
	}, ""), nil
}

// GetDownloadFileItemByIdQuery retrieves a file for download
func (u *FileItemUsecase) GetDownloadFileItemByIdQuery(params *fileitem.GetDownloadFileItemByIdQuery) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Check if file exists
	fileItem, err := u.FileItemRepo.GetByID(*params.ID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Check if user has access
	if fileItem.UserID != *userID {
		return nil, fmt.Errorf("access denied")
	}

	// Check if it's a directory (can't download directories)
	if fileItem.IsDirectory {
		return nil, fmt.Errorf("cannot download a directory")
	}

	// Get file from storage
	stream, err := u.storageService.DownloadFileOrDirectory(
		fileItem.ServerKey,
		fileItem.BucketName,
		getFullPath(fileItem))
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %v", err)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"fileItem":       fileItem,
		"downloadStream": stream,
	}, ""), nil
}

// GetTreeDirectoryQuery retrieves a directory tree
func (u *FileItemUsecase) GetTreeDirectoryQuery(params *fileitem.GetTreeDirectoryQuery) (*resp.Response, error) {
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت کاربر")
	}
	// Get tree
	items, err := u.FileItemRepo.GetTreeByUserIDAndParentID(*userID, params.ParentFileItemID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving directory tree: %v", err)
	}

	return resp.NewResponseData(resp.Success, resp.Data{
		"items": items,
	}, ""), nil
}

// Helper functions

// normalizeFileName removes invalid characters from a filename
func normalizeFileName(name string) string {
	reg := regexp.MustCompile(`[^\w\s/]`)
	name = reg.ReplaceAllString(name, "")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

// resolveNameConflict ensures a unique filename by adding a counter suffix if needed
func (u *FileItemUsecase) resolveNameConflict(parentPath, baseName, suffix, serverKey, bucketName string) (string, string) {
	fullPath := parentPath + baseName + suffix
	finalName := baseName + suffix
	counter := 1

	for {
		exists, _ := u.storageService.CheckFileOrDirectoryIsExists(serverKey, bucketName, fullPath)
		if !exists {
			break
		}

		newBaseName := fmt.Sprintf("%s_%d", baseName, counter)
		fullPath = parentPath + newBaseName + suffix
		finalName = newBaseName + suffix
		counter++
	}

	return fullPath, finalName
}

// getFullPath returns the full path for a file item
func getFullPath(fileItem *domain.FileItem) string {
	if fileItem.FilePath == "" {
		return fileItem.Name
	}
	return fileItem.FilePath + fileItem.Name
}
