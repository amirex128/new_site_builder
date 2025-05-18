package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/article"
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IArticleRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetAllByFilterAndSort(siteID int64, filters map[article.ArticleFilterEnum][]string, sort *article.ArticleSortEnum, paginationRequestDto common.PaginationRequestDto) ([]domain.Article, int64, error)
	GetByID(id int64) (domain.Article, error)
	GetBySlug(slug string) (domain.Article, error)
	GetBySlugAndSiteID(slug string, siteID int64) (domain.Article, error)
	Create(article domain.Article) error
	Update(article domain.Article) error
	Delete(id int64) error

	// Media relationship methods
	GetArticleMedia(articleID int64) ([]domain.Media, error)
	AddMediaToArticle(articleID int64, mediaID int64) error
	RemoveMediaFromArticle(articleID int64, mediaID int64) error
	RemoveAllMediaFromArticle(articleID int64) error

	// Category relationship methods
	GetArticleCategories(articleID int64) ([]domain.BlogCategory, error)
	AddCategoryToArticle(articleID int64, categoryID int64) error
	RemoveCategoryFromArticle(articleID int64, categoryID int64) error
	RemoveAllCategoriesFromArticle(articleID int64) error
}
