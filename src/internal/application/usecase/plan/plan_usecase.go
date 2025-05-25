package planusecase

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/plan"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"github.com/gin-gonic/gin"
)

type PlanUsecase struct {
	*usecase.BaseUsecase
	planRepo    repository.IPlanRepository
	userRepo    repository.IUserRepository
	authContext func(c *gin.Context) service.IAuthService
}

func NewPlanUsecase(c contract.IContainer) *PlanUsecase {
	return &PlanUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		planRepo:    c.GetPlanRepo(),
		userRepo:    c.GetUserRepo(),
		authContext: c.GetAuthTransientService(),
	}
}

func (u *PlanUsecase) CreatePlanCommand(params *plan.CreatePlanCommand) (*resp.Response, error) {
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "only admins can create plans")
	}
	newPlan := domain.Plan{
		Name:             *params.Name,
		ShowStatus:       *params.ShowStatus,
		Description:      *params.Description,
		Price:            int64(*params.Price),
		Duration:         *params.Duration,
		SmsCredits:       *params.SmsCredits,
		EmailCredits:     *params.EmailCredits,
		StorageMbCredits: *params.StorageMbCredits,
		AiCredits:        *params.AiCredits,
		AiImageCredits:   *params.AiImageCredits,
	}
	if params.DiscountType != nil {
		newPlan.DiscountType = string(*params.DiscountType)
	}
	if params.Discount != nil {
		discount := int64(*params.Discount)
		newPlan.Discount = &discount
	}
	if params.Feature != nil {
		newPlan.Feature = *params.Feature
	}
	err = u.planRepo.Create(newPlan)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Created, map[string]interface{}{"plan": newPlan}, "پلن با موفقیت ایجاد شد"), nil
}

func (u *PlanUsecase) UpdatePlanCommand(params *plan.UpdatePlanCommand) (*resp.Response, error) {
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "only admins can update plans")
	}
	existingPlan, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	if params.Name != nil {
		existingPlan.Name = *params.Name
	}
	if params.ShowStatus != nil {
		existingPlan.ShowStatus = *params.ShowStatus
	}
	if params.Description != nil {
		existingPlan.Description = *params.Description
	}
	if params.Price != nil {
		existingPlan.Price = int64(*params.Price)
	}
	if params.DiscountType != nil {
		existingPlan.DiscountType = string(*params.DiscountType)
	}
	if params.Discount != nil {
		discount := int64(*params.Discount)
		existingPlan.Discount = &discount
	}
	if params.Duration != nil {
		existingPlan.Duration = *params.Duration
	}
	if params.Feature != nil {
		existingPlan.Feature = *params.Feature
	}
	if params.SmsCredits != nil {
		existingPlan.SmsCredits = *params.SmsCredits
	}
	if params.EmailCredits != nil {
		existingPlan.EmailCredits = *params.EmailCredits
	}
	if params.StorageMbCredits != nil {
		existingPlan.StorageMbCredits = *params.StorageMbCredits
	}
	if params.AiCredits != nil {
		existingPlan.AiCredits = *params.AiCredits
	}
	if params.AiImageCredits != nil {
		existingPlan.AiImageCredits = *params.AiImageCredits
	}
	err = u.planRepo.Update(existingPlan)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Updated, map[string]interface{}{"plan": existingPlan}, "پلن با موفقیت بروزرسانی شد"), nil
}

func (u *PlanUsecase) DeletePlanCommand(params *plan.DeletePlanCommand) (*resp.Response, error) {
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "only admins can delete plans")
	}
	_, err = u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	err = u.planRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Deleted, map[string]interface{}{"success": true}, "پلن با موفقیت حذف شد"), nil
}

func (u *PlanUsecase) GetAllPlanQuery(params *plan.GetAllPlanQuery) (*resp.Response, error) {
	plansResult, err := u.planRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     plansResult.Items,
		"total":     plansResult.TotalCount,
		"page":      plansResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": plansResult.TotalPages,
	}, "لیست پلن ها با موفقیت دریافت شد"), nil
}

func (u *PlanUsecase) GetByIDPlanQuery(params *plan.GetByIDPlanQuery) (*resp.Response, error) {
	plan, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"plan": plan}, "پلن با موفقیت دریافت شد"), nil
}

func (u *PlanUsecase) CalculatePlanPriceQuery(params *plan.CalculatePlanPriceQuery) (*resp.Response, error) {
	plan, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	var finalPrice int64 = plan.Price
	var discountAmount int64 = 0
	if plan.Discount != nil && *plan.Discount > 0 {
		if plan.DiscountType == string(enums.FixedDiscountType) {
			discountAmount = *plan.Discount
			finalPrice = plan.Price - discountAmount
		} else if plan.DiscountType == string(enums.PercentageDiscountType) {
			discountAmount = (plan.Price * (*plan.Discount)) / 100
			finalPrice = plan.Price - discountAmount
		}
	}
	if finalPrice < 0 {
		finalPrice = 0
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"plan":           plan,
		"originalPrice":  plan.Price,
		"discountAmount": discountAmount,
		"finalPrice":     finalPrice,
	}, "قیمت پلن با موفقیت محاسبه شد"), nil
}
