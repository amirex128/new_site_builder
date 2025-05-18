package planusecase

import (
	"fmt"
	"strconv"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/plan"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type PlanUsecase struct {
	logger      sflogger.Logger
	planRepo    repository.IPlanRepository
	userRepo    repository.IUserRepository
	authContext common.IAuthContextService
}

func NewPlanUsecase(c contract.IContainer) *PlanUsecase {
	return &PlanUsecase{
		logger:      c.GetLogger(),
		planRepo:    c.GetPlanRepo(),
		userRepo:    c.GetUserRepo(),
		authContext: c.GetAuthContextService(),
	}
}

// CreatePlanCommand creates a new plan
func (u *PlanUsecase) CreatePlanCommand(params *plan.CreatePlanCommand) (any, error) {
	// Check admin access
	isAdmin, err := u.authContext.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("only admins can create plans")
	}

	// Create new plan
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

	// Set optional fields if provided
	if params.DiscountType != nil {
		newPlan.DiscountType = strconv.Itoa(int(*params.DiscountType))
	}

	if params.Discount != nil {
		discount := int64(*params.Discount)
		newPlan.Discount = &discount
	}

	if params.Feature != nil {
		newPlan.Feature = *params.Feature
	}

	// Create the plan
	err = u.planRepo.Create(newPlan)
	if err != nil {
		return nil, err
	}

	return newPlan, nil
}

// UpdatePlanCommand updates an existing plan
func (u *PlanUsecase) UpdatePlanCommand(params *plan.UpdatePlanCommand) (any, error) {
	// Check admin access
	isAdmin, err := u.authContext.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("only admins can update plans")
	}

	// Get the existing plan
	existingPlan, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
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
		existingPlan.DiscountType = strconv.Itoa(int(*params.DiscountType))
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

	// Update the plan
	err = u.planRepo.Update(existingPlan)
	if err != nil {
		return nil, err
	}

	return existingPlan, nil
}

// DeletePlanCommand deletes a plan
func (u *PlanUsecase) DeletePlanCommand(params *plan.DeletePlanCommand) (any, error) {
	// Check admin access
	isAdmin, err := u.authContext.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("only admins can delete plans")
	}

	// Check if the plan exists
	_, err = u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Delete the plan
	err = u.planRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Plan deleted successfully",
	}, nil
}

// GetAllPlanQuery gets all plans with pagination
func (u *PlanUsecase) GetAllPlanQuery(params *plan.GetAllPlanQuery) (any, error) {
	// Get all plans
	plans, count, err := u.planRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": plans,
		"total": count,
	}, nil
}

// GetByIDPlanQuery gets a plan by ID
func (u *PlanUsecase) GetByIDPlanQuery(params *plan.GetByIDPlanQuery) (any, error) {
	// Get the plan
	plan, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

// CalculatePlanPriceQuery calculates a plan's price with discounts
func (u *PlanUsecase) CalculatePlanPriceQuery(params *plan.CalculatePlanPriceQuery) (any, error) {
	// Get the plan
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
