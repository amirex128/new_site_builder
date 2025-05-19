package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	database *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{
		database: db,
	}
}

func (r *CustomerRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var count int64

	query := r.database.Model(&domain.Customer{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&customers)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return customers, count, nil
}

func (r *CustomerRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var count int64

	query := r.database.Model(&domain.Customer{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&customers)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return customers, count, nil
}

func (r *CustomerRepo) GetByID(id int64) (domain.Customer, error) {
	var customer domain.Customer
	result := r.database.First(&customer, id)
	if result.Error != nil {
		return customer, result.Error
	}
	return customer, nil
}

func (r *CustomerRepo) GetByEmail(email string) (domain.Customer, error) {
	var customer domain.Customer
	result := r.database.Where("email = ?", email).First(&customer)
	if result.Error != nil {
		return customer, result.Error
	}
	return customer, nil
}

func (r *CustomerRepo) GetByPhone(phone string) (domain.Customer, error) {
	var customer domain.Customer
	result := r.database.Where("phone = ?", phone).First(&customer)
	if result.Error != nil {
		return customer, result.Error
	}
	return customer, nil
}

func (r *CustomerRepo) Create(customer domain.Customer) error {
	result := r.database.Create(&customer)
	return result.Error
}

func (r *CustomerRepo) Update(customer domain.Customer) error {
	result := r.database.Save(&customer)
	return result.Error
}

func (r *CustomerRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Customer{}, id)
	return result.Error
}
