package orderusecase

import (
	"errors"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/order"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type OrderUsecase struct {
	*usecase.BaseUsecase
	orderRepo     repository.IOrderRepository
	basketRepo    repository.IBasketRepository
	orderItemRepo repository.IOrderItemRepository
	paymentRepo   repository.IPaymentRepository
	container     contract.IContainer
}

func NewOrderUsecase(c contract.IContainer) *OrderUsecase {
	return &OrderUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		orderRepo:     c.GetOrderRepo(),
		basketRepo:    c.GetBasketRepo(),
		orderItemRepo: c.GetOrderItemRepo(),
		paymentRepo:   c.GetPaymentRepo(),
		container:     c,
	}
}

func (u *OrderUsecase) CreateOrderRequestCommand(params *order.CreateOrderRequestCommand) (*resp.Response, error) {
	customerID, err := u.AuthContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, err
	}
	basket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(*customerID, *params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "سبد خرید وجود ندارد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if len(basket.Items) == 0 {
		return nil, resp.NewError(resp.BadRequest, "سبد خرید خالی است")
	}
	var courierPrice int64 = 100000
	totalWeight := 0
	for _, item := range basket.Items {
		if item.Product != nil {
			totalWeight += item.Quantity * item.Product.Weight
		}
	}
	totalFinalPrice := basket.TotalPriceWithCouponDiscount + courierPrice
	existingOrder, err := u.orderRepo.GetByID(0)
	var newOrder domain.Order
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
			newOrder.Courier = string(*params.Courier)
		} else {
			newOrder.Courier = "Post"
		}
		err = u.orderRepo.Create(&newOrder)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
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
			err = u.orderItemRepo.Create(&orderItem)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	} else if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	} else {
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
			existingOrder.Courier = string(*params.Courier)
		}
		err = u.orderRepo.Update(existingOrder)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		err = u.orderItemRepo.DeleteByOrderID(existingOrder.ID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
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
			err = u.orderItemRepo.Create(&orderItem)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
		newOrder = *existingOrder
	}
	orderData := map[string]string{
		"OrderId": strconv.FormatInt(newOrder.ID, 10),
	}
	paymentURL, err := u.paymentRepo.RequestPayment(
		newOrder.TotalFinalPrice,
		newOrder.ID,
		*customerID,
		string(*params.Gateway),
		orderData,
	)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Success, resp.Data{
		"success":    true,
		"paymentUrl": paymentURL,
		"IsDeleted":  false,
	}, "پرداخت با موفقیت انجام شد"), nil
}

func (u *OrderUsecase) CreateOrderVerifyCommand(params *order.CreateOrderVerifyCommand) (*resp.Response, error) {
	orderIDStr, ok := params.OrderData["OrderId"]
	if !ok {
		return nil, resp.NewError(resp.BadRequest, "شناسه سفارش در داده های پرداخت یافت نشد")
	}
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return nil, resp.NewError(resp.BadRequest, "شناسه سفارش نامعتبر است")
	}
	order, err := u.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	if order.OrderStatus == "Paid" {
		return nil, resp.NewError(resp.BadRequest, "سفارش قبلا پرداخت شده است")
	}
	paymentsResult, err := u.paymentRepo.GetAllByOrderID(orderID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if len(paymentsResult.Items) == 0 {
		return nil, resp.NewError(resp.NotFound, "پرداختی برای این سفارش یافت نشد")
	}
	payment := paymentsResult.Items[0]
	isVerified := *params.IsSuccess
	if isVerified {
		order.OrderStatus = "Paid"
		order.UpdatedAt = time.Now()
		err = u.orderRepo.Update(order)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		payment.PaymentStatusEnum = "Verified"
		payment.UpdatedAt = time.Now()
		err = u.paymentRepo.Update(&payment)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		return resp.NewResponseData(resp.Success, map[string]interface{}{
			"success": true,
			"message": "پرداخت با موفقیت انجام شد",
			"order":   order,
		}, "پرداخت با موفقیت انجام شد"), nil
	}
	payment.PaymentStatusEnum = "Failed"
	payment.UpdatedAt = time.Now()
	err = u.paymentRepo.Update(&payment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Success, resp.Data{
		"success": false,
		"message": "پرداخت ناموفق بود",
	}, "پرداخت ناموفق بود"), nil
}

func (u *OrderUsecase) GetAllOrderCustomerQuery(params *order.GetAllOrderCustomerQuery) (*resp.Response, error) {
	customerID, err := u.AuthContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, err
	}
	ordersResult, err := u.orderRepo.GetAllByCustomerID(*customerID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     ordersResult.Items,
		"total":     ordersResult.TotalCount,
		"page":      ordersResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": ordersResult.TotalPages,
	}, "لیست سفارشات با موفقیت دریافت شد"), nil
}

func (u *OrderUsecase) GetOrderCustomerDetailsQuery(params *order.GetOrderCustomerDetailsQuery) (*resp.Response, error) {
	customerID, err := u.AuthContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, err
	}
	order, err := u.orderRepo.GetByID(*params.OrderID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	if order.CustomerID != *customerID {
		return nil, resp.NewError(resp.Unauthorized, "سفارش متعلق به این کاربر نیست")
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"order": order}, "جزئیات سفارش با موفقیت دریافت شد"), nil
}

func (u *OrderUsecase) GetAllOrderUserQuery(params *order.GetAllOrderUserQuery) (*resp.Response, error) {
	ordersResult, err := u.orderRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     ordersResult.Items,
		"total":     ordersResult.TotalCount,
		"page":      ordersResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": ordersResult.TotalPages,
	}, "لیست سفارشات کاربر با موفقیت دریافت شد"), nil
}

func (u *OrderUsecase) GetOrderUserDetailsQuery(params *order.GetOrderUserDetailsQuery) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	order, err := u.orderRepo.GetByID(*params.OrderID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	siteRepo := u.container.GetSiteRepo()
	if siteRepo == nil {
		return nil, resp.NewError(resp.Internal, "site repository not available")
	}
	userSitesResult, err := siteRepo.GetAllByUserID(*userID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	hasAccess := false
	for _, site := range userSitesResult.Items {
		if site.ID == order.SiteID {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
		if err != nil || !isAdmin {
			return nil, resp.NewError(resp.Unauthorized, "شما به این سفارش دسترسی ندارید")
		}
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"order": order}, "جزئیات سفارش با موفقیت دریافت شد"), nil
}

func (u *OrderUsecase) AdminGetAllOrderUserQuery(params *order.AdminGetAllOrderUserQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, err
	}
	ordersResult, err := u.orderRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     ordersResult.Items,
		"total":     ordersResult.TotalCount,
		"page":      ordersResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": ordersResult.TotalPages,
	}, "لیست سفارشات ادمین با موفقیت دریافت شد"), nil
}
