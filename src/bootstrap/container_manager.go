package bootstrap

import (
	"context"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"github.com/amirex128/new_site_builder/src/internal/infra/repository"
	"github.com/amirex128/new_site_builder/src/internal/infra/service/auth"
	"github.com/amirex128/new_site_builder/src/internal/infra/service/cache"
	"github.com/amirex128/new_site_builder/src/internal/infra/service/payment"
	"github.com/amirex128/new_site_builder/src/internal/infra/service/storage"
	"github.com/gin-gonic/gin"

	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
	"github.com/amirex128/new_site_builder/src/config"
)

func ContainerProvider(ctx context.Context, cfg *config.Config, logger sflogger.Logger) *Container {
	mainDB := sform.MustDB("main")

	paymentRepo := repository.NewPaymentRepository(mainDB)
	gatewayRepo := repository.NewGatewayRepository(mainDB)

	identityService := auth.NewIdentityService(cfg.JwtSecretToken, cfg.JwtIssuer, cfg.JwtAudience, 24*time.Hour)
	return &Container{
		Config: cfg,
		Logger: logger,

		MainCache:       cache.NewRedis(sfredis.MustClient(ctx, "cache")),
		IdentityService: identityService,
		StorageService: storage.NewStorageService(
			storage.NewStorageClient(cfg.StorageS1Host, cfg.StorageS1AccessKey, cfg.StorageS1SecretKey),
			storage.NewStorageClient(cfg.StorageS2Host, cfg.StorageS2AccessKey, cfg.StorageS2SecretKey),
			storage.NewStorageClient(cfg.StorageS3Host, cfg.StorageS3AccessKey, cfg.StorageS3SecretKey),
		),
		PaymentService: payment.NewPaymentService(paymentRepo, gatewayRepo),
		// for transient
		AuthTransientService: func(c *gin.Context) service.IAuthService {
			return auth.NewAuthContextService(c, identityService)
		},
		stockCacheTransient: func() service.ICacheService {
			return cache.NewRedis(sfredis.MustClient(ctx, "stock"))
		},

		// Repositories
		ArticleRepo:               repository.NewArticleRepository(mainDB),
		BasketRepo:                repository.NewBasketRepository(mainDB),
		BasketItemRepo:            repository.NewBasketItemRepository(mainDB),
		ArticleCategoryRepo:       repository.NewArticleCategoryRepository(mainDB),
		CustomerRepo:              repository.NewCustomerRepository(mainDB),
		DiscountRepo:              repository.NewDiscountRepository(mainDB),
		HeaderFooterRepo:          repository.NewHeaderFooterRepository(mainDB),
		MediaRepo:                 repository.NewMediaRepository(mainDB),
		OrderRepo:                 repository.NewOrderRepository(mainDB),
		OrderItemRepo:             repository.NewOrderItemRepository(mainDB),
		PageRepo:                  repository.NewPageRepository(mainDB),
		PaymentRepo:               repository.NewPaymentRepository(mainDB),
		ProductRepo:               repository.NewProductRepository(mainDB),
		ProductCategoryRepo:       repository.NewProductCategoryRepository(mainDB),
		ProductReviewRepo:         repository.NewProductReviewRepository(mainDB),
		ProductVariantRepo:        repository.NewProductVariantRepository(mainDB),
		SettingRepo:               repository.NewSettingRepository(mainDB),
		SiteRepo:                  repository.NewSiteRepository(mainDB),
		TicketRepo:                repository.NewTicketRepository(mainDB),
		CustomerTicketRepo:        repository.NewCustomerTicketRepository(mainDB),
		UserRepo:                  repository.NewUserRepository(mainDB),
		UnitPriceRepo:             repository.NewUnitPriceRepository(mainDB),
		AddressRepo:               repository.NewAddressRepository(mainDB),
		CityRepo:                  repository.NewCityRepository(mainDB),
		ProvinceRepo:              repository.NewProvinceRepository(mainDB),
		PlanRepo:                  repository.NewPlanRepository(mainDB),
		RoleRepo:                  repository.NewRoleRepository(mainDB),
		PermissionRepo:            repository.NewPermissionRepository(mainDB),
		CommentRepo:               repository.NewCommentRepository(mainDB),
		CustomerCommentRepo:       repository.NewCustomerCommentRepository(mainDB),
		TicketMediaRepo:           repository.NewTicketMediaRepository(mainDB),
		CustomerTicketMediaRepo:   repository.NewCustomerTicketMediaRepository(mainDB),
		PageArticleUsageRepo:      repository.NewPageArticleUsageRepository(mainDB),
		PageProductUsageRepo:      repository.NewPageProductUsageRepository(mainDB),
		PageHeaderFooterUsageRepo: repository.NewPageHeaderFooterUsageRepository(mainDB),
	}
}
