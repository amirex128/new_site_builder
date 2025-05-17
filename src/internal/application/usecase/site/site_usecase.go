package siteusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/site"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type SiteUsecase struct {
	logger sflogger.Logger
	repo   repository.ISiteRepository
}

func NewSiteUsecase(c contract.IContainer) *SiteUsecase {
	return &SiteUsecase{
		logger: c.GetLogger(),
		repo:   c.GetSiteRepo(),
	}
}

func (u *SiteUsecase) CreateSiteCommand(params *site.CreateSiteCommand) (any, error) {
	// Implementation for creating a site
	fmt.Println(params)

	newSite := domain.Site{
		Domain:     *params.Domain,
		DomainType: strconv.Itoa(int(*params.DomainType)),
		Name:       *params.Name,
		Status:     strconv.Itoa(int(*params.Status)),
		SiteType:   strconv.Itoa(int(*params.SiteType)),
		UserID:     1, // Should come from auth context
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	err := u.repo.Create(newSite)
	if err != nil {
		return nil, err
	}

	return newSite, nil
}

func (u *SiteUsecase) UpdateSiteCommand(params *site.UpdateSiteCommand) (any, error) {
	// Implementation for updating a site
	fmt.Println(params)

	existingSite, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	if params.Domain != nil {
		existingSite.Domain = *params.Domain
	}

	if params.DomainType != nil {
		existingSite.DomainType = strconv.Itoa(int(*params.DomainType))
	}

	if params.Name != nil {
		existingSite.Name = *params.Name
	}

	if params.Status != nil {
		existingSite.Status = strconv.Itoa(int(*params.Status))
	}

	if params.SiteType != nil {
		existingSite.SiteType = strconv.Itoa(int(*params.SiteType))
	}

	existingSite.UpdatedAt = time.Now()

	err = u.repo.Update(existingSite)
	if err != nil {
		return nil, err
	}

	return existingSite, nil
}

func (u *SiteUsecase) DeleteSiteCommand(params *site.DeleteSiteCommand) (any, error) {
	// Implementation for deleting a site
	fmt.Println(params)

	err := u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": params.ID,
	}, nil
}

func (u *SiteUsecase) GetByIdSiteQuery(params *site.GetByIdSiteQuery) (any, error) {
	// Implementation to get site by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *SiteUsecase) GetAllSiteQuery(params *site.GetAllSiteQuery) (any, error) {
	// Implementation to get all sites for current user
	fmt.Println(params)

	// In a real implementation, get the user ID from the auth context
	userID := int64(1)

	result, count, err := u.repo.GetAllByUserID(userID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *SiteUsecase) AdminGetAllSiteQuery(params *site.AdminGetAllSiteQuery) (any, error) {
	// Implementation to get all sites for admin
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
