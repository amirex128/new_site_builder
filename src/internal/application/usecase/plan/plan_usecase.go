package planusecase

import (
	"fmt"
	"strconv"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/plan"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type PlanUsecase struct {
	logger   sflogger.Logger
	planRepo repository.IPlanRepository
}

func NewPlanUsecase(c contract.IContainer) *PlanUsecase {
	return &PlanUsecase{
		logger:   c.GetLogger(),
		planRepo: c.GetPlanRepo(),
	}
}

func (u *PlanUsecase) CreatePlanCommand(params *plan.CreatePlanCommand) (any, error) {
	// Implementation for creating a plan
	fmt.Println(params)

	var description string
	if params.Description != nil {
		description = *params.Description
	}

	var feature string
	if params.Feature != nil {
		feature = *params.Feature
	}

	var discount int64
	if params.Discount != nil {
		discount = *params.Discount
	}

	var discountType string
	if params.DiscountType != nil {
		discountType = strconv.Itoa(int(*params.DiscountType))
	}

	var smsCredits int
	if params.SmsCredits != nil {
		smsCredits = *params.SmsCredits
	}

	var emailCredits int
	if params.EmailCredits != nil {
		emailCredits = *params.EmailCredits
	}

	var storageMbCredits int
	if params.StorageCredits != nil {
		storageMbCredits = *params.StorageCredits
	}

	var aiCredits int
	if params.AiCredits != nil {
		aiCredits = *params.AiCredits
	}

	var aiImageCredits int
	if params.AiImageCredits != nil {
		aiImageCredits = *params.AiImageCredits
	}

	newPlan := domain.Plan{
		Name:             *params.Name,
		Description:      description,
		Price:            *params.Price,
		DiscountType:     discountType,
		Discount:         &discount,
		Duration:         *params.Duration,
		Feature:          feature,
		SmsCredits:       smsCredits,
		EmailCredits:     emailCredits,
		StorageMbCredits: storageMbCredits,
		AiCredits:        aiCredits,
		AiImageCredits:   aiImageCredits,
	}

	err := u.planRepo.Create(newPlan)
	if err != nil {
		return nil, err
	}

	return newPlan, nil
}

func (u *PlanUsecase) UpdatePlanCommand(params *plan.UpdatePlanCommand) (any, error) {
	// Implementation for updating a plan
	fmt.Println(params)

	existingPlan, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	if params.Name != nil {
		existingPlan.Name = *params.Name
	}

	if params.Description != nil {
		existingPlan.Description = *params.Description
	}

	if params.Price != nil {
		existingPlan.Price = *params.Price
	}

	if params.DiscountType != nil {
		existingPlan.DiscountType = strconv.Itoa(int(*params.DiscountType))
	}

	if params.Discount != nil {
		existingPlan.Discount = params.Discount
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

	if params.StorageCredits != nil {
		existingPlan.StorageMbCredits = *params.StorageCredits
	}

	if params.AiCredits != nil {
		existingPlan.AiCredits = *params.AiCredits
	}

	if params.AiImageCredits != nil {
		existingPlan.AiImageCredits = *params.AiImageCredits
	}

	err = u.planRepo.Update(existingPlan)
	if err != nil {
		return nil, err
	}

	return existingPlan, nil
}

func (u *PlanUsecase) DeletePlanCommand(params *plan.DeletePlanCommand) (any, error) {
	// Implementation for deleting a plan
	fmt.Println(params)

	err := u.planRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *PlanUsecase) GetByIdPlanQuery(params *plan.GetByIdPlanQuery) (any, error) {
	// Implementation to get plan by ID
	fmt.Println(params)

	result, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *PlanUsecase) GetAllPlanQuery(params *plan.GetAllPlanQuery) (any, error) {
	// Implementation to get all plans
	fmt.Println(params)

	result, count, err := u.planRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *PlanUsecase) CalculatePlanPriceQuery(params *plan.CalculatePlanPriceQuery) (any, error) {
	// Implementation to calculate a plan's price
	fmt.Println(params)

	plan, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, err
	}

	var finalPrice int64 = plan.Price
	var discountAmount int64 = 0

	// Calculate discount based on discount type
	if plan.Discount != nil && *plan.Discount > 0 {
		if plan.DiscountType == strconv.Itoa(int(user.Fixed)) {
			discountAmount = *plan.Discount
			finalPrice = plan.Price - discountAmount
		} else if plan.DiscountType == strconv.Itoa(int(user.Percentage)) {
			discountAmount = (plan.Price * (*plan.Discount)) / 100
			finalPrice = plan.Price - discountAmount
		}
	}

	// Ensure price doesn't go below zero
	if finalPrice < 0 {
		finalPrice = 0
	}

	return map[string]interface{}{
		"plan":           plan,
		"originalPrice":  plan.Price,
		"discountAmount": discountAmount,
		"finalPrice":     finalPrice,
	}, nil
}
