package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type AddressRepo struct {
	database *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepo {
	return &AddressRepo{
		database: db,
	}
}

func (r *AddressRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Address], error) {
	var addresses []domain.Address
	var count int64

	query := r.database.Model(&domain.Address{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&addresses)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(addresses, paginationRequestDto, count)
}

func (r *AddressRepo) GetAllByUserID(userID int64) ([]domain.Address, error) {
	var addresses []domain.Address
	var count int64

	query := r.database.Where("user_id = ?", userID).Model(&domain.Address{})
	query.Count(&count)

	result := query.Find(&addresses)
	if result.Error != nil {
		return nil, result.Error
	}

	return addresses, nil
}

func (r *AddressRepo) GetAllByCustomerID(customerID int64) ([]domain.Address, error) {
	var addresses []domain.Address
	var count int64

	query := r.database.Where("customer_id = ?", customerID).Model(&domain.Address{})
	query.Count(&count)

	result := query.Find(&addresses)
	if result.Error != nil {
		return nil, result.Error
	}

	return addresses, nil
}

func (r *AddressRepo) GetByID(id int64) (*domain.Address, error) {
	var address *domain.Address
	result := r.database.First(&address, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return address, nil
}

func (r *AddressRepo) Create(address *domain.Address) error {
	result := r.database.Create(address)
	return result.Error
}

func (r *AddressRepo) Update(address *domain.Address) error {
	result := r.database.Save(address)
	return result.Error
}

func (r *AddressRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Address{}, id)
	return result.Error
}

func (r *AddressRepo) AddAddressToUser(addressID int64, userID int64) error {
	// Update the user_id field directly in the database
	result := r.database.Model(&domain.Address{}).Where("id = ?", addressID).Update("user_id", userID)
	return result.Error
}

func (r *AddressRepo) RemoveAddressFromUser(addressID int64, userID int64) error {
	// Set user_id to NULL where matching both addressID and userID
	result := r.database.Model(&domain.Address{}).Where("id = ? AND user_id = ?", addressID, userID).Update("user_id", nil)
	return result.Error
}

func (r *AddressRepo) RemoveAllAddressesFromUser(userID int64) error {
	// Set user_id to NULL for all addresses with the given user_id
	result := r.database.Model(&domain.Address{}).Where("user_id = ?", userID).Update("user_id", nil)
	return result.Error
}

func (r *AddressRepo) AddAddressToCustomer(addressID int64, customerID int64) error {
	// Update the customer_id field directly in the database
	result := r.database.Model(&domain.Address{}).Where("id = ?", addressID).Update("customer_id", customerID)
	return result.Error
}

func (r *AddressRepo) RemoveAddressFromCustomer(addressID int64, customerID int64) error {
	// Set customer_id to NULL where matching both addressID and customerID
	result := r.database.Model(&domain.Address{}).Where("id = ? AND customer_id = ?", addressID, customerID).Update("customer_id", nil)
	return result.Error
}

func (r *AddressRepo) RemoveAllAddressesFromCustomer(customerID int64) error {
	// Set customer_id to NULL for all addresses with the given customer_id
	result := r.database.Model(&domain.Address{}).Where("customer_id = ?", customerID).Update("customer_id", nil)
	return result.Error
}
