package bootstrap

import (
	"context"

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

	// Create auth context service
	authContextService := auth.NewAuthContextService(ctx, mysql.NewUserRepository(mainDB), mysql.NewRoleRepository(mainDB))

	// Create identity service
	identityService := auth.NewIdentityService("your-jwt-secret", 24*time.Hour)

	// Get storage config
	storageConfig := cfg.Storage()

	return &Container{
		Config: cfg,
		Logger: logger,

		MainCache:          service.NewRedis(sfredis.MustClient(ctx, "cache")),
		AuthContextService: authContextService,
		IdentityService:    identityService,
		StorageService:     storage.NewStorageService(storageConfig.Bucket, storageConfig.Region, storageConfig.AccessKey, storageConfig.SecretKey),

		// for transient
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
		RoleRepo:                  mysql.NewRoleRepository(mainDB),
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
