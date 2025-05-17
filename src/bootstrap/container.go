package bootstrap

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	cache2 "github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
)

// Container
type Container struct {
	Config         contract.IConfig
	MemoryLoader   cache2.IMemoryLoader
	FoodPartyCache cache2.ICacheService

	stockCacheTransient func() cache2.ICacheService

	Logger                sflogger.Logger
	ArticleRepo           repository.IArticleRepository
	BasketRepo            repository.IBasketRepository
	BasketItemRepo        repository.IBasketItemRepository
	BlogCategoryRepo      repository.IBlogCategoryRepository
	CreditRepo            repository.ICreditRepository
	CouponRepo            repository.ICouponRepository
	CustomerRepo          repository.ICustomerRepository
	DefaultThemeRepo      repository.IDefaultThemeRepository
	DiscountRepo          repository.IDiscountRepository
	FileItemRepo          repository.IFileItemRepository
	GatewayRepo           repository.IGatewayRepository
	HeaderFooterRepo      repository.IHeaderFooterRepository
	MediaRepo             repository.IMediaRepository
	OrderRepo             repository.IOrderRepository
	OrderItemRepo         repository.IOrderItemRepository
	PageRepo              repository.IPageRepository
	ParbadPaymentRepo     repository.IParbadPaymentRepository
	ParbadTransactionRepo repository.IParbadTransactionRepository
	PaymentRepo           repository.IPaymentRepository
	ProductRepo           repository.IProductRepository
	ProductAttributeRepo  repository.IProductAttributeRepository
	ProductCategoryRepo   repository.IProductCategoryRepository
	ProductReviewRepo     repository.IProductReviewRepository
	ProductVariantRepo    repository.IProductVariantRepository
	ReturnItemRepo        repository.IReturnItemRepository
	SettingRepo           repository.ISettingRepository
	SiteRepo              repository.ISiteRepository
	StorageRepo           repository.IStorageRepository
	TicketRepo            repository.ITicketRepository
	CustomerTicketRepo    repository.ICustomerTicketRepository
	UserRepo              repository.IUserRepository
}

func (c *Container) GetConfig() contract.IConfig {
	return c.Config
}

func (c *Container) GetArticleRepo() repository.IArticleRepository {
	return c.ArticleRepo
}

func (c *Container) GetBasketRepo() repository.IBasketRepository {
	return c.BasketRepo
}

func (c *Container) GetBasketItemRepo() repository.IBasketItemRepository {
	return c.BasketItemRepo
}

func (c *Container) GetBlogCategoryRepo() repository.IBlogCategoryRepository {
	return c.BlogCategoryRepo
}

func (c *Container) GetCreditRepo() repository.ICreditRepository {
	return c.CreditRepo
}

func (c *Container) GetCouponRepo() repository.ICouponRepository {
	return c.CouponRepo
}

func (c *Container) GetCustomerRepo() repository.ICustomerRepository {
	return c.CustomerRepo
}

func (c *Container) GetDefaultThemeRepo() repository.IDefaultThemeRepository {
	return c.DefaultThemeRepo
}

func (c *Container) GetDiscountRepo() repository.IDiscountRepository {
	return c.DiscountRepo
}

func (c *Container) GetFileItemRepo() repository.IFileItemRepository {
	return c.FileItemRepo
}

func (c *Container) GetGatewayRepo() repository.IGatewayRepository {
	return c.GatewayRepo
}

func (c *Container) GetHeaderFooterRepo() repository.IHeaderFooterRepository {
	return c.HeaderFooterRepo
}

func (c *Container) GetMediaRepo() repository.IMediaRepository {
	return c.MediaRepo
}

func (c *Container) GetOrderRepo() repository.IOrderRepository {
	return c.OrderRepo
}

func (c *Container) GetOrderItemRepo() repository.IOrderItemRepository {
	return c.OrderItemRepo
}

func (c *Container) GetPageRepo() repository.IPageRepository {
	return c.PageRepo
}

func (c *Container) GetParbadPaymentRepo() repository.IParbadPaymentRepository {
	return c.ParbadPaymentRepo
}

func (c *Container) GetParbadTransactionRepo() repository.IParbadTransactionRepository {
	return c.ParbadTransactionRepo
}

func (c *Container) GetPaymentRepo() repository.IPaymentRepository {
	return c.PaymentRepo
}

func (c *Container) GetProductRepo() repository.IProductRepository {
	return c.ProductRepo
}

func (c *Container) GetProductAttributeRepo() repository.IProductAttributeRepository {
	return c.ProductAttributeRepo
}

func (c *Container) GetProductCategoryRepo() repository.IProductCategoryRepository {
	return c.ProductCategoryRepo
}

func (c *Container) GetProductReviewRepo() repository.IProductReviewRepository {
	return c.ProductReviewRepo
}

func (c *Container) GetProductVariantRepo() repository.IProductVariantRepository {
	return c.ProductVariantRepo
}

func (c *Container) GetReturnItemRepo() repository.IReturnItemRepository {
	return c.ReturnItemRepo
}

func (c *Container) GetSettingRepo() repository.ISettingRepository {
	return c.SettingRepo
}

func (c *Container) GetSiteRepo() repository.ISiteRepository {
	return c.SiteRepo
}

func (c *Container) GetStorageRepo() repository.IStorageRepository {
	return c.StorageRepo
}

func (c *Container) GetTicketRepo() repository.ITicketRepository {
	return c.TicketRepo
}

func (c *Container) GetCustomerTicketRepo() repository.ICustomerTicketRepository {
	return c.CustomerTicketRepo
}

func (c *Container) GetUserRepo() repository.IUserRepository {
	return c.UserRepo
}

func (c *Container) GetFoodPartyCash() cache2.ICacheService {
	return c.FoodPartyCache
}

func (c *Container) GetStockCacheTransient() cache2.ICacheService {
	return c.stockCacheTransient()
}

func (c *Container) GetLogger() sflogger.Logger {
	return c.Logger
}
