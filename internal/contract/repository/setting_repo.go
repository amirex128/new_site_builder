package repository

import (
	"github.com/amirex128/new_site_builder/internal/domain"
)

type ISettingRepository interface {
	GetBySiteID(siteID int64) (*domain.Setting, error)
	Create(setting *domain.Setting) error
	Update(setting *domain.Setting) error
}
