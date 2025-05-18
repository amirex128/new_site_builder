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
	roleRepo repository.IRoleRepository
	userRepo repository.IUserRepository
}

func NewPlanUsecase(c contract.IContainer) *PlanUsecase {
	return &PlanUsecase{
		logger:   c.GetLogger(),
		planRepo: c.GetPlanRepo(),
		roleRepo: c.GetRoleRepo(),
		userRepo: c.GetUserRepo(),
	}
}

// CreatePlanCommand creates a new plan
func (u *PlanUsecase) CreatePlanCommand(params *plan.CreatePlanCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()
	// In Golang, we would implement this with middleware or a service check

	// Prepare values with null handling
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

	// Create the plan entity
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

	// In .NET, there was an event added: entity.AddEvent(new PlanCreatedEventStore(entity))
	// In our monolith approach, we don't need to emit events

	// Save the plan
	err := u.planRepo.Create(newPlan)
	if err != nil {
		return nil, err
	}

	// Return the created plan
	return newPlan, nil
}

// UpdatePlanCommand updates an existing plan
func (u *PlanUsecase) UpdatePlanCommand(params *plan.UpdatePlanCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Get the existing plan
	existingPlan, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
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

	// In .NET, there was an event added: entity.AddEvent(new PlanUpdatedEventStore(entity))
	// In our monolith approach, we don't need to emit events

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
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Check if plan exists
	_, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Delete the plan
	err = u.planRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

// GetByIdPlanQuery gets a plan by ID
func (u *PlanUsecase) GetByIdPlanQuery(params *plan.GetByIdPlanQuery) (any, error) {
	// In .NET, there was a check for user access: gate.HasUserAccess(entity)
	// We would implement this with middleware or a service check

	result, err := u.planRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAllPlanQuery gets all plans with pagination
func (u *PlanUsecase) GetAllPlanQuery(params *plan.GetAllPlanQuery) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	result, count, err := u.planRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

// CalculatePlanPriceQuery calculates the price of a plan
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
