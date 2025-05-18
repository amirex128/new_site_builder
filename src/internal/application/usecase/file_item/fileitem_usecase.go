package fileitemusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/fileitem"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type FileItemUsecase struct {
	logger sflogger.Logger
	repo   repository.IFileItemRepository
}

func NewFileItemUsecase(c contract.IContainer) *FileItemUsecase {
	return &FileItemUsecase{
		logger: c.GetLogger(),
		repo:   c.GetFileItemRepo(),
	}
}

func (u *FileItemUsecase) CreateOrDirectoryItemCommand(params *fileitem.CreateOrDirectoryItemCommand) (any, error) {
	// Implementation for creating a file or directory
	fmt.Println(params)

	// This is a placeholder implementation
	var name string
	if params.Name != nil {
		name = *params.Name
	} else {
		name = "Unnamed"
	}

	fileItem := domain.FileItem{
		Name:        name,
		IsDirectory: *params.IsDirectory,
		Permission:  strconv.Itoa(int(*params.Permission)),
		BucketName:  "default-bucket",           // Placeholder
		ServerKey:   "default-key",              // Placeholder
		FilePath:    "/",                        // Placeholder
		Size:        0,                          // To be set with actual file size
		MimeType:    "application/octet-stream", // Default mime type
		UserID:      1,                          // Placeholder, should come from authenticated user
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	if params.ParentID != nil {
		fileItem.ParentID = params.ParentID
	}

	// TODO: Handle file upload if not a directory

	err := u.repo.Create(fileItem)
	if err != nil {
		return nil, err
	}

	return fileItem, nil
}

func (u *FileItemUsecase) DeleteFileItemCommand(params *fileitem.DeleteFileItemCommand) (any, error) {
	// Implementation for deleting a file item
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *FileItemUsecase) ForceDeleteFileItemCommand(params *fileitem.ForceDeleteFileItemCommand) (any, error) {
	// Implementation for force deleting a file item
	fmt.Println(params)

	// TODO: Implement permanent deletion logic
	// This might require custom repository methods beyond the standard interface

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *FileItemUsecase) UpdateFileItemCommand(params *fileitem.UpdateFileItemCommand) (any, error) {
	// Implementation for updating a file item
	fmt.Println(params)

	fileItem, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	if *params.IsChangePermission && params.Permission != nil {
		fileItem.Permission = strconv.Itoa(int(*params.Permission))
	}

	err = u.repo.Update(fileItem)
	if err != nil {
		return nil, err
	}

	return fileItem, nil
}

func (u *FileItemUsecase) FileOperationCommand(params *fileitem.FileOperationCommand) (any, error) {
	// Implementation for file operations (copy, move, rename)
	fmt.Println(params)

	fileItem, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	switch params.OperationType {
	case fileitem.Move:
		fileItem.ParentID = params.NewParentID
	case fileitem.Copy:
		// TODO: Implement copy logic
		// This might require custom repository methods beyond the standard interface
	case fileitem.Rename:
		if params.NewName != nil {
			fileItem.Name = *params.NewName
		}
	}

	err = u.repo.Update(fileItem)
	if err != nil {
		return nil, err
	}

	return fileItem, nil
}

func (u *FileItemUsecase) RestoreFileItemCommand(params *fileitem.RestoreFileItemCommand) (any, error) {
	// Implementation for restoring a deleted file item
	fmt.Println(params)

	// TODO: Implement actual restore logic
	// This would typically involve getting the file item marked as deleted and removing the deletion flag

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *FileItemUsecase) GetTreeDirectoryQuery(params *fileitem.GetTreeDirectoryQuery) (any, error) {
	// Implementation to get a directory tree
	fmt.Println(params)

	var parentID int64 = 0
	if params.ParentFileItemID != nil {
		parentID = *params.ParentFileItemID
	}

	// Using a basic pagination request for demonstration
	paginationRequest := common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	}

	result, _, err := u.repo.GetAllByParentID(parentID, paginationRequest)
	if err != nil {
		return nil, err
	}

	// TODO: In a real implementation, you would recursively fetch all child directories and files
	// to build a complete tree structure

	return map[string]interface{}{
		"items": result,
	}, nil
}

func (u *FileItemUsecase) GetDeletedTreeDirectoryQuery(params *fileitem.GetDeletedTreeDirectoryQuery) (any, error) {
	// Implementation to get a tree of deleted directories
	fmt.Println(params)

	// TODO: Implement actual logic to retrieve deleted items
	// This would typically involve querying for items with IsDeleted = true

	return map[string]interface{}{
		"items": []domain.FileItem{}, // Placeholder empty array
	}, nil
}

func (u *FileItemUsecase) GetDownloadFileItemByIdQuery(params *fileitem.GetDownloadFileItemByIdQuery) (any, error) {
	// Implementation to download a file by ID
	fmt.Println(params)

	fileItem, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if item is a directory (can't download directories)
	if fileItem.IsDirectory {
		return nil, fmt.Errorf("cannot download a directory")
	}

	// TODO: Implement actual file download logic
	// This would typically involve generating a download URL or stream from your storage system

	return map[string]interface{}{
		"fileItem":     fileItem,
		"downloadUrl":  "https://example.com/download/" + strconv.FormatInt(fileItem.ID, 10), // Placeholder URL
		"downloadName": fileItem.Name,
	}, nil
}
