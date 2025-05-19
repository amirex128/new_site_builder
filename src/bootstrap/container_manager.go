package bootstrap

import (
	"context"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/infra/service/payment"
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/infra/service/auth"
	"github.com/amirex128/new_site_builder/src/internal/infra/service/storage"

	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
	"github.com/amirex128/new_site_builder/src/config"
	"github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
	"github.com/amirex128/new_site_builder/src/internal/infra/repository/mysql"
	"github.com/amirex128/new_site_builder/src/internal/infra/service"
)

func ContainerProvider(ctx context.Context, cfg *config.Config, logger sflogger.Logger) *Container {
	mainDB := sform.MustDB("main")

	userRepo := mysql.NewUserRepository(mainDB)
	roleRepo := mysql.NewRoleRepository(mainDB)
	paymentRepo := mysql.NewPaymentRepository(mainDB)
	gatewayRepo := mysql.NewGatewayRepository(mainDB)

	return &Container{
		Config: cfg,
		Logger: logger,

		MainCache:       service.NewRedis(sfredis.MustClient(ctx, "cache")),
		IdentityService: auth.NewIdentityService(cfg.JwtSecretToken, cfg.JwtIssuer, cfg.JwtAudience, 24*time.Hour),
		StorageService: storage.NewStorageService(
			storage.NewStorageClient(cfg.StorageS1Host, cfg.StorageS1AccessKey, cfg.StorageS1SecretKey),
			storage.NewStorageClient(cfg.StorageS2Host, cfg.StorageS2AccessKey, cfg.StorageS2SecretKey),
			storage.NewStorageClient(cfg.StorageS3Host, cfg.StorageS3AccessKey, cfg.StorageS3SecretKey),
		),
		PaymentService: payment.NewPaymentService(paymentRepo, gatewayRepo),
		// for transient
		AuthContextTransientService: func(c *gin.Context) common.IAuthContextService {
			return auth.NewAuthContextService(ctx, request, userRepo, roleRepo)
		},
		stockCacheTransient: func() cache.ICacheService {
			return service.NewRedis(sfredis.MustClient(ctx, "stock"))
		},

		// Repositories
		ArticleRepo:               mysql.NewArticleRepository(mainDB),
		BasketRepo:                mysql.NewBasketRepository(mainDB),
		BasketItemRepo:            mysql.NewBasketItemRepository(mainDB),
		ArticleCategoryRepo:       mysql.NewArticleCategoryRepository(mainDB),
		CustomerRepo:              mysql.NewCustomerRepository(mainDB),
		DiscountRepo:              mysql.NewDiscountRepository(mainDB),
		HeaderFooterRepo:          mysql.NewHeaderFooterRepository(mainDB),
		MediaRepo:                 mysql.NewMediaRepository(mainDB),
		OrderRepo:                 mysql.NewOrderRepository(mainDB),
		OrderItemRepo:             mysql.NewOrderItemRepository(mainDB),
		PageRepo:                  mysql.NewPageRepository(mainDB),
		PaymentRepo:               mysql.NewPaymentRepository(mainDB),
		ProductRepo:               mysql.NewProductRepository(mainDB),
		ProductCategoryRepo:       mysql.NewProductCategoryRepository(mainDB),
		ProductReviewRepo:         mysql.NewProductReviewRepository(mainDB),
		ProductVariantRepo:        mysql.NewProductVariantRepository(mainDB),
		SettingRepo:               mysql.NewSettingRepository(mainDB),
		SiteRepo:                  mysql.NewSiteRepository(mainDB),
		TicketRepo:                mysql.NewTicketRepository(mainDB),
		CustomerTicketRepo:        mysql.NewCustomerTicketRepository(mainDB),
		UserRepo:                  mysql.NewUserRepository(mainDB),
		UnitPriceRepo:             mysql.NewUnitPriceRepository(mainDB),
		AddressRepo:               mysql.NewAddressRepository(mainDB),
		CityRepo:                  mysql.NewCityRepository(mainDB),
		ProvinceRepo:              mysql.NewProvinceRepository(mainDB),
		PlanRepo:                  mysql.NewPlanRepository(mainDB),
		RoleRepo:                  roleRepo,
		PermissionRepo:            mysql.NewPermissionRepository(mainDB),
		CommentRepo:               mysql.NewCommentRepository(mainDB),
		CustomerCommentRepo:       mysql.NewCustomerCommentRepository(mainDB),
		TicketMediaRepo:           mysql.NewTicketMediaRepository(mainDB),
		CustomerTicketMediaRepo:   mysql.NewCustomerTicketMediaRepository(mainDB),
		PageArticleUsageRepo:      mysql.NewPageArticleUsageRepository(mainDB),
		PageProductUsageRepo:      mysql.NewPageProductUsageRepository(mainDB),
		PageHeaderFooterUsageRepo: mysql.NewPageHeaderFooterUsageRepository(mainDB),
	}
}
