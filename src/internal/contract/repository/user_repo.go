package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IUserRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.User, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.User, int64, error)
	GetByID(id int64) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	GetByPhone(phone string) (domain.User, error)
	Create(user domain.User) error
	Update(user domain.User) error
	Delete(id int64) error
}
