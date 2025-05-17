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

	// TODO: Implement restore logic
	// This might require custom repository methods beyond the standard interface

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *FileItemUsecase) GetDeletedTreeDirectoryQuery(params *fileitem.GetDeletedTreeDirectoryQuery) (any, error) {
	// Implementation to get deleted tree directory
	fmt.Println(params)

	// TODO: Implement deleted tree retrieval logic
	// This might require custom repository methods beyond the standard interface

	return []interface{}{}, nil
}

func (u *FileItemUsecase) GetDownloadFileItemByIdQuery(params *fileitem.GetDownloadFileItemByIdQuery) (any, error) {
	// Implementation to download a file item by ID
	fmt.Println(params)

	fileItem, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Handle file download logic

	return fileItem, nil
}

func (u *FileItemUsecase) GetTreeDirectoryQuery(params *fileitem.GetTreeDirectoryQuery) (any, error) {
	// Implementation to get tree directory
	fmt.Println(params)

	// Empty pagination for now, can be enhanced later
	pagination := common.PaginationRequestDto{}

	// If parentID is provided, get children of that parent
	if params.ParentFileItemID != nil {
		// TODO: Implement proper tree fetching logic
		result, count, err := u.repo.GetAllByParentID(*params.ParentFileItemID, pagination)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"items": result,
			"total": count,
		}, nil
	}

	// Get root directories otherwise
	// TODO: Implement proper root directories fetching logic
	return []interface{}{}, nil
}
