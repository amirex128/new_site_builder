package repository

import (
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IConfigurationRepository interface {
	GetAll() ([]domain.Configuration, error)
	GetByKey(key string) (*domain.Configuration, error)
	SetKey(key string, value string) error
}
