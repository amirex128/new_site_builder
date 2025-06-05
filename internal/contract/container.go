package contract

import (
	"context"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	service2 "github.com/amirex128/new_site_builder/internal/contract/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IContainer provides methods to access all dependencies
type IContainer interface {
	GetCtx() context.Context
	GetLogger() sflogger.Logger

	// Services
	GetAuthTransientService() func(c *gin.Context) service2.IAuthService
	GetConfig() IConfig
	GetMainCache() service2.ICacheService
	GetStockCacheTransient() service2.ICacheService
	GetDB() *gorm.DB
	GetStorageService() service2.IStorageService
	GetIdentityService() service2.IIdentityService
	GetPaymentService() service2.IPaymentService
	GetMessageService() service2.IMessageService

	// Repositories
	GetAddressRepo() repository2.IAddressRepository
	GetArticleRepo() repository2.IArticleRepository
	GetArticleCategoryRepo() repository2.IArticleCategoryRepository
	GetBasketRepo() repository2.IBasketRepository
	GetBasketItemRepo() repository2.IBasketItemRepository
	GetCustomerRepo() repository2.ICustomerRepository
	GetCustomerTicketRepo() repository2.ICustomerTicketRepository
	GetCommentRepo() repository2.ICommentRepository
	GetCustomerCommentRepo() repository2.ICustomerCommentRepository
	GetTicketMediaRepo() repository2.ITicketMediaRepository
	GetCustomerTicketMediaRepo() repository2.ICustomerTicketMediaRepository
	GetDefaultThemeRepo() repository2.IDefaultThemeRepository
	GetDiscountRepo() repository2.IDiscountRepository
	GetFileItemRepo() repository2.IFileItemRepository
	GetHeaderFooterRepo() repository2.IHeaderFooterRepository
	GetOrderRepo() repository2.IOrderRepository
	GetOrderItemRepo() repository2.IOrderItemRepository
	GetPageRepo() repository2.IPageRepository
	GetPageArticleUsageRepo() repository2.IPageArticleUsageRepository
	GetPageProductUsageRepo() repository2.IPageProductUsageRepository
	GetPageHeaderFooterUsageRepo() repository2.IPageHeaderFooterUsageRepository
	GetPaymentRepo() repository2.IPaymentRepository
	GetPlanRepo() repository2.IPlanRepository
	GetProductRepo() repository2.IProductRepository
	GetProductCategoryRepo() repository2.IProductCategoryRepository
	GetProductReviewRepo() repository2.IProductReviewRepository
	GetRoleRepo() repository2.IRoleRepository
	GetSiteRepo() repository2.ISiteRepository
	GetTicketRepo() repository2.ITicketRepository
	GetUnitPriceRepo() repository2.IUnitPriceRepository
	GetUserRepo() repository2.IUserRepository

	// Additional repositories
	GetCreditRepo() repository2.ICreditRepository
	GetCouponRepo() repository2.ICouponRepository
	GetGatewayRepo() repository2.IGatewayRepository
	GetMediaRepo() repository2.IMediaRepository
	GetProductAttributeRepo() repository2.IProductAttributeRepository
	GetProductVariantRepo() repository2.IProductVariantRepository
	GetReturnItemRepo() repository2.IReturnItemRepository
	GetSettingRepo() repository2.ISettingRepository
	GetStorageRepo() repository2.IStorageRepository
	GetCityRepo() repository2.ICityRepository
	GetProvinceRepo() repository2.IProvinceRepository
	GetParbadPaymentRepo() repository2.IParbadPaymentRepository
	GetParbadTransactionRepo() repository2.IParbadTransactionRepository
	GetConfigurationRepo() repository2.IConfigurationRepository
}
