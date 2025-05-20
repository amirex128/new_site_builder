package fileitem

import (
	"mime/multipart"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateOrDirectoryItemCommand represents a command to create a file or directory
type CreateOrDirectoryItemCommand struct {
	File        *multipart.FileHeader         `json:"file,omitempty" validate:"required_if=IsDirectory false"`
	Name        *string                       `json:"name,omitempty" validate:"required_if=IsDirectory true,optional_text=0 255"`
	IsDirectory *bool                         `json:"isDirectory" validate:"required_bool"`
	Permission  *enums.FileItemPermissionEnum `json:"permission" validate:"required,enum"`
	ParentID    *int64                        `json:"parentId,omitempty" validate:"omitempty"`
}

// DeleteFileItemCommand represents a command to delete a file item
type DeleteFileItemCommand struct {
	ID *int64 `json:"id" validate:"required"`
}

// ForceDeleteFileItemCommand represents a command to permanently delete a file item
type ForceDeleteFileItemCommand struct {
	ID *int64 `json:"id" validate:"required"`
}

// RestoreFileItemCommand represents a command to restore a deleted file item
type RestoreFileItemCommand struct {
	ID *int64 `json:"id" validate:"required"`
}

// FileOperationCommand represents a command for file operations like copy, move, or rename
type FileOperationCommand struct {
	ID            *int64              `json:"id" validate:"required"`
	OperationType enums.OperationType `json:"operationType" validate:"required,enum"`
	NewName       *string             `json:"newName,omitempty" validate:"required_if=OperationType 2,optional_text=0 200"`
	NewParentID   *int64              `json:"newParentId,omitempty" validate:"required_if=OperationType 0 OperationType 1"`
}

// UpdateFileItemCommand represents a command to update a file item
type UpdateFileItemCommand struct {
	ID                 *int64                        `json:"id" validate:"required"`
	IsChangePermission *bool                         `json:"isChangePermission" validate:"required_bool"`
	Permission         *enums.FileItemPermissionEnum `json:"permission,omitempty" validate:"required_if=IsChangePermission true,enum_optional"`
}
