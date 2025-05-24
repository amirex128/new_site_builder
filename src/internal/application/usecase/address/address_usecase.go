package addressusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/contract/common"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/address"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type AddressUsecase struct {
	*usecase.BaseUsecase
	addressRepo  repository.IAddressRepository
	cityRepo     repository.ICityRepository
	provinceRepo repository.IProvinceRepository
	authContext  func(c *gin.Context) service.IAuthService
}

func NewAddressUsecase(c contract.IContainer) *AddressUsecase {
	return &AddressUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		addressRepo:  c.GetAddressRepo(),
		cityRepo:     c.GetCityRepo(),
		provinceRepo: c.GetProvinceRepo(),
		authContext:  c.GetAuthTransientService(),
	}
}

func (u *AddressUsecase) CreateAddressCommand(params *address.CreateAddressCommand) (*resp.Response, error) {
	var err error
	userID, customerID, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

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

	newAddress := &domain.Address{
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

	err = u.addressRepo.Create(newAddress)
	if err != nil {
		u.Logger.Error("Error creating address", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در ایجاد آدرس")
	}

	return resp.NewResponseData(resp.Created, resp.Data{
		"address": newAddress,
	}, "آدرس با موفقیت ایجاد شد"), nil
}

func (u *AddressUsecase) UpdateAddressCommand(params *address.UpdateAddressCommand) (*resp.Response, error) {
	userID, customerID, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	existingAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "آدرس یافت نشد")
	}

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
	err = u.addressRepo.Update(existingAddress)
	if err != nil {
		u.Logger.Error("Error updating address", map[string]interface{}{
			"error": err.Error(),
			"id":    *params.ID,
		})
		return nil, errors.New("خطا در بروزرسانی آدرس")
	}

	// Retrieve the updated address with relations
	fullAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		return resp.NewResponseData(resp.Updated, enhanceAddressResponse(existingAddress), "آدرس با موفقیت بروزرسانی شد"), nil // Return the basic address if can't retrieve with relations
	}

	return resp.NewResponseData(resp.Updated, enhanceAddressResponse(fullAddress), "آدرس با موفقیت بروزرسانی شد"), nil
}

// DeleteAddressCommand handles deleting an address
func (u *AddressUsecase) DeleteAddressCommand(params *address.DeleteAddressCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteAddressCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing address to check ownership
	existingAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("آدرس مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check ownership
	customerID, _ := u.authContext(u.Ctx).GetCustomerID()
	userID, _ := u.authContext(u.Ctx).GetUserID()

	isAdmin, _ := u.authContext(u.Ctx).IsAdmin()

	// Check if user has access to this address
	if !isAdmin && existingAddress.CustomerID != customerID && existingAddress.UserID != userID {
		return nil, errors.New("شما دسترسی به این آدرس ندارید")
	}

	// Delete address
	err = u.addressRepo.Delete(*params.ID)
	if err != nil {
		u.Logger.Error("Error deleting address", map[string]interface{}{
			"error": err.Error(),
			"id":    *params.ID,
		})
		return nil, errors.New("خطا در حذف آدرس")
	}

	return resp.NewResponse(resp.Deleted, "آدرس با موفقیت حذف شد"), nil
}

// GetByIdAddressQuery handles retrieving an address by ID
func (u *AddressUsecase) GetByIdAddressQuery(params *address.GetByIdAddressQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdAddressQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get address
	result, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("آدرس مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check ownership
	customerID, _ := u.authContext(u.Ctx).GetCustomerID()
	userID, _ := u.authContext(u.Ctx).GetUserID()

	isAdmin, _ := u.authContext(u.Ctx).IsAdmin()

	// Check if user has access to this address
	if !isAdmin && result.CustomerID != customerID && result.UserID != userID {
		return nil, errors.New("شما دسترسی به این آدرس ندارید")
	}

	return resp.NewResponseData(resp.Retrieved, enhanceAddressResponse(result), "آدرس با موفقیت دریافت شد"), nil
}

// GetAllAddressQuery handles retrieving all addresses for the current user/customer
func (u *AddressUsecase) GetAllAddressQuery(params *address.GetAllAddressQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllAddressQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	var results *common.PaginationResponseDto[domain.Address]
	var err error

	// Try to get customer ID first
	customerID, customerErr := u.authContext(u.Ctx).GetCustomerID()
	if customerErr == nil && customerID > 0 {
		// Get addresses for customer
		results, err = u.addressRepo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
		if err != nil {
			return nil, errors.New("خطا در دریافت آدرس ها")
		}
	} else {
		// Try user ID
		userID, userErr := u.authContext(u.Ctx).GetUserID()
		if userErr != nil {
			return nil, errors.New("خطا در احراز هویت کاربر")
		}

		// Get addresses for user
		results, err = u.addressRepo.GetAllByUserID(userID, params.PaginationRequestDto)
		if err != nil {
			return nil, errors.New("خطا در دریافت آدرس ها")
		}
	}

	// Enhance response with full details
	enhancedAddresses := make([]map[string]interface{}, 0, len(results.Items))
	for _, addr := range results.Items {
		enhancedAddresses = append(enhancedAddresses, enhanceAddressResponse(addr))
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     enhancedAddresses,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "آدرس ها با موفقیت دریافت شدند"), nil
}

// AdminGetAllAddressQuery handles retrieving all addresses for admin
func (u *AddressUsecase) AdminGetAllAddressQuery(params *address.AdminGetAllAddressQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllAddressQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	if !isAdmin {
		return nil, errors.New("شما دسترسی به این عملیات را ندارید")
	}

	// Get all addresses with pagination
	results, err := u.addressRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all addresses", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت آدرس ها")
	}

	// Enhance response with full details
	enhancedAddresses := make([]map[string]interface{}, 0, len(results.Items))
	for _, addr := range results.Items {
		enhancedAddresses = append(enhancedAddresses, enhanceAddressResponse(addr))
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     enhancedAddresses,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "آدرس ها با موفقیت دریافت شدند"), nil
}

// GetAllCityQuery handles retrieving all cities
func (u *AddressUsecase) GetAllCityQuery(params *address.GetAllCityQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllCityQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get all cities with pagination
	results, err := u.cityRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all cities", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت شهرها")
	}

	// Convert cities to simplified response format
	cities := make([]map[string]interface{}, 0, len(results.Items))
	for _, city := range results.Items {
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

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     cities,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "شهرها با موفقیت دریافت شدند"), nil
}

// GetAllProvinceQuery handles retrieving all provinces
func (u *AddressUsecase) GetAllProvinceQuery(params *address.GetAllProvinceQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllProvinceQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get all provinces with pagination
	results, err := u.provinceRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all provinces", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت استان ها")
	}

	// Convert provinces to simplified response format
	provinces := make([]map[string]interface{}, 0, len(results.Items))
	for _, province := range results.Items {
		provinceData := map[string]interface{}{
			"id":     province.ID,
			"name":   province.Name,
			"slug":   province.Slug,
			"status": province.Status,
		}

		provinces = append(provinces, provinceData)
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     provinces,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "استان ها با موفقیت دریافت شدند"), nil
}
