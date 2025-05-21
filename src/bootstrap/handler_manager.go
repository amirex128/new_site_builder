package bootstrap

import (
	handlermanager "github.com/amirex128/new_site_builder/src/bootstrap/handler_manager"
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

type HandlerManager struct {
	ArticleHandlerV1         *v1.ArticleHandler
	ArticleCategoryHandlerV1 *v1.ArticleCategoryHandler
	AddressHandlerV1         *v1.AddressHandler
	BasketHandlerV1          *v1.BasketHandler
	CustomerHandlerV1        *v1.CustomerHandler
	CustomerTicketHandlerV1  *v1.CustomerTicketHandler
	DefaultThemeHandlerV1    *v1.DefaultThemeHandler
	DiscountHandlerV1        *v1.DiscountHandler
	FileItemHandlerV1        *v1.FileItemHandler
	HeaderFooterHandlerV1    *v1.HeaderFooterHandler
	OrderHandlerV1           *v1.OrderHandler
	PageHandlerV1            *v1.PageHandler
	PaymentHandlerV1         *v1.PaymentHandler
	PlanHandlerV1            *v1.PlanHandler
	ProductHandlerV1         *v1.ProductHandler
	ProductCategoryHandlerV1 *v1.ProductCategoryHandler
	ProductReviewHandlerV1   *v1.ProductReviewHandler
	RoleHandlerV1            *v1.RoleHandler
	SiteHandlerV1            *v1.SiteHandler
	TicketHandlerV1          *v1.TicketHandler
	UnitPriceHandlerV1       *v1.UnitPriceHandler
	UserHandlerV1            *v1.UserHandler
	WebsiteHandlerV1         *v1.WebsiteHandler
}

func HttpHandlerBootstrap(c contract.IContainer) *HandlerManager {

	return &HandlerManager{
		ArticleHandlerV1:         handlermanager.ArticleInit(c),
		ArticleCategoryHandlerV1: handlermanager.ArticleCategoryInit(c),
		AddressHandlerV1:         handlermanager.AddressInit(c),
		BasketHandlerV1:          handlermanager.BasketInit(c),
		CustomerHandlerV1:        handlermanager.CustomerInit(c),
		CustomerTicketHandlerV1:  handlermanager.CustomerTicketInit(c),
		DefaultThemeHandlerV1:    handlermanager.DefaultThemeInit(c),
		DiscountHandlerV1:        handlermanager.DiscountInit(c),
		FileItemHandlerV1:        handlermanager.FileItemInit(c),
		HeaderFooterHandlerV1:    handlermanager.HeaderFooterInit(c),
		OrderHandlerV1:           handlermanager.OrderInit(c),
		PageHandlerV1:            handlermanager.PageInit(c),
		PaymentHandlerV1:         handlermanager.PaymentInit(c),
		PlanHandlerV1:            handlermanager.PlanInit(c),
		ProductHandlerV1:         handlermanager.ProductInit(c),
		ProductCategoryHandlerV1: handlermanager.ProductCategoryInit(c),
		ProductReviewHandlerV1:   handlermanager.ProductReviewInit(c),
		RoleHandlerV1:            handlermanager.RoleInit(c),
		SiteHandlerV1:            handlermanager.SiteInit(c),
		TicketHandlerV1:          handlermanager.TicketInit(c),
		UnitPriceHandlerV1:       handlermanager.UnitPriceInit(c),
		UserHandlerV1:            handlermanager.UserInit(c),
		WebsiteHandlerV1:         handlermanager.WebsiteInit(c),
	}
}
