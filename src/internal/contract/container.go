package contract

import (
	"context"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IContainer provides methods to access all dependencies
type IContainer interface {
	GetCtx() context.Context
	GetLogger() sflogger.Logger

	// Services
	GetAuthTransientService() func(c *gin.Context) service.IAuthService
	GetConfig() IConfig
	GetMainCache() service.ICacheService
	GetStockCacheTransient() service.ICacheService
	GetDB() *gorm.DB
	GetStorageService() service.IStorageService
	GetIdentityService() service.IIdentityService
	GetPaymentService() service.IPaymentService
	GetMessageService() service.IMessageService

	// Repositories
	GetAddressRepo() repository.IAddressRepository
	GetArticleRepo() repository.IArticleRepository
	GetArticleCategoryRepo() repository.IArticleCategoryRepository
	GetBasketRepo() repository.IBasketRepository
	GetBasketItemRepo() repository.IBasketItemRepository
	GetCustomerRepo() repository.ICustomerRepository
	GetCustomerTicketRepo() repository.ICustomerTicketRepository
	GetCommentRepo() repository.ICommentRepository
	GetCustomerCommentRepo() repository.ICustomerCommentRepository
	GetTicketMediaRepo() repository.ITicketMediaRepository
	GetCustomerTicketMediaRepo() repository.ICustomerTicketMediaRepository
	GetDefaultThemeRepo() repository.IDefaultThemeRepository
	GetDiscountRepo() repository.IDiscountRepository
	GetFileItemRepo() repository.IFileItemRepository
	GetHeaderFooterRepo() repository.IHeaderFooterRepository
	GetOrderRepo() repository.IOrderRepository
	GetOrderItemRepo() repository.IOrderItemRepository
	GetPageRepo() repository.IPageRepository
	GetPageArticleUsageRepo() repository.IPageArticleUsageRepository
	GetPageProductUsageRepo() repository.IPageProductUsageRepository
	GetPageHeaderFooterUsageRepo() repository.IPageHeaderFooterUsageRepository
	GetPaymentRepo() repository.IPaymentRepository
	GetPlanRepo() repository.IPlanRepository
	GetProductRepo() repository.IProductRepository
	GetProductCategoryRepo() repository.IProductCategoryRepository
	GetProductReviewRepo() repository.IProductReviewRepository
	GetRoleRepo() repository.IRoleRepository
	GetSiteRepo() repository.ISiteRepository
	GetTicketRepo() repository.ITicketRepository
	GetUnitPriceRepo() repository.IUnitPriceRepository
	GetUserRepo() repository.IUserRepository

	// Additional repositories
	GetCreditRepo() repository.ICreditRepository
	GetCouponRepo() repository.ICouponRepository
	GetGatewayRepo() repository.IGatewayRepository
	GetMediaRepo() repository.IMediaRepository
	GetProductAttributeRepo() repository.IProductAttributeRepository
	GetProductVariantRepo() repository.IProductVariantRepository
	GetReturnItemRepo() repository.IReturnItemRepository
	GetSettingRepo() repository.ISettingRepository
	GetStorageRepo() repository.IStorageRepository
	GetCityRepo() repository.ICityRepository
	GetProvinceRepo() repository.IProvinceRepository
	GetParbadPaymentRepo() repository.IParbadPaymentRepository
	GetParbadTransactionRepo() repository.IParbadTransactionRepository
}
