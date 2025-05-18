package mysql

import (
	"time"

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

	query := r.database.Model(&domain.FileItem{}).Where("is_deleted = ?", false)
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

	query := r.database.Model(&domain.FileItem{}).Where("user_id = ? AND is_deleted = ?", userID, false)
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

	query := r.database.Model(&domain.FileItem{}).Where("parent_id = ? AND is_deleted = ?", parentID, false)
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
	result := r.database.Where("is_deleted = ?", false).First(&fileItem, id)
	if result.Error != nil {
		return fileItem, result.Error
	}
	return fileItem, nil
}

func (r *FileItemRepo) GetByIDs(ids []int64) ([]domain.FileItem, error) {
	var fileItems []domain.FileItem
	result := r.database.Where("id IN ? AND is_deleted = ?", ids, false).Find(&fileItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return fileItems, nil
}

func (r *FileItemRepo) GetTreeByUserIDAndParentID(userID int64, parentID *int64) ([]domain.FileItem, error) {
	var fileItems []domain.FileItem

	query := r.database.Model(&domain.FileItem{}).Where("user_id = ? AND is_deleted = ?", userID, false)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	result := query.Find(&fileItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return fileItems, nil
}

func (r *FileItemRepo) GetDeletedItems(userID int64) ([]domain.FileItem, error) {
	var fileItems []domain.FileItem
	result := r.database.Where("user_id = ? AND is_deleted = ?", userID, true).Find(&fileItems)
	if result.Error != nil {
		return nil, result.Error
	}
	return fileItems, nil
}

func (r *FileItemRepo) SetDelete(id int64) error {
	now := time.Now()
	result := r.database.Model(&domain.FileItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": now,
	})
	return result.Error
}

func (r *FileItemRepo) SetRestore(id int64) error {
	result := r.database.Model(&domain.FileItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": false,
		"deleted_at": nil,
	})
	return result.Error
}

func (r *FileItemRepo) ForceDelete(id int64) error {
	// First get the file item to check if it's a directory
	var fileItem domain.FileItem
	if err := r.database.First(&fileItem, id).Error; err != nil {
		return err
	}

	// If it's a directory, also delete all its children recursively
	if fileItem.IsDirectory {
		// Get all children
		var children []domain.FileItem
		if err := r.database.Where("parent_id = ?", id).Find(&children).Error; err != nil {
			return err
		}

		// Delete each child recursively
		for _, child := range children {
			if err := r.ForceDelete(child.ID); err != nil {
				return err
			}
		}
	}

	// Delete the file item itself
	return r.database.Delete(&domain.FileItem{}, id).Error
}

func (r *FileItemRepo) Create(fileItem domain.FileItem) error {
	result := r.database.Create(&fileItem)
	return result.Error
}

func (r *FileItemRepo) Update(fileItem domain.FileItem) error {
	result := r.database.Save(&fileItem)
	return result.Error
}

func (r *FileItemRepo) UpdateFilePath(id int64, filePath string) error {
	result := r.database.Model(&domain.FileItem{}).Where("id = ?", id).Update("file_path", filePath)
	return result.Error
}

func (r *FileItemRepo) UpdateSize(id int64, sizeChange int64) error {
	// Use a raw SQL query to increment/decrement the size
	result := r.database.Exec("UPDATE file_items SET size = size + ? WHERE id = ?", sizeChange, id)
	return result.Error
}

func (r *FileItemRepo) UpdateParentID(id int64, parentID *int64) error {
	if parentID == nil {
		result := r.database.Model(&domain.FileItem{}).Where("id = ?", id).Update("parent_id", nil)
		return result.Error
	} else {
		result := r.database.Model(&domain.FileItem{}).Where("id = ?", id).Update("parent_id", *parentID)
		return result.Error
	}
}

func (r *FileItemRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.FileItem{}, id)
	return result.Error
}
