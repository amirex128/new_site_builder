package mysql

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type SettingRepo struct {
	database *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepo {
	return &SettingRepo{
		database: db,
	}
}

func (r *SettingRepo) GetBySiteID(siteID int64) (domain.Setting, error) {
	var setting domain.Setting
	result := r.database.Where("site_id = ?", siteID).First(&setting)
	if result.Error != nil {
		return setting, result.Error
	}
	return setting, nil
}

func (r *SettingRepo) Create(setting domain.Setting) error {
	result := r.database.Create(&setting)
	return result.Error
}

func (r *SettingRepo) Update(setting domain.Setting) error {
	result := r.database.Save(&setting)
	return result.Error
}
