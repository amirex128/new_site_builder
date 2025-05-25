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
	_, err := u.productRepo.GetByID(*params.ProductID)
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
	customerID, _ := u.authContext(u.Ctx).GetCustomerID()
	if customerID == nil || (customerID != nil && *customerID == 0) {
		customerID = userID
	}
	newReview := domain.ProductReview{
		Rating:     *params.Rating,
		Like:       *params.Like,
		Dislike:    *params.Dislike,
		Approved:   *params.Approved,
		ReviewText: *params.ReviewText,
		ProductID:  *params.ProductID,
		SiteID:     *params.SiteID,
		UserID:     *userID,
		CustomerID: *customerID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}
	err = u.repo.Create(newReview)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	createdReview, err := u.repo.GetByID(newReview.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Created, map[string]interface{}{"review": createdReview}, "نظر با موفقیت ایجاد شد"), nil
}

func (u *ProductReviewUsecase) UpdateProductReviewCommand(params *product_review.UpdateProductReviewCommand) (*resp.Response, error) {
	existingReview, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "نظر یافت نشد")
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
	if userID != nil && existingReview.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این نظر دسترسی ندارید")
	}
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
	err = u.repo.Update(existingReview)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	updatedReview, err := u.repo.GetByID(existingReview.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Updated, map[string]interface{}{"review": updatedReview}, "نظر با موفقیت بروزرسانی شد"), nil
}

func (u *ProductReviewUsecase) DeleteProductReviewCommand(params *product_review.DeleteProductReviewCommand) (*resp.Response, error) {
	existingReview, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "نظر یافت نشد")
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
	if userID != nil && existingReview.UserID != *userID && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما به این نظر دسترسی ندارید")
	}
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Deleted, map[string]interface{}{"success": true}, "نظر با موفقیت حذف شد"), nil
}

func (u *ProductReviewUsecase) GetByIdProductReviewQuery(params *product_review.GetByIdProductReviewQuery) (*resp.Response, error) {
	review, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if !review.Approved {
		userID, _ := u.authContext(u.Ctx).GetUserID()
		isAdmin, _ := u.authContext(u.Ctx).IsAdmin()
		if userID != nil && review.UserID != *userID && !isAdmin {
			return nil, resp.NewError(resp.Unauthorized, "شما به این نظر دسترسی ندارید")
		}
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"review": review}, "نظر با موفقیت دریافت شد"), nil
}

func (u *ProductReviewUsecase) GetAllProductReviewQuery(params *product_review.GetAllProductReviewQuery) (*resp.Response, error) {
	reviewsResult, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     reviewsResult.Items,
		"total":     reviewsResult.TotalCount,
		"page":      reviewsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": reviewsResult.TotalPages,
	}, "لیست نظرات با موفقیت دریافت شد"), nil
}

func (u *ProductReviewUsecase) AdminGetAllProductReviewQuery(params *product_review.AdminGetAllProductReviewQuery) (*resp.Response, error) {
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}
	reviewsResult, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     reviewsResult.Items,
		"total":     reviewsResult.TotalCount,
		"page":      reviewsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": reviewsResult.TotalPages,
	}, "لیست نظرات ادمین با موفقیت دریافت شد"), nil
}
