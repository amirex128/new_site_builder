package mysql

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type StorageRepo struct {
	database *gorm.DB
}

func NewStorageRepository(db *gorm.DB) *StorageRepo {
	return &StorageRepo{
		database: db,
	}
}

func (r *StorageRepo) GetByUserID(userID int64) (domain.Storage, error) {
	var storage domain.Storage
	result := r.database.Where("user_id = ?", userID).First(&storage)
	if result.Error != nil {
		return storage, result.Error
	}
	return storage, nil
}

func (r *StorageRepo) GetByID(id int64) (domain.Storage, error) {
	var storage domain.Storage
	result := r.database.First(&storage, id)
	if result.Error != nil {
		return storage, result.Error
	}
	return storage, nil
}

func (r *StorageRepo) Create(storage domain.Storage) error {
	result := r.database.Create(&storage)
	return result.Error
}

func (r *StorageRepo) Update(storage domain.Storage) error {
	result := r.database.Save(&storage)
	return result.Error
}

func (r *StorageRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Storage{}, id)
	return result.Error
}
