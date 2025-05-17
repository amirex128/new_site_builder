package pageusecase

import (
	"encoding/json"
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/page"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type PageUsecase struct {
	logger sflogger.Logger
	repo   repository.IPageRepository
}

func NewPageUsecase(c contract.IContainer) *PageUsecase {
	return &PageUsecase{
		logger: c.GetLogger(),
		repo:   c.GetPageRepo(),
	}
}

func (u *PageUsecase) CreatePageCommand(params *page.CreatePageCommand) (any, error) {
	// Implementation for creating a page
	fmt.Println(params)

	var description string
	if params.Description != nil {
		description = *params.Description
	}

	// Convert body to JSON string
	bodyJSON, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}

	// Convert SeoTags to JSON string
	seoTagsJSON, err := json.Marshal(params.SeoTags)
	if err != nil {
		return nil, err
	}

	newPage := domain.Page{
		SiteID:      *params.SiteID,
		HeaderID:    *params.HeaderID,
		FooterID:    *params.FooterID,
		Slug:        *params.Slug,
		Title:       *params.Title,
		Description: description,
		Body:        string(bodyJSON),
		SeoTags:     string(seoTagsJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
		UserID:      1, // Should come from auth context
	}

	err = u.repo.Create(newPage)
	if err != nil {
		return nil, err
	}

	// TODO: Handle media relations separately using the PageMedia join table

	return newPage, nil
}

func (u *PageUsecase) UpdatePageCommand(params *page.UpdatePageCommand) (any, error) {
	// Implementation for updating a page
	fmt.Println(params)

	existingPage, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	if params.SiteID != nil {
		existingPage.SiteID = *params.SiteID
	}

	if params.HeaderID != nil {
		existingPage.HeaderID = *params.HeaderID
	}

	if params.FooterID != nil {
		existingPage.FooterID = *params.FooterID
	}

	if params.Slug != nil {
		existingPage.Slug = *params.Slug
	}

	if params.Title != nil {
		existingPage.Title = *params.Title
	}

	if params.Description != nil {
		existingPage.Description = *params.Description
	}

	if params.Body != nil {
		bodyJSON, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		existingPage.Body = string(bodyJSON)
	}

	if params.SeoTags != nil {
		seoTagsJSON, err := json.Marshal(params.SeoTags)
		if err != nil {
			return nil, err
		}
		existingPage.SeoTags = string(seoTagsJSON)
	}

	existingPage.UpdatedAt = time.Now()

	err = u.repo.Update(existingPage)
	if err != nil {
		return nil, err
	}

	// TODO: Handle media relations separately using the PageMedia join table

	return existingPage, nil
}

func (u *PageUsecase) DeletePageCommand(params *page.DeletePageCommand) (any, error) {
	// Implementation for deleting a page
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *PageUsecase) GetByIdPageQuery(params *page.GetByIdPageQuery) (any, error) {
	// Implementation to get page by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *PageUsecase) GetAllPageQuery(params *page.GetAllPageQuery) (any, error) {
	// Implementation to get all pages
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

func (u *PageUsecase) AdminGetAllPageQuery(params *page.AdminGetAllPageQuery) (any, error) {
	// Implementation to get all pages for admin
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
