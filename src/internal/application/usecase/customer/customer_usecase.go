package customerusecase

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type CustomerUsecase struct {
	*usecase.BaseUsecase
	repo         repository.ICustomerRepository
	fileItemRepo repository.IFileItemRepository
	addressRepo  repository.IAddressRepository
	roleRepo     repository.IRoleRepository
	authContext  func(c *gin.Context) service.IAuthService
	identitySvc  service.IIdentityService
	messageSvc   service.IMessageService
}

func NewCustomerUsecase(c contract.IContainer) *CustomerUsecase {
	return &CustomerUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		repo:         c.GetCustomerRepo(),
		fileItemRepo: c.GetFileItemRepo(),
		addressRepo:  c.GetAddressRepo(),
		roleRepo:     c.GetRoleRepo(),
		authContext:  c.GetAuthTransientService(),
		identitySvc:  c.GetIdentityService(),
		messageSvc:   c.GetMessageService(),
	}
}

// LoginCustomerCommand handles customer login
func (u *CustomerUsecase) LoginCustomerCommand(params *customer.LoginCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("LoginCustomerCommand called", map[string]interface{}{
		"email": params.Email,
	})

	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.Unauthorized, "ایمیل یا رمز عبور اشتباه است")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	if !u.identitySvc.VerifyPassword(*params.Password, existingCustomer.Password, existingCustomer.Salt) {
		return nil, resp.NewError(resp.Unauthorized, "ایمیل یا رمز عبور اشتباه است")
	}

	roleNames := make([]string, 0, len(existingCustomer.Roles))
	for _, role := range existingCustomer.Roles {
		roleNames = append(roleNames, role.Name)
	}
	if len(roleNames) == 0 {
		roleNames = append(roleNames, "Customer")
	}
	token := u.identitySvc.TokenForCustomer(*existingCustomer).AddRoles(roleNames).Make()

	return resp.NewResponseData(resp.Created, map[string]interface{}{
		"token": token,
	}, "ورود با موفقیت انجام شد"), nil
}

// RegisterCustomerCommand handles customer registration
func (u *CustomerUsecase) RegisterCustomerCommand(params *customer.RegisterCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("RegisterCustomerCommand called", map[string]interface{}{
		"email":  params.Email,
		"siteId": params.SiteID,
	})

	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err == nil && existingCustomer.ID > 0 {
		return nil, resp.NewError(resp.BadRequest, "ایمیل قبلاً ثبت شده است")
	}

	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)
	newCustomer := domain.Customer{
		Email:     *params.Email,
		SiteID:    *params.SiteID,
		Password:  hashedPassword,
		Salt:      salt,
		IsActive:  enums.ActiveStatus,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}
	verificationCode := utils.GenerateVerificationCode()
	newCustomer.VerifyCode = &verificationCode

	err = u.repo.Create(&newCustomer)
	if err != nil {
		u.Logger.Error("Error creating customer", map[string]interface{}{
			"error": err.Error(),
			"email": *params.Email,
		})
		return nil, resp.NewError(resp.Internal, "خطا در ثبت نام مشتری")
	}

	createdCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات مشتری")
	}

	defaultRole, err := u.roleRepo.GetByName("Customer")
	if err == nil {
		err = u.roleRepo.AddRoleToCustomer(defaultRole.ID, createdCustomer.ID)
		if err != nil {
			u.Logger.Error("Error assigning default role to customer", map[string]interface{}{
				"error":      err.Error(),
				"customerId": createdCustomer.ID,
			})
		}
	}
	token := u.identitySvc.TokenForCustomer(*createdCustomer).AddRoles([]string{"Customer"}).Make()

	u.messageSvc.SendEmail(struct {
		To      string
		Subject string
		Body    string
	}{
		To:      *params.Email,
		Subject: "Verify Your Account",
		Body:    fmt.Sprintf("Your verification code is: %s", verificationCode),
	})

	return resp.NewResponseData(resp.Created, map[string]interface{}{
		"token": token,
	}, "ثبت نام با موفقیت انجام شد. لطفا حساب خود را تایید کنید."), nil
}

// RequestVerifyAndForgetCustomerCommand handles verification and password reset requests
func (u *CustomerUsecase) RequestVerifyAndForgetCustomerCommand(params *customer.RequestVerifyAndForgetCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("RequestVerifyAndForgetCustomerCommand called", map[string]interface{}{
		"email": params.Email,
		"phone": params.Phone,
		"type":  params.Type,
	})

	if *params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType {
		if params.Phone == nil {
			return nil, resp.NewError(resp.BadRequest, "شماره تلفن الزامی است")
		}
		existingCustomer, err := u.repo.GetByPhone(*params.Phone)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "شماره تلفن یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		verifyCode := utils.GenerateVerificationCode()
		expireAt := time.Now().Add(15 * time.Minute)
		existingCustomer.VerifyCode = &verifyCode
		existingCustomer.ExpireVerifyCodeAt = &expireAt
		existingCustomer.UpdatedAt = time.Now()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		u.messageSvc.SendSms(struct {
			To   string
			Body string
		}{
			To:   *params.Phone,
			Body: fmt.Sprintf("Your verification code is: %s", verifyCode),
		})
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success": true,
			"message": "کد تایید ارسال شد. لطفا تلفن خود را بررسی کنید.",
		}, "کد تایید ارسال شد. لطفا تلفن خود را بررسی کنید."), nil
	}

	if *params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType {
		if params.Email == nil {
			return nil, resp.NewError(resp.BadRequest, "ایمیل الزامی است")
		}
		existingCustomer, err := u.repo.GetByEmail(*params.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "ایمیل یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		verifyCode := utils.GenerateVerificationCode()
		expireAt := time.Now().Add(24 * time.Hour)
		existingCustomer.VerifyCode = &verifyCode
		existingCustomer.ExpireVerifyCodeAt = &expireAt
		existingCustomer.UpdatedAt = time.Now()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		gatewayURL := os.Getenv("HTTP_GATEWAY_URL")
		if gatewayURL == "" {
			gatewayURL = "http://localhost"
		}
		resetLink := generatePasswordResetLink(gatewayURL, existingCustomer.ID, verifyCode, string(*params.Type))
		u.messageSvc.SendEmail(struct {
			To      string
			Subject string
			Body    string
		}{
			To:      *params.Email,
			Subject: "Your Verification Code",
			Body:    fmt.Sprintf("Your verification code is: %s\nOr use this link: %s", verifyCode, resetLink),
		})
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success": true,
			"message": "کد تایید ارسال شد. لطفا ایمیل خود را بررسی کنید.",
		}, "کد تایید ارسال شد. لطفا ایمیل خود را بررسی کنید."), nil
	}

	return nil, resp.NewError(resp.BadRequest, "نوع تأیید نامعتبر است")
}

// UpdateProfileCustomerCommand handles updating customer profile
func (u *CustomerUsecase) UpdateProfileCustomerCommand(params *customer.UpdateProfileCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateProfileCustomerCommand called", map[string]interface{}{
		"firstName": params.FirstName,
		"lastName":  params.LastName,
	})

	_, customerID, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت مشتری")
	}
	if customerID == nil {
		return nil, resp.NewError(resp.Unauthorized, "شناسه مشتری یافت نشد")
	}

	existingCustomer, err := u.repo.GetByID(*customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مشتری یافت نشد")
		}
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
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی پروفایل مشتری")
	}

	if len(params.AddressIDs) > 0 {
		err = u.addressRepo.RemoveAllAddressesFromCustomer(*customerID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی آدرس‌های مشتری")
		}
		for _, addressID := range params.AddressIDs {
			err = u.addressRepo.AddAddressToCustomer(addressID, *customerID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در افزودن آدرس به مشتری")
			}
		}
	}

	return resp.NewResponseData(resp.Updated, map[string]interface{}{
		"customer": existingCustomer,
	}, "پروفایل مشتری با موفقیت بروزرسانی شد"), nil
}

// VerifyCustomerQuery handles customer verification
func (u *CustomerUsecase) VerifyCustomerQuery(params *customer.VerifyCustomerQuery) (*resp.Response, error) {
	u.Logger.Info("VerifyCustomerQuery called", map[string]interface{}{
		"email": params.Email,
		"code":  params.Code,
		"type":  params.Type,
	})

	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "ایمیل یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	if existingCustomer.VerifyCode == nil || *existingCustomer.VerifyCode != *params.Code {
		return nil, resp.NewError(resp.BadRequest, "کد تایید نامعتبر است")
	}
	if existingCustomer.ExpireVerifyCodeAt == nil || time.Now().After(*existingCustomer.ExpireVerifyCodeAt) {
		return nil, resp.NewError(resp.BadRequest, "کد تایید منقضی شده است")
	}

	existingCustomer.VerifyCode = nil
	existingCustomer.ExpireVerifyCodeAt = nil
	existingCustomer.UpdatedAt = time.Now()

	var resetToken string

	switch *params.Type {
	case enums.VerifyEmailType:
		existingCustomer.VerifyEmail = enums.ActiveStatus
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success": true,
			"message": "تایید ایمیل با موفقیت انجام شد.",
		}, "تایید ایمیل با موفقیت انجام شد."), nil
	case enums.VerifyPhoneType:
		existingCustomer.VerifyPhone = enums.ActiveStatus
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success": true,
			"message": "تایید تلفن با موفقیت انجام شد.",
		}, "تایید تلفن با موفقیت انجام شد."), nil
	case enums.ForgetPasswordEmailType:
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("customer_id", strconv.FormatInt(existingCustomer.ID, 10)).Make()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		resetURL := os.Getenv("RESET_PASSWORD_URL")
		if resetURL == "" {
			resetURL = "https://example.com/reset-password"
		}
		queryParams := url.Values{}
		queryParams.Add("token", resetToken)
		resetURL = resetURL + "?" + queryParams.Encode()
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success":     true,
			"message":     "تایید بازیابی رمز عبور با موفقیت انجام شد.",
			"reset_token": resetToken,
			"reset_url":   resetURL,
		}, "تایید بازیابی رمز عبور با موفقیت انجام شد."), nil
	case enums.ForgetPasswordPhoneType:
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("customer_id", strconv.FormatInt(existingCustomer.ID, 10)).Make()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success":     true,
			"message":     "تایید بازیابی رمز عبور با موفقیت انجام شد.",
			"reset_token": resetToken,
		}, "تایید بازیابی رمز عبور با موفقیت انجام شد."), nil
	default:
		return nil, resp.NewError(resp.BadRequest, "نوع تأیید نامعتبر است")
	}
}

// GetProfileCustomerQuery handles getting customer profile
func (u *CustomerUsecase) GetProfileCustomerQuery(params *customer.GetProfileCustomerQuery) (*resp.Response, error) {
	u.Logger.Info("GetProfileCustomerQuery called", nil)

	_, customerID, _, err := u.authContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت مشتری")
	}
	if customerID == nil {
		return nil, resp.NewError(resp.Unauthorized, "شناسه مشتری یافت نشد")
	}

	existingCustomer, err := u.repo.GetByID(*customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مشتری یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	addressesResult, err := u.addressRepo.GetAllByCustomerID(*customerID, common.PaginationRequestDto{Page: 1, PageSize: 100})
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"customer": existingCustomer,
		"address":  addressesResult.Items,
	}, "پروفایل مشتری با موفقیت دریافت شد"), nil
}

// AdminGetAllCustomerQuery handles admin getting all customers
func (u *CustomerUsecase) AdminGetAllCustomerQuery(params *customer.AdminGetAllCustomerQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllCustomerQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	customersResult, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت لیست مشتریان")
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     customersResult.Items,
		"total":     customersResult.TotalCount,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (customersResult.TotalCount + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "لیست مشتریان با موفقیت دریافت شد (ادمین)"), nil
}

// Helper function to generate a password reset link
func generatePasswordResetLink(baseURL string, customerID int64, code string, verifyType string) string {
	return fmt.Sprintf("%s/reset-password?id=%d&code=%s&type=%s", baseURL, customerID, code, verifyType)
}
