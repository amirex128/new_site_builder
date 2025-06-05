package repository

import (
	"errors"
	"github.com/amirex128/new_site_builder/internal/domain"
	"time"

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

func (r *StorageRepo) GetByUserID(userID int64) (*domain.Storage, error) {
	var storage domain.Storage
	result := r.database.Where("user_id = ?", userID).First(&storage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create a new storage record if none exists
			newStorage := domain.Storage{
				UsedSpaceKb: 0,
				QuotaKb:     0,
				ChargedAt:   time.Now(),
				ExpireAt:    time.Now(),
				UserID:      userID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				IsDeleted:   false,
			}
			if err := r.Create(&newStorage); err != nil {
				return nil, err
			}
			return &newStorage, nil
		}
		return nil, result.Error
	}
	return &storage, nil
}

func (r *StorageRepo) GetByID(id int64) (*domain.Storage, error) {
	var storage domain.Storage
	result := r.database.First(&storage, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &storage, nil
}

func (r *StorageRepo) Create(storage *domain.Storage) error {
	result := r.database.Create(storage)
	return result.Error
}

func (r *StorageRepo) Update(storage *domain.Storage) error {
	result := r.database.Save(storage)
	return result.Error
}

func (r *StorageRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Storage{}, id)
	return result.Error
}

func (r *StorageRepo) SetIncreaseUsedSpaceKb(id int64, sizeKb int64) error {
	result := r.database.Exec("UPDATE storages SET used_space_kb = used_space_kb + ? WHERE id = ?", sizeKb, id)
	return result.Error
}

func (r *StorageRepo) CheckQuotaExceeded(id int64, sizeBytes int64) (bool, error) {
	var storage domain.Storage
	if err := r.database.First(&storage, id).Error; err != nil {
		return true, err
	}

	// Convert bytes to KB for comparison
	sizeKb := sizeBytes / 1024

	// Check if usage would exceed quota
	return (storage.UsedSpaceKb + sizeKb) > storage.QuotaKb, nil
}

func (r *StorageRepo) CheckHasExpired(id int64) (bool, error) {
	var storage domain.Storage
	if err := r.database.First(&storage, id).Error; err != nil {
		return true, err
	}

	now := time.Now()
	return now.After(storage.ExpireAt), nil
}

func (r *StorageRepo) ChargeStorage(id int64, quotaKb int64, expireAt time.Time) error {
	result := r.database.Model(&domain.Storage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"quota_kb":   quotaKb,
		"charged_at": time.Now(),
		"expire_at":  expireAt,
	})
	return result.Error
}
