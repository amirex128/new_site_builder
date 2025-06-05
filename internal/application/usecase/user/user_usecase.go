package userusecase

import (
	"fmt"
	user2 "github.com/amirex128/new_site_builder/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	service2 "github.com/amirex128/new_site_builder/internal/contract/service"
	"github.com/amirex128/new_site_builder/internal/domain"
	enums2 "github.com/amirex128/new_site_builder/internal/domain/enums"
	"strconv"
	"time"
)

type UserUsecase struct {
	*usecase.BaseUsecase
	userRepo      repository2.IUserRepository
	planRepo      repository2.IPlanRepository
	addressRepo   repository2.IAddressRepository
	paymentRepo   repository2.IPaymentRepository
	identitySvc   service2.IIdentityService
	siteRepo      repository2.ISiteRepository
	messageSvc    service2.IMessageService
	unitPriceRepo repository2.IUnitPriceRepository
}

func NewUserUsecase(c contract.IContainer) *UserUsecase {
	return &UserUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		siteRepo:      c.GetSiteRepo(),
		userRepo:      c.GetUserRepo(),
		planRepo:      c.GetPlanRepo(),
		addressRepo:   c.GetAddressRepo(),
		paymentRepo:   c.GetPaymentRepo(),
		identitySvc:   c.GetIdentityService(),
		messageSvc:    c.GetMessageService(),
		unitPriceRepo: c.GetUnitPriceRepo(),
	}
}

func (u *UserUsecase) UpdateProfileUserCommand(params *user2.UpdateProfileUserCommand) (*resp.Response, error) {
	userId, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userId == nil {
		return nil, err
	}
	existingUser, err := u.userRepo.GetByID(*userId)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.FirstName != nil {
		existingUser.FirstName = *params.FirstName
	}
	if params.LastName != nil {
		existingUser.LastName = *params.LastName
	}
	if params.Email != nil {
		existingUser.Email = *params.Email
	}
	if params.Password != nil {
		hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)
		existingUser.Password = hashedPassword
		existingUser.Salt = salt
	}
	if params.NationalCode != nil {
		existingUser.NationalCode = *params.NationalCode
	}
	if params.Phone != nil {
		existingUser.Phone = *params.Phone
	}
	if params.AiTypeEnum != nil {
		existingUser.AiTypeEnum = *params.AiTypeEnum
	}
	if params.UseCustomEmailSmtp != nil {
		existingUser.UseCustomEmailSmtp = string(*params.UseCustomEmailSmtp)
	}
	if params.Smtp != nil {
		existingUser.SmtpHost = params.Smtp.Host
		existingUser.SmtpPort = &params.Smtp.Port
		existingUser.SmtpUsername = params.Smtp.Username
		existingUser.SmtpPassword = params.Smtp.Password
	}
	existingUser.UpdatedAt = time.Now()
	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if len(params.AddressIDs) > 0 {
		err = u.addressRepo.RemoveAllAddressesFromUser(*userId)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		for _, addressID := range params.AddressIDs {
			err = u.addressRepo.AddAddressToUser(addressID, *userId)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	}
	return resp.NewResponse(resp.Updated, "اطلاعات کاربری با موفقیت به روز شد"), nil
}

func (u *UserUsecase) GetProfileUserQuery(params *user2.GetProfileUserQuery) (*resp.Response, error) {
	userId, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil || userId == nil {
		return nil, err
	}
	existingUser, err := u.userRepo.GetByID(*userId)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "کاربر یافت نشد")
	}
	addresses, err := u.addressRepo.GetAllByUserID(*userId)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{"user": existingUser, "addresses": addresses}, "اطلاعات کاربری با موفقیت دریافت شد"), nil
}

func (u *UserUsecase) RegisterUserCommand(params *user2.RegisterUserCommand) (*resp.Response, error) {
	_, err := u.userRepo.GetByEmail(*params.Email)
	if err == nil {
		return nil, resp.NewError(resp.BadRequest, "ایمیل وارد شده قبلا استفاده شده است")
	}

	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)

	newUser := domain.User{
		Email:        *params.Email,
		Password:     hashedPassword,
		Salt:         salt,
		VerifyEmail:  enums2.InactiveStatus,
		VerifyPhone:  enums2.InactiveStatus,
		IsActive:     enums2.ActiveStatus,
		AiTypeEnum:   enums2.GPT35Type,
		UserTypeEnum: enums2.UserTypeValue,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	verificationCode := utils.GenerateVerificationCode()
	newUser.VerifyCode = &verificationCode

	err = u.userRepo.Create(&newUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ثبت نام کاربر")
	}

	token := u.identitySvc.TokenForUser(newUser).Make()
	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": "Bearer " + token,
		},
		"ثبت نام با موفقیت انجام شد. لطفا حساب خود را فعال کنید.",
	), nil
}

func (u *UserUsecase) LoginUserCommand(params *user2.LoginUserCommand) (*resp.Response, error) {
	existingUser, err := u.userRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "اطلاعات وارد شده صحیح نمی باشد")
	}

	if existingUser.IsActive == enums2.InactiveStatus {
		return nil, resp.NewError(resp.Unauthorized, "حساب کاربری غیرفعال است")
	}

	if !u.identitySvc.VerifyPassword(*params.Password, existingUser.Password, existingUser.Salt) {
		return nil, resp.NewError(resp.Unauthorized, "اطلاعات وارد شده صحیح نمی باشد")
	}

	token := u.identitySvc.TokenForUser(*existingUser).Make()

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": "Bearer " + token,
		},
		"ورود با موفقیت انجام شد",
	), nil
}

func (u *UserUsecase) RequestVerifyAndForgetUserCommand(params *user2.RequestVerifyAndForgetUserCommand) (*resp.Response, error) {
	var existingUser *domain.User
	var err error

	if params.Type != nil && (*params.Type == enums2.VerifyEmailType || *params.Type == enums2.ForgetPasswordEmailType) {
		if params.Email == nil {
			return nil, resp.NewError(resp.BadRequest, "ایمیل الزامی است")
		}
		existingUser, err = u.userRepo.GetByEmail(*params.Email)
	} else if params.Type != nil && (*params.Type == enums2.VerifyPhoneType || *params.Type == enums2.ForgetPasswordPhoneType) {
		if params.Phone == nil {
			return nil, resp.NewError(resp.BadRequest, "شماره تلفن الزامی است")
		}
		existingUser, err = u.userRepo.GetByPhone(*params.Phone)
	} else {
		return nil, resp.NewError(resp.BadRequest, "نوع احراز هویت صحیح نمی باشد")
	}

	if err != nil {
		return nil, resp.NewError(resp.NotFound, "کاربر یافت نشد")
	}

	verificationCode := utils.GenerateVerificationCode()

	if params.Type != nil && (*params.Type == enums2.VerifyEmailType || *params.Type == enums2.ForgetPasswordEmailType) {
		existingUser.VerifyCode = &verificationCode
	} else if params.Type != nil && (*params.Type == enums2.VerifyPhoneType || *params.Type == enums2.ForgetPasswordPhoneType) {
		existingUser.VerifyCode = &verificationCode
	}

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	if params.Type != nil && (*params.Type == enums2.VerifyEmailType || *params.Type == enums2.ForgetPasswordEmailType) && params.Email != nil {
		msg := struct {
			To      string
			Subject string
			Body    string
		}{
			To:      *params.Email,
			Subject: "Your Verification Code",
			Body:    fmt.Sprintf("Your verification code is: %d", verificationCode),
		}
		u.messageSvc.SendEmail(msg)
	} else if params.Type != nil && (*params.Type == enums2.VerifyPhoneType || *params.Type == enums2.ForgetPasswordPhoneType) && params.Phone != nil {
		msg := struct {
			To   string
			Body string
		}{
			To:   *params.Phone,
			Body: fmt.Sprintf("Your verification code is: %d", verificationCode),
		}
		u.messageSvc.SendSms(msg)
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"success": true,
			"message": "Verification code sent. Please check your email or phone.",
		},
		"Verification code sent. Please check your email or phone.",
	), nil
}

func (u *UserUsecase) VerifyUserQuery(params *user2.VerifyUserQuery) (*resp.Response, error) {
	existingUser, err := u.userRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "کاربر یافت نشد")
	}

	codeStr := *params.Code

	var resetToken string

	if params.Type != nil && *params.Type == enums2.VerifyEmailType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "کد احراز هویت صحیح نمی باشد")
		}
		existingUser.IsActive = enums2.ActiveStatus // Activate the user
		existingUser.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums2.VerifyPhoneType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "کد احراز هویت صحیح نمی باشد")
		}
		existingUser.IsActive = enums2.ActiveStatus // Activate the user
		existingUser.VerifyPhone = ""
	} else if params.Type != nil && *params.Type == enums2.ForgetPasswordEmailType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "کد احراز هویت صحیح نمی باشد")
		}
		// Provide a token for password reset
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("user_id", strconv.FormatInt(existingUser.ID, 10)).Make()
		existingUser.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums2.ForgetPasswordPhoneType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "کد احراز هویت صحیح نمی باشد")
		}
		// Provide a token for password reset
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("user_id", strconv.FormatInt(existingUser.ID, 10)).Make()
		existingUser.VerifyPhone = ""
	} else {
		return nil, resp.NewError(resp.BadRequest, "نوع احراز هویت صحیح نمی باشد")
	}

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	respData := resp.Data{
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

func (u *UserUsecase) ChargeCreditRequestUserCommand(params *user2.ChargeCreditRequestUserCommand) (*resp.Response, error) {
	// Get the  user ID
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Calculate total amount
	var totalAmount int64 = 0
	orderData := make(map[string]string)

	for i, unitPrice := range params.UnitPrices {
		unitPriceObj, err := u.unitPriceRepo.GetByName(string(*unitPrice.UnitPriceName))
		if err != nil {
			return nil, resp.NewError(resp.BadRequest, "قیمت واحد یافت نشد")
		}
		var itemPrice int64 = unitPriceObj.Price * int64(*unitPrice.UnitPriceCount)

		if string(*unitPrice.UnitPriceName) == "storage_mb_credits" && unitPrice.UnitPriceDay != nil {
			itemPrice = itemPrice * int64(*unitPrice.UnitPriceDay)
		}

		totalAmount += itemPrice

		orderData[fmt.Sprintf("unitPrice_%d_name", i)] = string(*unitPrice.UnitPriceName)
		orderData[fmt.Sprintf("unitPrice_%d_count", i)] = strconv.Itoa(*unitPrice.UnitPriceCount)
		if unitPrice.UnitPriceDay != nil {
			orderData[fmt.Sprintf("unitPrice_%d_days", i)] = strconv.Itoa(*unitPrice.UnitPriceDay)
		}
	}
	orderID := time.Now().UnixNano()

	orderData["userId"] = strconv.FormatInt(*userID, 10)
	orderData["totalAmount"] = strconv.FormatInt(totalAmount, 10)
	orderData["finalFrontReturnUrl"] = *params.FinalFrontReturnUrl

	paymentUrl, err := u.paymentRepo.RequestPayment(totalAmount, orderID, *userID, string(*params.Gateway), orderData)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"paymentUrl": paymentUrl,
			"orderId":    orderID,
		},
		"لینک سفارش با موفقیت ایجاد شد",
	), nil
}

func (u *UserUsecase) UpgradePlanRequestUserCommand(params *user2.UpgradePlanRequestUserCommand) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	plan, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "پلن یافت نشد")
	}

	var finalPrice int64 = plan.Price
	var discountAmount int64 = 0

	if plan.Discount != nil && *plan.Discount > 0 {
		if plan.DiscountType == string(enums2.FixedDiscountType) {
			discountAmount = *plan.Discount
			finalPrice = plan.Price - discountAmount
		} else if plan.DiscountType == string(enums2.PercentageDiscountType) {
			discountAmount = (plan.Price * (*plan.Discount)) / 100
			finalPrice = plan.Price - discountAmount
		}
	}

	if finalPrice < 0 {
		finalPrice = 0
	}

	orderData := make(map[string]string)
	orderData["userId"] = strconv.FormatInt(*userID, 10)
	orderData["planId"] = strconv.FormatInt(*params.PlanID, 10)
	orderData["originalPrice"] = strconv.FormatInt(plan.Price, 10)
	orderData["discountAmount"] = strconv.FormatInt(discountAmount, 10)
	orderData["finalPrice"] = strconv.FormatInt(finalPrice, 10)
	orderData["finalFrontReturnUrl"] = *params.FinalFrontReturnUrl

	orderID := time.Now().UnixNano()

	paymentUrl, err := u.paymentRepo.RequestPayment(finalPrice, orderID, *userID, string(*params.Gateway), orderData)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"paymentUrl": paymentUrl,
			"orderId":    orderID,
			"plan":       plan,
		},
		"لینک سفارش با موفقیت ایجاد شد",
	), nil
}

func (u *UserUsecase) AdminGetAllUserQuery(params *user2.AdminGetAllUserQuery) (*resp.Response, error) {
	err := u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	result, err := u.userRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت کاربران")
	}

	return resp.NewResponseData(resp.Success, result, "کاربران با موفقیت دریافت شدند"), nil
}
