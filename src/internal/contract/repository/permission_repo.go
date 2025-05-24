package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPermissionRepository interface {
	GetAll() ([]domain.Permission, error)
	GetByID(id int64) (domain.Permission, error)
	GetByName(name string) (domain.Permission, error)
	Create(permission domain.Permission) error
	Update(permission domain.Permission) error
	Delete(id int64) error
}
