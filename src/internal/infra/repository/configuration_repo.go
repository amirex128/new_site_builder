package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type ConfigurationRepo struct {
	database *gorm.DB
}

func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepo {
	return &ConfigurationRepo{
		database: db,
	}
}

func (r *ConfigurationRepo) GetAll() ([]domain.Configuration, error) {
	var configurations []domain.Configuration
	result := r.database.Order("id DESC").Find(&configurations)
	if result.Error != nil {
		return nil, result.Error
	}
	return configurations, nil
}

func (r *ConfigurationRepo) GetByKey(key string) (*domain.Configuration, error) {
	var configuration domain.Configuration
	result := r.database.Where("`key` = ?", key).First(&configuration)
	if result.Error != nil {
		return nil, result.Error
	}
	return &configuration, nil
}

func (r *ConfigurationRepo) SetKey(key string, value string) error {
	configuration := domain.Configuration{
		Key:   key,
		Value: value,
	}
	result := r.database.Where("`key` = ?", key).FirstOrCreate(&configuration)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
