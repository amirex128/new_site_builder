package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/sample/entity"
)

type ISampleRepository interface {
	GetAllSample() []entity.SimpleEntity
}
