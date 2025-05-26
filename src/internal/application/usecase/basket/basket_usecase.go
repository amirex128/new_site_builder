package basketusecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/basket"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type BasketUsecase struct {
	*usecase.BaseUsecase
	basketRepo         repository.IBasketRepository
	basketItemRepo     repository.IBasketItemRepository
	productRepo        repository.IProductRepository
	productVariantRepo repository.IProductVariantRepository
	discountRepo       repository.IDiscountRepository
	authContext        func(c *gin.Context) service.IAuthService
}

func NewBasketUsecase(c contract.IContainer) *BasketUsecase {
	return &BasketUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		basketRepo:         c.GetBasketRepo(),
		basketItemRepo:     c.GetBasketItemRepo(),
		productRepo:        c.GetProductRepo(),
		productVariantRepo: c.GetProductVariantRepo(),
		discountRepo:       c.GetDiscountRepo(),
		authContext:        c.GetAuthTransientService(),
	}
}

func (u *BasketUsecase) UpdateBasketCommand(params *basket.UpdateBasketCommand) (*resp.Response, error) {
	if params.BasketItems == nil || len(params.BasketItems) == 0 {
		return nil, resp.NewError(resp.BadRequest, "آیتم‌های سبد خرید الزامی هستند")
	}
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	hasBasket := false
	existingBasket, err := u.basketRepo.GetBasketByCustomerIDAndSiteID(*customerID, *params.SiteID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.Internal, "خطا در دریافت سبد خرید")
		}
	} else {
		hasBasket = true
	}

	if params.SimpleAdd != nil && *params.SimpleAdd {
		if !hasBasket {
			existingBasket = &domain.Basket{
				SiteID:                       *params.SiteID,
				CustomerID:                   *customerID,
				TotalRawPrice:                0,
				TotalCouponDiscount:          0,
				TotalPriceWithCouponDiscount: 0,
				CreatedAt:                    time.Now(),
				UpdatedAt:                    time.Now(),
			}
			if err := u.basketRepo.Create(existingBasket); err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در ایجاد سبد خرید")
			}
		}
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
				basketItem.ID = *item.BasketItemID
				if err := u.basketItemRepo.Update(&basketItem); err != nil {
					return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی آیتم سبد خرید")
				}
			} else {
				if err := u.basketItemRepo.Create(&basketItem); err != nil {
					return nil, resp.NewError(resp.Internal, "خطا در افزودن آیتم به سبد خرید")
				}
			}
		}
		updatedBasket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(*customerID, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در دریافت سبد خرید")
		}
		return resp.NewResponseData(resp.Updated, updatedBasket, "سبد خرید با موفقیت بروزرسانی شد"), nil
	} else {
		if !hasBasket {
			return nil, resp.NewError(resp.NotFound, "سبد خرید وجود ندارد")
		}
		if err := u.basketRepo.DeleteBasketItems(existingBasket.ID); err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در حذف آیتم‌های سبد خرید")
		}
		var totalRawPrice int64 = 0
		var totalCouponDiscount int64 = 0
		var totalPriceWithCouponDiscount int64 = 0
		for _, item := range params.BasketItems {
			product, err := u.productRepo.GetByID(*item.ProductID)
			if err != nil {
				return nil, resp.NewError(resp.NotFound, "محصول یافت نشد")
			}
			var variantPrice int64 = 0
			var variant *domain.ProductVariant
			if item.ProductVariantID != nil {
				variant, err = u.productVariantRepo.GetByID(*item.ProductVariantID)
				if err != nil {
					return nil, resp.NewError(resp.NotFound, "تنوع محصول یافت نشد")
				}
				variantPrice = variant.Price
			} else {
				variants, err := u.productVariantRepo.GetAllByProductID(product.ID)
				if err == nil && len(variants) > 0 {
					variantPrice = variants[0].Price
					variant = &variants[0]
				}
			}
			if variant == nil || variantPrice == 0 {
				return nil, resp.NewError(resp.BadRequest, "قیمت محصول یا تنوع آن یافت نشد")
			}
			if variant.Stock < *item.Quantity {
				return nil, resp.NewError(resp.BadRequest, fmt.Sprintf("موجودی محصول %s با تنوع %s کافی نیست", product.Name, variant.Name))
			}
			rawPrice := variantPrice
			finalRawPrice := rawPrice * int64(*item.Quantity)
			totalRawPrice += finalRawPrice
			var justCouponPrice int64 = 0
			var justDiscountPrice int64 = 0
			var finalPriceWithCouponDiscount int64 = finalRawPrice
			if params.Code != nil && *params.Code != "" {
				discount, err := u.discountRepo.GetByCode(*params.Code)
				if err == nil && discount.ID > 0 {
					if discount.ExpiryDate.After(time.Now()) && discount.Quantity > 0 {
						if discount.Type == "percentage" {
							discountAmount := (finalRawPrice * discount.Value) / 100
							justDiscountPrice = discountAmount
							finalPriceWithCouponDiscount = finalRawPrice - discountAmount
						} else if discount.Type == "fixed" {
							if discount.Value > finalRawPrice {
								justDiscountPrice = finalRawPrice
								finalPriceWithCouponDiscount = 0
							} else {
								justDiscountPrice = discount.Value
								finalPriceWithCouponDiscount = finalRawPrice - discount.Value
							}
						}
						existingBasket.DiscountID = &discount.ID
					}
				}
			}
			totalCouponDiscount += justCouponPrice + justDiscountPrice
			totalPriceWithCouponDiscount += finalPriceWithCouponDiscount
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
				IsDeleted:                    false,
			}
			if item.ProductVariantID != nil {
				basketItem.ProductVariantID = *item.ProductVariantID
			}
			if err := u.basketItemRepo.Create(&basketItem); err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در افزودن آیتم به سبد خرید")
			}
		}
		existingBasket.TotalRawPrice = totalRawPrice
		existingBasket.TotalCouponDiscount = totalCouponDiscount
		existingBasket.TotalPriceWithCouponDiscount = totalPriceWithCouponDiscount
		existingBasket.UpdatedAt = time.Now()
		if err := u.basketRepo.Update(existingBasket); err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی سبد خرید")
		}
		updatedBasket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(*customerID, *params.SiteID)
		if err != nil {
			return nil, resp.NewError(resp.NotFound, "خطا در دریافت سبد خرید")
		}
		return resp.NewResponseData(resp.Updated, updatedBasket, "سبد خرید با موفقیت بروزرسانی شد"), nil
	}
}

func (u *BasketUsecase) GetBasketQuery(params *basket.GetBasketQuery) (*resp.Response, error) {
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	basket, err := u.basketRepo.GetBasketWithItemsByCustomerIDAndSiteID(*customerID, *params.SiteID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, basket, "سبد خرید با موفقیت دریافت شد"), nil
}

func (u *BasketUsecase) GetAllBasketUserQuery(params *basket.GetAllBasketUserQuery) (*resp.Response, error) {
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	basketsResult, err := u.basketRepo.GetAllByCustomerID(*customerID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت سبدهای خرید")
	}
	return resp.NewResponseData(resp.Retrieved, basketsResult, "سبدهای خرید با موفقیت دریافت شدند"), nil
}

func (u *BasketUsecase) AdminGetAllBasketUserQuery(params *basket.AdminGetAllBasketUserQuery) (*resp.Response, error) {
	result, err := u.basketRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت سبدهای خرید")
	}
	return resp.NewResponseData(resp.Retrieved, result, "سبدهای خرید با موفقیت دریافت شدند (مدیر)"), nil
}
