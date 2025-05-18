package router

import (
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/gin-gonic/gin"
)

type RouterV1 struct {
	h *bootstrap.HandlerManager
}

func (v RouterV1) Routes(route *gin.RouterGroup) {
	// NOTE: HandlerManager in bootstrap/handler_manager.go needs to be updated with the following fields:
	// - FileItemHandlerV1
	// - BasketHandlerV1
	// - OrderHandlerV1
	// - PaymentHandlerV1
	// - ProductCategoryHandlerV1
	// - ProductHandlerV1
	// - WebsiteHandlerV1
	// - DefaultThemeHandlerV1
	// - PageHandlerV1
	// - HeaderFooterHandlerV1
	// - SiteHandlerV1
	// - TicketHandlerV1
	// - CustomerTicketHandlerV1
	// - DiscountHandlerV1
	// - ProductReviewHandlerV1
	// - AddressHandlerV1
	// - PlanHandlerV1
	// - RoleHandlerV1
	// - UnitPriceHandlerV1
	// - UserHandlerV1
	// - CustomerHandlerV1
	// - ArticleCategoryHandlerV1

	// Article routes
	articleRoute := route.Group("/article")
	articleRoute.POST("", v.h.ArticleHandlerV1.ArticleCreate)
	articleRoute.PUT("", v.h.ArticleHandlerV1.ArticleUpdate)
	articleRoute.DELETE("", v.h.ArticleHandlerV1.ArticleDelete)
	articleRoute.GET("", v.h.ArticleHandlerV1.ArticleGet)
	articleRoute.GET("/all", v.h.ArticleHandlerV1.ArticleGetAll)
	articleRoute.POST("/filters-sort", v.h.ArticleHandlerV1.ArticleGetByFiltersSort)
	articleRoute.GET("/admin/all", v.h.ArticleHandlerV1.AdminArticleGetAll)

	// Article Category routes
	articleCategoryRoute := route.Group("/article-category")
	articleCategoryRoute.POST("", v.h.ArticleCategoryHandlerV1.CategoryCreate)
	articleCategoryRoute.PUT("", v.h.ArticleCategoryHandlerV1.CategoryUpdate)
	articleCategoryRoute.DELETE("", v.h.ArticleCategoryHandlerV1.CategoryDelete)
	articleCategoryRoute.GET("", v.h.ArticleCategoryHandlerV1.CategoryGet)
	articleCategoryRoute.GET("/all", v.h.ArticleCategoryHandlerV1.CategoryGetAll)
	articleCategoryRoute.GET("/admin/all", v.h.ArticleCategoryHandlerV1.AdminCategoryGetAll)

	// Address routes
	addressRoute := route.Group("/address")
	addressRoute.POST("", v.h.AddressHandlerV1.CreateAddress)
	addressRoute.PUT("", v.h.AddressHandlerV1.UpdateAddress)
	addressRoute.DELETE("", v.h.AddressHandlerV1.DeleteAddress)
	addressRoute.GET("", v.h.AddressHandlerV1.GetByIdAddress)
	addressRoute.GET("/all", v.h.AddressHandlerV1.GetAllAddress)
	addressRoute.GET("/city/all", v.h.AddressHandlerV1.GetAllCity)
	addressRoute.GET("/province/all", v.h.AddressHandlerV1.GetAllProvince)
	addressRoute.GET("/admin/all", v.h.AddressHandlerV1.AdminGetAllAddress)

	// Plan routes
	planRoute := route.Group("/plan")
	planRoute.POST("", v.h.PlanHandlerV1.CreatePlan)
	planRoute.PUT("", v.h.PlanHandlerV1.UpdatePlan)
	planRoute.DELETE("", v.h.PlanHandlerV1.DeletePlan)
	planRoute.GET("", v.h.PlanHandlerV1.GetByIdPlan)
	planRoute.GET("/all", v.h.PlanHandlerV1.GetAllPlan)
	planRoute.GET("/calculate", v.h.PlanHandlerV1.CalculatePlanPrice)

	// Role routes
	roleRoute := route.Group("/role")
	roleRoute.POST("", v.h.RoleHandlerV1.CreateRole)
	roleRoute.PUT("", v.h.RoleHandlerV1.UpdateRole)
	roleRoute.PUT("/customer", v.h.RoleHandlerV1.SetRoleToCustomer)
	roleRoute.PUT("/user", v.h.RoleHandlerV1.SetRoleToUser)
	roleRoute.PUT("/plan", v.h.RoleHandlerV1.SetRoleToPlan)
	roleRoute.GET("/permission/all", v.h.RoleHandlerV1.GetAllPermission)
	roleRoute.GET("/all", v.h.RoleHandlerV1.GetAllRole)
	roleRoute.GET("/permissions", v.h.RoleHandlerV1.GetRolePermissions)

	// FileItem routes
	fileItemRoute := route.Group("/file-item")
	fileItemRoute.POST("", v.h.FileItemHandlerV1.CreateOrDirectoryItem)
	fileItemRoute.PUT("", v.h.FileItemHandlerV1.UpdateFileItem)
	fileItemRoute.DELETE("", v.h.FileItemHandlerV1.DeleteFileItem)
	fileItemRoute.DELETE("/force", v.h.FileItemHandlerV1.ForceDeleteFileItem)
	fileItemRoute.PUT("/restore", v.h.FileItemHandlerV1.RestoreFileItem)
	fileItemRoute.POST("/operation", v.h.FileItemHandlerV1.FileOperation)
	fileItemRoute.GET("/tree", v.h.FileItemHandlerV1.GetTreeDirectory)
	fileItemRoute.GET("/tree/deleted", v.h.FileItemHandlerV1.GetDeletedTreeDirectory)
	fileItemRoute.GET("/download", v.h.FileItemHandlerV1.GetDownloadFileItemById)

	// Basket routes
	basketRoute := route.Group("/basket")
	basketRoute.PUT("", v.h.BasketHandlerV1.UpdateBasket)
	basketRoute.GET("", v.h.BasketHandlerV1.GetBasket)
	basketRoute.GET("/user/all", v.h.BasketHandlerV1.GetAllBasketUser)
	basketRoute.GET("/admin/all", v.h.BasketHandlerV1.AdminGetAllBasketUser)

	// Order routes
	orderRoute := route.Group("/order")
	orderRoute.POST("", v.h.OrderHandlerV1.CreateOrderRequest)
	orderRoute.GET("/customer/all", v.h.OrderHandlerV1.GetAllOrderCustomer)
	orderRoute.GET("/customer/details", v.h.OrderHandlerV1.GetOrderCustomerDetails)
	orderRoute.GET("/user/all", v.h.OrderHandlerV1.GetAllOrderUser)
	orderRoute.GET("/user/details", v.h.OrderHandlerV1.GetOrderUserDetails)
	orderRoute.GET("/admin/all", v.h.OrderHandlerV1.AdminGetAllOrderUser)

	// Payment routes
	paymentRoute := route.Group("/payment")
	paymentRoute.POST("/verify", v.h.PaymentHandlerV1.VerifyPayment)
	paymentRoute.GET("/admin/all", v.h.PaymentHandlerV1.AdminGetAllPayment)

	// Gateway routes
	gatewayRoute := route.Group("/gateway")
	gatewayRoute.POST("", v.h.PaymentHandlerV1.CreateOrUpdateGateway)
	gatewayRoute.GET("", v.h.PaymentHandlerV1.GetByIdGateway)
	gatewayRoute.GET("/admin/all", v.h.PaymentHandlerV1.AdminGetAllGateway)

	// Product Category routes
	productCategoryRoute := route.Group("/product-category")
	productCategoryRoute.POST("", v.h.ProductCategoryHandlerV1.CreateCategory)
	productCategoryRoute.PUT("", v.h.ProductCategoryHandlerV1.UpdateCategory)
	productCategoryRoute.DELETE("", v.h.ProductCategoryHandlerV1.DeleteCategory)
	productCategoryRoute.GET("", v.h.ProductCategoryHandlerV1.GetByIdCategory)
	productCategoryRoute.GET("/all", v.h.ProductCategoryHandlerV1.GetAllCategory)
	productCategoryRoute.GET("/admin/all", v.h.ProductCategoryHandlerV1.AdminGetAllCategory)

	// Product routes
	productRoute := route.Group("/product")
	productRoute.POST("", v.h.ProductHandlerV1.CreateProduct)
	productRoute.PUT("", v.h.ProductHandlerV1.UpdateProduct)
	productRoute.DELETE("", v.h.ProductHandlerV1.DeleteProduct)
	productRoute.GET("", v.h.ProductHandlerV1.GetByIdProduct)
	productRoute.GET("/all", v.h.ProductHandlerV1.GetAllProduct)
	productRoute.POST("/filters-sort", v.h.ProductHandlerV1.GetByFiltersSortProduct)
	productRoute.GET("/admin/all", v.h.ProductHandlerV1.AdminGetAllProduct)

	// Discount routes
	discountRoute := route.Group("/discount")
	discountRoute.POST("", v.h.DiscountHandlerV1.CreateDiscount)
	discountRoute.PUT("", v.h.DiscountHandlerV1.UpdateDiscount)
	discountRoute.DELETE("", v.h.DiscountHandlerV1.DeleteDiscount)
	discountRoute.GET("", v.h.DiscountHandlerV1.GetByIdDiscount)
	discountRoute.GET("/all", v.h.DiscountHandlerV1.GetAllDiscount)
	discountRoute.GET("/admin/all", v.h.DiscountHandlerV1.AdminGetAllDiscount)

	// Product Review routes
	productReviewRoute := route.Group("/product-review")
	productReviewRoute.POST("", v.h.ProductReviewHandlerV1.CreateProductReview)
	productReviewRoute.PUT("", v.h.ProductReviewHandlerV1.UpdateProductReview)
	productReviewRoute.DELETE("", v.h.ProductReviewHandlerV1.DeleteProductReview)
	productReviewRoute.GET("", v.h.ProductReviewHandlerV1.GetByIdProductReview)
	productReviewRoute.GET("/all", v.h.ProductReviewHandlerV1.GetAllProductReview)
	productReviewRoute.GET("/admin/all", v.h.ProductReviewHandlerV1.AdminGetAllProductReview)

	// Website routes
	websiteRoute := route.Group("/website")
	websiteRoute.GET("/page", v.h.WebsiteHandlerV1.GetByDomainPage)
	websiteRoute.GET("/header-footer", v.h.WebsiteHandlerV1.GetByDomainHeaderFooter)
	websiteRoute.GET("/product/search", v.h.WebsiteHandlerV1.ProductSearchList)
	websiteRoute.POST("/article/filters-sort", v.h.WebsiteHandlerV1.GetFiltersSortArticle)
	websiteRoute.POST("/product/filters-sort", v.h.WebsiteHandlerV1.GetFiltersSortProduct)
	websiteRoute.GET("/article/category", v.h.WebsiteHandlerV1.GetArticlesByCategorySlug)
	websiteRoute.GET("/product/category", v.h.WebsiteHandlerV1.GetProductsByCategorySlug)
	websiteRoute.GET("/article", v.h.WebsiteHandlerV1.GetSingleArticleBySlug)
	websiteRoute.GET("/product", v.h.WebsiteHandlerV1.GetSingleProductBySlug)

	// DefaultTheme routes
	defaultThemeRoute := route.Group("/default-theme")
	defaultThemeRoute.POST("", v.h.DefaultThemeHandlerV1.CreateDefaultTheme)
	defaultThemeRoute.PUT("", v.h.DefaultThemeHandlerV1.UpdateDefaultTheme)
	defaultThemeRoute.DELETE("", v.h.DefaultThemeHandlerV1.DeleteDefaultTheme)
	defaultThemeRoute.GET("", v.h.DefaultThemeHandlerV1.GetByIdDefaultTheme)
	defaultThemeRoute.GET("/all", v.h.DefaultThemeHandlerV1.GetAllDefaultTheme)

	// Page routes
	pageRoute := route.Group("/page")
	pageRoute.POST("", v.h.PageHandlerV1.CreatePage)
	pageRoute.PUT("", v.h.PageHandlerV1.UpdatePage)
	pageRoute.DELETE("", v.h.PageHandlerV1.DeletePage)
	pageRoute.GET("", v.h.PageHandlerV1.GetByIdPage)
	pageRoute.GET("/all", v.h.PageHandlerV1.GetAllPage)
	pageRoute.GET("/admin/all", v.h.PageHandlerV1.AdminGetAllPage)

	// HeaderFooter routes
	headerFooterRoute := route.Group("/header-footer")
	headerFooterRoute.POST("", v.h.HeaderFooterHandlerV1.CreateHeaderFooter)
	headerFooterRoute.PUT("", v.h.HeaderFooterHandlerV1.UpdateHeaderFooter)
	headerFooterRoute.DELETE("", v.h.HeaderFooterHandlerV1.DeleteHeaderFooter)
	headerFooterRoute.GET("", v.h.HeaderFooterHandlerV1.GetByIdHeaderFooter)
	headerFooterRoute.GET("/all", v.h.HeaderFooterHandlerV1.GetAllHeaderFooter)
	headerFooterRoute.GET("/admin/all", v.h.HeaderFooterHandlerV1.AdminGetAllHeaderFooter)

	// Site routes
	siteRoute := route.Group("/site")
	siteRoute.POST("", v.h.SiteHandlerV1.CreateSite)
	siteRoute.PUT("", v.h.SiteHandlerV1.UpdateSite)
	siteRoute.DELETE("", v.h.SiteHandlerV1.DeleteSite)
	siteRoute.GET("", v.h.SiteHandlerV1.GetByIdSite)
	siteRoute.GET("/all", v.h.SiteHandlerV1.GetAllSite)
	siteRoute.GET("/admin/all", v.h.SiteHandlerV1.AdminGetAllSite)

	// Ticket routes
	ticketRoute := route.Group("/ticket")
	ticketRoute.POST("", v.h.TicketHandlerV1.CreateTicket)
	ticketRoute.PUT("", v.h.TicketHandlerV1.ReplayTicket)
	ticketRoute.PUT("/admin", v.h.TicketHandlerV1.AdminReplayTicket)
	ticketRoute.GET("", v.h.TicketHandlerV1.GetByIdTicket)
	ticketRoute.GET("/all", v.h.TicketHandlerV1.GetAllTicket)
	ticketRoute.GET("/admin/all", v.h.TicketHandlerV1.AdminGetAllTicket)

	// CustomerTicket routes
	customerTicketRoute := route.Group("/customer-ticket")
	customerTicketRoute.POST("", v.h.CustomerTicketHandlerV1.CreateCustomerTicket)
	customerTicketRoute.PUT("", v.h.CustomerTicketHandlerV1.ReplayCustomerTicket)
	customerTicketRoute.PUT("/admin", v.h.CustomerTicketHandlerV1.AdminReplayCustomerTicket)
	customerTicketRoute.GET("", v.h.CustomerTicketHandlerV1.GetByIdCustomerTicket)
	customerTicketRoute.GET("/all", v.h.CustomerTicketHandlerV1.GetAllCustomerTicket)
	customerTicketRoute.GET("/admin/all", v.h.CustomerTicketHandlerV1.AdminGetAllCustomerTicket)

	// UnitPrice routes
	unitPriceRoute := route.Group("/unit-price")
	unitPriceRoute.PUT("", v.h.UnitPriceHandlerV1.UpdateUnitPrice)
	unitPriceRoute.GET("/calculate", v.h.UnitPriceHandlerV1.CalculateUnitPrice)
	unitPriceRoute.GET("/all", v.h.UnitPriceHandlerV1.GetAllUnitPrice)

	// User routes
	userRoute := route.Group("/user")
	userRoute.POST("/register", v.h.UserHandlerV1.RegisterUser)
	userRoute.POST("/login", v.h.UserHandlerV1.LoginUser)
	userRoute.POST("/verify-forget", v.h.UserHandlerV1.RequestVerifyAndForgetUser)
	userRoute.GET("/verify", v.h.UserHandlerV1.VerifyUser)
	userRoute.PUT("/profile", v.h.UserHandlerV1.UpdateProfileUser)
	userRoute.GET("/profile", v.h.UserHandlerV1.GetProfileUser)
	userRoute.POST("/charge-credit", v.h.UserHandlerV1.ChargeCreditRequestUser)
	userRoute.POST("/upgrade-plan", v.h.UserHandlerV1.UpgradePlanRequestUser)
	userRoute.GET("/admin/all", v.h.UserHandlerV1.AdminGetAllUser)

	// Customer routes
	customerRoute := route.Group("/customer")
	customerRoute.POST("/register", v.h.CustomerHandlerV1.RegisterCustomer)
	customerRoute.POST("/login", v.h.CustomerHandlerV1.LoginCustomer)
	customerRoute.POST("/verify-forget", v.h.CustomerHandlerV1.RequestVerifyAndForgetCustomer)
	customerRoute.GET("/verify", v.h.CustomerHandlerV1.VerifyCustomer)
	customerRoute.PUT("/profile", v.h.CustomerHandlerV1.UpdateProfileCustomer)
	customerRoute.GET("/profile", v.h.CustomerHandlerV1.GetProfileCustomer)
	customerRoute.GET("/admin/all", v.h.CustomerHandlerV1.AdminGetAllCustomer)
}
