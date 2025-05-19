package repository

import (
	"strings"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/article"
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

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

func (r *ArticleRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	query := r.database.Model(&domain.Article{}).Where("is_deleted = ?", false)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	query := r.database.Model(&domain.Article{}).
		Where("site_id = ?", siteID).
		Where("is_deleted = ?", false)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	// For many-to-many relationship using the join table
	query := r.database.Model(&domain.Article{}).
		Joins("JOIN article_category ON article_category.article_id = articles.id").
		Where("article_category.category_id = ?", categoryID).
		Where("articles.is_deleted = ?", false)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetAllByFilterAndSort(
	siteID int64,
	filters map[article.ArticleFilterEnum][]string,
	sort *article.ArticleSortEnum,
	paginationRequestDto common.PaginationRequestDto,
) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var count int64

	// Start building the query
	query := r.database.Model(&domain.Article{}).
		Where("site_id = ?", siteID).
		Where("is_deleted = ?", false)

	// Apply filters if they exist
	if filters != nil {
		for filterType, values := range filters {
			switch filterType {
			case article.RateRange:
				// Expects min,max format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("rate BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case article.ReviewRange:
				// Expects min,max format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("review_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case article.VisitedRange:
				// Expects min,max format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("visited_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case article.AddedRange:
				// Expects start,end date format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("created_at BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case article.UpdatedRange:
				// Expects start,end date format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("updated_at BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case article.CategoryIds:
				// Join with category table to filter by categoryIds
				if len(values) > 0 {
					categoryIds := strings.Join(values, ",")
					query = query.Joins("JOIN article_category ac ON ac.article_id = articles.id").
						Where("ac.category_id IN (?)", categoryIds)
				}
			case article.ArticleIds:
				// Filter by article IDs
				if len(values) > 0 {
					articleIds := strings.Join(values, ",")
					query = query.Where("id IN (?)", articleIds)
				}
			case article.Badges:
				// Filter by badges (which are comma-separated in a single field)
				if len(values) > 0 {
					for _, badge := range values {
						query = query.Where("badges LIKE ?", "%"+badge+"%")
					}
				}
			}
		}
	}

	// Apply sorting if specified
	if sort != nil {
		switch *sort {
		case article.TitleAZ:
			query = query.Order("title ASC")
		case article.TitleZA:
			query = query.Order("title DESC")
		case article.RecentlyAdded:
			query = query.Order("created_at DESC")
		case article.RecentlyUpdated:
			query = query.Order("updated_at DESC")
		case article.MostVisited:
			query = query.Order("visited_count DESC")
		case article.LeastVisited:
			query = query.Order("visited_count ASC")
		case article.MostRated:
			query = query.Order("rate DESC")
		case article.LeastRated:
			query = query.Order("rate ASC")
		case article.MostReviewed:
			query = query.Order("review_count DESC")
		case article.LeastReviewed:
			query = query.Order("review_count ASC")
		}
	} else {
		// Default sort by recently updated
		query = query.Order("updated_at DESC")
	}

	// Get count
	query.Count(&count)

	// Apply pagination
	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	// Get results
	result := query.Limit(limit).Offset(offset).Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func (r *ArticleRepo) GetByID(id int64) (domain.Article, error) {
	var article domain.Article

	result := r.database.
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&article)

	if result.Error != nil {
		return domain.Article{}, result.Error
	}

	return article, nil
}

func (r *ArticleRepo) GetBySlug(slug string) (domain.Article, error) {
	var article domain.Article

	result := r.database.
		Where("slug = ?", slug).
		Where("is_deleted = ?", false).
		First(&article)

	if result.Error != nil {
		return domain.Article{}, result.Error
	}

	return article, nil
}

func (r *ArticleRepo) GetBySlugAndSiteID(slug string, siteID int64) (domain.Article, error) {
	var article domain.Article

	result := r.database.
		Where("slug = ?", slug).
		Where("site_id = ?", siteID).
		Where("is_deleted = ?", false).
		First(&article)

	if result.Error != nil {
		return domain.Article{}, result.Error
	}

	return article, nil
}

func (r *ArticleRepo) Create(article domain.Article) error {
	return r.database.Create(&article).Error
}

func (r *ArticleRepo) Update(article domain.Article) error {
	return r.database.Save(&article).Error
}

func (r *ArticleRepo) Delete(id int64) error {
	// Soft delete
	return r.database.Model(&domain.Article{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// Media relationship methods

func (r *ArticleRepo) GetArticleMedia(articleID int64) ([]domain.Media, error) {
	var mediaItems []domain.Media

	err := r.database.
		Joins("JOIN article_media ON article_media.media_id = media.id").
		Where("article_media.article_id = ?", articleID).
		Find(&mediaItems).Error

	if err != nil {
		return nil, err
	}

	return mediaItems, nil
}

func (r *ArticleRepo) AddMediaToArticle(articleID int64, mediaID int64) error {
	articleMedia := domain.ArticleMedia{
		ArticleID: articleID,
		MediaID:   mediaID,
	}

	return r.database.Create(&articleMedia).Error
}

func (r *ArticleRepo) RemoveMediaFromArticle(articleID int64, mediaID int64) error {
	return r.database.
		Where("article_id = ? AND media_id = ?", articleID, mediaID).
		Delete(&domain.ArticleMedia{}).Error
}

func (r *ArticleRepo) RemoveAllMediaFromArticle(articleID int64) error {
	return r.database.
		Where("article_id = ?", articleID).
		Delete(&domain.ArticleMedia{}).Error
}

// Category relationship methods

func (r *ArticleRepo) GetArticleCategories(articleID int64) ([]domain.BlogCategory, error) {
	var categories []domain.BlogCategory

	err := r.database.
		Joins("JOIN article_category ON article_category.category_id = blog_categories.id").
		Where("article_category.article_id = ?", articleID).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *ArticleRepo) AddCategoryToArticle(articleID int64, categoryID int64) error {
	articleCategory := domain.ArticleArticleCategory{
		ArticleID:  articleID,
		CategoryID: categoryID,
	}

	return r.database.Create(&articleCategory).Error
}

func (r *ArticleRepo) RemoveCategoryFromArticle(articleID int64, categoryID int64) error {
	return r.database.
		Where("article_id = ? AND category_id = ?", articleID, categoryID).
		Delete(&domain.ArticleArticleCategory{}).Error
}

func (r *ArticleRepo) RemoveAllCategoriesFromArticle(articleID int64) error {
	return r.database.
		Where("article_id = ?", articleID).
		Delete(&domain.ArticleArticleCategory{}).Error
}
