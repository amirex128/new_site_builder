package productreviewusecase

import (
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product_review"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ProductReviewUsecase struct {
	logger sflogger.Logger
	repo   repository.IProductReviewRepository
}

func NewProductReviewUsecase(c contract.IContainer) *ProductReviewUsecase {
	return &ProductReviewUsecase{
		logger: c.GetLogger(),
		repo:   c.GetProductReviewRepo(),
	}
}

func (u *ProductReviewUsecase) CreateProductReviewCommand(params *product_review.CreateProductReviewCommand) (any, error) {
	// Implementation for creating a product review
	fmt.Println(params)

	newReview := domain.ProductReview{
		Rating:     *params.Rating,
		Like:       *params.Like,
		Dislike:    *params.Dislike,
		Approved:   *params.Approved,
		ReviewText: *params.ReviewText,
		ProductID:  *params.ProductID,
		SiteID:     *params.SiteID,
		CustomerID: 1, // Should come from auth context
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	err := u.repo.Create(newReview)
	if err != nil {
		return nil, err
	}

	return newReview, nil
}

func (u *ProductReviewUsecase) UpdateProductReviewCommand(params *product_review.UpdateProductReviewCommand) (any, error) {
	// Implementation for updating a product review
	fmt.Println(params)

	existingReview, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return existingReview, nil
}

func (u *ProductReviewUsecase) DeleteProductReviewCommand(params *product_review.DeleteProductReviewCommand) (any, error) {
	// Implementation for deleting a product review
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *ProductReviewUsecase) GetByIdProductReviewQuery(params *product_review.GetByIdProductReviewQuery) (any, error) {
	// Implementation to get product review by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *ProductReviewUsecase) GetAllProductReviewQuery(params *product_review.GetAllProductReviewQuery) (any, error) {
	// Implementation to get all product reviews by site ID
	fmt.Println(params)

	result, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *ProductReviewUsecase) AdminGetAllProductReviewQuery(params *product_review.AdminGetAllProductReviewQuery) (any, error) {
	// Implementation to get all product reviews for admin
	fmt.Println(params)

	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
