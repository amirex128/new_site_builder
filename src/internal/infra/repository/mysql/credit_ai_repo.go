package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type CreditRepo struct {
	database *gorm.DB
}

func NewCreditRepository(db *gorm.DB) *CreditRepo {
	return &CreditRepo{
		database: db,
	}
}

func (r *CreditRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Credit, int64, error) {
	var credits []domain.Credit
	var count int64

	query := r.database.Model(&domain.Credit{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&credits)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return credits, count, nil
}

func (r *CreditRepo) GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Credit, int64, error) {
	var credits []domain.Credit
	var count int64

	query := r.database.Model(&domain.Credit{}).Where("user_id = ?", userID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&credits)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return credits, count, nil
}

func (r *CreditRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Credit, int64, error) {
	var credits []domain.Credit
	var count int64

	query := r.database.Model(&domain.Credit{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&credits)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return credits, count, nil
}

func (r *CreditRepo) GetByID(id int64) (domain.Credit, error) {
	var credit domain.Credit
	result := r.database.First(&credit, id)
	if result.Error != nil {
		return credit, result.Error
	}
	return credit, nil
}

func (r *CreditRepo) Create(credit domain.Credit) error {
	result := r.database.Create(&credit)
	return result.Error
}

func (r *CreditRepo) Update(credit domain.Credit) error {
	result := r.database.Save(&credit)
	return result.Error
}

func (r *CreditRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Credit{}, id)
	return result.Error
}
