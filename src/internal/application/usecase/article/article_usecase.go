package articleusecase

import (
	"strings"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/article"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ArticleUsecase struct {
	*usecase.BaseUsecase
	articleRepo  repository.IArticleRepository
	categoryRepo repository.IArticleCategoryRepository
	mediaRepo    repository.IMediaRepository
}

func NewArticleUsecase(c contract.IContainer) *ArticleUsecase {
	return &ArticleUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		articleRepo:  c.GetArticleRepo(),
		categoryRepo: c.GetArticleCategoryRepo(),
		mediaRepo:    c.GetMediaRepo(),
	}
}

func (u *ArticleUsecase) CreateArticleCommand(params *article.CreateArticleCommand) (*resp.Response, error) {
	// Implementation for creating an article based on .NET CreateArticleCommand
	u.Logger.Info("Creating new article", map[string]interface{}{"title": *params.Title})

	// Convert SeoTags slice to string (comma-separated)
	var seoTags string
	if params.SeoTags != nil && len(params.SeoTags) > 0 {
		seoTags = strings.Join(params.SeoTags, ",")
	}

	// Create new article
	newArticle := domain.Article{
		Title:        *params.Title,
		Description:  *params.Description,
		Body:         *params.Body,
		Slug:         *params.Slug,
		SiteID:       *params.SiteID,
		SeoTags:      seoTags,
		UserID:       1, // Should come from auth context in real impl
		VisitedCount: 0,
		ReviewCount:  0,
		Rate:         0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}

	// Save the article to get its ID
	err := u.articleRepo.Create(newArticle)
	if err != nil {
		return nil, err
	}

	// Handle optional media relationships
	if params.MediaIDs != nil && len(params.MediaIDs) > 0 {
		for _, mediaID := range params.MediaIDs {
			err = u.articleRepo.AddMediaToArticle(newArticle.ID, mediaID)
			if err != nil {
				u.Logger.Errorf("Failed to add media %d to article %d: %v", mediaID, newArticle.ID, err)
				// Continue with other media instead of failing completely
			}
		}
	}

	// Handle optional category relationships
	if params.CategoryIDs != nil && len(params.CategoryIDs) > 0 {
		for _, categoryID := range params.CategoryIDs {
			err = u.articleRepo.AddCategoryToArticle(newArticle.ID, categoryID)
			if err != nil {
				u.Logger.Errorf("Failed to add category %d to article %d: %v", categoryID, newArticle.ID, err)
				// Continue with other categories instead of failing completely
			}
		}
	}

	// Return the created article
	// Convert domain.Article to map for response data
	articleMap := map[string]interface{}{
		"id":           newArticle.ID,
		"title":        newArticle.Title,
		"description":  newArticle.Description,
		"body":         newArticle.Body,
		"slug":         newArticle.Slug,
		"siteId":       newArticle.SiteID,
		"seoTags":      newArticle.SeoTags,
		"userId":       newArticle.UserID,
		"visitedCount": newArticle.VisitedCount,
		"reviewCount":  newArticle.ReviewCount,
		"rate":         newArticle.Rate,
		"createdAt":    newArticle.CreatedAt,
		"updatedAt":    newArticle.UpdatedAt,
		"isDeleted":    newArticle.IsDeleted,
	}
	return resp.NewResponseData(resp.Created, articleMap, "Article created successfully"), nil
}

func (u *ArticleUsecase) UpdateArticleCommand(params *article.UpdateArticleCommand) (*resp.Response, error) {
	// Implementation for updating an article based on .NET UpdateArticleCommand
	u.Logger.Info("Updating article", map[string]interface{}{"id": *params.ID})

	// Get existing article
	existingArticle, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this article
	// In a real implementation, check if the current user has rights to edit this article
	// This is equivalent to the gate.HasUserAccess(entity) call in .NET

	// Update fields if provided
	if params.Title != nil {
		existingArticle.Title = *params.Title
	}

	if params.Description != nil {
		existingArticle.Description = *params.Description
	}

	if params.Body != nil {
		existingArticle.Body = *params.Body
	}

	if params.Slug != nil {
		existingArticle.Slug = *params.Slug
	}

	// Update SeoTags if provided
	if params.SeoTags != nil {
		existingArticle.SeoTags = strings.Join(params.SeoTags, ",")
	}

	existingArticle.UpdatedAt = time.Now()

	// Save changes
	err = u.articleRepo.Update(existingArticle)
	if err != nil {
		return nil, err
	}

	// Handle media relationships if provided
	if params.MediaIDs != nil {
		// First remove all existing media associations
		err = u.articleRepo.RemoveAllMediaFromArticle(existingArticle.ID)
		if err != nil {
			u.Logger.Errorf("Failed to remove media from article %d: %v", existingArticle.ID, err)
		}

		// Then add the new ones
		for _, mediaID := range params.MediaIDs {
			err = u.articleRepo.AddMediaToArticle(existingArticle.ID, mediaID)
			if err != nil {
				u.Logger.Errorf("Failed to add media %d to article %d: %v", mediaID, existingArticle.ID, err)
			}
		}
	}

	// Handle category relationships if provided
	if params.CategoryIDs != nil {
		// First remove all existing category associations
		err = u.articleRepo.RemoveAllCategoriesFromArticle(existingArticle.ID)
		if err != nil {
			u.Logger.Errorf("Failed to remove categories from article %d: %v", existingArticle.ID, err)
		}

		// Then add the new ones
		for _, categoryID := range params.CategoryIDs {
			err = u.articleRepo.AddCategoryToArticle(existingArticle.ID, categoryID)
			if err != nil {
				u.Logger.Errorf("Failed to add category %d to article %d: %v", categoryID, existingArticle.ID, err)
			}
		}
	}

	// Convert domain.Article to map for response data
	articleMap := map[string]interface{}{
		"id":           existingArticle.ID,
		"title":        existingArticle.Title,
		"description":  existingArticle.Description,
		"body":         existingArticle.Body,
		"slug":         existingArticle.Slug,
		"siteId":       existingArticle.SiteID,
		"seoTags":      existingArticle.SeoTags,
		"userId":       existingArticle.UserID,
		"visitedCount": existingArticle.VisitedCount,
		"reviewCount":  existingArticle.ReviewCount,
		"rate":         existingArticle.Rate,
		"createdAt":    existingArticle.CreatedAt,
		"updatedAt":    existingArticle.UpdatedAt,
		"isDeleted":    existingArticle.IsDeleted,
	}
	return resp.NewResponseData(resp.Updated, articleMap, "Article updated successfully"), nil
}

func (u *ArticleUsecase) DeleteArticleCommand(params *article.DeleteArticleCommand) (*resp.Response, error) {
	// Implementation for deleting an article based on .NET DeleteArticleCommand
	u.Logger.Info("Deleting article", map[string]interface{}{"id": *params.ID})

	// Get the article first to ensure it exists
	_, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this article
	// In a real implementation, check if the current user has rights to delete this article
	// This is equivalent to the gate.HasUserAccess(entity) call in .NET

	// Delete the article
	err = u.articleRepo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return resp.NewResponse(resp.Deleted, "Article deleted successfully"), nil
}

func (u *ArticleUsecase) GetByIdArticleQuery(params *article.GetByIdArticleQuery) (*resp.Response, error) {
	// Implementation to get article by ID based on .NET GetByIdArticleQuery
	u.Logger.Info("Getting article by ID", map[string]interface{}{"id": *params.ID})

	// Get the article
	result, err := u.articleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this article
	// In a real implementation, check if the current user has rights to view this article
	// This is equivalent to the gate.HasUserAccess(entity) call in .NET

	// Get media information
	mediaItems, err := u.articleRepo.GetArticleMedia(result.ID)
	if err != nil {
		u.Logger.Errorf("Failed to get media for article %d: %v", result.ID, err)
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"article": result,
		"media":   mediaItems,
	}, "Article retrieved successfully"), nil
}

func (u *ArticleUsecase) GetSingleArticleQuery(params *article.GetSingleArticleQuery) (*resp.Response, error) {
	// Implementation to get article by slug based on .NET GetSingleArticleQuery
	u.Logger.Info("Getting article by slug", map[string]interface{}{
		"slug":   *params.Slug,
		"siteID": *params.SiteID,
	})

	// Get the article by slug and site ID
	result, err := u.articleRepo.GetBySlugAndSiteID(*params.Slug, *params.SiteID)
	if err != nil {
		return nil, err
	}

	// Get media information
	mediaItems, err := u.articleRepo.GetArticleMedia(result.ID)
	if err != nil {
		u.Logger.Errorf("Failed to get media for article %d: %v", result.ID, err)
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"article": result,
		"media":   mediaItems,
	}, "Article retrieved successfully"), nil
}

func (u *ArticleUsecase) GetAllArticleQuery(params *article.GetAllArticleQuery) (*resp.Response, error) {
	// Implementation to get all articles by site ID, based on .NET GetAllArticleQuery
	u.Logger.Info("Getting all articles by site ID", map[string]interface{}{"siteID": *params.SiteID})

	// Check if user has access to this site
	// In a real implementation, check if the current user has rights to view articles for this site
	// This is equivalent to the gate.HasSiteAccess(request.SiteId) call in .NET

	result, count, err := u.articleRepo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// For each article, get media information
	// In a more efficient implementation, this would be done in a single query
	articlesWithMedia := make([]map[string]interface{}, len(result))
	for i, article := range result {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			u.Logger.Errorf("Failed to get media for article %d: %v", article.ID, err)
			media = []domain.Media{}
		}

		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items": articlesWithMedia,
		"total": count,
	}, "Articles retrieved successfully"), nil
}

func (u *ArticleUsecase) GetArticleByCategoryQuery(params *article.GetArticleByCategoryQuery) (*resp.Response, error) {
	// Implementation to get articles by category, based on .NET GetArticleByCategoryQuery
	u.Logger.Info("Getting articles by category slug", map[string]interface{}{
		"slug":   *params.Slug,
		"siteID": *params.SiteID,
	})

	// First get the category by slug
	category, err := u.categoryRepo.GetBySlugAndSiteID(*params.Slug, *params.SiteID)
	if err != nil {
		return nil, err
	}

	// Then get articles for this category
	result, count, err := u.articleRepo.GetAllByCategoryID(category.ID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// For each article, get media information
	articlesWithMedia := make([]map[string]interface{}, len(result))
	for i, article := range result {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			u.Logger.Errorf("Failed to get media for article %d: %v", article.ID, err)
			media = []domain.Media{}
		}

		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":    articlesWithMedia,
		"total":    count,
		"category": category,
	}, "Articles retrieved successfully"), nil
}

func (u *ArticleUsecase) GetByFiltersSortArticleQuery(params *article.GetByFiltersSortArticleQuery) (*resp.Response, error) {
	// Implementation to get articles with filtering and sorting, based on .NET GetByFiltersSortArticleQuery
	u.Logger.Info("Getting articles with filters and sorting", map[string]interface{}{"siteID": *params.SiteID})

	// This is a more complex query that would need special handling
	// For now, we'll implement a basic version that just calls through to a repository method
	result, count, err := u.articleRepo.GetAllByFilterAndSort(
		*params.SiteID,
		params.SelectedFilters,
		params.SelectedSort,
		params.PaginationRequestDto,
	)

	if err != nil {
		return nil, err
	}

	// For each article, get media information
	articlesWithMedia := make([]map[string]interface{}, len(result))
	for i, article := range result {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			u.Logger.Errorf("Failed to get media for article %d: %v", article.ID, err)
			media = []domain.Media{}
		}

		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items": articlesWithMedia,
		"total": count,
	}, "Articles retrieved successfully"), nil
}

func (u *ArticleUsecase) AdminGetAllArticleQuery(params *article.AdminGetAllArticleQuery) (*resp.Response, error) {
	// Implementation to get all articles for admin, based on .NET AdminGetAllArticleQuery
	u.Logger.Info("Admin getting all articles", map[string]interface{}{})

	// Check if user has admin access
	// In a real implementation, check if the current user has admin rights
	// This is equivalent to the gate.HasSiteAccess() call in .NET

	result, count, err := u.articleRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// For each article, get media information
	articlesWithMedia := make([]map[string]interface{}, len(result))
	for i, article := range result {
		media, err := u.articleRepo.GetArticleMedia(article.ID)
		if err != nil {
			u.Logger.Errorf("Failed to get media for article %d: %v", article.ID, err)
			media = []domain.Media{}
		}

		articlesWithMedia[i] = map[string]interface{}{
			"article": article,
			"media":   media,
		}
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items": articlesWithMedia,
		"total": count,
	}, "Articles retrieved successfully"), nil
}
