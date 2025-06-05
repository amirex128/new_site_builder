package handlermanager

import (
	v2 "github.com/amirex128/new_site_builder/internal/api/handler/http/v1"
	"github.com/amirex128/new_site_builder/internal/application/usecase/address"
	"github.com/amirex128/new_site_builder/internal/application/usecase/article"
	"github.com/amirex128/new_site_builder/internal/application/usecase/article_category"
	"github.com/amirex128/new_site_builder/internal/application/usecase/basket"
	"github.com/amirex128/new_site_builder/internal/application/usecase/customer"
	"github.com/amirex128/new_site_builder/internal/application/usecase/customer_ticket"
	"github.com/amirex128/new_site_builder/internal/application/usecase/default_theme"
	"github.com/amirex128/new_site_builder/internal/application/usecase/discount"
	"github.com/amirex128/new_site_builder/internal/application/usecase/file_item"
	"github.com/amirex128/new_site_builder/internal/application/usecase/header_footer"
	"github.com/amirex128/new_site_builder/internal/application/usecase/order"
	"github.com/amirex128/new_site_builder/internal/application/usecase/page"
	"github.com/amirex128/new_site_builder/internal/application/usecase/payment"
	"github.com/amirex128/new_site_builder/internal/application/usecase/plan"
	"github.com/amirex128/new_site_builder/internal/application/usecase/product"
	"github.com/amirex128/new_site_builder/internal/application/usecase/product_category"
	"github.com/amirex128/new_site_builder/internal/application/usecase/product_review"
	"github.com/amirex128/new_site_builder/internal/application/usecase/role"
	"github.com/amirex128/new_site_builder/internal/application/usecase/site"
	"github.com/amirex128/new_site_builder/internal/application/usecase/ticket"
	"github.com/amirex128/new_site_builder/internal/application/usecase/unit_price"
	"github.com/amirex128/new_site_builder/internal/application/usecase/user"
	"github.com/amirex128/new_site_builder/internal/application/usecase/website"
	"github.com/amirex128/new_site_builder/internal/contract"
)

func AddressInit(c contract.IContainer) *v2.AddressHandler {
	use := addressusecase.NewAddressUsecase(c)
	handler := v2.NewAddressHandler(use)

	return handler
}

func ArticleCategoryInit(c contract.IContainer) *v2.ArticleCategoryHandler {
	use := articlecategoryusecase.NewArticleCategoryUsecase(c)
	handler := v2.NewBlogCategoryHandler(use)

	return handler
}

func ArticleInit(c contract.IContainer) *v2.ArticleHandler {
	use := articleusecase.NewArticleUsecase(c)
	productList := v2.NewArticleHandler(use)

	return productList
}

func BasketInit(c contract.IContainer) *v2.BasketHandler {
	use := basketusecase.NewBasketUsecase(c)
	handler := v2.NewBasketHandler(use)

	return handler
}

func CustomerInit(c contract.IContainer) *v2.CustomerHandler {
	use := customerusecase.NewCustomerUsecase(c)
	handler := v2.NewCustomerHandler(use)

	return handler
}

func CustomerTicketInit(c contract.IContainer) *v2.CustomerTicketHandler {
	use := customerticketusecase.NewCustomerTicketUsecase(c)
	handler := v2.NewCustomerTicketHandler(use)

	return handler
}

func DefaultThemeInit(c contract.IContainer) *v2.DefaultThemeHandler {
	use := defaultthemeusecase.NewDefaultThemeUsecase(c)
	handler := v2.NewDefaultThemeHandler(use)

	return handler
}

func DiscountInit(c contract.IContainer) *v2.DiscountHandler {
	use := discountusecase.NewDiscountUsecase(c)
	handler := v2.NewDiscountHandler(use)

	return handler
}

func FileItemInit(c contract.IContainer) *v2.FileItemHandler {
	use := fileitemusecase.NewFileItemUsecase(c)
	handler := v2.NewFileItemHandler(use)

	return handler
}

func HeaderFooterInit(c contract.IContainer) *v2.HeaderFooterHandler {
	use := headerfooterusecase.NewHeaderFooterUsecase(c)
	handler := v2.NewHeaderFooterHandler(use)

	return handler
}

func OrderInit(c contract.IContainer) *v2.OrderHandler {
	use := orderusecase.NewOrderUsecase(c)
	handler := v2.NewOrderHandler(use)

	return handler
}

func PageInit(c contract.IContainer) *v2.PageHandler {
	use := pageusecase.NewPageUsecase(c)
	handler := v2.NewPageHandler(use)

	return handler
}

func PaymentInit(c contract.IContainer) *v2.PaymentHandler {
	use := paymentusecase.NewPaymentUsecase(c)
	handler := v2.NewPaymentHandler(use)

	return handler
}

func PlanInit(c contract.IContainer) *v2.PlanHandler {
	use := planusecase.NewPlanUsecase(c)
	handler := v2.NewPlanHandler(use)

	return handler
}

func ProductCategoryInit(c contract.IContainer) *v2.ProductCategoryHandler {
	use := productcategoryusecase.NewProductCategoryUsecase(c)
	handler := v2.NewProductCategoryHandler(use)

	return handler
}

func ProductInit(c contract.IContainer) *v2.ProductHandler {
	use := productusecase.NewProductUsecase(c)
	handler := v2.NewProductHandler(use)

	return handler
}

func ProductReviewInit(c contract.IContainer) *v2.ProductReviewHandler {
	use := productreviewusecase.NewProductReviewUsecase(c)
	handler := v2.NewProductReviewHandler(use)

	return handler
}

func RoleInit(c contract.IContainer) *v2.RoleHandler {
	use := roleusecase.NewRoleUsecase(c)
	handler := v2.NewRoleHandler(use)

	return handler
}

func SiteInit(c contract.IContainer) *v2.SiteHandler {
	use := siteusecase.NewSiteUsecase(c)
	handler := v2.NewSiteHandler(use)

	return handler
}

func TicketInit(c contract.IContainer) *v2.TicketHandler {
	use := ticketusecase.NewTicketUsecase(c)
	handler := v2.NewTicketHandler(use)

	return handler
}
func UnitPriceInit(c contract.IContainer) *v2.UnitPriceHandler {
	use := unitpriceusecase.NewUnitPriceUsecase(c)
	handler := v2.NewUnitPriceHandler(use)

	return handler
}

func UserInit(c contract.IContainer) *v2.UserHandler {
	use := userusecase.NewUserUsecase(c)
	handler := v2.NewUserHandler(use)

	return handler
}

func WebsiteInit(c contract.IContainer) *v2.WebsiteHandler {
	use := websiteusecase.NewWebsiteUsecase(c)
	handler := v2.NewWebsiteHandler(use)

	return handler
}
