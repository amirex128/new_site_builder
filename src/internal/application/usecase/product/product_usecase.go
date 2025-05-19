package productusecase

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type ProductUsecase struct {
	logger               sflogger.Logger
	repo                 repository.IProductRepository
	productCategoryRepo  repository.IProductCategoryRepository
	discountRepo         repository.IDiscountRepository
	mediaRepo            repository.IMediaRepository
	productVariantRepo   repository.IProductVariantRepository
	productAttributeRepo repository.IProductAttributeRepository
	couponRepo           repository.ICouponRepository
	authContextSvc       common.IAuthContextService
}

func NewProductUsecase(c contract.IContainer) *ProductUsecase {
	return &ProductUsecase{
		logger:               c.GetLogger(),
		repo:                 c.GetProductRepo(),
		productCategoryRepo:  c.GetProductCategoryRepo(),
		discountRepo:         c.GetDiscountRepo(),
		mediaRepo:            c.GetMediaRepo(),
		productVariantRepo:   c.GetProductVariantRepo(),
		productAttributeRepo: c.GetProductAttributeRepo(),
		couponRepo:           c.GetCouponRepo(),
		authContextSvc:       c.GetAuthContextTransientService(),
	}
}

func (u *ProductUsecase) CreateProductCommand(params *product.CreateProductCommand) (any, error) {
	u.logger.Info("CreateProductCommand called", map[string]interface{}{
		"name":   *params.Name,
		"siteId": *params.SiteID,
	})

	// Validate slug uniqueness
	existingProductBySlug, err := u.repo.GetBySlug(*params.Slug)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil && existingProductBySlug.SiteID == *params.SiteID {
			return nil, errors.New("نامک تکراری است")
		}
	}

	// Get user ID from auth context
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Prepare SEO tags
	var seoTagsStr string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTagsStr = strings.Join(params.SeoTags, ",")
	}

	// Convert status enum to string
	statusStr := strconv.Itoa(int(*params.Status))

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
	err = u.repo.Create(newProduct)
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
				u.logger.Error("Failed to create product variant", map[string]interface{}{
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
				Type:      strconv.Itoa(int(*attrDTO.Type)),
				ProductID: newProduct.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = u.productAttributeRepo.Create(attr)
			if err != nil {
				u.logger.Error("Failed to create product attribute", map[string]interface{}{
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
			Type:       strconv.Itoa(int(*params.Coupon.Type)),
			Value:      *params.Coupon.Value,
			ExpiryDate: *params.Coupon.ExpiryDate,
			ProductID:  newProduct.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = u.couponRepo.Create(coupon)
		if err != nil {
			u.logger.Error("Failed to create coupon", map[string]interface{}{
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
			u.logger.Info("Adding product to category", map[string]interface{}{
				"productId":  newProduct.ID,
				"categoryId": categoryID,
			})
		}
	}

	// Handle discount associations
	if params.DiscountIDs != nil && len(params.DiscountIDs) > 0 {
		// Similar to categories, add to discount_product table
		for _, discountID := range params.DiscountIDs {
			u.logger.Info("Adding discount to product", map[string]interface{}{
				"productId":  newProduct.ID,
				"discountId": discountID,
			})
		}
	}

	// Handle media associations
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		for _, mediaID := range params.MediaIDs {
			u.logger.Info("Adding media to product", map[string]interface{}{
				"productId": newProduct.ID,
				"mediaId":   mediaID,
			})
		}
	}

	// Fetch the product with related data
	createdProduct, err := u.repo.GetByID(newProduct.ID)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (u *ProductUsecase) UpdateProductCommand(params *product.UpdateProductCommand) (any, error) {
	u.logger.Info("UpdateProductCommand called", map[string]interface{}{
		"id":     *params.ID,
		"siteId": *params.SiteID,
	})

	// Get existing product
	existingProduct, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// In our monolithic approach, we check directly
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingProduct.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این محصول دسترسی ندارید")
	}

	// Validate slug uniqueness if changed
	if params.Slug != nil && *params.Slug != existingProduct.Slug {
		existingProductBySlug, err := u.repo.GetBySlug(*params.Slug)
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
		existingProduct.Status = strconv.Itoa(int(*params.Status))
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
	err = u.repo.Update(existingProduct)
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
	updatedProduct, err := u.repo.GetByID(existingProduct.ID)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (u *ProductUsecase) DeleteProductCommand(params *product.DeleteProductCommand) (any, error) {
	u.logger.Info("DeleteProductCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing product
	existingProduct, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Check if user has access to this product
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingProduct.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این محصول دسترسی ندارید")
	}

	// Delete the product
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *ProductUsecase) GetByIdProductQuery(params *product.GetByIdProductQuery) (any, error) {
	u.logger.Info("GetByIdProductQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get product by ID
	product, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Check user access - typically products are visible to all
	// But we might want to log access
	userID, _ := u.authContextSvc.GetUserID()
	if userID > 0 {
		u.logger.Info("Product accessed by user", map[string]interface{}{
			"productId": product.ID,
			"userId":    userID,
		})
	}

	// Load media URLs if any
	// In a real implementation, we would fetch the media information
	// and add it to the response

	return product, nil
}

func (u *ProductUsecase) GetAllProductQuery(params *product.GetAllProductQuery) (any, error) {
	u.logger.Info("GetAllProductQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check site access - simplified in monolithic app
	// In a real implementation, we would check if the user has access to this site

	// Get all products for the site
	products, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
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

func (u *ProductUsecase) GetByFiltersSortProductQuery(params *product.GetByFiltersSortProductQuery) (any, error) {
	u.logger.Info("GetByFiltersSortProductQuery called", map[string]interface{}{
		"siteId":       *params.SiteID,
		"page":         params.Page,
		"pageSize":     params.PageSize,
		"selectedSort": params.SelectedSort,
		"filterCount":  len(params.SelectedFilters),
	})

	// Build query filters based on selected filters
	queryBuilder := &filterQueryBuilder{
		siteID: *params.SiteID,
		logger: u.logger,
	}

	// Apply filters if any are provided
	if params.SelectedFilters != nil && len(params.SelectedFilters) > 0 {
		// Process all filter types
		for filterType, values := range params.SelectedFilters {
			switch filterType {
			case product.PriceRange:
				if len(values) == 2 {
					minPrice, errMin := strconv.ParseInt(values[0], 10, 64)
					maxPrice, errMax := strconv.ParseInt(values[1], 10, 64)
					if errMin == nil && errMax == nil {
						queryBuilder.addPriceRangeFilter(minPrice, maxPrice)
					}
				}
			case product.RatingRange:
				if len(values) == 2 {
					minRating, errMin := strconv.Atoi(values[0])
					maxRating, errMax := strconv.Atoi(values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addRatingRangeFilter(minRating, maxRating)
					}
				}
			case product.SellingRange:
				if len(values) == 2 {
					minSelling, errMin := strconv.Atoi(values[0])
					maxSelling, errMax := strconv.Atoi(values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addSellingRangeFilter(minSelling, maxSelling)
					}
				}
			case product.VisitedRange:
				if len(values) == 2 {
					minVisited, errMin := strconv.Atoi(values[0])
					maxVisited, errMax := strconv.Atoi(values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addVisitedRangeFilter(minVisited, maxVisited)
					}
				}
			case product.ReviewRange:
				if len(values) == 2 {
					minReview, errMin := strconv.Atoi(values[0])
					maxReview, errMax := strconv.Atoi(values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addReviewRangeFilter(minReview, maxReview)
					}
				}
			case product.WeightRange:
				if len(values) == 2 {
					minWeight, errMin := strconv.Atoi(values[0])
					maxWeight, errMax := strconv.Atoi(values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addWeightRangeFilter(minWeight, maxWeight)
					}
				}
			case product.CategoryIds:
				if len(values) > 0 {
					var categoryIDs []int64
					for _, val := range values {
						id, err := strconv.ParseInt(val, 10, 64)
						if err == nil {
							categoryIDs = append(categoryIDs, id)
						}
					}
					if len(categoryIDs) > 0 {
						queryBuilder.addCategoryIDsFilter(categoryIDs)
					}
				}
			case product.ProductIds:
				if len(values) > 0 {
					var productIDs []int64
					for _, val := range values {
						id, err := strconv.ParseInt(val, 10, 64)
						if err == nil {
							productIDs = append(productIDs, id)
						}
					}
					if len(productIDs) > 0 {
						queryBuilder.addProductIDsFilter(productIDs)
					}
				}
			case product.FreeSend:
				if len(values) > 0 {
					freeSend, err := strconv.ParseBool(values[0])
					if err == nil {
						queryBuilder.addFreeSendFilter(freeSend)
					}
				}
			case product.UpdatedRange:
				if len(values) == 2 {
					minUpdated, errMin := time.Parse(time.RFC3339, values[0])
					maxUpdated, errMax := time.Parse(time.RFC3339, values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addUpdatedRangeFilter(minUpdated, maxUpdated)
					}
				}
			case product.AddedRange:
				if len(values) == 2 {
					minAdded, errMin := time.Parse(time.RFC3339, values[0])
					maxAdded, errMax := time.Parse(time.RFC3339, values[1])
					if errMin == nil && errMax == nil {
						queryBuilder.addAddedRangeFilter(minAdded, maxAdded)
					}
				}
			}
		}
	}

	// Apply sorting if specified
	if params.SelectedSort != nil {
		queryBuilder.addSorting(*params.SelectedSort)
	}

	// In a real implementation, we would use the query builder to construct a complex SQL query
	// For now, we'll query the repository and then filter/sort in memory for demonstration
	products, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Apply filters and sorting (simplistic implementation for demonstration)
	// In a real implementation, this would be handled by the database query
	var filteredProducts []domain.Product
	for _, p := range products {
		// In a real implementation, we would filter based on queryBuilder parameters
		// For this example, we'll include all products
		filteredProducts = append(filteredProducts, p)
	}

	// Sort products based on selected sort (if any)
	if params.SelectedSort != nil {
		// In a real implementation, this would be handled by the database ORDER BY
		switch *params.SelectedSort {
		case product.NameAZ:
			sort.Slice(filteredProducts, func(i, j int) bool {
				return filteredProducts[i].Name < filteredProducts[j].Name
			})
		case product.NameZA:
			sort.Slice(filteredProducts, func(i, j int) bool {
				return filteredProducts[i].Name > filteredProducts[j].Name
			})
		case product.RecentlyAdded:
			sort.Slice(filteredProducts, func(i, j int) bool {
				return filteredProducts[i].CreatedAt.After(filteredProducts[j].CreatedAt)
			})
		case product.MostVisited:
			sort.Slice(filteredProducts, func(i, j int) bool {
				return filteredProducts[i].VisitedCount > filteredProducts[j].VisitedCount
			})
		}
	}

	// Enhancement: Load associated data (variants, categories, etc.)
	var enhancedProducts []map[string]interface{}
	for _, p := range filteredProducts {
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

	// Apply pagination to the enhanced products
	// In a real implementation, this would be handled by the database LIMIT/OFFSET
	startIndex := int((params.Page - 1) * params.PageSize)
	endIndex := int(params.Page * params.PageSize)
	totalProducts := len(enhancedProducts)

	// Adjust indices if they are out of bounds
	if startIndex >= totalProducts {
		startIndex = 0
		endIndex = 0
	}
	if endIndex > totalProducts {
		endIndex = totalProducts
	}

	// Get paginated slice of products
	var paginatedProducts []map[string]interface{}
	if startIndex < totalProducts {
		paginatedProducts = enhancedProducts[startIndex:endIndex]
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     paginatedProducts,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// filterQueryBuilder is a helper struct to build SQL queries for product filtering
type filterQueryBuilder struct {
	siteID       int64
	priceRange   []int64
	ratingRange  []int
	sellingRange []int
	visitedRange []int
	reviewRange  []int
	weightRange  []int
	categoryIDs  []int64
	productIDs   []int64
	freeSend     *bool
	updatedRange []time.Time
	addedRange   []time.Time
	selectedSort *product.ProductSortEnum
	logger       sflogger.Logger
}

func (qb *filterQueryBuilder) addPriceRangeFilter(min, max int64) {
	qb.priceRange = []int64{min, max}
	qb.logger.Info("Added price range filter", map[string]interface{}{
		"min": min,
		"max": max,
	})
}

func (qb *filterQueryBuilder) addRatingRangeFilter(min, max int) {
	qb.ratingRange = []int{min, max}
}

func (qb *filterQueryBuilder) addSellingRangeFilter(min, max int) {
	qb.sellingRange = []int{min, max}
}

func (qb *filterQueryBuilder) addVisitedRangeFilter(min, max int) {
	qb.visitedRange = []int{min, max}
}

func (qb *filterQueryBuilder) addReviewRangeFilter(min, max int) {
	qb.reviewRange = []int{min, max}
}

func (qb *filterQueryBuilder) addWeightRangeFilter(min, max int) {
	qb.weightRange = []int{min, max}
}

func (qb *filterQueryBuilder) addCategoryIDsFilter(categoryIDs []int64) {
	qb.categoryIDs = categoryIDs
}

func (qb *filterQueryBuilder) addProductIDsFilter(productIDs []int64) {
	qb.productIDs = productIDs
}

func (qb *filterQueryBuilder) addFreeSendFilter(freeSend bool) {
	qb.freeSend = &freeSend
}

func (qb *filterQueryBuilder) addUpdatedRangeFilter(min, max time.Time) {
	qb.updatedRange = []time.Time{min, max}
}

func (qb *filterQueryBuilder) addAddedRangeFilter(min, max time.Time) {
	qb.addedRange = []time.Time{min, max}
}

func (qb *filterQueryBuilder) addSorting(sortType product.ProductSortEnum) {
	qb.selectedSort = &sortType
}

func (u *ProductUsecase) GetProductByCategoryQuery(params *product.GetProductByCategoryQuery) (any, error) {
	u.logger.Info("GetProductByCategoryQuery called", map[string]interface{}{
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
	products, count, err := u.repo.GetAllByCategoryID(category.ID, params.PaginationRequestDto)
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

func (u *ProductUsecase) GetSingleProductQuery(params *product.GetSingleProductQuery) (any, error) {
	u.logger.Info("GetSingleProductQuery called", map[string]interface{}{
		"slug":   *params.Slug,
		"siteId": *params.SiteID,
	})

	// Get product by slug
	product, err := u.repo.GetBySlug(*params.Slug)
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
			_ = u.repo.Update(product) // Ignore error as this is a background operation
		}
	}()

	return response, nil
}

func (u *ProductUsecase) CalculateProductsPriceQuery(params *product.CalculateProductsPriceQuery) (any, error) {
	u.logger.Info("CalculateProductsPriceQuery called", map[string]interface{}{
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
	// Note: In an actual implementation, we would need a repo method to get products with variants in a single query
	// For now, we'll get each product individually
	for _, productID := range productIDs {
		product, err := u.repo.GetByID(productID)
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
			discountType, _ := strconv.Atoi(discount.Type)
			if discountType == 0 { // Percentage
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

func (u *ProductUsecase) AdminGetAllProductQuery(params *product.AdminGetAllProductQuery) (any, error) {
	u.logger.Info("AdminGetAllProductQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all products across all sites for admin
	products, count, err := u.repo.GetAll(params.PaginationRequestDto)
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
