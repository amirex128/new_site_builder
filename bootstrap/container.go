package bootstrap

import (
	"context"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	service2 "github.com/amirex128/new_site_builder/internal/contract/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Container
type Container struct {
	Ctx                  context.Context
	Config               contract.IConfig
	MemoryLoader         service2.IMemoryLoaderService
	MainCache            service2.ICacheService
	AuthTransientService func(c *gin.Context) service2.IAuthService
	IdentityService      service2.IIdentityService
	StorageService       service2.IStorageService
	stockCacheTransient  func() service2.ICacheService
	PaymentService       service2.IPaymentService
	MessageService       service2.IMessageService

	Logger                    sflogger.Logger
	ArticleRepo               repository2.IArticleRepository
	BasketRepo                repository2.IBasketRepository
	BasketItemRepo            repository2.IBasketItemRepository
	ArticleCategoryRepo       repository2.IArticleCategoryRepository
	CreditRepo                repository2.ICreditRepository
	CouponRepo                repository2.ICouponRepository
	CustomerRepo              repository2.ICustomerRepository
	DefaultThemeRepo          repository2.IDefaultThemeRepository
	DiscountRepo              repository2.IDiscountRepository
	FileItemRepo              repository2.IFileItemRepository
	GatewayRepo               repository2.IGatewayRepository
	HeaderFooterRepo          repository2.IHeaderFooterRepository
	MediaRepo                 repository2.IMediaRepository
	OrderRepo                 repository2.IOrderRepository
	OrderItemRepo             repository2.IOrderItemRepository
	PageRepo                  repository2.IPageRepository
	ParbadPaymentRepo         repository2.IParbadPaymentRepository
	ParbadTransactionRepo     repository2.IParbadTransactionRepository
	PaymentRepo               repository2.IPaymentRepository
	ProductRepo               repository2.IProductRepository
	ProductAttributeRepo      repository2.IProductAttributeRepository
	ProductCategoryRepo       repository2.IProductCategoryRepository
	ProductReviewRepo         repository2.IProductReviewRepository
	ProductVariantRepo        repository2.IProductVariantRepository
	ReturnItemRepo            repository2.IReturnItemRepository
	SettingRepo               repository2.ISettingRepository
	SiteRepo                  repository2.ISiteRepository
	StorageRepo               repository2.IStorageRepository
	TicketRepo                repository2.ITicketRepository
	CustomerTicketRepo        repository2.ICustomerTicketRepository
	UserRepo                  repository2.IUserRepository
	UnitPriceRepo             repository2.IUnitPriceRepository
	AddressRepo               repository2.IAddressRepository
	CityRepo                  repository2.ICityRepository
	ProvinceRepo              repository2.IProvinceRepository
	PlanRepo                  repository2.IPlanRepository
	RoleRepo                  repository2.IRoleRepository
	PermissionRepo            repository2.IPermissionRepository
	CommentRepo               repository2.ICommentRepository
	CustomerCommentRepo       repository2.ICustomerCommentRepository
	TicketMediaRepo           repository2.ITicketMediaRepository
	CustomerTicketMediaRepo   repository2.ICustomerTicketMediaRepository
	PageArticleUsageRepo      repository2.IPageArticleUsageRepository
	PageProductUsageRepo      repository2.IPageProductUsageRepository
	PageHeaderFooterUsageRepo repository2.IPageHeaderFooterUsageRepository
	ConfigurationRepo         repository2.IConfigurationRepository
}

func (c *Container) GetCtx() context.Context {
	return c.Ctx
}
func (c *Container) GetCommentRepo() repository2.ICommentRepository {
	return c.CommentRepo
}

func (c *Container) GetCustomerCommentRepo() repository2.ICustomerCommentRepository {
	return c.CustomerCommentRepo
}

func (c *Container) GetTicketMediaRepo() repository2.ITicketMediaRepository {
	return c.TicketMediaRepo
}

func (c *Container) GetCustomerTicketMediaRepo() repository2.ICustomerTicketMediaRepository {
	return c.CustomerTicketMediaRepo
}

func (c *Container) GetPageArticleUsageRepo() repository2.IPageArticleUsageRepository {
	return c.PageArticleUsageRepo
}

func (c *Container) GetPageProductUsageRepo() repository2.IPageProductUsageRepository {
	return c.PageProductUsageRepo
}

func (c *Container) GetPageHeaderFooterUsageRepo() repository2.IPageHeaderFooterUsageRepository {
	return c.PageHeaderFooterUsageRepo
}

func (c *Container) GetStorageService() service2.IStorageService {
	return c.StorageService
}

func (c *Container) GetIdentityService() service2.IIdentityService {
	return c.IdentityService
}

func (c *Container) GetConfig() contract.IConfig {
	return c.Config
}

func (c *Container) GetArticleRepo() repository2.IArticleRepository {
	return c.ArticleRepo
}

func (c *Container) GetBasketRepo() repository2.IBasketRepository {
	return c.BasketRepo
}

func (c *Container) GetBasketItemRepo() repository2.IBasketItemRepository {
	return c.BasketItemRepo
}

func (c *Container) GetArticleCategoryRepo() repository2.IArticleCategoryRepository {
	return c.ArticleCategoryRepo
}

func (c *Container) GetCreditRepo() repository2.ICreditRepository {
	return c.CreditRepo
}

func (c *Container) GetCouponRepo() repository2.ICouponRepository {
	return c.CouponRepo
}

func (c *Container) GetCustomerRepo() repository2.ICustomerRepository {
	return c.CustomerRepo
}

func (c *Container) GetDefaultThemeRepo() repository2.IDefaultThemeRepository {
	return c.DefaultThemeRepo
}

func (c *Container) GetDiscountRepo() repository2.IDiscountRepository {
	return c.DiscountRepo
}

func (c *Container) GetFileItemRepo() repository2.IFileItemRepository {
	return c.FileItemRepo
}

func (c *Container) GetGatewayRepo() repository2.IGatewayRepository {
	return c.GatewayRepo
}

func (c *Container) GetHeaderFooterRepo() repository2.IHeaderFooterRepository {
	return c.HeaderFooterRepo
}

func (c *Container) GetMediaRepo() repository2.IMediaRepository {
	return c.MediaRepo
}

func (c *Container) GetOrderRepo() repository2.IOrderRepository {
	return c.OrderRepo
}

func (c *Container) GetOrderItemRepo() repository2.IOrderItemRepository {
	return c.OrderItemRepo
}

func (c *Container) GetPageRepo() repository2.IPageRepository {
	return c.PageRepo
}

func (c *Container) GetParbadPaymentRepo() repository2.IParbadPaymentRepository {
	return c.ParbadPaymentRepo
}

func (c *Container) GetParbadTransactionRepo() repository2.IParbadTransactionRepository {
	return c.ParbadTransactionRepo
}

func (c *Container) GetPaymentRepo() repository2.IPaymentRepository {
	return c.PaymentRepo
}

func (c *Container) GetProductRepo() repository2.IProductRepository {
	return c.ProductRepo
}

func (c *Container) GetProductAttributeRepo() repository2.IProductAttributeRepository {
	return c.ProductAttributeRepo
}

func (c *Container) GetProductCategoryRepo() repository2.IProductCategoryRepository {
	return c.ProductCategoryRepo
}

func (c *Container) GetProductReviewRepo() repository2.IProductReviewRepository {
	return c.ProductReviewRepo
}

func (c *Container) GetProductVariantRepo() repository2.IProductVariantRepository {
	return c.ProductVariantRepo
}

func (c *Container) GetReturnItemRepo() repository2.IReturnItemRepository {
	return c.ReturnItemRepo
}

func (c *Container) GetSettingRepo() repository2.ISettingRepository {
	return c.SettingRepo
}

func (c *Container) GetSiteRepo() repository2.ISiteRepository {
	return c.SiteRepo
}

func (c *Container) GetStorageRepo() repository2.IStorageRepository {
	return c.StorageRepo
}

func (c *Container) GetTicketRepo() repository2.ITicketRepository {
	return c.TicketRepo
}

func (c *Container) GetCustomerTicketRepo() repository2.ICustomerTicketRepository {
	return c.CustomerTicketRepo
}

func (c *Container) GetUserRepo() repository2.IUserRepository {
	return c.UserRepo
}

func (c *Container) GetUnitPriceRepo() repository2.IUnitPriceRepository {
	return c.UnitPriceRepo
}

func (c *Container) GetAddressRepo() repository2.IAddressRepository {
	return c.AddressRepo
}

func (c *Container) GetCityRepo() repository2.ICityRepository {
	return c.CityRepo
}

func (c *Container) GetProvinceRepo() repository2.IProvinceRepository {
	return c.ProvinceRepo
}

func (c *Container) GetPlanRepo() repository2.IPlanRepository {
	return c.PlanRepo
}

func (c *Container) GetRoleRepo() repository2.IRoleRepository {
	return c.RoleRepo
}

func (c *Container) GetPermissionRepo() repository2.IPermissionRepository {
	return c.PermissionRepo
}

func (c *Container) GetAuthTransientService() func(c *gin.Context) service2.IAuthService {
	return c.AuthTransientService
}

func (c *Container) GetMainCache() service2.ICacheService {
	return c.MainCache
}

func (c *Container) GetStockCacheTransient() service2.ICacheService {
	return c.stockCacheTransient()
}

func (c *Container) GetLogger() sflogger.Logger {
	return c.Logger
}

func (c *Container) GetDB() *gorm.DB {
	return sform.MustDB("main")
}

func (c *Container) GetPaymentService() service2.IPaymentService {
	return c.PaymentService
}

func (c *Container) GetMessageService() service2.IMessageService {
	return c.MessageService
}

func (c *Container) GetConfigurationRepo() repository2.IConfigurationRepository {
	return c.ConfigurationRepo
}
