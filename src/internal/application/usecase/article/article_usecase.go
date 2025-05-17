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
