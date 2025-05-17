package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type FileItemRepo struct {
	database *gorm.DB
}

func NewFileItemRepository(db *gorm.DB) *FileItemRepo {
	return &FileItemRepo{
		database: db,
	}
}

func (r *FileItemRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.FileItem, int64, error) {
	var fileItems []domain.FileItem
	var count int64

	query := r.database.Model(&domain.FileItem{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&fileItems)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return fileItems, count, nil
}

func (r *FileItemRepo) GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.FileItem, int64, error) {
	var fileItems []domain.FileItem
	var count int64

	query := r.database.Model(&domain.FileItem{}).Where("user_id = ?", userID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&fileItems)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return fileItems, count, nil
}

func (r *FileItemRepo) GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.FileItem, int64, error) {
	var fileItems []domain.FileItem
	var count int64

	query := r.database.Model(&domain.FileItem{}).Where("parent_id = ?", parentID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&fileItems)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return fileItems, count, nil
}

func (r *FileItemRepo) GetByID(id int64) (domain.FileItem, error) {
	var fileItem domain.FileItem
	result := r.database.First(&fileItem, id)
	if result.Error != nil {
		return fileItem, result.Error
	}
	return fileItem, nil
}

func (r *FileItemRepo) Create(fileItem domain.FileItem) error {
	result := r.database.Create(&fileItem)
	return result.Error
}

func (r *FileItemRepo) Update(fileItem domain.FileItem) error {
	result := r.database.Save(&fileItem)
	return result.Error
}

func (r *FileItemRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.FileItem{}, id)
	return result.Error
}
