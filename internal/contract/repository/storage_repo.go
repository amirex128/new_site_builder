package repository

import (
	"github.com/amirex128/new_site_builder/internal/domain"
	"time"
)

type IStorageRepository interface {
	GetByUserID(userID int64) (*domain.Storage, error)
	GetByID(id int64) (*domain.Storage, error)
	Create(storage *domain.Storage) error
	Update(storage *domain.Storage) error
	Delete(id int64) error
	SetIncreaseUsedSpaceKb(id int64, sizeKb int64) error
	CheckQuotaExceeded(id int64, sizeBytes int64) (bool, error)
	CheckHasExpired(id int64) (bool, error)
	ChargeStorage(id int64, quotaKb int64, expireAt time.Time) error
}
