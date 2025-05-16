package mysql

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	database *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{
		database: db,
	}
}

func (r *PaymentRepo) GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetAllByOrderID(orderID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{}).Where("order_id = ?", orderID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetByID(id int64) (domain.Payment, error) {
	var payment domain.Payment
	result := r.database.First(&payment, id)
	if result.Error != nil {
		return payment, result.Error
	}
	return payment, nil
}

func (r *PaymentRepo) GetByTrackingNumber(trackingNumber string) (domain.Payment, error) {
	var payment domain.Payment
	result := r.database.Where("tracking_number = ?", trackingNumber).First(&payment)
	if result.Error != nil {
		return payment, result.Error
	}
	return payment, nil
}

func (r *PaymentRepo) Create(payment domain.Payment) error {
	result := r.database.Create(&payment)
	return result.Error
}

func (r *PaymentRepo) Update(payment domain.Payment) error {
	result := r.database.Save(&payment)
	return result.Error
}

func (r *PaymentRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Payment{}, id)
	return result.Error
}
