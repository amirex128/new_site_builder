package bootstrap

import (
	"context"

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

	return &Container{
		Config: cfg,
		Logger: logger,

		//todo: create name constant
		FoodPartyCache: service.NewRedis(sfredis.MustClient(ctx, "foodparty")),

		// for transient
		stockCacheTransient: func() cache.ICacheService {
			return service.NewRedis(sfredis.MustClient(ctx, "stock"))
		},

		// Repositories
		ArticleRepo:           mysql.NewArticleRepository(mainDB),
		BasketRepo:            mysql.NewBasketRepository(mainDB),
		ArticleCategoryRepo:   mysql.NewArticleCategoryRepository(mainDB),
		CustomerRepo:          mysql.NewCustomerRepository(mainDB),
		DiscountRepo:          mysql.NewDiscountRepository(mainDB),
		HeaderFooterRepo:      mysql.NewHeaderFooterRepository(mainDB),
		MediaRepo:             mysql.NewMediaRepository(mainDB),
		OrderRepo:             mysql.NewOrderRepository(mainDB),
		PageRepo:              mysql.NewPageRepository(mainDB),
		PaymentRepo:           mysql.NewPaymentRepository(mainDB),
		ProductRepo:           mysql.NewProductRepository(mainDB),
		ProductCategoryRepo:   mysql.NewProductCategoryRepository(mainDB),
		ProductReviewRepo:     mysql.NewProductReviewRepository(mainDB),
		SettingRepo:           mysql.NewSettingRepository(mainDB),
		SiteRepo:              mysql.NewSiteRepository(mainDB),
		TicketRepo:            mysql.NewTicketRepository(mainDB),
		CustomerTicketRepo:    mysql.NewCustomerTicketRepository(mainDB),
		UserRepo:              mysql.NewUserRepository(mainDB),
		BasketItemRepo:        mysql.NewBasketItemRepository(mainDB),
		CreditRepo:            mysql.NewCreditRepository(mainDB),
		CouponRepo:            mysql.NewCouponRepository(mainDB),
		DefaultThemeRepo:      mysql.NewDefaultThemeRepository(mainDB),
		FileItemRepo:          mysql.NewFileItemRepository(mainDB),
		GatewayRepo:           mysql.NewGatewayRepository(mainDB),
		OrderItemRepo:         mysql.NewOrderItemRepository(mainDB),
		ParbadPaymentRepo:     mysql.NewParbadPaymentRepository(mainDB),
		ParbadTransactionRepo: mysql.NewParbadTransactionRepository(mainDB),
		ProductAttributeRepo:  mysql.NewProductAttributeRepository(mainDB),
		ProductVariantRepo:    mysql.NewProductVariantRepository(mainDB),
		ReturnItemRepo:        mysql.NewReturnItemRepository(mainDB),
		StorageRepo:           mysql.NewStorageRepository(mainDB),

		// New repositories
		UnitPriceRepo:  mysql.NewUnitPriceRepository(mainDB),
		AddressRepo:    mysql.NewAddressRepository(mainDB),
		CityRepo:       mysql.NewCityRepository(mainDB),
		ProvinceRepo:   mysql.NewProvinceRepository(mainDB),
		PlanRepo:       mysql.NewPlanRepository(mainDB),
		RoleRepo:       mysql.NewRoleRepository(mainDB),
		PermissionRepo: mysql.NewPermissionRepository(mainDB),
	}
}
