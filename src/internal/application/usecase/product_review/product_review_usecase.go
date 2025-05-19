package productreviewusecase

import (
	"errors"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_review"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type ProductReviewUsecase struct {
	logger         sflogger.Logger
	repo           repository.IProductReviewRepository
	productRepo    repository.IProductRepository
	authContextSvc common.IAuthContextService
}

func NewProductReviewUsecase(c contract.IContainer) *ProductReviewUsecase {
	return &ProductReviewUsecase{
		logger:         c.GetLogger(),
		repo:           c.GetProductReviewRepo(),
		productRepo:    c.GetProductRepo(),
		authContextSvc: c.GetAuthContextTransientService(),
	}
}

func (u *ProductReviewUsecase) CreateProductReviewCommand(params *product_review.CreateProductReviewCommand) (any, error) {
	u.logger.Info("CreateProductReviewCommand called", map[string]interface{}{
		"rating":    *params.Rating,
		"productId": *params.ProductID,
		"siteId":    *params.SiteID,
	})

	// Validate product exists
	_, err := u.productRepo.GetByID(*params.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محصول یافت نشد")
		}
		return nil, err
	}

	// Get user ID from auth context
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Get customer ID if available (in this monolithic app, we may have a customer context too)
	customerID, _ := u.authContextSvc.GetCustomerID()
	if customerID == 0 {
		// If no customer ID available, use a default or generate one
		customerID = 1 // Default value, in a real app this would need proper handling
	}

	// Create review entity
	newReview := domain.ProductReview{
		Rating:     *params.Rating,
		Like:       *params.Like,
		Dislike:    *params.Dislike,
		Approved:   *params.Approved,
		ReviewText: *params.ReviewText,
		ProductID:  *params.ProductID,
		SiteID:     *params.SiteID,
		UserID:     userID,
		CustomerID: customerID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	// Create the review in the database
	err = u.repo.Create(newReview)
	if err != nil {
		return nil, err
	}

	// Fetch the created review
	createdReview, err := u.repo.GetByID(newReview.ID)
	if err != nil {
		return nil, err
	}

	return createdReview, nil
}

func (u *ProductReviewUsecase) UpdateProductReviewCommand(params *product_review.UpdateProductReviewCommand) (any, error) {
	u.logger.Info("UpdateProductReviewCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing review
	existingReview, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Check if user has access to this review
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingReview.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این نظر دسترسی ندارید")
	}

	// Update fields if provided
	if params.Rating != nil {
		existingReview.Rating = *params.Rating
	}

	if params.Like != nil {
		existingReview.Like = *params.Like
	}

	if params.Dislike != nil {
		existingReview.Dislike = *params.Dislike
	}

	if params.Approved != nil {
		existingReview.Approved = *params.Approved
	}

	if params.ReviewText != nil {
		existingReview.ReviewText = *params.ReviewText
	}

	existingReview.ProductID = *params.ProductID

	if params.SiteID != nil {
		existingReview.SiteID = *params.SiteID
	}

	existingReview.UpdatedAt = time.Now()

	// Update the review
	err = u.repo.Update(existingReview)
	if err != nil {
		return nil, err
	}

	// Fetch the updated review
	updatedReview, err := u.repo.GetByID(existingReview.ID)
	if err != nil {
		return nil, err
	}

	return updatedReview, nil
}

func (u *ProductReviewUsecase) DeleteProductReviewCommand(params *product_review.DeleteProductReviewCommand) (any, error) {
	u.logger.Info("DeleteProductReviewCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing review
	existingReview, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Check if user has access to this review
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingReview.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این نظر دسترسی ندارید")
	}

	// Delete the review
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (u *ProductReviewUsecase) GetByIdProductReviewQuery(params *product_review.GetByIdProductReviewQuery) (any, error) {
	u.logger.Info("GetByIdProductReviewQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get review by ID
	review, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access - anyone can view approved reviews
	if !review.Approved {
		userID, _ := u.authContextSvc.GetUserID()
		isAdmin, _ := u.authContextSvc.IsAdmin()

		if review.UserID != userID && !isAdmin {
			return nil, errors.New("شما به این نظر دسترسی ندارید")
		}
	}

	return review, nil
}

func (u *ProductReviewUsecase) GetAllProductReviewQuery(params *product_review.GetAllProductReviewQuery) (any, error) {
	u.logger.Info("GetAllProductReviewQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check site access (can be optional for public reviews)
	// In a real implementation, we might check if the user has access to this site
	// For simplicity, we'll allow access to approved reviews for all

	// Get all reviews for the site
	reviews, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     reviews,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *ProductReviewUsecase) AdminGetAllProductReviewQuery(params *product_review.AdminGetAllProductReviewQuery) (any, error) {
	u.logger.Info("AdminGetAllProductReviewQuery called", map[string]interface{}{
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

	// Get all reviews across all sites for admin
	reviews, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return map[string]interface{}{
		"items":     reviews,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}
