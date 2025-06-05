package bootstrap

import (
	"context"
	"github.com/amirex128/new_site_builder/config"
	service3 "github.com/amirex128/new_site_builder/internal/contract/service"
	repository2 "github.com/amirex128/new_site_builder/internal/infra/repository"
	"github.com/amirex128/new_site_builder/internal/infra/service"
	"github.com/gin-gonic/gin"

	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

func ContainerProvider(ctx context.Context, cfg *config.Config, logger sflogger.Logger) *Container {
	mainDB := sform.MustDB("main")

	paymentRepo := repository2.NewPaymentRepository(mainDB)
	gatewayRepo := repository2.NewGatewayRepository(mainDB)

	identityService := service.NewIdentityService(cfg.JwtSecretToken, cfg.JwtIssuer, cfg.JwtAudience, 24*time.Hour)
	return &Container{
		Ctx:    ctx,
		Config: cfg,
		Logger: logger,

		MainCache:       service.NewRedis(sfredis.MustClient(ctx, "cache")),
		IdentityService: identityService,
		StorageService: service.NewStorageService(
			service.NewStorageClient(cfg.StorageS1Host, cfg.StorageS1AccessKey, cfg.StorageS1SecretKey),
			service.NewStorageClient(cfg.StorageS2Host, cfg.StorageS2AccessKey, cfg.StorageS2SecretKey),
			service.NewStorageClient(cfg.StorageS3Host, cfg.StorageS3AccessKey, cfg.StorageS3SecretKey),
		),
		PaymentService: service.NewPaymentService(paymentRepo, gatewayRepo),
		MessageService: service.NewRabbitMqService(ctx, logger),
		// for transient
		AuthTransientService: func(c *gin.Context) service3.IAuthService {
			return service.NewAuthContextService(c, identityService)
		},
		stockCacheTransient: func() service3.ICacheService {
			return service.NewRedis(sfredis.MustClient(ctx, "stock"))
		},

		// Repositories
		ArticleRepo:               repository2.NewArticleRepository(mainDB),
		BasketRepo:                repository2.NewBasketRepository(mainDB),
		BasketItemRepo:            repository2.NewBasketItemRepository(mainDB),
		ArticleCategoryRepo:       repository2.NewArticleCategoryRepository(mainDB),
		CustomerRepo:              repository2.NewCustomerRepository(mainDB),
		DiscountRepo:              repository2.NewDiscountRepository(mainDB),
		HeaderFooterRepo:          repository2.NewHeaderFooterRepository(mainDB),
		MediaRepo:                 repository2.NewMediaRepository(mainDB),
		OrderRepo:                 repository2.NewOrderRepository(mainDB),
		OrderItemRepo:             repository2.NewOrderItemRepository(mainDB),
		PageRepo:                  repository2.NewPageRepository(mainDB),
		PaymentRepo:               repository2.NewPaymentRepository(mainDB),
		ProductRepo:               repository2.NewProductRepository(mainDB),
		ProductCategoryRepo:       repository2.NewProductCategoryRepository(mainDB),
		ProductReviewRepo:         repository2.NewProductReviewRepository(mainDB),
		ProductVariantRepo:        repository2.NewProductVariantRepository(mainDB),
		SettingRepo:               repository2.NewSettingRepository(mainDB),
		SiteRepo:                  repository2.NewSiteRepository(mainDB),
		TicketRepo:                repository2.NewTicketRepository(mainDB),
		CustomerTicketRepo:        repository2.NewCustomerTicketRepository(mainDB),
		UserRepo:                  repository2.NewUserRepository(mainDB),
		UnitPriceRepo:             repository2.NewUnitPriceRepository(mainDB),
		AddressRepo:               repository2.NewAddressRepository(mainDB),
		CityRepo:                  repository2.NewCityRepository(mainDB),
		ProvinceRepo:              repository2.NewProvinceRepository(mainDB),
		PlanRepo:                  repository2.NewPlanRepository(mainDB),
		RoleRepo:                  repository2.NewRoleRepository(mainDB),
		PermissionRepo:            repository2.NewPermissionRepository(mainDB),
		CommentRepo:               repository2.NewCommentRepository(mainDB),
		CustomerCommentRepo:       repository2.NewCustomerCommentRepository(mainDB),
		TicketMediaRepo:           repository2.NewTicketMediaRepository(mainDB),
		CustomerTicketMediaRepo:   repository2.NewCustomerTicketMediaRepository(mainDB),
		PageArticleUsageRepo:      repository2.NewPageArticleUsageRepository(mainDB),
		PageProductUsageRepo:      repository2.NewPageProductUsageRepository(mainDB),
		PageHeaderFooterUsageRepo: repository2.NewPageHeaderFooterUsageRepository(mainDB),
		DefaultThemeRepo:          repository2.NewDefaultThemeRepository(mainDB),
		FileItemRepo:              repository2.NewFileItemRepository(mainDB),
		GatewayRepo:               repository2.NewGatewayRepository(mainDB),
		ConfigurationRepo:         repository2.NewConfigurationRepository(mainDB),
	}
}
