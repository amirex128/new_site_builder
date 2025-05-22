package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"strings"

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
	filters map[enums.ArticleFilterEnum][]string,
	sort *enums.ArticleSortEnum,
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
			case "rate_range":
				// Expects min,max format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("rate BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "review_range":
				// Expects min,max format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("review_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "visited_range":
				// Expects min,max format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("visited_count BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "added_range":
				// Expects start,end date format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("created_at BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "updated_range":
				// Expects start,end date format in the first value
				if len(values) > 0 {
					parts := strings.Split(values[0], ",")
					if len(parts) == 2 {
						query = query.Where("updated_at BETWEEN ? AND ?", parts[0], parts[1])
					}
				}
			case "category_ids":
				// Join with category table to filter by categoryIds
				if len(values) > 0 {
					categoryIds := strings.Join(values, ",")
					query = query.Joins("JOIN article_category ac ON ac.article_id = articles.id").
						Where("ac.category_id IN (?)", categoryIds)
				}
			case "article_ids":
				// Filter by article IDs
				if len(values) > 0 {
					articleIds := strings.Join(values, ",")
					query = query.Where("id IN (?)", articleIds)
				}
			case "badges":
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
		case "title_az":
			query = query.Order("title ASC")
		case "title_za":
			query = query.Order("title DESC")
		case "recently_added":
			query = query.Order("created_at DESC")
		case "recently_updated":
			query = query.Order("updated_at DESC")
		case "most_visited":
			query = query.Order("visited_count DESC")
		case "least_visited":
			query = query.Order("visited_count ASC")
		case "most_rated":
			query = query.Order("rate DESC")
		case "least_rated":
			query = query.Order("rate ASC")
		case "most_reviewed":
			query = query.Order("review_count DESC")
		case "least_reviewed":
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

func (r *ArticleRepo) GetArticleCategories(articleID int64) ([]domain.ArticleCategory, error) {
	var categories []domain.ArticleCategory

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
