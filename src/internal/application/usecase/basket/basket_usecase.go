package basketusecase

import (
	"errors"
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/basket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type BasketUsecase struct {
	logger             sflogger.Logger
	basketRepo         repository.IBasketRepository
	basketItemRepo     repository.IBasketItemRepository
	productRepo        repository.IProductRepository
	productVariantRepo repository.IProductVariantRepository
	discountRepo       repository.IDiscountRepository
	authContextSvc     common.IAuthContextService
}

func NewBasketUsecase(c contract.IContainer) *BasketUsecase {
	return &BasketUsecase{
		logger:             c.GetLogger(),
		basketRepo:         c.GetBasketRepo(),
		basketItemRepo:     c.GetBasketItemRepo(),
		productRepo:        c.GetProductRepo(),
		productVariantRepo: c.GetProductVariantRepo(),
		discountRepo:       c.GetDiscountRepo(),
		authContextSvc:     c.GetAuthContextTransientService(),
	}
}

func (u *BasketUsecase) UpdateBasketCommand(params *basket.UpdateBasketCommand) (any, error) {
	u.logger.Info("UpdateBasketCommand called", map[string]interface{}{
		"params": params,
	})

	if params.BasketItems == nil || len(params.BasketItems) == 0 {
		return nil, errors.New("آیتم‌های سبد خرید الزامی هستند")
	}

	customerID, err := u.authContextSvc.GetCustomerID()
	if err != nil {
		return nil, err
	}

	siteID := *params.SiteID

	// Get or create the customer's basket for this site
	existingBasket, err := u.basketRepo.GetBasketByCustomerIDAndSiteID(customerID, siteID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Simple add mode just adds/updates the basket items without price calculation
	if params.SimpleAdd != nil && *params.SimpleAdd {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new basket
			newBasket := domain.Basket{
				SiteID:                       siteID,
				CustomerID:                   customerID,
				TotalRawPrice:                0,
				TotalCouponDiscount:          0,
				TotalPriceWithCouponDiscount: 0,
				CreatedAt:                    time.Now(),
				UpdatedAt:                    time.Now(),
			}

			if err := u.basketRepo.Create(newBasket); err != nil {
				return nil, err
			}

			// Get the newly created basket
			existingBasket, err = u.basketRepo.GetBasketByCustomerIDAndSiteID(customerID, siteID)
			if err != nil {
				return nil, err
			}
		}

		// Handle basket items
		for _, item := range params.BasketItems {
			basketItem := domain.BasketItem{
				BasketID:  existingBasket.ID,
				ProductID: *item.ProductID,
				Quantity:  *item.Quantity,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if item.ProductVariantID != nil {
				basketItem.ProductVariantID = *item.ProductVariantID
			}

			if item.BasketItemID != nil && *item.BasketItemID > 0 {
				// Update existing item
				basketItem.ID = *item.BasketItemID
				if err := u.basketItemRepo.Update(basketItem); err != nil {
					return nil, err
				}
			} else {
				// Create new item
				if err := u.basketItemRepo.Create(basketItem); err != nil {
					return nil, err
				}
			}
		}

		// Get the updated basket with items
		updatedBasket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(customerID, siteID)
		if err != nil {
			return nil, err
		}

		return updatedBasket, nil
	} else {
		// Complex mode with price calculation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سبد خرید وجود ندارد")
		}

		// Implement CalculateProductsPrice logic
		// Delete existing basket items
		if err := u.basketRepo.DeleteBasketItems(existingBasket.ID); err != nil {
			return nil, err
		}

		var totalRawPrice int64 = 0
		var totalCouponDiscount int64 = 0
		var totalPriceWithCouponDiscount int64 = 0

		// Process each basket item
		for _, item := range params.BasketItems {
			// Get product and variant to calculate prices
			product, err := u.productRepo.GetByID(*item.ProductID)
			if err != nil {
				return nil, err
			}

			var variantPrice int64 = 0
			var variant domain.ProductVariant
			if item.ProductVariantID != nil {
				// Get the product variant
				variant, err = u.productVariantRepo.GetByID(*item.ProductVariantID)
				if err != nil {
					return nil, err
				}
				variantPrice = variant.Price
			} else {
				// If no variant is specified, try to get the first variant for this product
				variants, count, err := u.productVariantRepo.GetAllByProductID(product.ID, common.PaginationRequestDto{
					Page:     1,
					PageSize: 1,
				})
				if err == nil && count > 0 {
					variantPrice = variants[0].Price
					variant = variants[0]
				}
			}

			if variantPrice == 0 {
				return nil, errors.New("قیمت محصول یا تنوع آن یافت نشد")
			}

			// Check if there's enough stock
			if variant.Stock < *item.Quantity {
				return nil, errors.New(fmt.Sprintf("موجودی محصول %s با تنوع %s کافی نیست", product.Name, variant.Name))
			}

			rawPrice := variantPrice
			finalRawPrice := rawPrice * int64(*item.Quantity)
			totalRawPrice += finalRawPrice

			// Apply discount if available
			var justCouponPrice int64 = 0
			var justDiscountPrice int64 = 0
			var finalPriceWithCouponDiscount int64 = finalRawPrice

			// Apply discount code if provided
			if params.Code != nil && *params.Code != "" {
				// Get discount by code
				discount, err := u.discountRepo.GetByCode(*params.Code)
				if err == nil && discount.ID > 0 {
					// Check if discount is valid
					if discount.ExpiryDate.After(time.Now()) && discount.Quantity > 0 {
						// Apply discount logic
						if discount.Type == "percentage" {
							discountAmount := (finalRawPrice * discount.Value) / 100
							justDiscountPrice = discountAmount
							finalPriceWithCouponDiscount = finalRawPrice - discountAmount
						} else if discount.Type == "fixed" {
							// Ensure discount doesn't exceed the price
							if discount.Value > finalRawPrice {
								justDiscountPrice = finalRawPrice
								finalPriceWithCouponDiscount = 0
							} else {
								justDiscountPrice = discount.Value
								finalPriceWithCouponDiscount = finalRawPrice - discount.Value
							}
						}

						// Set discount ID in the basket
						existingBasket.DiscountID = &discount.ID
					}
				}
			}

			// Apply product coupon if available
			// This would require additional implementation to fetch product coupons
			// For now, we'll skip this part as it's not clear from the current Go code structure
			// how product coupons are stored and related to products

			totalCouponDiscount += justCouponPrice + justDiscountPrice
			totalPriceWithCouponDiscount += finalPriceWithCouponDiscount

			// Create the basket item
			basketItem := domain.BasketItem{
				BasketID:                     existingBasket.ID,
				ProductID:                    *item.ProductID,
				Quantity:                     *item.Quantity,
				RawPrice:                     rawPrice,
				FinalRawPrice:                finalRawPrice,
				JustCouponPrice:              justCouponPrice,
				JustDiscountPrice:            justDiscountPrice,
				FinalPriceWithCouponDiscount: finalPriceWithCouponDiscount,
				CreatedAt:                    time.Now(),
				UpdatedAt:                    time.Now(),
			}

			if item.ProductVariantID != nil {
				basketItem.ProductVariantID = *item.ProductVariantID
			}

			if err := u.basketItemRepo.Create(basketItem); err != nil {
				return nil, err
			}
		}

		// Update the basket with calculated values
		existingBasket.TotalRawPrice = totalRawPrice
		existingBasket.TotalCouponDiscount = totalCouponDiscount
		existingBasket.TotalPriceWithCouponDiscount = totalPriceWithCouponDiscount
		existingBasket.UpdatedAt = time.Now()

		if err := u.basketRepo.Update(existingBasket); err != nil {
			return nil, err
		}

		// Get the updated basket with items
		updatedBasket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(customerID, siteID)
		if err != nil {
			return nil, err
		}

		return updatedBasket, nil
	}
}

func (u *BasketUsecase) GetBasketQuery(params *basket.GetBasketQuery) (any, error) {
	u.logger.Info("GetBasketQuery called", map[string]interface{}{
		"params": params,
	})

	customerID, err := u.authContextSvc.GetCustomerID()
	if err != nil {
		return nil, err
	}

	basket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(customerID, *params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return an empty basket
			return map[string]interface{}{
				"id":                           0,
				"siteId":                       *params.SiteID,
				"customerId":                   customerID,
				"totalRawPrice":                0,
				"totalCouponDiscount":          0,
				"totalPriceWithCouponDiscount": 0,
				"items":                        []interface{}{},
			}, nil
		}
		return nil, err
	}

	return basket, nil
}

func (u *BasketUsecase) GetAllBasketUserQuery(params *basket.GetAllBasketUserQuery) (any, error) {
	u.logger.Info("GetAllBasketUserQuery called", map[string]interface{}{
		"params": params,
	})

	customerID, err := u.authContextSvc.GetCustomerID()
	if err != nil {
		return nil, err
	}

	baskets, count, err := u.basketRepo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": baskets,
		"total": count,
	}, nil
}

func (u *BasketUsecase) AdminGetAllBasketUserQuery(params *basket.AdminGetAllBasketUserQuery) (any, error) {
	u.logger.Info("AdminGetAllBasketUserQuery called", map[string]interface{}{
		"params": params,
	})

	result, count, err := u.basketRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
