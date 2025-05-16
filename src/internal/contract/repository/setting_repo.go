package repository

import (
	"go-boilerplate/src/internal/domain"
)

type ISettingRepository interface {
	GetBySiteID(siteID int64) (domain.Setting, error)
	Create(setting domain.Setting) error
	Update(setting domain.Setting) error
}
