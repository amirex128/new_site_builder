package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	enums2 "github.com/amirex128/new_site_builder/internal/domain/enums"
)

type IArticleRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Article], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Article], error)
	GetAllByCategoryID(categoryID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Article], error)
	GetAllByFilterAndSort(siteID int64, filters map[enums2.ArticleFilterEnum][]string, sort *enums2.ArticleSortEnum, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Article], error)
	GetByID(id int64) (*domain2.Article, error)
	GetBySlug(slug string) (*domain2.Article, error)
	GetBySlugAndSiteID(slug string, siteID int64) (*domain2.Article, error)
	Create(article *domain2.Article) error
	Update(article *domain2.Article) error
	Delete(id int64) error

	// Media relationship methods
	GetArticleMedia(articleID int64) ([]domain2.Media, error)
	AddMediaToArticle(articleID int64, mediaID int64) error
	RemoveMediaFromArticle(articleID int64, mediaID int64) error
	RemoveAllMediaFromArticle(articleID int64) error

	// Category relationship methods
	GetArticleCategories(articleID int64) ([]domain2.ArticleCategory, error)
	AddCategoryToArticle(articleID int64, categoryID int64) error
	RemoveCategoryFromArticle(articleID int64, categoryID int64) error
	RemoveAllCategoriesFromArticle(articleID int64) error
}
