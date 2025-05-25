package productusecase

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"gorm.io/gorm"
)

type ProductUsecase struct {
	*usecase.BaseUsecase
	productRepo          repository.IProductRepository
	productCategoryRepo  repository.IProductCategoryRepository
	discountRepo         repository.IDiscountRepository
	mediaRepo            repository.IMediaRepository
	productVariantRepo   repository.IProductVariantRepository
	productAttributeRepo repository.IProductAttributeRepository
	couponRepo           repository.ICouponRepository
	authContext          func(c *gin.Context) service.IAuthService
}

func NewProductUsecase(c contract.IContainer) *ProductUsecase {
	return &ProductUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		productRepo:          c.GetProductRepo(),
		productCategoryRepo:  c.GetProductCategoryRepo(),
		discountRepo:         c.GetDiscountRepo(),
		mediaRepo:            c.GetMediaRepo(),
		productVariantRepo:   c.GetProductVariantRepo(),
		productAttributeRepo: c.GetProductAttributeRepo(),
		couponRepo:           c.GetCouponRepo(),
		authContext:          c.GetAuthTransientService(),
	}
}

func (u *ProductUsecase) CreateProductCommand(params *product.CreateProductCommand) (*resp.Response, error) {
	existingProductBySlug, err := u.productRepo.GetBySlug(*params.Slug)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil && existingProductBySlug.SiteID == *params.SiteID {
			return nil, resp.NewError(resp.BadRequest, "نامک تکراری است")
		}
	}
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	var seoTagsStr string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTagsStr = strings.Join(params.SeoTags, ",")
	}
	statusStr := *params.Status
	newProduct := domain.Product{
		Name:            *params.Name,
		Description:     getStringValueOrEmpty(params.Description),
		Status:          statusStr,
		Weight:          *params.Weight,
		FreeSend:        *params.FreeSend,
		LongDescription: getStringValueOrEmpty(params.LongDescription),
		Slug:            *params.Slug,
		SeoTags:         seoTagsStr,
		SiteID:          *params.SiteID,
		UserID:          *userID,
		SellingCount:    0,
		VisitedCount:    0,
		ReviewCount:     0,
		Rate:            0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsDeleted:       false,
	}
	err = u.productRepo.Create(newProduct)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.ProductVariants != nil && len(params.ProductVariants) > 0 {
		for _, variantDTO := range params.ProductVariants {
			variant := domain.ProductVariant{
				Name:      *variantDTO.Name,
				Price:     *variantDTO.Price,
				Stock:     *variantDTO.Stock,
				ProductID: newProduct.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = u.productVariantRepo.Create(variant)
			if err != nil {
				continue
			}
		}
	}
	if params.ProductAttributes != nil && len(params.ProductAttributes) > 0 {
		for _, attrDTO := range params.ProductAttributes {
			attr := domain.ProductAttribute{
				Name:      *attrDTO.Name,
				Value:     *attrDTO.Value,
				Type:      *attrDTO.Type,
				ProductID: newProduct.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = u.productAttributeRepo.Create(attr)
			if err != nil {
				continue
			}
		}
	}
	if params.Coupon != nil {
		coupon := domain.Coupon{
			Quantity:   *params.Coupon.Quantity,
			Type:       *params.Coupon.Type,
			Value:      *params.Coupon.Value,
			ExpiryDate: *params.Coupon.ExpiryDate,
			ProductID:  newProduct.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = u.couponRepo.Create(coupon)
		if err != nil {
			// continue
		}
	}
	createdProduct, err := u.productRepo.GetByID(newProduct.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Created, resp.Data{
		"product": createdProduct,
	}, "محصول با موفقیت ایجاد شد"), nil
}

func (u *ProductUsecase) UpdateProductCommand(params *product.UpdateProductCommand) (*resp.Response, error) {
	existingProduct, err := u.productRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "محصول یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if userID != nil && existingProduct.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این محصول دسترسی ندارید")
	}
	if params.Slug != nil && *params.Slug != existingProduct.Slug {
		existingProductBySlug, err := u.productRepo.GetBySlug(*params.Slug)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			if err == nil && existingProductBySlug.SiteID == *params.SiteID && existingProductBySlug.ID != *params.ID {
				return nil, resp.NewError(resp.BadRequest, "نامک تکراری است")
			}
		}
	}
	if params.Name != nil {
		existingProduct.Name = *params.Name
	}
	if params.Description != nil {
		existingProduct.Description = *params.Description
	}
	if params.Status != nil {
		existingProduct.Status = *params.Status
	}
	if params.Weight != nil {
		existingProduct.Weight = *params.Weight
	}
	if params.FreeSend != nil {
		existingProduct.FreeSend = *params.FreeSend
	}
	if params.LongDescription != nil {
		existingProduct.LongDescription = *params.LongDescription
	}
	if params.Slug != nil {
		existingProduct.Slug = *params.Slug
	}
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		existingProduct.SeoTags = strings.Join(params.SeoTags, ",")
	}
	existingProduct.SiteID = *params.SiteID
	existingProduct.UpdatedAt = time.Now()
	existingProduct.UserID = *userID
	err = u.productRepo.Update(existingProduct)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	updatedProduct, err := u.productRepo.GetByID(existingProduct.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Updated, resp.Data{
		"product": updatedProduct,
	}, "محصول با موفقیت بروزرسانی شد"), nil
}

func (u *ProductUsecase) DeleteProductCommand(params *product.DeleteProductCommand) (*resp.Response, error) {
	existingProduct, err := u.productRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "محصول یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if userID != nil && existingProduct.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این محصول دسترسی ندارید")
	}
	err = u.productRepo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Deleted, resp.Data{
		"success": true,
	}, "محصول با موفقیت حذف شد"), nil
}

func (u *ProductUsecase) GetByIdProductQuery(params *product.GetByIdProductQuery) (*resp.Response, error) {
	product, err := u.productRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "محصول یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"product": product,
	}, "محصول با موفقیت دریافت شد"), nil
}

func (u *ProductUsecase) GetAllProductQuery(params *product.GetAllProductQuery) (*resp.Response, error) {
	result, err := u.productRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	products := result.Items
	count := result.TotalCount
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     products,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "محصولات با موفقیت دریافت شدند"), nil
}

func (u *ProductUsecase) GetByFiltersSortProductQuery(params *product.GetByFiltersSortProductQuery) (*resp.Response, error) {
	filters := make(map[enums.ProductFilterEnum][]string)
	for k, v := range params.SelectedFilters {
		filters[enums.ProductFilterEnum(k)] = v
	}
	var sortStr *string
	if params.SelectedSort != nil {
		s := string(*params.SelectedSort)
		sortStr = &s
	}
	result, err := u.productRepo.GetAllByFilterAndSort(
		*params.SiteID,
		filters,
		sortStr,
		params.PaginationRequestDto,
	)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	products := result.Items
	count := result.TotalCount
	var enhancedProducts []map[string]interface{}
	for _, p := range products {
		productData := map[string]interface{}{
			"id":              p.ID,
			"name":            p.Name,
			"description":     p.Description,
			"status":          p.Status,
			"weight":          p.Weight,
			"sellingCount":    p.SellingCount,
			"visitedCount":    p.VisitedCount,
			"reviewCount":     p.ReviewCount,
			"rate":            p.Rate,
			"freeSend":        p.FreeSend,
			"longDescription": p.LongDescription,
			"slug":            p.Slug,
			"siteId":          p.SiteID,
			"createdAt":       p.CreatedAt,
			"updatedAt":       p.UpdatedAt,
		}
		if p.SeoTags != "" {
			productData["seoTags"] = strings.Split(p.SeoTags, ",")
		}
		if len(p.Media) > 0 {
			var mediaItems []map[string]interface{}
			for _, media := range p.Media {
				mediaItems = append(mediaItems, map[string]interface{}{
					"id":  media.ID,
					"url": "/api/media/" + strconv.FormatInt(media.ID, 10),
				})
			}
			productData["media"] = mediaItems
		}
		enhancedProducts = append(enhancedProducts, productData)
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     enhancedProducts,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "محصولات با موفقیت دریافت شدند"), nil
}

func (u *ProductUsecase) GetProductByCategoryQuery(params *product.GetProductByCategoryQuery) (*resp.Response, error) {
	category, err := u.productCategoryRepo.GetBySlug(*params.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "دسته‌بندی یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if category.SiteID != *params.SiteID {
		return nil, resp.NewError(resp.NotFound, "دسته‌بندی یافت نشد")
	}
	result, err := u.productRepo.GetAllByCategoryID(category.ID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	products := result.Items
	count := result.TotalCount
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     products,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "محصولات با موفقیت دریافت شدند"), nil
}

func (u *ProductUsecase) GetSingleProductQuery(params *product.GetSingleProductQuery) (*resp.Response, error) {
	product, err := u.productRepo.GetBySlug(*params.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "محصول یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if product.SiteID != *params.SiteID {
		return nil, resp.NewError(resp.NotFound, "محصول یافت نشد")
	}
	response := map[string]interface{}{
		"id":              product.ID,
		"name":            product.Name,
		"description":     product.Description,
		"status":          product.Status,
		"weight":          product.Weight,
		"sellingCount":    product.SellingCount,
		"visitedCount":    product.VisitedCount,
		"reviewCount":     product.ReviewCount,
		"rate":            product.Rate,
		"freeSend":        product.FreeSend,
		"longDescription": product.LongDescription,
		"slug":            product.Slug,
		"siteId":          product.SiteID,
		"createdAt":       product.CreatedAt,
		"updatedAt":       product.UpdatedAt,
	}
	if product.SeoTags != "" {
		response["seoTags"] = strings.Split(product.SeoTags, ",")
	}
	if len(product.Media) > 0 {
		var mediaItems []map[string]interface{}
		for _, media := range product.Media {
			mediaItems = append(mediaItems, map[string]interface{}{
				"id":  media.ID,
				"url": "/api/media/" + strconv.FormatInt(media.ID, 10),
			})
		}
		response["media"] = mediaItems
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data(response), "محصول با موفقیت دریافت شد"), nil
}

func (u *ProductUsecase) CalculateProductsPriceQuery(params *product.CalculateProductsPriceQuery) (*resp.Response, error) {
	response := struct {
		CalculatedPrices             []map[string]interface{} `json:"calculatedPrices"`
		TotalRawPrice                int64                    `json:"totalRawPrice"`
		TotalCouponDiscount          int64                    `json:"totalCouponDiscount"`
		TotalPriceWithCouponDiscount int64                    `json:"totalPriceWithCouponDiscount"`
		DiscountID                   *int64                   `json:"discountId"`
		ResponseStatus               struct {
			IsSuccess bool   `json:"isSuccess"`
			Message   string `json:"message"`
		} `json:"responseStatus"`
	}{
		CalculatedPrices: []map[string]interface{}{},
		ResponseStatus: struct {
			IsSuccess bool   `json:"isSuccess"`
			Message   string `json:"message"`
		}{
			IsSuccess: true,
			Message:   "Products price calculation successful.",
		},
	}
	if params.OrderBasketItems == nil || len(params.OrderBasketItems) == 0 {
		response.ResponseStatus.IsSuccess = false
		response.ResponseStatus.Message = "سبد خرید خالی است"
		return resp.NewResponseData(resp.Retrieved, resp.Data{
			"calculatedPrices":             response.CalculatedPrices,
			"totalRawPrice":                response.TotalRawPrice,
			"totalCouponDiscount":          response.TotalCouponDiscount,
			"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
			"discountId":                   response.DiscountID,
			"responseStatus":               response.ResponseStatus,
		}, response.ResponseStatus.Message), nil
	}
	var discount domain.Discount
	var discountFound bool
	var isDiscountUsed bool
	if params.Code != nil && *params.Code != "" {
		var err error
		if *params.IsOrderVerify {
			discountID, parseErr := strconv.ParseInt(*params.Code, 10, 64)
			if parseErr == nil {
				discount, err = u.discountRepo.GetByID(discountID)
				if err == nil {
					discountFound = true
				}
			}
		} else {
			discount, err = u.discountRepo.GetByCode(*params.Code)
			if err == nil {
				discountFound = true
			}
		}
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		} else if discountFound {
			if discount.SiteID != *params.SiteID {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "کد تخفیف معتبر نیست"
				return resp.NewResponseData(resp.Retrieved, resp.Data{
					"calculatedPrices":             response.CalculatedPrices,
					"totalRawPrice":                response.TotalRawPrice,
					"totalCouponDiscount":          response.TotalCouponDiscount,
					"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
					"discountId":                   response.DiscountID,
					"responseStatus":               response.ResponseStatus,
				}, response.ResponseStatus.Message), nil
			}
			if discount.Quantity <= 0 {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "ظرفیت استفاده از این کد تخفیف به پایان رسیده است"
				return resp.NewResponseData(resp.Retrieved, resp.Data{
					"calculatedPrices":             response.CalculatedPrices,
					"totalRawPrice":                response.TotalRawPrice,
					"totalCouponDiscount":          response.TotalCouponDiscount,
					"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
					"discountId":                   response.DiscountID,
					"responseStatus":               response.ResponseStatus,
				}, response.ResponseStatus.Message), nil
			}
			if discount.ExpiryDate.Before(time.Now()) {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "این کد تخفیف منقضی شده است"
				return resp.NewResponseData(resp.Retrieved, resp.Data{
					"calculatedPrices":             response.CalculatedPrices,
					"totalRawPrice":                response.TotalRawPrice,
					"totalCouponDiscount":          response.TotalCouponDiscount,
					"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
					"discountId":                   response.DiscountID,
					"responseStatus":               response.ResponseStatus,
				}, response.ResponseStatus.Message), nil
			}
			isDiscountUsed = false
		}
	}
	var productIDs []int64
	productVariantMap := make(map[int64]int64)
	quantityMap := make(map[int64]int)
	for _, item := range params.OrderBasketItems {
		productIDs = append(productIDs, *item.ProductID)
		productVariantMap[*item.ProductID] = *item.ProductVariationID
		quantityMap[*item.ProductID] = *item.Quantity
	}
	for _, productID := range productIDs {
		product, err := u.productRepo.GetByID(productID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "یکی از محصولات موجود نیست"
				return resp.NewResponseData(resp.Retrieved, resp.Data{
					"calculatedPrices":             response.CalculatedPrices,
					"totalRawPrice":                response.TotalRawPrice,
					"totalCouponDiscount":          response.TotalCouponDiscount,
					"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
					"discountId":                   response.DiscountID,
					"responseStatus":               response.ResponseStatus,
				}, response.ResponseStatus.Message), nil
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		variantID := productVariantMap[productID]
		var variant domain.ProductVariant
		var variantFound bool
		paginationRequest := common.PaginationRequestDto{
			Page:     1,
			PageSize: 100,
		}
		variantsResult, err := u.productVariantRepo.GetAllByProductID(productID, paginationRequest)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		variants := variantsResult.Items
		for _, v := range variants {
			if v.ID == variantID {
				variant = v
				variantFound = true
				break
			}
		}
		if !variantFound {
			response.ResponseStatus.IsSuccess = false
			response.ResponseStatus.Message = "یکی از تنوع های محصول موجود نیست"
			return resp.NewResponseData(resp.Retrieved, resp.Data{
				"calculatedPrices":             response.CalculatedPrices,
				"totalRawPrice":                response.TotalRawPrice,
				"totalCouponDiscount":          response.TotalCouponDiscount,
				"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
				"discountId":                   response.DiscountID,
				"responseStatus":               response.ResponseStatus,
			}, response.ResponseStatus.Message), nil
		}
		quantity := quantityMap[productID]
		if variant.Stock < quantity {
			response.ResponseStatus.IsSuccess = false
			response.ResponseStatus.Message = "موجودی محصول " + product.Name + " کافی نیست"
			return resp.NewResponseData(resp.Retrieved, resp.Data{
				"calculatedPrices":             response.CalculatedPrices,
				"totalRawPrice":                response.TotalRawPrice,
				"totalCouponDiscount":          response.TotalCouponDiscount,
				"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
				"discountId":                   response.DiscountID,
				"responseStatus":               response.ResponseStatus,
			}, response.ResponseStatus.Message), nil
		}
		rawPrice := variant.Price * int64(quantity)
		response.TotalRawPrice += rawPrice
		var coupon domain.Coupon
		var couponFound bool
		discountValue := int64(0)
		couponValue := int64(0)
		if couponFound {
			if coupon.Type == "0" {
				couponValue = (rawPrice * coupon.Value) / 100
			} else {
				couponValue = coupon.Value * int64(quantity)
				if couponValue > rawPrice {
					couponValue = rawPrice
				}
			}
		}
		if discountFound && !isDiscountUsed {
			if discount.Type == enums.PercentageDiscountType {
				discountValue = (rawPrice * discount.Value) / 100
			} else {
				discountValue = discount.Value
				if discountValue > rawPrice {
					discountValue = rawPrice
				}
			}
		}
		finalPriceWithDiscounts := rawPrice - (couponValue + discountValue)
		if finalPriceWithDiscounts < 0 {
			finalPriceWithDiscounts = 0
		}
		response.TotalCouponDiscount += (couponValue + discountValue)
		response.TotalPriceWithCouponDiscount += finalPriceWithDiscounts
		calculatedPrice := map[string]interface{}{
			"basketItemId":                 params.OrderBasketItems[0].BasketItemID,
			"rawPrice":                     variant.Price,
			"quantity":                     quantity,
			"finalRawPrice":                rawPrice,
			"justCouponPrice":              couponValue,
			"justDiscountPrice":            discountValue,
			"finalPriceWithCouponDiscount": finalPriceWithDiscounts,
			"productId":                    product.ID,
			"productVariantId":             variant.ID,
			"freeSend":                     product.FreeSend,
			"weight":                       product.Weight,
		}
		response.CalculatedPrices = append(response.CalculatedPrices, calculatedPrice)
	}
	if discountFound && !isDiscountUsed {
		response.DiscountID = &discount.ID
	}
	responseData := resp.Data{
		"calculatedPrices":             response.CalculatedPrices,
		"totalRawPrice":                response.TotalRawPrice,
		"totalCouponDiscount":          response.TotalCouponDiscount,
		"totalPriceWithCouponDiscount": response.TotalPriceWithCouponDiscount,
		"discountId":                   response.DiscountID,
		"responseStatus":               response.ResponseStatus,
	}
	return resp.NewResponseData(resp.Success, responseData, "محاسبه قیمت محصولات با موفقیت انجام شد"), nil
}

func (u *ProductUsecase) AdminGetAllProductQuery(params *product.AdminGetAllProductQuery) (*resp.Response, error) {
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}
	result, err := u.productRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	products := result.Items
	count := result.TotalCount
	return resp.NewResponseData(resp.Retrieved, resp.Data{
		"items":     products,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, "محصولات با موفقیت دریافت شدند"), nil
}

func getStringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
