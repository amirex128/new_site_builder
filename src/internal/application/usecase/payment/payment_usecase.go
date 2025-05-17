package paymentusecase

import (
	"fmt"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/payment"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type PaymentUsecase struct {
	logger      sflogger.Logger
	paymentRepo repository.IPaymentRepository
	gatewayRepo repository.IGatewayRepository
}

func NewPaymentUsecase(c contract.IContainer) *PaymentUsecase {
	return &PaymentUsecase{
		logger:      c.GetLogger(),
		paymentRepo: c.GetPaymentRepo(),
		gatewayRepo: c.GetGatewayRepo(),
	}
}

func (u *PaymentUsecase) VerifyPaymentCommand(params *payment.VerifyPaymentCommand) (any, error) {
	// Implementation for verifying a payment
	fmt.Println(params)

	// TODO: Implement proper payment verification logic
	// This might involve:
	// 1. Getting the payment from the database
	// 2. Calling the payment gateway's API to verify the payment
	// 3. Updating the payment status in the database

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *PaymentUsecase) CreateOrUpdateGatewayCommand(params *payment.CreateOrUpdateGatewayCommand) (any, error) {
	// Implementation for creating or updating a gateway
	fmt.Println(params)

	// TODO: Implement proper gateway creation/update logic
	// This would typically involve creating or updating a gateway record in the database

	// Placeholder response
	return map[string]interface{}{
		"success": true,
		"id":      1, // Placeholder ID
	}, nil
}

func (u *PaymentUsecase) GetByIdGatewayQuery(params *payment.GetByIdGatewayQuery) (any, error) {
	// Implementation to get a gateway by ID
	fmt.Println(params)

	gateway, err := u.gatewayRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return gateway, nil
}

func (u *PaymentUsecase) AdminGetAllGatewayQuery(params *payment.AdminGetAllGatewayQuery) (any, error) {
	// Implementation to get all gateways for admin
	fmt.Println(params)

	// Use the GetAll method from the gateway repository
	// Note: IGatewayRepository interface might not have a GetAll method defined
	// You may need to implement it in the repository
	// For now, using a placeholder response

	return map[string]interface{}{
		"items": []interface{}{},
		"total": 0,
	}, nil
}

func (u *PaymentUsecase) AdminGetAllPaymentQuery(params *payment.AdminGetAllPaymentQuery) (any, error) {
	// Implementation to get all payments for admin
	fmt.Println(params)

	result, count, err := u.paymentRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
