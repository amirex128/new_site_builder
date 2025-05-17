package productusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ProductUsecase struct {
	logger sflogger.Logger
	repo   repository.IProductRepository
}

func NewProductUsecase(c contract.IContainer) *ProductUsecase {
	return &ProductUsecase{
		logger: c.GetLogger(),
		repo:   c.GetProductRepo(),
	}
}

func (u *ProductUsecase) CreateProductCommand(params *product.CreateProductCommand) (any, error) {
	// Implementation for creating a product
	fmt.Println(params)

	// Convert nullable string fields to non-nullable
	var description string
	if params.Description != nil {
		description = *params.Description
	}

	var longDescription string
	if params.LongDescription != nil {
		longDescription = *params.LongDescription
	}

	// Convert status enum to string
	statusStr := strconv.Itoa(int(*params.Status))

	newProduct := domain.Product{
		Name:            *params.Name,
		Description:     description,
		Status:          statusStr,
		Weight:          *params.Weight,
		FreeSend:        *params.FreeSend,
		LongDescription: longDescription,
		Slug:            *params.Slug,
		SiteID:          *params.SiteID,
		SellingCount:    0,
		VisitedCount:    0,
		ReviewCount:     0,
		Rate:            0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsDeleted:       false,
		UserID:          1, // Should come from auth context
	}

	err := u.repo.Create(newProduct)
	if err != nil {
		return nil, err
	}

	// TODO: Handle related entities like media, categories, variants, attributes, etc.

	return newProduct, nil
}

func (u *ProductUsecase) UpdateProductCommand(params *product.UpdateProductCommand) (any, error) {
	// Implementation for updating a product
	fmt.Println(params)

	existingProduct, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

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

	existingProduct.SiteID = *params.SiteID
	existingProduct.UpdatedAt = time.Now()

	err = u.repo.Update(existingProduct)
	if err != nil {
		return nil, err
	}

	// TODO: Handle related entities like media, categories, variants, attributes, etc.

	return existingProduct, nil
}

func (u *ProductUsecase) DeleteProductCommand(params *product.DeleteProductCommand) (any, error) {
	// Implementation for deleting a product
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *ProductUsecase) GetByIdProductQuery(params *product.GetByIdProductQuery) (any, error) {
	// Implementation to get product by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) GetAllProductQuery(params *product.GetAllProductQuery) (any, error) {
	// Implementation to get all products
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

func (u *ProductUsecase) GetByFiltersSortProductQuery(params *product.GetByFiltersSortProductQuery) (any, error) {
	// Implementation to get products with filtering and sorting
	fmt.Println(params)

	// TODO: Implement proper filtering and sorting logic
	// For now, just return all products by site ID
	result, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *ProductUsecase) AdminGetAllProductQuery(params *product.AdminGetAllProductQuery) (any, error) {
	// Implementation to get all products for admin
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
