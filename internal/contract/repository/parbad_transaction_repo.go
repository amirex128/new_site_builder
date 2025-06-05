package repository

import (
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IParbadTransactionRepository interface {
	GetByPaymentID(paymentID int64) ([]domain.ParbadTransaction, error)
	GetByID(id int64) (*domain.ParbadTransaction, error)
	Create(transaction *domain.ParbadTransaction) error
	Update(transaction *domain.ParbadTransaction) error
}
