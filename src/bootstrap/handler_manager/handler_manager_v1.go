package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	addressusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/address"
	articleusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article"
	articlecategoryusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/article_category"
	basketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/basket"
	customerusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/customer"
	customerticketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/customer_ticket"
	defaultthemeusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/default_theme"
	discountusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/discount"
	fileitemusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/file_item"
	headerfooterusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/header_footer"
	orderusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/order"
	pageusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/page"
	paymentusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/payment"
	planusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/plan"
	productusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product"
	productcategoryusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product_category"
	productreviewusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product_review"
	roleusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/role"
	siteusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/site"
	ticketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/ticket"
	unitpriceusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/unit_price"
	userusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/user"
	websiteusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/website"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func AddressInit(c contract.IContainer) *v1.AddressHandler {
	use := addressusecase.NewAddressUsecase(c)
	handler := v1.NewAddressHandler(use)

	return handler
}

func ArticleCategoryInit(c contract.IContainer) *v1.ArticleCategoryHandler {
	use := articlecategoryusecase.NewArticleCategoryUsecase(c)
	handler := v1.NewBlogCategoryHandler(use)

	return handler
}

func ArticleInit(c contract.IContainer) *v1.ArticleHandler {
	use := articleusecase.NewArticleUsecase(c)
	productList := v1.NewArticleHandler(use)

	return productList
}

func BasketInit(c contract.IContainer) *v1.BasketHandler {
	use := basketusecase.NewBasketUsecase(c)
	handler := v1.NewBasketHandler(use)

	return handler
}

func CustomerInit(c contract.IContainer) *v1.CustomerHandler {
	use := customerusecase.NewCustomerUsecase(c)
	handler := v1.NewCustomerHandler(use)

	return handler
}

func CustomerTicketInit(c contract.IContainer) *v1.CustomerTicketHandler {
	use := customerticketusecase.NewCustomerTicketUsecase(c)
	handler := v1.NewCustomerTicketHandler(use)

	return handler
}

func DefaultThemeInit(c contract.IContainer) *v1.DefaultThemeHandler {
	use := defaultthemeusecase.NewDefaultThemeUsecase(c)
	handler := v1.NewDefaultThemeHandler(use)

	return handler
}

func DiscountInit(c contract.IContainer) *v1.DiscountHandler {
	use := discountusecase.NewDiscountUsecase(c)
	handler := v1.NewDiscountHandler(use)

	return handler
}

func FileItemInit(c contract.IContainer) *v1.FileItemHandler {
	use := fileitemusecase.NewFileItemUsecase(c)
	handler := v1.NewFileItemHandler(use)

	return handler
}

func HeaderFooterInit(c contract.IContainer) *v1.HeaderFooterHandler {
	use := headerfooterusecase.NewHeaderFooterUsecase(c)
	handler := v1.NewHeaderFooterHandler(use)

	return handler
}

func OrderInit(c contract.IContainer) *v1.OrderHandler {
	use := orderusecase.NewOrderUsecase(c)
	handler := v1.NewOrderHandler(use)

	return handler
}

func PageInit(c contract.IContainer) *v1.PageHandler {
	use := pageusecase.NewPageUsecase(c)
	handler := v1.NewPageHandler(use)

	return handler
}

func PaymentInit(c contract.IContainer) *v1.PaymentHandler {
	use := paymentusecase.NewPaymentUsecase(c)
	handler := v1.NewPaymentHandler(use)

	return handler
}

func PlanInit(c contract.IContainer) *v1.PlanHandler {
	use := planusecase.NewPlanUsecase(c)
	handler := v1.NewPlanHandler(use)

	return handler
}

func ProductCategoryInit(c contract.IContainer) *v1.ProductCategoryHandler {
	use := productcategoryusecase.NewProductCategoryUsecase(c)
	handler := v1.NewProductCategoryHandler(use)

	return handler
}

func ProductInit(c contract.IContainer) *v1.ProductHandler {
	use := productusecase.NewProductUsecase(c)
	handler := v1.NewProductHandler(use)

	return handler
}

func ProductReviewInit(c contract.IContainer) *v1.ProductReviewHandler {
	use := productreviewusecase.NewProductReviewUsecase(c)
	handler := v1.NewProductReviewHandler(use)

	return handler
}

func RoleInit(c contract.IContainer) *v1.RoleHandler {
	use := roleusecase.NewRoleUsecase(c)
	handler := v1.NewRoleHandler(use)

	return handler
}

func SiteInit(c contract.IContainer) *v1.SiteHandler {
	use := siteusecase.NewSiteUsecase(c)
	handler := v1.NewSiteHandler(use)

	return handler
}

func TicketInit(c contract.IContainer) *v1.TicketHandler {
	use := ticketusecase.NewTicketUsecase(c)
	handler := v1.NewTicketHandler(use)

	return handler
}
func UnitPriceInit(c contract.IContainer) *v1.UnitPriceHandler {
	use := unitpriceusecase.NewUnitPriceUsecase(c)
	handler := v1.NewUnitPriceHandler(use)

	return handler
}

func UserInit(c contract.IContainer) *v1.UserHandler {
	use := userusecase.NewUserUsecase(c)
	handler := v1.NewUserHandler(use)

	return handler
}

func WebsiteInit(c contract.IContainer) *v1.WebsiteHandler {
	use := websiteusecase.NewWebsiteUsecase(c)
	handler := v1.NewWebsiteHandler(use)

	return handler
}
