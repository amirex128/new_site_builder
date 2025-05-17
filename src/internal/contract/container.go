package contract

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
)

type IContainer interface {
	GetArticleRepo() repository.IArticleRepository
	GetBasketRepo() repository.IBasketRepository
	GetBasketItemRepo() repository.IBasketItemRepository
	GetArticleCategoryRepo() repository.IArticleCategoryRepository
	GetCreditRepo() repository.ICreditRepository
	GetCouponRepo() repository.ICouponRepository
	GetCustomerRepo() repository.ICustomerRepository
	GetDefaultThemeRepo() repository.IDefaultThemeRepository
	GetDiscountRepo() repository.IDiscountRepository
	GetFileItemRepo() repository.IFileItemRepository
	GetGatewayRepo() repository.IGatewayRepository
	GetHeaderFooterRepo() repository.IHeaderFooterRepository
	GetMediaRepo() repository.IMediaRepository
	GetOrderRepo() repository.IOrderRepository
	GetOrderItemRepo() repository.IOrderItemRepository
	GetPageRepo() repository.IPageRepository
	GetParbadPaymentRepo() repository.IParbadPaymentRepository
	GetParbadTransactionRepo() repository.IParbadTransactionRepository
	GetPaymentRepo() repository.IPaymentRepository
	GetProductRepo() repository.IProductRepository
	GetProductAttributeRepo() repository.IProductAttributeRepository
	GetProductCategoryRepo() repository.IProductCategoryRepository
	GetProductReviewRepo() repository.IProductReviewRepository
	GetProductVariantRepo() repository.IProductVariantRepository
	GetReturnItemRepo() repository.IReturnItemRepository
	GetSettingRepo() repository.ISettingRepository
	GetSiteRepo() repository.ISiteRepository
	GetStorageRepo() repository.IStorageRepository
	GetTicketRepo() repository.ITicketRepository
	GetCustomerTicketRepo() repository.ICustomerTicketRepository
	GetUserRepo() repository.IUserRepository
	GetConfig() IConfig
	GetFoodPartyCash() cache.ICacheService
	GetStockCacheTransient() cache.ICacheService
	GetLogger() sflogger.Logger
}
