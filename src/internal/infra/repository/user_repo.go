package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type UserRepo struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{
		database: db,
	}
}

func (r *UserRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.User, int64, error) {
	var users []domain.User
	var count int64

	query := r.database.Model(&domain.User{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, count, nil
}

func (r *UserRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.User, int64, error) {
	var users []domain.User
	var count int64

	// Users related to sites might be through Site relation, adjust query as needed
	query := r.database.Model(&domain.User{}).
		Joins("JOIN sites ON sites.user_id = users.id").
		Where("sites.id = ?", siteID)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, count, nil
}

func (r *UserRepo) GetByID(id int64) (domain.User, error) {
	var user domain.User
	result := r.database.First(&user, id)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *UserRepo) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	result := r.database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *UserRepo) GetByPhone(phone string) (domain.User, error) {
	var user domain.User
	result := r.database.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *UserRepo) Create(user domain.User) error {
	result := r.database.Create(&user)
	return result.Error
}

func (r *UserRepo) Update(user domain.User) error {
	result := r.database.Save(&user)
	return result.Error
}

func (r *UserRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.User{}, id)
	return result.Error
}
