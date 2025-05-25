package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type ParbadTransactionRepo struct {
	database *gorm.DB
}

func NewParbadTransactionRepository(db *gorm.DB) *ParbadTransactionRepo {
	return &ParbadTransactionRepo{
		database: db,
	}
}

func (r *ParbadTransactionRepo) GetByPaymentID(paymentID int64) ([]domain.ParbadTransaction, error) {
	var transactions []domain.ParbadTransaction
	result := r.database.Where("payment_id = ?", paymentID).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (r *ParbadTransactionRepo) GetByID(id int64) (*domain.ParbadTransaction, error) {
	var transaction domain.ParbadTransaction
	result := r.database.First(&transaction, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func (r *ParbadTransactionRepo) Create(transaction *domain.ParbadTransaction) error {
	result := r.database.Create(transaction)
	return result.Error
}

func (r *ParbadTransactionRepo) Update(transaction *domain.ParbadTransaction) error {
	result := r.database.Save(transaction)
	return result.Error
}
