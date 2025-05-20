package fileitem

import (
	"mime/multipart"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateOrDirectoryItemCommand represents a command to create a file or directory
type CreateOrDirectoryItemCommand struct {
	File        *multipart.FileHeader         `json:"file,omitempty" nameFa:"فایل" validate:"required_if=IsDirectory false"`
	Name        *string                       `json:"name,omitempty" nameFa:"نام" validate:"required_if=IsDirectory true,optional_text=0 255"`
	IsDirectory *bool                         `json:"isDirectory" nameFa:"دایرکتوری" validate:"required_bool"`
	Permission  *enums.FileItemPermissionEnum `json:"permission" nameFa:"دسترسی" validate:"required,enum"`
	ParentID    *int64                        `json:"parentId,omitempty" nameFa:"شناسه والد" validate:"omitempty"`
}

// DeleteFileItemCommand represents a command to delete a file item
type DeleteFileItemCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required"`
}

// ForceDeleteFileItemCommand represents a command to permanently delete a file item
type ForceDeleteFileItemCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required"`
}

// RestoreFileItemCommand represents a command to restore a deleted file item
type RestoreFileItemCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required"`
}

// FileOperationCommand represents a command for file operations like copy, move, or rename
type FileOperationCommand struct {
	ID            *int64              `json:"id" nameFa:"شناسه" validate:"required"`
	OperationType enums.OperationType `json:"operationType" nameFa:"نوع عملیات" validate:"required,enum"`
	NewName       *string             `json:"newName,omitempty" nameFa:"نام جدید" validate:"required_if=OperationType 2,optional_text=0 200"`
	NewParentID   *int64              `json:"newParentId,omitempty" nameFa:"شناسه والد جدید" validate:"required_if=OperationType 0 OperationType 1"`
}

// UpdateFileItemCommand represents a command to update a file item
type UpdateFileItemCommand struct {
	ID                 *int64                        `json:"id" nameFa:"شناسه" validate:"required"`
	IsChangePermission *bool                         `json:"isChangePermission" nameFa:"تغییر دسترسی" validate:"required_bool"`
	Permission         *enums.FileItemPermissionEnum `json:"permission,omitempty" nameFa:"دسترسی" validate:"required_if=IsChangePermission true,enum_optional"`
}
