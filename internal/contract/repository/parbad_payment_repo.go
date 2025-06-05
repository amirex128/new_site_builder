package repository

import (
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IParbadPaymentRepository interface {
	GetByTrackingNumber(trackingNumber int64) (*domain.ParbadPayment, error)
	GetByID(id int64) (*domain.ParbadPayment, error)
	Create(payment *domain.ParbadPayment) error
	Update(payment *domain.ParbadPayment) error
}
