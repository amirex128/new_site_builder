package headerfooterusecase

import (
	"encoding/json"
	"fmt"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/header_footer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type HeaderFooterUsecase struct {
	logger sflogger.Logger
	repo   repository.IHeaderFooterRepository
}

func NewHeaderFooterUsecase(c contract.IContainer) *HeaderFooterUsecase {
	return &HeaderFooterUsecase{
		logger: c.GetLogger(),
		repo:   c.GetHeaderFooterRepo(),
	}
}

func (u *HeaderFooterUsecase) CreateHeaderFooterCommand(params *header_footer.CreateHeaderFooterCommand) (any, error) {
	// Implementation for creating a header/footer
	fmt.Println(params)

	// Convert body to JSON string
	bodyJSON, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}

	headerFooter := domain.HeaderFooter{
		SiteID:    *params.SiteID,
		Title:     *params.Title,
		IsMain:    *params.IsMain,
		Body:      string(bodyJSON),
		Type:      int(*params.Type),
		UserID:    1, // Should come from auth context
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	err = u.repo.Create(headerFooter)
	if err != nil {
		return nil, err
	}

	return headerFooter, nil
}

func (u *HeaderFooterUsecase) UpdateHeaderFooterCommand(params *header_footer.UpdateHeaderFooterCommand) (any, error) {
	// Implementation for updating a header/footer
	fmt.Println(params)

	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	existingHeaderFooter.SiteID = *params.SiteID

	if params.Title != nil {
		existingHeaderFooter.Title = *params.Title
	}

	existingHeaderFooter.IsMain = *params.IsMain
	existingHeaderFooter.Type = int(*params.Type)

	if params.Body != nil {
		bodyJSON, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		existingHeaderFooter.Body = string(bodyJSON)
	}

	existingHeaderFooter.UpdatedAt = time.Now()

	err = u.repo.Update(existingHeaderFooter)
	if err != nil {
		return nil, err
	}

	return existingHeaderFooter, nil
}

func (u *HeaderFooterUsecase) DeleteHeaderFooterCommand(params *header_footer.DeleteHeaderFooterCommand) (any, error) {
	// Implementation for deleting a header/footer
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *HeaderFooterUsecase) GetByIdHeaderFooterQuery(params *header_footer.GetByIdHeaderFooterQuery) (any, error) {
	// Implementation to get header/footer by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *HeaderFooterUsecase) GetAllHeaderFooterQuery(params *header_footer.GetAllHeaderFooterQuery) (any, error) {
	// Implementation to get all header/footers
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

func (u *HeaderFooterUsecase) AdminGetAllHeaderFooterQuery(params *header_footer.AdminGetAllHeaderFooterQuery) (any, error) {
	// Implementation for admin to get all header/footers
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
