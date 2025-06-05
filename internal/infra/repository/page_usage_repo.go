package repository

import (
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

// PageArticleUsageRepo implements IPageArticleUsageRepository
type PageArticleUsageRepo struct {
	database *gorm.DB
}

func NewPageArticleUsageRepository(db *gorm.DB) *PageArticleUsageRepo {
	return &PageArticleUsageRepo{
		database: db,
	}
}

func (r *PageArticleUsageRepo) GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageArticleUsage, error) {
	var usages []domain.PageArticleUsage
	result := r.database.
		Joins("JOIN pages ON page_article_usages.page_id = pages.id").
		Where("page_article_usages.page_id = ? AND pages.site_id = ?", pageID, siteID).
		Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *PageArticleUsageRepo) GetByArticleIDsAndSiteID(articleIDs []int64, siteID int64) ([]domain.PageArticleUsage, error) {
	var usages []domain.PageArticleUsage
	result := r.database.
		Joins("JOIN pages ON page_article_usages.page_id = pages.id").
		Where("page_article_usages.article_id IN ? AND pages.site_id = ?", articleIDs, siteID).
		Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *PageArticleUsageRepo) DeleteByPageIDAndSiteID(pageID, siteID int64) error {
	return r.database.
		Exec("DELETE pau FROM page_article_usages pau JOIN pages p ON pau.page_id = p.id WHERE pau.page_id = ? AND p.site_id = ?", pageID, siteID).
		Error
}

func (r *PageArticleUsageRepo) Create(usage *domain.PageArticleUsage) error {
	result := r.database.Create(usage)
	return result.Error
}

func (r *PageArticleUsageRepo) CreateBatch(usages []domain.PageArticleUsage) error {
	result := r.database.Create(usages)
	return result.Error
}

func (r *PageArticleUsageRepo) AddArticleToPage(pageID int64, articleID int64, position string) error {
	usage := domain.PageArticleUsage{
		PageID:    pageID,
		ArticleID: articleID,
	}
	result := r.database.Create(&usage)
	return result.Error
}

func (r *PageArticleUsageRepo) RemoveArticleFromPage(pageID int64, articleID int64) error {
	result := r.database.Where("page_id = ? AND article_id = ?", pageID, articleID).Delete(&domain.PageArticleUsage{})
	return result.Error
}

func (r *PageArticleUsageRepo) GetArticlesByPageID(pageID int64) ([]domain.PageArticleUsage, error) {
	var pageArticleUsages []domain.PageArticleUsage
	result := r.database.Where("page_id = ?", pageID).Preload("Article").Find(&pageArticleUsages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pageArticleUsages, nil
}

// PageProductUsageRepo implements IPageProductUsageRepository
type PageProductUsageRepo struct {
	database *gorm.DB
}

func NewPageProductUsageRepository(db *gorm.DB) *PageProductUsageRepo {
	return &PageProductUsageRepo{
		database: db,
	}
}

func (r *PageProductUsageRepo) GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageProductUsage, error) {
	var usages []domain.PageProductUsage
	result := r.database.
		Joins("JOIN pages ON page_product_usages.page_id = pages.id").
		Where("page_product_usages.page_id = ? AND pages.site_id = ?", pageID, siteID).
		Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *PageProductUsageRepo) GetByProductIDsAndSiteID(productIDs []int64, siteID int64) ([]domain.PageProductUsage, error) {
	var usages []domain.PageProductUsage
	result := r.database.
		Joins("JOIN pages ON page_product_usages.page_id = pages.id").
		Where("page_product_usages.product_id IN ? AND pages.site_id = ?", productIDs, siteID).
		Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *PageProductUsageRepo) DeleteByPageIDAndSiteID(pageID, siteID int64) error {
	return r.database.
		Exec("DELETE ppu FROM page_product_usages ppu JOIN pages p ON ppu.page_id = p.id WHERE ppu.page_id = ? AND p.site_id = ?", pageID, siteID).
		Error
}

func (r *PageProductUsageRepo) Create(usage *domain.PageProductUsage) error {
	result := r.database.Create(usage)
	return result.Error
}

func (r *PageProductUsageRepo) CreateBatch(usages []domain.PageProductUsage) error {
	result := r.database.Create(usages)
	return result.Error
}

func (r *PageProductUsageRepo) AddProductToPage(pageID int64, productID int64, position string) error {
	usage := domain.PageProductUsage{
		PageID:    pageID,
		ProductID: productID,
	}
	result := r.database.Create(&usage)
	return result.Error
}

func (r *PageProductUsageRepo) RemoveProductFromPage(pageID int64, productID int64) error {
	result := r.database.Where("page_id = ? AND product_id = ?", pageID, productID).Delete(&domain.PageProductUsage{})
	return result.Error
}

func (r *PageProductUsageRepo) GetProductsByPageID(pageID int64) ([]domain.PageProductUsage, error) {
	var pageProductUsages []domain.PageProductUsage
	result := r.database.Where("page_id = ?", pageID).Preload("Product").Find(&pageProductUsages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pageProductUsages, nil
}

// PageHeaderFooterUsageRepo implements IPageHeaderFooterUsageRepository
type PageHeaderFooterUsageRepo struct {
	database *gorm.DB
}

func NewPageHeaderFooterUsageRepository(db *gorm.DB) *PageHeaderFooterUsageRepo {
	return &PageHeaderFooterUsageRepo{
		database: db,
	}
}

func (r *PageHeaderFooterUsageRepo) GetByPageIDAndSiteID(pageID, siteID int64) ([]domain.PageHeaderFooterUsage, error) {
	var usages []domain.PageHeaderFooterUsage
	result := r.database.
		Joins("JOIN pages ON page_header_footer_usages.page_id = pages.id").
		Where("page_header_footer_usages.page_id = ? AND pages.site_id = ?", pageID, siteID).
		Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *PageHeaderFooterUsageRepo) GetByHeaderFooterIDsAndSiteID(headerFooterIDs []int64, siteID int64) ([]domain.PageHeaderFooterUsage, error) {
	var usages []domain.PageHeaderFooterUsage
	result := r.database.
		Joins("JOIN pages ON page_header_footer_usages.page_id = pages.id").
		Where("page_header_footer_usages.header_footer_id IN ? AND pages.site_id = ?", headerFooterIDs, siteID).
		Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *PageHeaderFooterUsageRepo) DeleteByPageIDAndSiteID(pageID, siteID int64) error {
	return r.database.
		Exec("DELETE phfu FROM page_header_footer_usages phfu JOIN pages p ON phfu.page_id = p.id WHERE phfu.page_id = ? AND p.site_id = ?", pageID, siteID).
		Error
}

func (r *PageHeaderFooterUsageRepo) Create(usage *domain.PageHeaderFooterUsage) error {
	result := r.database.Create(usage)
	return result.Error
}

func (r *PageHeaderFooterUsageRepo) CreateBatch(usages []domain.PageHeaderFooterUsage) error {
	result := r.database.Create(usages)
	return result.Error
}

func (r *PageHeaderFooterUsageRepo) AddHeaderFooterToPage(pageID int64, headerFooterID int64, position string) error {
	usage := domain.PageHeaderFooterUsage{
		PageID:         pageID,
		HeaderFooterID: headerFooterID,
	}
	result := r.database.Create(&usage)
	return result.Error
}

func (r *PageHeaderFooterUsageRepo) RemoveHeaderFooterFromPage(pageID int64, headerFooterID int64) error {
	result := r.database.Where("page_id = ? AND header_footer_id = ?", pageID, headerFooterID).Delete(&domain.PageHeaderFooterUsage{})
	return result.Error
}

func (r *PageHeaderFooterUsageRepo) GetHeaderFootersByPageID(pageID int64) ([]domain.PageHeaderFooterUsage, error) {
	var pageHeaderFooterUsages []domain.PageHeaderFooterUsage
	result := r.database.Where("page_id = ?", pageID).Preload("HeaderFooter").Find(&pageHeaderFooterUsages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pageHeaderFooterUsages, nil
}
