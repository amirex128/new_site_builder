package addressusecase

import (
	"errors"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/address"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type AddressUsecase struct {
	logger         sflogger.Logger
	repo           repository.IAddressRepository
	cityRepo       repository.ICityRepository
	provinceRepo   repository.IProvinceRepository
	authContextSvc common.IAuthContextService
}

func NewAddressUsecase(c contract.IContainer) *AddressUsecase {
	return &AddressUsecase{
		logger:         c.GetLogger(),
		repo:           c.GetAddressRepo(),
		cityRepo:       c.GetCityRepo(),
		provinceRepo:   c.GetProvinceRepo(),
		authContextSvc: c.GetAuthContextTransientService()(),
	}
}

// CreateAddressCommand handles the creation of a new address
func (u *AddressUsecase) CreateAddressCommand(params *address.CreateAddressCommand) (any, error) {
	u.logger.Info("CreateAddressCommand called", map[string]interface{}{
		"title": *params.Title,
	})

	var customerID, userID int64
	var err error

	// If customer ID is provided, use it
	if params.CustomerID != nil {
		customerID = *params.CustomerID
	} else {
		// Otherwise try to get from auth context
		customerID, err = u.authContextSvc.GetCustomerID()
		if err != nil {
			u.logger.Info("No customer ID in auth context, trying user ID", nil)
			// Not a customer, try as a user
			customerID = 0
		}
	}

	// If user ID is provided, use it
	if params.UserID != nil {
		userID = *params.UserID
	} else if customerID == 0 {
		// If no customer ID, try to get user ID from auth context
		userID, err = u.authContextSvc.GetUserID()
		if err != nil {
			return nil, errors.New("خطا در احراز هویت کاربر")
		}
	}

	// Validate city and province exist
	_, err = u.cityRepo.GetByID(*params.CityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("شهر مورد نظر یافت نشد")
		}
		return nil, err
	}

	_, err = u.provinceRepo.GetByID(*params.ProvinceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("استان مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Create new address
	newAddress := domain.Address{
		Title:       *params.Title,
		Latitude:    params.Latitude,
		Longitude:   params.Longitude,
		AddressLine: *params.AddressLine,
		PostalCode:  *params.PostalCode,
		CityID:      *params.CityID,
		ProvinceID:  *params.ProvinceID,
		UserID:      userID,
		CustomerID:  customerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	// Save to repository
	err = u.repo.Create(newAddress)
	if err != nil {
		u.logger.Error("Error creating address", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در ایجاد آدرس")
	}

	// If we have a customer or user ID, add the relationship
	if customerID > 0 {
		// This would typically involve updating a many-to-many relationship
		// For simplicity, we'll just log that this would happen
		u.logger.Info("Would add address to customer", map[string]interface{}{
			"addressId":  newAddress.ID,
			"customerId": customerID,
		})
	} else if userID > 0 {
		err = u.repo.AddAddressToUser(newAddress.ID, userID)
		if err != nil {
			u.logger.Error("Error adding address to user", map[string]interface{}{
				"error":     err.Error(),
				"addressId": newAddress.ID,
				"userId":    userID,
			})
			// Continue despite error - the address was created
		}
	}

	// Retrieve the address with relations to return
	fullAddress, err := u.repo.GetByID(newAddress.ID)
	if err != nil {
		return newAddress, nil // Return the basic address if can't retrieve with relations
	}

	return enhanceAddressResponse(fullAddress), nil
}

// UpdateAddressCommand handles updating an existing address
func (u *AddressUsecase) UpdateAddressCommand(params *address.UpdateAddressCommand) (any, error) {
	u.logger.Info("UpdateAddressCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing address
	existingAddress, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("آدرس مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check ownership
	customerID, _ := u.authContextSvc.GetCustomerID()
	userID, _ := u.authContextSvc.GetUserID()

	isAdmin, _ := u.authContextSvc.IsAdmin()

	// Check if user has access to this address
	if !isAdmin && existingAddress.CustomerID != customerID && existingAddress.UserID != userID {
		return nil, errors.New("شما دسترسی به این آدرس ندارید")
	}

	// Validate city and province exist
	_, err = u.cityRepo.GetByID(*params.CityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("شهر مورد نظر یافت نشد")
		}
		return nil, err
	}

	_, err = u.provinceRepo.GetByID(*params.ProvinceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("استان مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Update fields if provided
	if params.Title != nil {
		existingAddress.Title = *params.Title
	}

	if params.Latitude != nil {
		existingAddress.Latitude = params.Latitude
	}

	if params.Longitude != nil {
		existingAddress.Longitude = params.Longitude
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

	// Update in repository
	err = u.repo.Update(existingAddress)
	if err != nil {
		u.logger.Error("Error updating address", map[string]interface{}{
			"error": err.Error(),
			"id":    *params.ID,
		})
		return nil, errors.New("خطا در بروزرسانی آدرس")
	}

	// Retrieve the updated address with relations
	fullAddress, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return existingAddress, nil // Return the basic address if can't retrieve with relations
	}

	return enhanceAddressResponse(fullAddress), nil
}

// DeleteAddressCommand handles deleting an address
func (u *AddressUsecase) DeleteAddressCommand(params *address.DeleteAddressCommand) (any, error) {
	u.logger.Info("DeleteAddressCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing address to check ownership
	existingAddress, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("آدرس مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check ownership
	customerID, _ := u.authContextSvc.GetCustomerID()
	userID, _ := u.authContextSvc.GetUserID()

	isAdmin, _ := u.authContextSvc.IsAdmin()

	// Check if user has access to this address
	if !isAdmin && existingAddress.CustomerID != customerID && existingAddress.UserID != userID {
		return nil, errors.New("شما دسترسی به این آدرس ندارید")
	}

	// Delete address
	err = u.repo.Delete(*params.ID)
	if err != nil {
		u.logger.Error("Error deleting address", map[string]interface{}{
			"error": err.Error(),
			"id":    *params.ID,
		})
		return nil, errors.New("خطا در حذف آدرس")
	}

	return map[string]interface{}{
		"success": true,
		"message": "آدرس با موفقیت حذف شد",
	}, nil
}

// GetByIdAddressQuery handles retrieving an address by ID
func (u *AddressUsecase) GetByIdAddressQuery(params *address.GetByIdAddressQuery) (any, error) {
	u.logger.Info("GetByIdAddressQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get address
	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("آدرس مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check ownership
	customerID, _ := u.authContextSvc.GetCustomerID()
	userID, _ := u.authContextSvc.GetUserID()

	isAdmin, _ := u.authContextSvc.IsAdmin()

	// Check if user has access to this address
	if !isAdmin && result.CustomerID != customerID && result.UserID != userID {
		return nil, errors.New("شما دسترسی به این آدرس ندارید")
	}

	return enhanceAddressResponse(result), nil
}

// GetAllAddressQuery handles retrieving all addresses for the current user/customer
func (u *AddressUsecase) GetAllAddressQuery(params *address.GetAllAddressQuery) (any, error) {
	u.logger.Info("GetAllAddressQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	var results []domain.Address
	var count int64
	var err error

	// Try to get customer ID first
	customerID, customerErr := u.authContextSvc.GetCustomerID()
	if customerErr == nil && customerID > 0 {
		// Get addresses for customer
		results, err = u.repo.GetAllByCustomerID(customerID)
		if err != nil {
			return nil, errors.New("خطا در دریافت آدرس ها")
		}
		count = int64(len(results))
	} else {
		// Try user ID
		userID, userErr := u.authContextSvc.GetUserID()
		if userErr != nil {
			return nil, errors.New("خطا در احراز هویت کاربر")
		}

		// Get addresses for user
		results, err = u.repo.GetAllByUserID(userID)
		if err != nil {
			return nil, errors.New("خطا در دریافت آدرس ها")
		}
		count = int64(len(results))
	}

	// Enhance response with full details
	enhancedAddresses := make([]map[string]interface{}, 0, len(results))
	for _, addr := range results {
		enhancedAddresses = append(enhancedAddresses, enhanceAddressResponse(addr))
	}

	return map[string]interface{}{
		"items":     enhancedAddresses,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// AdminGetAllAddressQuery handles retrieving all addresses for admin
func (u *AddressUsecase) AdminGetAllAddressQuery(params *address.AdminGetAllAddressQuery) (any, error) {
	u.logger.Info("AdminGetAllAddressQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check if user is admin
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	if !isAdmin {
		return nil, errors.New("شما دسترسی به این عملیات را ندارید")
	}

	// Get all addresses with pagination
	results, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.logger.Error("Error getting all addresses", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت آدرس ها")
	}

	// Enhance response with full details
	enhancedAddresses := make([]map[string]interface{}, 0, len(results))
	for _, addr := range results {
		enhancedAddresses = append(enhancedAddresses, enhanceAddressResponse(addr))
	}

	return map[string]interface{}{
		"items":     enhancedAddresses,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// GetAllCityQuery handles retrieving all cities
func (u *AddressUsecase) GetAllCityQuery(params *address.GetAllCityQuery) (any, error) {
	u.logger.Info("GetAllCityQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get all cities with pagination
	results, count, err := u.cityRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.logger.Error("Error getting all cities", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت شهرها")
	}

	// Convert cities to simplified response format
	cities := make([]map[string]interface{}, 0, len(results))
	for _, city := range results {
		cityData := map[string]interface{}{
			"id":         city.ID,
			"name":       city.Name,
			"slug":       city.Slug,
			"status":     city.Status,
			"provinceId": city.ProvinceID,
		}

		// Add province info if available
		if city.Province != nil {
			cityData["province"] = map[string]interface{}{
				"id":   city.Province.ID,
				"name": city.Province.Name,
			}
		}

		cities = append(cities, cityData)
	}

	return map[string]interface{}{
		"items":     cities,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// GetAllProvinceQuery handles retrieving all provinces
func (u *AddressUsecase) GetAllProvinceQuery(params *address.GetAllProvinceQuery) (any, error) {
	u.logger.Info("GetAllProvinceQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get all provinces with pagination
	results, count, err := u.provinceRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.logger.Error("Error getting all provinces", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت استان ها")
	}

	// Convert provinces to simplified response format
	provinces := make([]map[string]interface{}, 0, len(results))
	for _, province := range results {
		provinceData := map[string]interface{}{
			"id":     province.ID,
			"name":   province.Name,
			"slug":   province.Slug,
			"status": province.Status,
		}

		provinces = append(provinces, provinceData)
	}

	return map[string]interface{}{
		"items":     provinces,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// Helper function to enhance address response with structured data
func enhanceAddressResponse(a domain.Address) map[string]interface{} {
	response := map[string]interface{}{
		"id":          a.ID,
		"title":       a.Title,
		"latitude":    a.Latitude,
		"longitude":   a.Longitude,
		"addressLine": a.AddressLine,
		"postalCode":  a.PostalCode,
		"cityId":      a.CityID,
		"provinceId":  a.ProvinceID,
		"userId":      a.UserID,
		"customerId":  a.CustomerID,
		"createdAt":   a.CreatedAt,
		"updatedAt":   a.UpdatedAt,
	}

	// Add city info if available
	if a.City != nil {
		response["city"] = map[string]interface{}{
			"id":   a.City.ID,
			"name": a.City.Name,
		}
	}

	// Add province info if available
	if a.Province != nil {
		response["province"] = map[string]interface{}{
			"id":   a.Province.ID,
			"name": a.Province.Name,
		}
	}

	return response
}
