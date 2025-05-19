package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	database *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{
		database: db,
	}
}

func (r *PermissionRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Permission, int64, error) {
	var permissions []domain.Permission
	var count int64

	query := r.database.Model(&domain.Permission{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&permissions)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return permissions, count, nil
}

func (r *PermissionRepo) GetByID(id int64) (domain.Permission, error) {
	var permission domain.Permission
	result := r.database.First(&permission, id)
	if result.Error != nil {
		return permission, result.Error
	}
	return permission, nil
}

func (r *PermissionRepo) GetByName(name string) (domain.Permission, error) {
	var permission domain.Permission
	result := r.database.Where("name = ?", name).First(&permission)
	if result.Error != nil {
		return permission, result.Error
	}
	return permission, nil
}

func (r *PermissionRepo) Create(permission domain.Permission) error {
	result := r.database.Create(&permission)
	return result.Error
}

func (r *PermissionRepo) Update(permission domain.Permission) error {
	result := r.database.Save(&permission)
	return result.Error
}

func (r *PermissionRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Permission{}, id)
	return result.Error
}
