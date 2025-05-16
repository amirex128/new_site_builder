package repository

import (
	"go-boilerplate/src/internal/domain/sample/entity"
)

type ISampleRepository interface {
	GetAllSample() []entity.SimpleEntity
}
