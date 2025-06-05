package addressusecase

import (
	"errors"
	address2 "github.com/amirex128/new_site_builder/internal/application/dto/address"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	"github.com/amirex128/new_site_builder/internal/domain"
	"github.com/amirex128/new_site_builder/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type AddressUsecase struct {
	*usecase.BaseUsecase
	addressRepo  repository2.IAddressRepository
	cityRepo     repository2.ICityRepository
	provinceRepo repository2.IProvinceRepository
}

func NewAddressUsecase(c contract.IContainer) *AddressUsecase {
	return &AddressUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		addressRepo:  c.GetAddressRepo(),
		cityRepo:     c.GetCityRepo(),
		provinceRepo: c.GetProvinceRepo(),
	}
}

func (u *AddressUsecase) CreateAddressCommand(params *address2.CreateAddressCommand) (*resp.Response, error) {
	var err error
	userID, customerID, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, err
	}

	_, err = u.cityRepo.GetByID(*params.CityID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "شهر مورد نظر یافت نشد")
	}

	_, err = u.provinceRepo.GetByID(*params.ProvinceID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "استان مورد نظر یافت نشد")
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
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد آدرس")
	}

	return resp.NewResponseData(resp.Created, newAddress, "آدرس با موفقیت ایجاد شد"), nil
}

func (u *AddressUsecase) UpdateAddressCommand(params *address2.UpdateAddressCommand) (*resp.Response, error) {
	existingAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "آدرس یافت نشد")
	}

	userType, err := u.AuthContext(u.Ctx).GetUserType()
	if err != nil {
		return nil, err
	}
	if *userType == enums.UserTypeValue {
		err = u.CheckAccessUserModel(existingAddress)
		if err != nil {
			return nil, err
		}
	}else if *userType == enums.CustomerTypeValue{
		err = u.CheckAccessCustomerModel(existingAddress)
		if err != nil {
			return nil, err
		}
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

	err = u.addressRepo.Update(existingAddress)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Updated, existingAddress, "آدرس با موفقیت بروزرسانی شد"), nil
}

func (u *AddressUsecase) DeleteAddressCommand(params *address2.DeleteAddressCommand) (*resp.Response, error) {
	existingAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت آدرس")
	}

	userType, err := u.AuthContext(u.Ctx).GetUserType()
	if err != nil {
		return nil, err
	}
	if *userType == enums.UserTypeValue {
		err = u.CheckAccessUserModel(existingAddress)
		if err != nil {
			return nil, err
		}
	}else if *userType == enums.CustomerTypeValue{
		err = u.CheckAccessCustomerModel(existingAddress)
		if err != nil {
			return nil, err
		}
	}

	err = u.addressRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در حذف آدرس")
	}

	return resp.NewResponse(resp.Deleted, "آدرس با موفقیت حذف شد"), nil
}

func (u *AddressUsecase) GetByIdAddressQuery(params *address2.GetByIdAddressQuery) (*resp.Response, error) {
	existingAddress, err := u.addressRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "آدرس مورد نظر یافت نشد")
		}
		return nil, err
	}
	userType, err := u.AuthContext(u.Ctx).GetUserType()
	if err != nil {
		return nil, err
	}
	if *userType == enums.UserTypeValue {
		err = u.CheckAccessUserModel(existingAddress)
		if err != nil {
			return nil, err
		}
	}else if *userType == enums.CustomerTypeValue{
		err = u.CheckAccessCustomerModel(existingAddress)
		if err != nil {
			return nil, err
		}
	}

	return resp.NewResponseData(resp.Retrieved, existingAddress, "آدرس با موفقیت دریافت شد"), nil
}

func (u *AddressUsecase) GetAllAddressQuery(params *address2.GetAllAddressQuery) (*resp.Response, error) {

	var results []domain.Address
	var err error

	userID, customerID, userType, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, err
	}

	if *userType == enums.CustomerTypeValue {
		results, err = u.addressRepo.GetAllByCustomerID(*customerID)
		if err != nil {
			return nil, errors.New("خطا در دریافت آدرس ها")
		}
	}
	if *userType == enums.UserTypeValue {
		results, err = u.addressRepo.GetAllByUserID(*userID)
		if err != nil {
			return nil, errors.New("خطا در دریافت آدرس ها")
		}
	}

	return resp.NewResponseData(resp.Retrieved, results, "آدرس ها با موفقیت دریافت شدند"), nil
}

func (u *AddressUsecase) AdminGetAllAddressQuery(params *address2.AdminGetAllAddressQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	results, err := u.addressRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت آدرس ها")
	}

	return resp.NewResponseData(resp.Retrieved, results, "آدرس ها با موفقیت دریافت شدند"), nil
}

func (u *AddressUsecase) GetAllCityQuery(params *address2.GetAllCityQuery) (*resp.Response, error) {
	results, err := u.cityRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت شهر ها")
	}

	return resp.NewResponseData(resp.Retrieved, results, "شهرها با موفقیت دریافت شدند"), nil
}

func (u *AddressUsecase) GetAllProvinceQuery(params *address2.GetAllProvinceQuery) (*resp.Response, error) {
	results, err := u.provinceRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت استان ها")
	}

	return resp.NewResponseData(resp.Retrieved, results, "استان ها با موفقیت دریافت شدند"), nil
}
