package customerusecase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type CustomerUsecase struct {
	*usecase.BaseUsecase
	repo         repository.ICustomerRepository
	fileItemRepo repository.IFileItemRepository
	addressRepo  repository.IAddressRepository
	roleRepo     repository.IRoleRepository
	identitySvc  service.IIdentityService
	messageSvc   service.IMessageService
}

func NewCustomerUsecase(c contract.IContainer) *CustomerUsecase {
	return &CustomerUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		repo:         c.GetCustomerRepo(),
		fileItemRepo: c.GetFileItemRepo(),
		addressRepo:  c.GetAddressRepo(),
		roleRepo:     c.GetRoleRepo(),
		identitySvc:  c.GetIdentityService(),
		messageSvc:   c.GetMessageService(),
	}
}

func (u *CustomerUsecase) LoginCustomerCommand(params *customer.LoginCustomerCommand) (*resp.Response, error) {
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "اطلاعات وارد شده صحیح نمی باشد")
	}

	if existingCustomer.IsActive == enums.InactiveStatus {
		return nil, resp.NewError(resp.Unauthorized, "حساب کاربری غیرفعال است")
	}

	if !u.identitySvc.VerifyPassword(*params.Password, existingCustomer.Password, existingCustomer.Salt) {
		return nil, resp.NewError(resp.Unauthorized, "اطلاعات وارد شده صحیح نمی باشد")
	}

	token := u.identitySvc.TokenForCustomer(*existingCustomer).Make()

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": token,
		},
		"ورود با موفقیت انجام شد",
	), nil
}

func (u *CustomerUsecase) RegisterCustomerCommand(params *customer.RegisterCustomerCommand) (*resp.Response, error) {
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err == nil && existingCustomer.ID > 0 {
		return nil, resp.NewError(resp.BadRequest, "ایمیل وارد شده قبلا استفاده شده است")
	}

	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)
	newCustomer := domain.Customer{
		Email:     *params.Email,
		SiteID:    *params.SiteID,
		Password:  hashedPassword,
		Salt:      salt,
		IsActive:  enums.InactiveStatus,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}
	verificationCode := utils.GenerateVerificationCode()
	newCustomer.VerifyCode = &verificationCode

	err = u.repo.Create(&newCustomer)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ثبت نام مشتری")
	}
	token := u.identitySvc.TokenForCustomer(newCustomer).Make()
	return resp.NewResponseData(
		resp.Created,
		map[string]interface{}{
			"token": token,
		},
		"ثبت نام با موفقیت انجام شد. لطفا حساب خود را فعال کنید.",
	), nil
}

func (u *CustomerUsecase) RequestVerifyAndForgetCustomerCommand(params *customer.RequestVerifyAndForgetCustomerCommand) (*resp.Response, error) {
	var existingCustomer *domain.Customer
	var err error

	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) {
		if params.Email == nil {
			return nil, resp.NewError(resp.BadRequest, "ایمیل الزامی است")
		}
		existingCustomer, err = u.repo.GetByEmail(*params.Email)
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) {
		if params.Phone == nil {
			return nil, resp.NewError(resp.BadRequest, "شماره تلفن الزامی است")
		}
		existingCustomer, err = u.repo.GetByPhone(*params.Phone)
	} else {
		return nil, resp.NewError(resp.BadRequest, "نوع احراز هویت صحیح نمی باشد")
	}

	if err != nil {
		return nil, resp.NewError(resp.NotFound, "مشتری یافت نشد")
	}

	verificationCode := utils.GenerateVerificationCode()
	existingCustomer.VerifyCode = &verificationCode

	err = u.repo.Update(existingCustomer)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) && params.Email != nil {
		u.messageSvc.SendEmail(struct {
			To      string
			Subject string
			Body    string
		}{
			To:      *params.Email,
			Subject: "Your Verification Code",
			Body:    fmt.Sprintf("Your verification code is: %s", verificationCode),
		})
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) && params.Phone != nil {
		u.messageSvc.SendSms(struct {
			To   string
			Body string
		}{
			To:   *params.Phone,
			Body: fmt.Sprintf("Your verification code is: %s", verificationCode),
		})
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"success": true,
			"message": "کد احراز هویت با موفقیت ارسال شد",
		},
		"کد احراز هویت با موفقیت ارسال شد",
	), nil
}

func (u *CustomerUsecase) UpdateProfileCustomerCommand(params *customer.UpdateProfileCustomerCommand) (*resp.Response, error) {
	_, customerID, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil || customerID == nil {
		return nil, err
	}
	existingCustomer, err := u.repo.GetByID(*customerID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.FirstName != nil {
		existingCustomer.FirstName = *params.FirstName
	}
	if params.LastName != nil {
		existingCustomer.LastName = *params.LastName
	}
	if params.Email != nil {
		existingCustomer.Email = *params.Email
	}
	if params.Password != nil {
		hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)
		existingCustomer.Password = hashedPassword
		existingCustomer.Salt = salt
	}
	if params.NationalCode != nil {
		existingCustomer.NationalCode = *params.NationalCode
	}
	if params.Phone != nil {
		existingCustomer.Phone = *params.Phone
	}
	existingCustomer.UpdatedAt = time.Now()
	err = u.repo.Update(existingCustomer)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if len(params.AddressIDs) > 0 {
		err = u.addressRepo.RemoveAllAddressesFromCustomer(*customerID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		for _, addressID := range params.AddressIDs {
			err = u.addressRepo.AddAddressToCustomer(addressID, *customerID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	}
	return resp.NewResponse(resp.Updated, "اطلاعات مشتری با موفقیت به روز شد"), nil
}

func (u *CustomerUsecase) VerifyCustomerQuery(params *customer.VerifyCustomerQuery) (*resp.Response, error) {
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "مشتری یافت نشد")
	}

	codeStr := *params.Code
	if existingCustomer.VerifyCode == nil || *existingCustomer.VerifyCode != codeStr {
		return nil, resp.NewError(resp.BadRequest, "کد احراز هویت صحیح نمی باشد")
	}

	var resetToken string

	if params.Type != nil && *params.Type == enums.VerifyEmailType {
		existingCustomer.IsActive = enums.ActiveStatus
		existingCustomer.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums.VerifyPhoneType {
		existingCustomer.IsActive = enums.ActiveStatus
		existingCustomer.VerifyPhone = ""
	} else if params.Type != nil && *params.Type == enums.ForgetPasswordEmailType {
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("customer_id", strconv.FormatInt(existingCustomer.ID, 10)).Make()
		existingCustomer.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums.ForgetPasswordPhoneType {
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("customer_id", strconv.FormatInt(existingCustomer.ID, 10)).Make()
		existingCustomer.VerifyPhone = ""
	} else {
		return nil, resp.NewError(resp.BadRequest, "نوع احراز هویت صحیح نمی باشد")
	}

	err = u.repo.Update(existingCustomer)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	respData := map[string]interface{}{
		"success": true,
		"message": "احراز هویت با موفقیت انجام شد",
	}
	if resetToken != "" {
		respData["reset_token"] = resetToken
	}

	return resp.NewResponseData(
		resp.Success,
		respData,
		"احراز هویت با موفقیت انجام شد",
	), nil
}

func (u *CustomerUsecase) GetProfileCustomerQuery(params *customer.GetProfileCustomerQuery) (*resp.Response, error) {
	_, customerID, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil || customerID == nil {
		return nil, err
	}
	existingCustomer, err := u.repo.GetByID(*customerID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	addresses, err := u.addressRepo.GetAllByCustomerID(*customerID, common.PaginationRequestDto{Page: 1, PageSize: 100})
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{"customer": existingCustomer, "addresses": addresses}, "اطلاعات مشتری با موفقیت دریافت شد"), nil
}

func (u *CustomerUsecase) AdminGetAllCustomerQuery(params *customer.AdminGetAllCustomerQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}
	customersResult, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت لیست مشتریان")
	}

	return resp.NewResponseData(resp.Retrieved, customersResult, "لیست مشتریان با موفقیت دریافت شد (ادمین)"), nil
}
