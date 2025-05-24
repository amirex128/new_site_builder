package productreviewusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_review"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type ProductReviewUsecase struct {
	*usecase.BaseUsecase
	logger      sflogger.Logger
	repo        repository.IProductReviewRepository
	productRepo repository.IProductRepository
	authContext func(c *gin.Context) service.IAuthService
}

func NewProductReviewUsecase(c contract.IContainer) *ProductReviewUsecase {
	return &ProductReviewUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		repo:        c.GetProductReviewRepo(),
		productRepo: c.GetProductRepo(),
		authContext: c.GetAuthTransientService(),
	}
}

func (u *ProductReviewUsecase) CreateProductReviewCommand(params *product_review.CreateProductReviewCommand) (*resp.Response, error) {
	u.Logger.Info("CreateProductReviewCommand called", map[string]interface{}{
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
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Get customer ID if available (in this monolithic app, we may have a customer context too)
	customerID, _ := u.authContext(u.Ctx).GetCustomerID()
	// Use userId if customerID is 0
	if customerID == 0 {
		customerID = userID
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

	return resp.NewResponseData(resp.Created, map[string]interface{}{"review": createdReview}, "نظر با موفقیت ایجاد شد"), nil
}

func (u *ProductReviewUsecase) UpdateProductReviewCommand(params *product_review.UpdateProductReviewCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateProductReviewCommand called", map[string]interface{}{
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
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Check if user has access to this review
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
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

	return resp.NewResponseData(resp.Updated, map[string]interface{}{"review": updatedReview}, "نظر با موفقیت بروزرسانی شد"), nil
}

func (u *ProductReviewUsecase) DeleteProductReviewCommand(params *product_review.DeleteProductReviewCommand) (*resp.Response, error) {
	u.Logger.Info("DeleteProductReviewCommand called", map[string]interface{}{
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
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Check if user has access to this review
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
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

	return resp.NewResponseData(resp.Deleted, map[string]interface{}{
		"success": true,
	}, "نظر با موفقیت حذف شد"), nil
}

func (u *ProductReviewUsecase) GetByIdProductReviewQuery(params *product_review.GetByIdProductReviewQuery) (*resp.Response, error) {
	u.Logger.Info("GetByIdProductReviewQuery called", map[string]interface{}{
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
		userID, _ := u.authContext(u.Ctx).GetUserID()
		isAdmin, _ := u.authContext(u.Ctx).IsAdmin()

		if review.UserID != userID && !isAdmin {
			return nil, errors.New("شما به این نظر دسترسی ندارید")
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"review": review}, "نظر با موفقیت دریافت شد"), nil
}

func (u *ProductReviewUsecase) GetAllProductReviewQuery(params *product_review.GetAllProductReviewQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllProductReviewQuery called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check site access (can be optional for public reviews)
	// In a real implementation, we might check if the user has access to this site
	// For simplicity, we'll allow access to approved reviews for all

	// Get all reviews for the site
	reviewsResult, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     reviewsResult.Items,
		"total":     reviewsResult.TotalCount,
		"page":      reviewsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": reviewsResult.TotalPages,
	}, "لیست نظرات با موفقیت دریافت شد"), nil
}

func (u *ProductReviewUsecase) AdminGetAllProductReviewQuery(params *product_review.AdminGetAllProductReviewQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllProductReviewQuery called", map[string]interface{}{
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

	// Get all reviews across all sites for admin
	reviewsResult, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Return paginated result
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     reviewsResult.Items,
		"total":     reviewsResult.TotalCount,
		"page":      reviewsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": reviewsResult.TotalPages,
	}, "لیست نظرات ادمین با موفقیت دریافت شد"), nil
}
