package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type ParbadPaymentRepo struct {
	database *gorm.DB
}

func NewParbadPaymentRepository(db *gorm.DB) *ParbadPaymentRepo {
	return &ParbadPaymentRepo{
		database: db,
	}
}

func (r *ParbadPaymentRepo) GetByTrackingNumber(trackingNumber int64) (*domain.ParbadPayment, error) {
	var payment domain.ParbadPayment
	result := r.database.Where("tracking_number = ?", trackingNumber).First(&payment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (r *ParbadPaymentRepo) GetByID(id int64) (*domain.ParbadPayment, error) {
	var payment domain.ParbadPayment
	result := r.database.First(&payment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (r *ParbadPaymentRepo) Create(payment *domain.ParbadPayment) error {
	result := r.database.Create(payment)
	return result.Error
}

func (r *ParbadPaymentRepo) Update(payment *domain.ParbadPayment) error {
	result := r.database.Save(payment)
	return result.Error
}
