package articleusecase

import (
	"fmt"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/article"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
)

type ArticleUsecase struct {
	logger sflogger.Logger
	repo   repository.IArticleRepository
}

func NewArticleUsecase(c contract.IContainer) *ArticleUsecase {
	return &ArticleUsecase{
		logger: c.GetLogger(),
		repo:   c.GetArticleRepo(),
	}
}

func (u *ArticleUsecase) CreateArticleCommand(params *article.CreateArticleCommand) ([]any, error) {

	fmt.Println(params)

	return nil, nil
}

func (u *ArticleUsecase) UpdateArticleCommand(params *article.UpdateArticleCommand) ([]any, error) {

	fmt.Println(params)

	return nil, nil
}

func (u *ArticleUsecase) DeleteArticleCommand(params *article.DeleteArticleCommand) ([]any, error) {

	fmt.Println(params)

	return nil, nil
}

func (u *ArticleUsecase) GetByIdArticleQuery(params *article.GetByIdArticleQuery) (any, error) {
	// Implementation to get article by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *ArticleUsecase) GetAllArticleQuery(params *article.GetAllArticleQuery) (any, error) {
	// Implementation to get all articles
	fmt.Println(params)

	result, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *ArticleUsecase) GetByFiltersSortArticleQuery(params *article.GetByFiltersSortArticleQuery) (any, error) {
	// Implementation to get articles with filtering and sorting
	fmt.Println(params)

	// TODO: Add implementation with filtering and sorting
	// For now, we are simply reusing the GetAllBySiteID method as a placeholder
	// In a real implementation, you would apply the filters and sorting based on params.SelectedFilters and params.SelectedSort
	result, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *ArticleUsecase) AdminGetAllArticleQuery(params *article.AdminGetAllArticleQuery) (any, error) {
	// Implementation to get all articles for admin
	fmt.Println(params)

	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
