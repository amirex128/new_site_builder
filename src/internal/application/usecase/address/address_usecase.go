package addressusecase

import (
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/address"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type AddressUsecase struct {
	logger       sflogger.Logger
	addressRepo  repository.IAddressRepository
	cityRepo     repository.ICityRepository
	provinceRepo repository.IProvinceRepository
}

func NewAddressUsecase(c contract.IContainer) *AddressUsecase {
	return &AddressUsecase{
		logger:       c.GetLogger(),
		addressRepo:  c.GetAddressRepo(),
		cityRepo:     c.GetCityRepo(),
		provinceRepo: c.GetProvinceRepo(),
	}
}

func (u *AddressUsecase) CreateAddressCommand(params *address.CreateAddressCommand) (any, error) {
	// Implementation for creating an address
	fmt.Println(params)

	var customerID, userID int64
	if params.CustomerID != nil {
		customerID = *params.CustomerID
	}
	if params.UserID != nil {
		userID = *params.UserID
	}

	// Create copies of the float values since we need pointers
	latitude := *params.Latitude
	longitude := *params.Longitude

	newAddress := domain.Address{
		CustomerID:  customerID,
		UserID:      userID,
		Title:       *params.Title,
		Latitude:    &latitude,
		Longitude:   &longitude,
		AddressLine: *params.AddressLine,
		PostalCode:  *params.PostalCode,
		CityID:      *params.CityID,
		ProvinceID:  *params.ProvinceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	err := u.addressRepo.Create(newAddress)
	if err != nil {
		return nil, err
	}

	return newAddress, nil
}

func (u *AddressUsecase) UpdateAddressCommand(params *address.UpdateAddressCommand) (any, error) {
	// Implementation for updating an address
	fmt.Println(params)

	existingAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	if params.Title != nil {
		existingAddress.Title = *params.Title
	}

	if params.Latitude != nil {
		latitude := *params.Latitude
		existingAddress.Latitude = &latitude
	}

	if params.Longitude != nil {
		longitude := *params.Longitude
		existingAddress.Longitude = &longitude
	}

	if params.AddressLine != nil {
		existingAddress.AddressLine = *params.AddressLine
	}

	if params.PostalCode != nil {
		existingAddress.PostalCode = *params.PostalCode
	}

	existingAddress.CityID = *params.CityID
	existingAddress.ProvinceID = *params.ProvinceID
	existingAddress.UpdatedAt = time.Now()

	err = u.addressRepo.Update(existingAddress)
	if err != nil {
		return nil, err
	}

	return existingAddress, nil
}

func (u *AddressUsecase) DeleteAddressCommand(params *address.DeleteAddressCommand) (any, error) {
	// Implementation for deleting an address
	fmt.Println(params)

	err := u.addressRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *AddressUsecase) GetByIdAddressQuery(params *address.GetByIdAddressQuery) (any, error) {
	// Implementation to get address by ID
	fmt.Println(params)

	result, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *AddressUsecase) GetAllAddressQuery(params *address.GetAllAddressQuery) (any, error) {
	// Implementation to get all addresses
	fmt.Println(params)

	result, count, err := u.addressRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *AddressUsecase) GetAllCityQuery(params *address.GetAllCityQuery) (any, error) {
	// Implementation to get all cities
	fmt.Println(params)

	result, count, err := u.cityRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *AddressUsecase) GetAllProvinceQuery(params *address.GetAllProvinceQuery) (any, error) {
	// Implementation to get all provinces
	fmt.Println(params)

	result, count, err := u.provinceRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *AddressUsecase) AdminGetAllAddressQuery(params *address.AdminGetAllAddressQuery) (any, error) {
	// Implementation for admin to get all addresses
	fmt.Println(params)

	result, count, err := u.addressRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
