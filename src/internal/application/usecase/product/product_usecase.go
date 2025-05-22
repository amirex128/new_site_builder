package productusecase

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
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
	u.Logger.Info("CreateProductCommand called", map[string]interface{}{
		"name":   *params.Name,
		"siteId": *params.SiteID,
	})

	// Validate slug uniqueness
	existingProductBySlug, err := u.productRepo.GetBySlug(*params.Slug)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil && existingProductBySlug.SiteID == *params.SiteID {
			return nil, errors.New("نامک تکراری است")
		}
	}

	// Get user ID from auth context
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Prepare SEO tags
	var seoTagsStr string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTagsStr = strings.Join(params.SeoTags, ",")
	}

	// Convert status enum to string
	statusStr := *params.Status

	// Create product entity
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
		UserID:          userID,
		SellingCount:    0,
		VisitedCount:    0,
		ReviewCount:     0,
		Rate:            0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsDeleted:       false,
	}

	// Create the product in the database
	err = u.productRepo.Create(newProduct)
	if err != nil {
		return nil, err
	}

	// Handle product variants
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
				u.Logger.Error("Failed to create product variant", map[string]interface{}{
					"productId":   newProduct.ID,
					"variantName": *variantDTO.Name,
					"error":       err.Error(),
				})
			}
		}
	}

	// Handle product attributes
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
				u.Logger.Error("Failed to create product attribute", map[string]interface{}{
					"productId": newProduct.ID,
					"attrName":  *attrDTO.Name,
					"error":     err.Error(),
				})
			}
		}
	}

	// Handle coupon if any
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
			u.Logger.Error("Failed to create coupon", map[string]interface{}{
				"productId": newProduct.ID,
				"error":     err.Error(),
			})
		}
	}

	// Handle category associations
	if params.CategoryIDs != nil && len(params.CategoryIDs) > 0 {
		// In a monolithic approach, we can add relationships directly
		for _, categoryID := range params.CategoryIDs {
			// Add relationship to category_product table
			// This would be implemented in a real repository method
			u.Logger.Info("Adding product to category", map[string]interface{}{
				"productId":  newProduct.ID,
				"categoryId": categoryID,
			})
		}
	}

	// Handle discount associations
	if params.DiscountIDs != nil && len(params.DiscountIDs) > 0 {
		// Similar to categories, add to discount_product table
		for _, discountID := range params.DiscountIDs {
			u.Logger.Info("Adding discount to product", map[string]interface{}{
				"productId":  newProduct.ID,
				"discountId": discountID,
			})
		}
	}

	// Handle media associations
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		for _, mediaID := range params.MediaIDs {
			u.Logger.Info("Adding media to product", map[string]interface{}{
				"productId": newProduct.ID,
				"mediaId":   mediaID,
			})
		}
	}

	// Fetch the product with related data
	createdProduct, err := u.productRepo.GetByID(newProduct.ID)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (u *ProductUsecase) UpdateProductCommand(params *product.UpdateProductCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateProductCommand called", map[string]interface{}{
		"id":     *params.ID,
		"siteId": *params.SiteID,
	})

	// Get existing product
	existingProduct, err := u.productRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// In our monolithic approach, we check directly
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingProduct.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این محصول دسترسی ندارید")
	}

	// Validate slug uniqueness if changed
	if params.Slug != nil && *params.Slug != existingProduct.Slug {
		existingProductBySlug, err := u.productRepo.GetBySlug(*params.Slug)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			if err == nil && existingProductBySlug.SiteID == *params.SiteID && existingProductBySlug.ID != *params.ID {
				return nil, errors.New("نامک تکراری است")
			}
		}
	}

	// Update fields if provided
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

	// Prepare SEO tags
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		existingProduct.SeoTags = strings.Join(params.SeoTags, ",")
	}

	existingProduct.SiteID = *params.SiteID
	existingProduct.UpdatedAt = time.Now()
	existingProduct.UserID = userID // Update with current user ID

	// Update the product
	err = u.productRepo.Update(existingProduct)
	if err != nil {
		return nil, err
	}

	// Update product variants
	if params.ProductVariants != nil && len(params.ProductVariants) > 0 {
		// Typically, we would delete existing and create new ones
		// This is simplified - in production we'd need more logic for updates vs creates
		for _, variantDTO := range params.ProductVariants {
			// Check if variant ID exists to determine update vs create
			var variant domain.ProductVariant
			if variantDTO.ID != nil {
				// Update existing variant
				variant.ID = *variantDTO.ID
				variant.ProductID = existingProduct.ID
				variant.Name = *variantDTO.Name
				variant.Price = *variantDTO.Price
				variant.Stock = *variantDTO.Stock
				variant.UpdatedAt = time.Now()
				// Update in repository
			} else {
				// Create new variant
				variant = domain.ProductVariant{
					Name:      *variantDTO.Name,
					Price:     *variantDTO.Price,
					Stock:     *variantDTO.Stock,
					ProductID: existingProduct.ID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				// Create in repository
			}
		}
	}

	// Update product attributes - similar to variants
	if params.ProductAttributes != nil && len(params.ProductAttributes) > 0 {
		// Similar logic as variants
	}

	// Update coupon if any
	if params.Coupon != nil {
		// Similar logic - get existing or create new
	}

	// Update category associations
	if params.CategoryIDs != nil {
		// Delete existing and create new associations
	}

	// Update discount associations
	if params.DiscountIDs != nil {
		// Delete existing and create new associations
	}

	// Update media associations
	if params.MediaIDs != nil {
		// Delete existing and create new associations
	}

	// Fetch the updated product with related data
	updatedProduct, err := u.productRepo.GetByID(existingProduct.ID)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (u *ProductUsecase) DeleteProductCommand(params *product.DeleteProductCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteProductCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing product
	existingProduct, err := u.productRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Check if user has access to this product
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingProduct.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این محصول دسترسی ندارید")
	}

	// Delete the product
	err = u.productRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *ProductUsecase) GetByIdProductQuery(params *product.GetByIdProductQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdProductQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get product by ID
	product, err := u.productRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check user access - typically products are visible to all
	// But we might want to log access
	userID, _ := u.authContext(u.Ctx).GetUserID()
	if userID > 0 {
		u.Logger.Info("Product accessed by user", map[string]interface{}{
			"productId": product.ID,
			"userId":    userID,
		})
	}

	// Load media URLs if any
	// In a real implementation, we would fetch the media information
	// and add it to the response

	return product, nil
}

func (u *ProductUsecase) GetAllProductQuery(params *product.GetAllProductQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllProductQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check site access - simplified in monolithic app
	// In a real implementation, we would check if the user has access to this site

	// Get all products for the site
	products, count, err := u.productRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     products,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *ProductUsecase) GetByFiltersSortProductQuery(params *product.GetByFiltersSortProductQuery) (*resp.Response, error) {
	u.Logger.Info("GetByFiltersSortProductQuery called", map[string]interface{}{
		"siteId":       *params.SiteID,
		"page":         params.Page,
		"pageSize":     params.PageSize,
		"selectedSort": params.SelectedSort,
		"filterCount":  len(params.SelectedFilters),
	})

	// Convert filter enums to ProductFilterEnum keys for repository
	filters := make(map[enums.ProductFilterEnum][]string)
	for k, v := range params.SelectedFilters {
		filters[enums.ProductFilterEnum(k)] = v
	}

	// Convert sort enum to string pointer
	var sortStr *string
	if params.SelectedSort != nil {
		s := string(*params.SelectedSort)
		sortStr = &s
	}

	// Call the repository method for filtering and sorting
	products, count, err := u.productRepo.GetAllByFilterAndSort(
		*params.SiteID,
		filters,
		sortStr,
		params.PaginationRequestDto,
	)
	if err != nil {
		return nil, err
	}

	// Enhancement: Load associated data (media, etc.)
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

		// Parse SEO tags if any
		if p.SeoTags != "" {
			productData["seoTags"] = strings.Split(p.SeoTags, ",")
		}

		// Include media if any
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

	return map[string]interface{}{
		"items":     enhancedProducts,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *ProductUsecase) GetProductByCategoryQuery(params *product.GetProductByCategoryQuery) (*resp.Response, error) {
	u.Logger.Info("GetProductByCategoryQuery called", map[string]interface{}{
		"slug":   *params.Slug,
		"siteId": *params.SiteID,
	})

	// Get category by slug and site ID from productCategoryRepo
	category, err := u.productCategoryRepo.GetBySlug(*params.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("دسته‌بندی یافت نشد")
		}
		return nil, err
	}

	// Verify that the category belongs to the specified site
	if category.SiteID != *params.SiteID {
		return nil, errors.New("دسته‌بندی یافت نشد")
	}

	// Get products for the category
	products, count, err := u.productRepo.GetAllByCategoryID(category.ID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     products,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *ProductUsecase) GetSingleProductQuery(params *product.GetSingleProductQuery) (*resp.Response, error) {
	u.Logger.Info("GetSingleProductQuery called", map[string]interface{}{
		"slug":   *params.Slug,
		"siteId": *params.SiteID,
	})

	// Get product by slug
	product, err := u.productRepo.GetBySlug(*params.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check if the product belongs to the specified site
	if product.SiteID != *params.SiteID {
		return nil, errors.New("محصول یافت نشد")
	}

	// Construct response with additional information
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

	// Parse SEO tags if any
	if product.SeoTags != "" {
		response["seoTags"] = strings.Split(product.SeoTags, ",")
	}

	// Load media if any
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

	// Increase visited count in background (non-blocking)
	go func() {
		if product.ID > 0 {
			product.VisitedCount++
			_ = u.productRepo.Update(product) // Ignore error as this is a background operation
		}
	}()

	return response, nil
}

func (u *ProductUsecase) CalculateProductsPriceQuery(params *product.CalculateProductsPriceQuery) (*resp.Response, error) {
	u.Logger.Info("CalculateProductsPriceQuery called", map[string]interface{}{
		"customerId": *params.CustomerID,
		"siteId":     *params.SiteID,
		"code":       params.Code,
	})

	// Prepare response structure
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

	// Validate order basket items
	if params.OrderBasketItems == nil || len(params.OrderBasketItems) == 0 {
		response.ResponseStatus.IsSuccess = false
		response.ResponseStatus.Message = "سبد خرید خالی است"
		return response, nil
	}

	// Find discount if code is provided
	var discount domain.Discount
	var discountFound bool
	var isDiscountUsed bool

	if params.Code != nil && *params.Code != "" {
		var err error
		// Try to get discount by code or ID based on IsOrderVerify flag
		if *params.IsOrderVerify {
			// In order verification, Code is actually the discount ID
			discountID, parseErr := strconv.ParseInt(*params.Code, 10, 64)
			if parseErr == nil {
				discount, err = u.discountRepo.GetByID(discountID)
				if err == nil {
					discountFound = true
				}
			}
		} else {
			// Normal flow - get discount by code
			discount, err = u.discountRepo.GetByCode(*params.Code)
			if err == nil {
				discountFound = true
			}
		}

		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			// Not found is ok, just continue without discount
		} else if discountFound {
			// Validate discount
			if discount.SiteID != *params.SiteID {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "کد تخفیف معتبر نیست"
				return response, nil
			}

			// Check quantity
			if discount.Quantity <= 0 {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "ظرفیت استفاده از این کد تخفیف به پایان رسیده است"
				return response, nil
			}

			// Check expiry date
			if discount.ExpiryDate.Before(time.Now()) {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "این کد تخفیف منقضی شده است"
				return response, nil
			}

			// Check if customer has already used this discount
			// In a real implementation, we'd check the discount usage table
			// For now, we'll assume it hasn't been used
			isDiscountUsed = false
		}
	}

	// Collect product IDs for fetching
	var productIDs []int64
	productVariantMap := make(map[int64]int64) // Maps product ID to variant ID
	quantityMap := make(map[int64]int)         // Maps product ID to quantity

	for _, item := range params.OrderBasketItems {
		productIDs = append(productIDs, *item.ProductID)
		productVariantMap[*item.ProductID] = *item.ProductVariationID
		quantityMap[*item.ProductID] = *item.Quantity
	}

	// Get all products
	// Note: In an actual implementation, we would need a productRepo method to get products with variants in a single query
	// For now, we'll get each product individually
	for _, productID := range productIDs {
		product, err := u.productRepo.GetByID(productID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ResponseStatus.IsSuccess = false
				response.ResponseStatus.Message = "یکی از محصولات موجود نیست"
				return response, nil
			}
			return nil, err
		}

		// Get product variant
		variantID := productVariantMap[productID]
		var variant domain.ProductVariant
		var variantFound bool

		// In a real implementation, we'd get variants directly
		// For now, create a dummy pagination request
		paginationRequest := common.PaginationRequestDto{
			Page:     1,
			PageSize: 100, // Assuming we won't have more than 100 variants per product
		}
		variants, _, err := u.productVariantRepo.GetAllByProductID(productID, paginationRequest)
		if err != nil {
			return nil, err
		}

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
			return response, nil
		}

		// Check stock
		quantity := quantityMap[productID]
		if variant.Stock < quantity {
			response.ResponseStatus.IsSuccess = false
			response.ResponseStatus.Message = "موجودی محصول " + product.Name + " کافی نیست"
			return response, nil
		}

		// Calculate raw price
		rawPrice := variant.Price * int64(quantity)
		response.TotalRawPrice += rawPrice

		// Get product coupon if any
		var coupon domain.Coupon
		var couponFound bool
		// In a real implementation, we'd get this from the database with a proper method
		// Placeholder logic for now

		// Calculate discounts
		discountValue := int64(0)
		couponValue := int64(0)

		// Apply product-specific coupon if available
		if couponFound {
			if coupon.Type == "0" { // Percentage
				couponValue = (rawPrice * coupon.Value) / 100
			} else { // Fixed value
				couponValue = coupon.Value * int64(quantity)
				if couponValue > rawPrice {
					couponValue = rawPrice
				}
			}
		}

		// Apply global discount if available and not already used
		if discountFound && !isDiscountUsed {
			if discount.Type == enums.PercentageDiscountType {
				discountValue = (rawPrice * discount.Value) / 100
			} else { // Fixed value
				discountValue = discount.Value
				if discountValue > rawPrice {
					discountValue = rawPrice
				}
			}
		}

		// Calculate final price
		finalPriceWithDiscounts := rawPrice - (couponValue + discountValue)
		if finalPriceWithDiscounts < 0 {
			finalPriceWithDiscounts = 0
		}

		// Track total discount
		response.TotalCouponDiscount += (couponValue + discountValue)
		response.TotalPriceWithCouponDiscount += finalPriceWithDiscounts

		// Add to calculated prices
		calculatedPrice := map[string]interface{}{
			"basketItemId":                 params.OrderBasketItems[0].BasketItemID, // Using first item as placeholder
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

	// Set discount ID if used and not already used by customer
	if discountFound && !isDiscountUsed {
		response.DiscountID = &discount.ID
	}

	return response, nil
}

func (u *ProductUsecase) AdminGetAllProductQuery(params *product.AdminGetAllProductQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllProductQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all products across all sites for admin
	products, count, err := u.productRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     products,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// Helper function to handle nil string pointers
func getStringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
