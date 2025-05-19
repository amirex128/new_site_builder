package headerfooterusecase

import (
	"encoding/json"
	"errors"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/header_footer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type HeaderFooterUsecase struct {
	logger         sflogger.Logger
	repo           repository.IHeaderFooterRepository
	siteRepo       repository.ISiteRepository
	authContextSvc common.IAuthContextService
}

func NewHeaderFooterUsecase(c contract.IContainer) *HeaderFooterUsecase {
	return &HeaderFooterUsecase{
		logger:         c.GetLogger(),
		repo:           c.GetHeaderFooterRepo(),
		siteRepo:       c.GetSiteRepo(),
		authContextSvc: c.GetAuthContextTransiantService(),
	}
}

func (u *HeaderFooterUsecase) CreateHeaderFooterCommand(params *header_footer.CreateHeaderFooterCommand) (any, error) {
	u.logger.Info("CreateHeaderFooterCommand called", map[string]interface{}{
		"siteId": *params.SiteID,
		"title":  *params.Title,
		"type":   *params.Type,
	})

	// Check if site exists
	_, err := u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Get user ID from auth context
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

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
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Create in repository
	err = u.repo.Create(headerFooter)
	if err != nil {
		return nil, err
	}

	// Get created header/footer with relations
	createdHeaderFooter, err := u.repo.GetByID(headerFooter.ID)
	if err != nil {
		return nil, err
	}

	return createdHeaderFooter, nil
}

func (u *HeaderFooterUsecase) UpdateHeaderFooterCommand(params *header_footer.UpdateHeaderFooterCommand) (any, error) {
	u.logger.Info("UpdateHeaderFooterCommand called", map[string]interface{}{
		"id":     *params.ID,
		"siteId": *params.SiteID,
		"type":   *params.Type,
	})

	// Get existing header/footer
	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("هدر/فوتر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingHeaderFooter.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این هدر/فوتر دسترسی ندارید")
	}

	// Check if site exists
	_, err = u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("سایت مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Update fields
	existingHeaderFooter.SiteID = *params.SiteID
	existingHeaderFooter.IsMain = *params.IsMain
	existingHeaderFooter.Type = int(*params.Type)

	if params.Title != nil {
		existingHeaderFooter.Title = *params.Title
	}

	if params.Body != nil {
		bodyJSON, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		existingHeaderFooter.Body = string(bodyJSON)
	}

	existingHeaderFooter.UpdatedAt = time.Now()

	// Update in repository
	err = u.repo.Update(existingHeaderFooter)
	if err != nil {
		return nil, err
	}

	// Get updated header/footer with relations
	updatedHeaderFooter, err := u.repo.GetByID(existingHeaderFooter.ID)
	if err != nil {
		return nil, err
	}

	return updatedHeaderFooter, nil
}

func (u *HeaderFooterUsecase) DeleteHeaderFooterCommand(params *header_footer.DeleteHeaderFooterCommand) (any, error) {
	u.logger.Info("DeleteHeaderFooterCommand called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get existing header/footer
	existingHeaderFooter, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("هدر/فوتر مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if existingHeaderFooter.UserID != userID && !isAdmin {
		return nil, errors.New("شما به این هدر/فوتر دسترسی ندارید")
	}

	// Delete header/footer
	err = u.repo.Delete(*params.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": existingHeaderFooter.ID,
	}, nil
}

func (u *HeaderFooterUsecase) GetByIdHeaderFooterQuery(params *header_footer.GetByIdHeaderFooterQuery) (any, error) {
	u.logger.Info("GetByIdHeaderFooterQuery called", map[string]interface{}{
		"id":     params.ID,
		"siteId": *params.SiteID,
	})

	if params.ID != nil {
		// Get by ID
		headerFooter, err := u.repo.GetByID(*params.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("هدر/فوتر مورد نظر یافت نشد")
			}
			return nil, err
		}

		// Check if header/footer belongs to the specified site
		if headerFooter.SiteID != *params.SiteID {
			return nil, errors.New("هدر/فوتر مورد نظر متعلق به سایت مشخص شده نیست")
		}

		// Check if type matches if specified
		if params.Type != nil && headerFooter.Type != int(*params.Type) {
			return nil, errors.New("نوع هدر/فوتر مطابقت ندارد")
		}

		return headerFooter, nil
	} else if len(params.IDs) > 0 {
		// Get by multiple IDs
		// In a monolithic approach, we would need to implement a method to get by IDs in the repository
		// For now, we'll use a simple approach of getting each one individually
		var result []domain.HeaderFooter
		for _, id := range params.IDs {
			headerFooter, err := u.repo.GetByID(id)
			if err == nil && headerFooter.SiteID == *params.SiteID {
				result = append(result, headerFooter)
			}
		}
		return result, nil
	}

	return nil, errors.New("شناسه هدر/فوتر الزامی است")
}

func (u *HeaderFooterUsecase) GetAllHeaderFooterQuery(params *header_footer.GetAllHeaderFooterQuery) (any, error) {
	u.logger.Info("GetAllHeaderFooterQuery called", map[string]interface{}{
		"siteId": *params.SiteID,
		"type":   params.Type,
	})

	// Check site access (in a real implementation, we would verify user has access to the site)
	// We don't need the userID here, so remove the unused variable

	// Get all header/footers for the site
	headerFooters, count, err := u.repo.GetAllBySiteID(*params.SiteID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	// Filter by type if specified
	if params.Type != nil {
		var filteredHeaderFooters []domain.HeaderFooter
		for _, hf := range headerFooters {
			if hf.Type == int(*params.Type) {
				filteredHeaderFooters = append(filteredHeaderFooters, hf)
			}
		}
		headerFooters = filteredHeaderFooters
		count = int64(len(filteredHeaderFooters))
	}

	return map[string]interface{}{
		"items":     headerFooters,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *HeaderFooterUsecase) AdminGetAllHeaderFooterQuery(params *header_footer.AdminGetAllHeaderFooterQuery) (any, error) {
	u.logger.Info("AdminGetAllHeaderFooterQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all header/footers across all sites for admin
	headerFooters, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items":     headerFooters,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *HeaderFooterUsecase) GetHeaderFooterByDomainOrSiteIdQuery(params *header_footer.GetHeaderFooterByDomainOrSiteIdQuery) (any, error) {
	u.logger.Info("GetHeaderFooterByDomainOrSiteIdQuery called", map[string]interface{}{
		"siteId": params.SiteID,
		"domain": params.Domain,
	})

	var siteID int64

	if params.SiteID != nil {
		// Use provided site ID
		siteID = *params.SiteID
	} else if params.Domain != nil {
		// Get site by domain
		site, err := u.siteRepo.GetByDomain(*params.Domain)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("سایت مورد نظر یافت نشد")
			}
			return nil, err
		}
		siteID = site.ID
	} else {
		return nil, errors.New("شناسه سایت یا دامنه الزامی است")
	}

	// Get main header/footer for site
	// In a real implementation, we would have a dedicated repository method for this
	// For now, we'll get all and filter for the main one
	headerFooters, _, err := u.repo.GetAllBySiteID(siteID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return nil, err
	}

	// Find main header/footer
	for _, hf := range headerFooters {
		if hf.IsMain {
			return hf, nil
		}
	}

	return nil, errors.New("هدر/فوتر اصلی برای سایت مورد نظر یافت نشد")
}
