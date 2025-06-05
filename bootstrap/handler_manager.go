package bootstrap

import (
	"github.com/amirex128/new_site_builder/bootstrap/handler_manager"
	v2 "github.com/amirex128/new_site_builder/internal/api/handler/http/v1"
	"github.com/amirex128/new_site_builder/internal/contract"
)

type HandlerManager struct {
	ArticleHandlerV1         *v2.ArticleHandler
	ArticleCategoryHandlerV1 *v2.ArticleCategoryHandler
	AddressHandlerV1         *v2.AddressHandler
	BasketHandlerV1          *v2.BasketHandler
	CustomerHandlerV1        *v2.CustomerHandler
	CustomerTicketHandlerV1  *v2.CustomerTicketHandler
	DefaultThemeHandlerV1    *v2.DefaultThemeHandler
	DiscountHandlerV1        *v2.DiscountHandler
	FileItemHandlerV1        *v2.FileItemHandler
	HeaderFooterHandlerV1    *v2.HeaderFooterHandler
	OrderHandlerV1           *v2.OrderHandler
	PageHandlerV1            *v2.PageHandler
	PaymentHandlerV1         *v2.PaymentHandler
	PlanHandlerV1            *v2.PlanHandler
	ProductHandlerV1         *v2.ProductHandler
	ProductCategoryHandlerV1 *v2.ProductCategoryHandler
	ProductReviewHandlerV1   *v2.ProductReviewHandler
	RoleHandlerV1            *v2.RoleHandler
	SiteHandlerV1            *v2.SiteHandler
	TicketHandlerV1          *v2.TicketHandler
	UnitPriceHandlerV1       *v2.UnitPriceHandler
	UserHandlerV1            *v2.UserHandler
	WebsiteHandlerV1         *v2.WebsiteHandler
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
