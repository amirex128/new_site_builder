package orderusecase

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/order"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type OrderUsecase struct {
	logger         sflogger.Logger
	orderRepo      repository.IOrderRepository
	basketRepo     repository.IBasketRepository
	orderItemRepo  repository.IOrderItemRepository
	paymentRepo    repository.IPaymentRepository
	authContextSvc common.IAuthContextService
}

func NewOrderUsecase(c contract.IContainer) *OrderUsecase {
	return &OrderUsecase{
		logger:         c.GetLogger(),
		orderRepo:      c.GetOrderRepo(),
		basketRepo:     c.GetBasketRepo(),
		orderItemRepo:  c.GetOrderItemRepo(),
		paymentRepo:    c.GetPaymentRepo(),
		authContextSvc: c.GetAuthContextService(),
	}
}

func (u *OrderUsecase) CreateOrderRequestCommand(params *order.CreateOrderRequestCommand) (any, error) {
	customerID, err := u.authContextSvc.GetCustomerID()
	if err != nil {
		return nil, err
	}

	siteID := *params.SiteID

	// Get the customer's basket
	basket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(customerID, siteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سبد خرید وجود ندارد")
		}
		return nil, err
	}

	if len(basket.Items) == 0 {
		return nil, errors.New("سبد خرید خالی است")
	}

	// Calculate courier price and total weight
	var courierPrice int64 = 100000 // Default courier price
	var totalWeight int = 0

	// Calculate total weight
	for _, item := range basket.Items {
		if item.Product != nil {
			totalWeight += item.Quantity * item.Product.Weight
		}
	}

	totalFinalPrice := basket.TotalPriceWithCouponDiscount + courierPrice

	// Check if an order already exists for this basket
	existingOrder, err := u.orderRepo.GetByID(0) // Need to add method to get by basket ID

	var newOrder domain.Order

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create a new order
		newOrder = domain.Order{
			SiteID:                       basket.SiteID,
			CustomerID:                   basket.CustomerID,
			TotalRawPrice:                basket.TotalRawPrice,
			TotalCouponDiscount:          basket.TotalCouponDiscount,
			TotalPriceWithCouponDiscount: basket.TotalPriceWithCouponDiscount,
			CourierPrice:                 courierPrice,
			OrderStatus:                  "WaitForPay",
			TotalFinalPrice:              totalFinalPrice,
			Description:                  "",
			TotalWeight:                  totalWeight,
			BasketID:                     basket.ID,
			DiscountID:                   basket.DiscountID,
			AddressID:                    *params.AddressID,
			CreatedAt:                    time.Now(),
			UpdatedAt:                    time.Now(),
		}

		if params.Description != nil {
			newOrder.Description = *params.Description
		}

		if params.Courier != nil {
			var courierName string
			switch *params.Courier {
			case order.Post:
				courierName = "Post"
			case order.Tipax:
				courierName = "Tipax"
			default:
				courierName = "Post"
			}
			newOrder.Courier = courierName
		} else {
			newOrder.Courier = "Post"
		}

		if err := u.orderRepo.Create(newOrder); err != nil {
			return nil, err
		}

		// Create order items
		for _, basketItem := range basket.Items {
			orderItem := domain.OrderItem{
				OrderID:                      newOrder.ID,
				ProductID:                    basketItem.ProductID,
				ProductVariantID:             basketItem.ProductVariantID,
				Quantity:                     basketItem.Quantity,
				RawPrice:                     basketItem.RawPrice,
				FinalRawPrice:                basketItem.FinalRawPrice,
				FinalPriceWithCouponDiscount: basketItem.FinalPriceWithCouponDiscount,
				JustCouponPrice:              basketItem.JustCouponPrice,
				JustDiscountPrice:            basketItem.JustDiscountPrice,
				CreatedAt:                    time.Now(),
				UpdatedAt:                    time.Now(),
			}

			if err := u.orderItemRepo.Create(orderItem); err != nil {
				return nil, err
			}
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update existing order
		existingOrder.TotalRawPrice = basket.TotalRawPrice
		existingOrder.TotalCouponDiscount = basket.TotalCouponDiscount
		existingOrder.TotalPriceWithCouponDiscount = basket.TotalPriceWithCouponDiscount
		existingOrder.CourierPrice = courierPrice
		existingOrder.OrderStatus = "WaitForPay"
		existingOrder.TotalFinalPrice = totalFinalPrice
		existingOrder.TotalWeight = totalWeight
		existingOrder.DiscountID = basket.DiscountID
		existingOrder.AddressID = *params.AddressID
		existingOrder.UpdatedAt = time.Now()

		if params.Description != nil {
			existingOrder.Description = *params.Description
		}

		if params.Courier != nil {
			var courierName string
			switch *params.Courier {
			case order.Post:
				courierName = "Post"
			case order.Tipax:
				courierName = "Tipax"
			default:
				courierName = "Post"
			}
			existingOrder.Courier = courierName
		}

		if err := u.orderRepo.Update(existingOrder); err != nil {
			return nil, err
		}

		// Delete existing order items and create new ones
		// TODO: Add delete order items by order ID method

		// Create order items
		for _, basketItem := range basket.Items {
			orderItem := domain.OrderItem{
				OrderID:                      existingOrder.ID,
				ProductID:                    basketItem.ProductID,
				ProductVariantID:             basketItem.ProductVariantID,
				Quantity:                     basketItem.Quantity,
				RawPrice:                     basketItem.RawPrice,
				FinalRawPrice:                basketItem.FinalRawPrice,
				FinalPriceWithCouponDiscount: basketItem.FinalPriceWithCouponDiscount,
				JustCouponPrice:              basketItem.JustCouponPrice,
				JustDiscountPrice:            basketItem.JustDiscountPrice,
				CreatedAt:                    time.Now(),
				UpdatedAt:                    time.Now(),
			}

			if err := u.orderItemRepo.Create(orderItem); err != nil {
				return nil, err
			}
		}

		newOrder = existingOrder
	}

	// Create order data for payment
	orderData := map[string]string{
		"OrderId": fmt.Sprintf("%d", newOrder.ID),
	}

	// Request payment gateway
	paymentURL, err := u.paymentRepo.RequestPayment(
		newOrder.TotalFinalPrice,
		newOrder.ID,
		customerID,
		strconv.Itoa(int(*params.Gateway)),
		orderData,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":    true,
		"paymentUrl": paymentURL,
	}, nil
}

func (u *OrderUsecase) CreateOrderVerifyCommand(params *order.CreateOrderVerifyCommand) (any, error) {
	// Extract order ID from order data
	orderIDStr, ok := params.OrderData["OrderId"]
	if !ok {
		return nil, errors.New("شناسه سفارش در داده های پرداخت یافت نشد")
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("شناسه سفارش نامعتبر است")
	}

	// Get order details
	order, err := u.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if order is already paid
	if order.OrderStatus == "Paid" {
		return nil, errors.New("سفارش قبلا پرداخت شده است")
	}

	// Get latest payment for this order
	payments, count, err := u.paymentRepo.GetAllByOrderID(orderID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("پرداختی برای این سفارش یافت نشد")
	}

	payment := payments[0]

	// Verify payment with gateway
	// This is a simplified implementation
	// In a real application, you would call the payment gateway's API
	isVerified := *params.IsSuccess

	if isVerified {
		// Update order status
		order.OrderStatus = "Paid"
		order.UpdatedAt = time.Now()

		if err := u.orderRepo.Update(order); err != nil {
			return nil, err
		}

		// Update payment status
		payment.PaymentStatusEnum = "Verified"
		payment.UpdatedAt = time.Now()

		if err := u.paymentRepo.Update(payment); err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"success": true,
			"message": "پرداخت با موفقیت انجام شد",
			"order":   order,
		}, nil
	}

	// Payment failed
	payment.PaymentStatusEnum = "Failed"
	payment.UpdatedAt = time.Now()

	if err := u.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": false,
		"message": "پرداخت ناموفق بود",
	}, nil
}

func (u *OrderUsecase) GetAllOrderCustomerQuery(params *order.GetAllOrderCustomerQuery) (any, error) {
	customerID, err := u.authContextSvc.GetCustomerID()
	if err != nil {
		return nil, err
	}

	orders, count, err := u.orderRepo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": orders,
		"total": count,
	}, nil
}

func (u *OrderUsecase) GetOrderCustomerDetailsQuery(params *order.GetOrderCustomerDetailsQuery) (any, error) {
	customerID, err := u.authContextSvc.GetCustomerID()
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.GetByID(*params.OrderID)
	if err != nil {
		return nil, err
	}

	// Verify the order belongs to the current customer
	if order.CustomerID != customerID {
		return nil, errors.New("سفارش متعلق به این کاربر نیست")
	}

	return order, nil
}

func (u *OrderUsecase) GetAllOrderUserQuery(params *order.GetAllOrderUserQuery) (any, error) {
	orders, count, err := u.orderRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": orders,
		"total": count,
	}, nil
}

func (u *OrderUsecase) GetOrderUserDetailsQuery(params *order.GetOrderUserDetailsQuery) (any, error) {
	order, err := u.orderRepo.GetByID(*params.OrderID)
	if err != nil {
		return nil, err
	}

	// TODO: Verify the order belongs to the user's site
	// In a real implementation, we would get the user's site ID and compare it with the order's site ID

	return order, nil
}

func (u *OrderUsecase) AdminGetAllOrderUserQuery(params *order.AdminGetAllOrderUserQuery) (any, error) {
	orders, count, err := u.orderRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": orders,
		"total": count,
	}, nil
}
