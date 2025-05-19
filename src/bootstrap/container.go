package bootstrap

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
	"github.com/amirex128/new_site_builder/src/internal/contract/service/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Container
type Container struct {
	Config                      contract.IConfig
	MemoryLoader                cache.IMemoryLoader
	MainCache                   cache.ICacheService
	AuthContextTransientService func(c *gin.Context) common.IAuthContextService
	IdentityService             common.IIdentityService
	StorageService              storage.IStorageService
	stockCacheTransient         func() cache.ICacheService
	PaymentService              service.IPaymentService

	Logger                    sflogger.Logger
	ArticleRepo               repository.IArticleRepository
	BasketRepo                repository.IBasketRepository
	BasketItemRepo            repository.IBasketItemRepository
	ArticleCategoryRepo       repository.IArticleCategoryRepository
	CreditRepo                repository.ICreditRepository
	CouponRepo                repository.ICouponRepository
	CustomerRepo              repository.ICustomerRepository
	DefaultThemeRepo          repository.IDefaultThemeRepository
	DiscountRepo              repository.IDiscountRepository
	FileItemRepo              repository.IFileItemRepository
	GatewayRepo               repository.IGatewayRepository
	HeaderFooterRepo          repository.IHeaderFooterRepository
	MediaRepo                 repository.IMediaRepository
	OrderRepo                 repository.IOrderRepository
	OrderItemRepo             repository.IOrderItemRepository
	PageRepo                  repository.IPageRepository
	ParbadPaymentRepo         repository.IParbadPaymentRepository
	ParbadTransactionRepo     repository.IParbadTransactionRepository
	PaymentRepo               repository.IPaymentRepository
	ProductRepo               repository.IProductRepository
	ProductAttributeRepo      repository.IProductAttributeRepository
	ProductCategoryRepo       repository.IProductCategoryRepository
	ProductReviewRepo         repository.IProductReviewRepository
	ProductVariantRepo        repository.IProductVariantRepository
	ReturnItemRepo            repository.IReturnItemRepository
	SettingRepo               repository.ISettingRepository
	SiteRepo                  repository.ISiteRepository
	StorageRepo               repository.IStorageRepository
	TicketRepo                repository.ITicketRepository
	CustomerTicketRepo        repository.ICustomerTicketRepository
	UserRepo                  repository.IUserRepository
	UnitPriceRepo             repository.IUnitPriceRepository
	AddressRepo               repository.IAddressRepository
	CityRepo                  repository.ICityRepository
	ProvinceRepo              repository.IProvinceRepository
	PlanRepo                  repository.IPlanRepository
	RoleRepo                  repository.IRoleRepository
	PermissionRepo            repository.IPermissionRepository
	CommentRepo               repository.ICommentRepository
	CustomerCommentRepo       repository.ICustomerCommentRepository
	TicketMediaRepo           repository.ITicketMediaRepository
	CustomerTicketMediaRepo   repository.ICustomerTicketMediaRepository
	PageArticleUsageRepo      repository.IPageArticleUsageRepository
	PageProductUsageRepo      repository.IPageProductUsageRepository
	PageHeaderFooterUsageRepo repository.IPageHeaderFooterUsageRepository
}

func (c *Container) GetCommentRepo() repository.ICommentRepository {
	return c.CommentRepo
}

func (c *Container) GetCustomerCommentRepo() repository.ICustomerCommentRepository {
	return c.CustomerCommentRepo
}

func (c *Container) GetTicketMediaRepo() repository.ITicketMediaRepository {
	return c.TicketMediaRepo
}

func (c *Container) GetCustomerTicketMediaRepo() repository.ICustomerTicketMediaRepository {
	return c.CustomerTicketMediaRepo
}

func (c *Container) GetPageArticleUsageRepo() repository.IPageArticleUsageRepository {
	return c.PageArticleUsageRepo
}

func (c *Container) GetPageProductUsageRepo() repository.IPageProductUsageRepository {
	return c.PageProductUsageRepo
}

func (c *Container) GetPageHeaderFooterUsageRepo() repository.IPageHeaderFooterUsageRepository {
	return c.PageHeaderFooterUsageRepo
}

func (c *Container) GetStorageService() storage.IStorageService {
	return c.StorageService
}

func (c *Container) GetIdentityService() common.IIdentityService {
	return c.IdentityService
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

func (c *Container) GetArticleCategoryRepo() repository.IArticleCategoryRepository {
	return c.ArticleCategoryRepo
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

func (c *Container) GetUnitPriceRepo() repository.IUnitPriceRepository {
	return c.UnitPriceRepo
}

func (c *Container) GetAddressRepo() repository.IAddressRepository {
	return c.AddressRepo
}

func (c *Container) GetCityRepo() repository.ICityRepository {
	return c.CityRepo
}

func (c *Container) GetProvinceRepo() repository.IProvinceRepository {
	return c.ProvinceRepo
}

func (c *Container) GetPlanRepo() repository.IPlanRepository {
	return c.PlanRepo
}

func (c *Container) GetRoleRepo() repository.IRoleRepository {
	return c.RoleRepo
}

func (c *Container) GetPermissionRepo() repository.IPermissionRepository {
	return c.PermissionRepo
}

func (c *Container) GetAuthContextTransientService() func(c *gin.Context) common.IAuthContextService {
	return c.AuthContextTransientService
}

func (c *Container) GetMainCache() cache.ICacheService {
	return c.MainCache
}

func (c *Container) GetStockCacheTransient() cache.ICacheService {
	return c.stockCacheTransient()
}

func (c *Container) GetLogger() sflogger.Logger {
	return c.Logger
}

func (c *Container) GetDB() *gorm.DB {
	return sform.MustDB("main")
}

func (c *Container) GetPaymentService() service.IPaymentService {
	return c.PaymentService
}
