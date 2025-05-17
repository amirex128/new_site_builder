package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IStorageRepository interface {
	GetByUserID(userID int64) (domain.Storage, error)
	GetByID(id int64) (domain.Storage, error)
	Create(storage domain.Storage) error
	Update(storage domain.Storage) error
	Delete(id int64) error
}
