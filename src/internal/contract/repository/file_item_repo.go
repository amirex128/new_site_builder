package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IFileItemRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.FileItem], error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.FileItem], error)
	GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.FileItem], error)
	GetByID(id int64) (domain.FileItem, error)
	GetByIDs(ids []int64) ([]domain.FileItem, error)
	GetTreeByUserIDAndParentID(userID int64, parentID *int64) ([]domain.FileItem, error)
	GetDeletedItems(userID int64) ([]domain.FileItem, error)
	SetDelete(id int64) error
	SetRestore(id int64) error
	ForceDelete(id int64) error
	Create(fileItem domain.FileItem) error
	Update(fileItem domain.FileItem) error
	UpdateFilePath(id int64, filePath string) error
	UpdateSize(id int64, sizeChange int64) error
	UpdateParentID(id int64, parentID *int64) error
	Delete(id int64) error
}
