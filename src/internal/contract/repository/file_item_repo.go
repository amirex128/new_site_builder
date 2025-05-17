package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IFileItemRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.FileItem, int64, error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.FileItem, int64, error)
	GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.FileItem, int64, error)
	GetByID(id int64) (domain.FileItem, error)
	Create(fileItem domain.FileItem) error
	Update(fileItem domain.FileItem) error
	Delete(id int64) error
}
