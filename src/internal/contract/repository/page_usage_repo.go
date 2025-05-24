package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPageArticleUsageRepository interface {
	GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageArticleUsage, error)
	GetByArticleIDsAndSiteID(articleIDs []int64, siteID int64) ([]domain.PageArticleUsage, error)
	DeleteByPageIDAndSiteID(pageID, siteID int64) error
	Create(usage domain.PageArticleUsage) error
	CreateBatch(usages []domain.PageArticleUsage) error
}

type IPageProductUsageRepository interface {
	GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageProductUsage, error)
	GetByProductIDsAndSiteID(productIDs []int64, siteID int64) ([]domain.PageProductUsage, error)
	DeleteByPageIDAndSiteID(pageID, siteID int64) error
	Create(usage domain.PageProductUsage) error
	CreateBatch(usages []domain.PageProductUsage) error
}

type IPageHeaderFooterUsageRepository interface {
	GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageHeaderFooterUsage, error)
	GetByHeaderFooterIDsAndSiteID(headerFooterIDs []int64, siteID int64) ([]domain.PageHeaderFooterUsage, error)
	DeleteByPageIDAndSiteID(pageID, siteID int64) error
	Create(usage domain.PageHeaderFooterUsage) error
	CreateBatch(usages []domain.PageHeaderFooterUsage) error
}
