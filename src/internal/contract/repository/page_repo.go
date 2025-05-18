package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPageRepository interface {
	// GetAll returns all pages with pagination
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Page, int64, error)

	// GetAllBySiteID returns pages for a specific site with pagination
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Page, int64, error)

	// GetByID returns a page by its ID
	GetByID(id int64) (domain.Page, error)

	// GetByIDAndSiteID returns a page by ID and site ID
	GetByIDAndSiteID(id, siteID int64) (domain.Page, error)

	// GetBySlug returns a page by slug and site ID
	GetBySlug(slug string, siteID int64) (domain.Page, error)

	// GetByPaths returns pages by their paths (slugs) and site ID
	GetByPaths(paths []string, siteID int64) ([]domain.Page, error)

	// GetByIDs returns pages by their IDs and site ID
	GetByIDs(ids []int64, siteID int64) ([]domain.Page, error)

	// Create creates a new page
	Create(page domain.Page) error

	// Update updates an existing page
	Update(page domain.Page) error

	// Delete deletes a page by ID
	Delete(id int64) error

	// AddMediaToPage adds media to a page
	AddMediaToPage(pageID int64, mediaIDs []int64) error

	// RemoveAllMediaFromPage removes all media from a page
	RemoveAllMediaFromPage(pageID int64) error
}
