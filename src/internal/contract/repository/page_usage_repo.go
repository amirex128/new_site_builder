package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

// IPageArticleUsageRepository defines the interface for page-article usage operations
type IPageArticleUsageRepository interface {
	// GetByPageIDAndSiteID returns usages by page ID and site ID
	GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageArticleUsage, error)
	// GetByArticleIDsAndSiteID returns usages by article IDs and site ID
	GetByArticleIDsAndSiteID(articleIDs []int64, siteID int64) ([]domain.PageArticleUsage, error)
	// DeleteByPageIDAndSiteID deletes usages by page ID and site ID
	DeleteByPageIDAndSiteID(pageID, siteID int64) error
	// Create creates a new page-article usage
	Create(usage domain.PageArticleUsage) error
	// CreateBatch creates multiple page-article usages
	CreateBatch(usages []domain.PageArticleUsage) error
}

// IPageProductUsageRepository defines the interface for page-product usage operations
type IPageProductUsageRepository interface {
	// GetByPageIDAndSiteID returns usages by page ID and site ID
	GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageProductUsage, error)
	// GetByProductIDsAndSiteID returns usages by product IDs and site ID
	GetByProductIDsAndSiteID(productIDs []int64, siteID int64) ([]domain.PageProductUsage, error)
	// DeleteByPageIDAndSiteID deletes usages by page ID and site ID
	DeleteByPageIDAndSiteID(pageID, siteID int64) error
	// Create creates a new page-product usage
	Create(usage domain.PageProductUsage) error
	// CreateBatch creates multiple page-product usages
	CreateBatch(usages []domain.PageProductUsage) error
}

// IPageHeaderFooterUsageRepository defines the interface for page-header/footer usage operations
type IPageHeaderFooterUsageRepository interface {
	// GetByPageIDAndSiteID returns usages by page ID and site ID
	GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageHeaderFooterUsage, error)
	// GetByHeaderFooterIDsAndSiteID returns usages by header/footer IDs and site ID
	GetByHeaderFooterIDsAndSiteID(headerFooterIDs []int64, siteID int64) ([]domain.PageHeaderFooterUsage, error)
	// DeleteByPageIDAndSiteID deletes usages by page ID and site ID
	DeleteByPageIDAndSiteID(pageID, siteID int64) error
	// Create creates a new page-header/footer usage
	Create(usage domain.PageHeaderFooterUsage) error
	// CreateBatch creates multiple page-header/footer usages
	CreateBatch(usages []domain.PageHeaderFooterUsage) error
}
