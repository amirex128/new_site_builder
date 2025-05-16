package mysql

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"

	"gorm.io/gorm"
)

type ArticleRepo struct {
	database *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{
		database: db,
	}
}

func (r *ArticleRepo) GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	query := r.database.Model(&domain.Article{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	query := r.database.Model(&domain.Article{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetAllByCategoryID(categoryID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	// For many-to-many relationship using the join table
	query := r.database.Model(&domain.Article{}).
		Joins("JOIN article_category ON article_category.article_id = articles.id").
		Where("article_category.category_id = ?", categoryID)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetByID(id int64) (domain.Article, error) {
	var article domain.Article
	result := r.database.First(&article, id)
	if result.Error != nil {
		return article, result.Error
	}
	return article, nil
}

func (r *ArticleRepo) GetBySlug(slug string) (domain.Article, error) {
	var article domain.Article
	result := r.database.Where("slug = ?", slug).First(&article)
	if result.Error != nil {
		return article, result.Error
	}
	return article, nil
}

func (r *ArticleRepo) Create(article domain.Article) error {
	result := r.database.Create(&article)
	return result.Error
}

func (r *ArticleRepo) Update(article domain.Article) error {
	result := r.database.Save(&article)
	return result.Error
}

func (r *ArticleRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Article{}, id)
	return result.Error
}
