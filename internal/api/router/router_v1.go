package router

import (
	"github.com/amirex128/new_site_builder/bootstrap"
	"github.com/amirex128/new_site_builder/internal/api/middleware"
	"github.com/amirex128/new_site_builder/internal/contract"
	"github.com/gin-gonic/gin"
)

type RouterV1 struct {
	h         *bootstrap.HandlerManager
	container contract.IContainer
}

func (v RouterV1) Routes(route *gin.RouterGroup) {
	auth := middleware.NewAuthenticator(v.container.GetAuthTransientService(), v.container.GetIdentityService())

	// Article routes
	articleRoute := route.Group("/article")
	articleRoute.POST("", auth.Authenticate(), v.h.ArticleHandlerV1.ArticleCreate)
	articleRoute.PUT("", auth.Authenticate(), v.h.ArticleHandlerV1.ArticleUpdate)
	articleRoute.DELETE("", auth.Authenticate(), v.h.ArticleHandlerV1.ArticleDelete)
	articleRoute.GET("", auth.Authenticate(), v.h.ArticleHandlerV1.ArticleGet)
	articleRoute.GET("/all", auth.Authenticate(), v.h.ArticleHandlerV1.ArticleGetAll)
	articleRoute.POST("/filters-sort", auth.Authenticate(), v.h.ArticleHandlerV1.ArticleGetByFiltersSort)
	articleRoute.GET("/admin/all", auth.Authenticate(), v.h.ArticleHandlerV1.AdminArticleGetAll)

	// Article Category routes
	articleCategoryRoute := route.Group("/article-category")
	articleCategoryRoute.POST("", auth.Authenticate(), v.h.ArticleCategoryHandlerV1.CategoryCreate)
	articleCategoryRoute.PUT("", auth.Authenticate(), v.h.ArticleCategoryHandlerV1.CategoryUpdate)
	articleCategoryRoute.DELETE("", auth.Authenticate(), v.h.ArticleCategoryHandlerV1.CategoryDelete)
	articleCategoryRoute.GET("", auth.Authenticate(), v.h.ArticleCategoryHandlerV1.CategoryGet)
	articleCategoryRoute.GET("/all", auth.Authenticate(), v.h.ArticleCategoryHandlerV1.CategoryGetAll)
	articleCategoryRoute.GET("/admin/all", auth.Authenticate(), v.h.ArticleCategoryHandlerV1.AdminCategoryGetAll)

	// Address routes
	addressRoute := route.Group("/address")
	addressRoute.POST("", auth.Authenticate(), v.h.AddressHandlerV1.CreateAddress)
	addressRoute.PUT("", auth.Authenticate(), v.h.AddressHandlerV1.UpdateAddress)
	addressRoute.DELETE("", auth.Authenticate(), v.h.AddressHandlerV1.DeleteAddress)
	addressRoute.GET("", auth.Authenticate(), v.h.AddressHandlerV1.GetByIdAddress)
	addressRoute.GET("/all", auth.Authenticate(), v.h.AddressHandlerV1.GetAllAddress)
	addressRoute.GET("/city/all", auth.Authenticate(), v.h.AddressHandlerV1.GetAllCity)
	addressRoute.GET("/province/all", auth.Authenticate(), v.h.AddressHandlerV1.GetAllProvince)
	addressRoute.GET("/admin/all", auth.Authenticate(), v.h.AddressHandlerV1.AdminGetAllAddress)

	// Plan routes
	planRoute := route.Group("/plan")
	planRoute.POST("", auth.Authenticate(), v.h.PlanHandlerV1.CreatePlan)
	planRoute.PUT("", auth.Authenticate(), v.h.PlanHandlerV1.UpdatePlan)
	planRoute.DELETE("", auth.Authenticate(), v.h.PlanHandlerV1.DeletePlan)
	planRoute.GET("", auth.Authenticate(), v.h.PlanHandlerV1.GetByIdPlan)
	planRoute.GET("/all", auth.Authenticate(), v.h.PlanHandlerV1.GetAllPlan)
	planRoute.GET("/calculate", auth.Authenticate(), v.h.PlanHandlerV1.CalculatePlanPrice)

	// Role routes
	roleRoute := route.Group("/role")
	roleRoute.POST("", auth.Authenticate(), v.h.RoleHandlerV1.CreateRole)
	roleRoute.PUT("", auth.Authenticate(), v.h.RoleHandlerV1.UpdateRole)
	roleRoute.PUT("/customer", auth.Authenticate(), v.h.RoleHandlerV1.SetRoleToCustomer)
	roleRoute.PUT("/user", auth.Authenticate(), v.h.RoleHandlerV1.SetRoleToUser)
	roleRoute.PUT("/plan", auth.Authenticate(), v.h.RoleHandlerV1.SetRoleToPlan)
	roleRoute.GET("/permission/all", auth.Authenticate(), v.h.RoleHandlerV1.GetAllPermission)
	roleRoute.GET("/all", auth.Authenticate(), v.h.RoleHandlerV1.GetAllRole)
	roleRoute.GET("/permissions", auth.Authenticate(), v.h.RoleHandlerV1.GetRolePermissions)

	// FileItem routes
	fileItemRoute := route.Group("/file-item")
	fileItemRoute.POST("", auth.Authenticate(), v.h.FileItemHandlerV1.CreateOrDirectoryItem)
	fileItemRoute.PUT("", auth.Authenticate(), v.h.FileItemHandlerV1.UpdateFileItem)
	fileItemRoute.DELETE("", auth.Authenticate(), v.h.FileItemHandlerV1.DeleteFileItem)
	fileItemRoute.DELETE("/force", auth.Authenticate(), v.h.FileItemHandlerV1.ForceDeleteFileItem)
	fileItemRoute.PUT("/restore", auth.Authenticate(), v.h.FileItemHandlerV1.RestoreFileItem)
	fileItemRoute.POST("/operation", auth.Authenticate(), v.h.FileItemHandlerV1.FileOperation)
	fileItemRoute.GET("/tree", auth.Authenticate(), v.h.FileItemHandlerV1.GetTreeDirectory)
	fileItemRoute.GET("/tree/deleted", auth.Authenticate(), v.h.FileItemHandlerV1.GetDeletedTreeDirectory)
	fileItemRoute.GET("/download", auth.Authenticate(), v.h.FileItemHandlerV1.GetDownloadFileItemById)

	// Basket routes
	basketRoute := route.Group("/basket")
	basketRoute.PUT("", auth.Authenticate(), v.h.BasketHandlerV1.UpdateBasket)
	basketRoute.GET("", auth.Authenticate(), v.h.BasketHandlerV1.GetBasket)
	basketRoute.GET("/user/all", auth.Authenticate(), v.h.BasketHandlerV1.GetAllBasketUser)
	basketRoute.GET("/admin/all", auth.Authenticate(), v.h.BasketHandlerV1.AdminGetAllBasketUser)

	// Order routes
	orderRoute := route.Group("/order")
	orderRoute.POST("", auth.Authenticate(), v.h.OrderHandlerV1.CreateOrderRequest)
	orderRoute.GET("/customer/all", auth.Authenticate(), v.h.OrderHandlerV1.GetAllOrderCustomer)
	orderRoute.GET("/customer/details", auth.Authenticate(), v.h.OrderHandlerV1.GetOrderCustomerDetails)
	orderRoute.GET("/user/all", auth.Authenticate(), v.h.OrderHandlerV1.GetAllOrderUser)
	orderRoute.GET("/user/details", auth.Authenticate(), v.h.OrderHandlerV1.GetOrderUserDetails)
	orderRoute.GET("/admin/all", auth.Authenticate(), v.h.OrderHandlerV1.AdminGetAllOrderUser)

	// Payment routes
	paymentRoute := route.Group("/payment")
	paymentRoute.POST("/verify", auth.Authenticate(), v.h.PaymentHandlerV1.VerifyPayment)
	paymentRoute.GET("/admin/all", auth.Authenticate(), v.h.PaymentHandlerV1.AdminGetAllPayment)

	// Gateway routes
	gatewayRoute := route.Group("/gateway")
	gatewayRoute.POST("", auth.Authenticate(), v.h.PaymentHandlerV1.CreateOrUpdateGateway)
	gatewayRoute.GET("", auth.Authenticate(), v.h.PaymentHandlerV1.GetByIdGateway)
	gatewayRoute.GET("/admin/all", auth.Authenticate(), v.h.PaymentHandlerV1.AdminGetAllGateway)

	// Product Category routes
	productCategoryRoute := route.Group("/product-category")
	productCategoryRoute.POST("", auth.Authenticate(), v.h.ProductCategoryHandlerV1.CreateCategory)
	productCategoryRoute.PUT("", auth.Authenticate(), v.h.ProductCategoryHandlerV1.UpdateCategory)
	productCategoryRoute.DELETE("", auth.Authenticate(), v.h.ProductCategoryHandlerV1.DeleteCategory)
	productCategoryRoute.GET("", auth.Authenticate(), v.h.ProductCategoryHandlerV1.GetByIdCategory)
	productCategoryRoute.GET("/all", auth.Authenticate(), v.h.ProductCategoryHandlerV1.GetAllCategory)
	productCategoryRoute.GET("/admin/all", auth.Authenticate(), v.h.ProductCategoryHandlerV1.AdminGetAllCategory)

	// Product routes
	productRoute := route.Group("/product")
	productRoute.POST("", auth.Authenticate(), v.h.ProductHandlerV1.CreateProduct)
	productRoute.PUT("", auth.Authenticate(), v.h.ProductHandlerV1.UpdateProduct)
	productRoute.DELETE("", auth.Authenticate(), v.h.ProductHandlerV1.DeleteProduct)
	productRoute.GET("", auth.Authenticate(), v.h.ProductHandlerV1.GetByIdProduct)
	productRoute.GET("/all", auth.Authenticate(), v.h.ProductHandlerV1.GetAllProduct)
	productRoute.POST("/filters-sort", auth.Authenticate(), v.h.ProductHandlerV1.GetByFiltersSortProduct)
	productRoute.GET("/admin/all", auth.Authenticate(), v.h.ProductHandlerV1.AdminGetAllProduct)

	// Discount routes
	discountRoute := route.Group("/discount")
	discountRoute.POST("", auth.Authenticate(), v.h.DiscountHandlerV1.CreateDiscount)
	discountRoute.PUT("", auth.Authenticate(), v.h.DiscountHandlerV1.UpdateDiscount)
	discountRoute.DELETE("", auth.Authenticate(), v.h.DiscountHandlerV1.DeleteDiscount)
	discountRoute.GET("", auth.Authenticate(), v.h.DiscountHandlerV1.GetByIdDiscount)
	discountRoute.GET("/all", auth.Authenticate(), v.h.DiscountHandlerV1.GetAllDiscount)
	discountRoute.GET("/admin/all", auth.Authenticate(), v.h.DiscountHandlerV1.AdminGetAllDiscount)

	// Product Review routes
	productReviewRoute := route.Group("/product-review")
	productReviewRoute.POST("", auth.Authenticate(), v.h.ProductReviewHandlerV1.CreateProductReview)
	productReviewRoute.PUT("", auth.Authenticate(), v.h.ProductReviewHandlerV1.UpdateProductReview)
	productReviewRoute.DELETE("", auth.Authenticate(), v.h.ProductReviewHandlerV1.DeleteProductReview)
	productReviewRoute.GET("", auth.Authenticate(), v.h.ProductReviewHandlerV1.GetByIdProductReview)
	productReviewRoute.GET("/all", auth.Authenticate(), v.h.ProductReviewHandlerV1.GetAllProductReview)
	productReviewRoute.GET("/admin/all", auth.Authenticate(), v.h.ProductReviewHandlerV1.AdminGetAllProductReview)

	// Website routes
	websiteRoute := route.Group("/website")
	websiteRoute.GET("/page", auth.Authenticate(), v.h.WebsiteHandlerV1.GetByDomainPage)
	websiteRoute.GET("/header-footer", auth.Authenticate(), v.h.WebsiteHandlerV1.GetByDomainHeaderFooter)
	websiteRoute.GET("/product/search", auth.Authenticate(), v.h.WebsiteHandlerV1.ProductSearchList)
	websiteRoute.POST("/article/filters-sort", auth.Authenticate(), v.h.WebsiteHandlerV1.GetFiltersSortArticle)
	websiteRoute.POST("/product/filters-sort", auth.Authenticate(), v.h.WebsiteHandlerV1.GetFiltersSortProduct)
	websiteRoute.GET("/article/category", auth.Authenticate(), v.h.WebsiteHandlerV1.GetArticlesByCategorySlug)
	websiteRoute.GET("/product/category", auth.Authenticate(), v.h.WebsiteHandlerV1.GetProductsByCategorySlug)
	websiteRoute.GET("/article", auth.Authenticate(), v.h.WebsiteHandlerV1.GetSingleArticleBySlug)
	websiteRoute.GET("/product", auth.Authenticate(), v.h.WebsiteHandlerV1.GetSingleProductBySlug)

	// DefaultTheme routes
	defaultThemeRoute := route.Group("/default-theme")
	defaultThemeRoute.POST("", auth.Authenticate(), v.h.DefaultThemeHandlerV1.CreateDefaultTheme)
	defaultThemeRoute.PUT("", auth.Authenticate(), v.h.DefaultThemeHandlerV1.UpdateDefaultTheme)
	defaultThemeRoute.DELETE("", auth.Authenticate(), v.h.DefaultThemeHandlerV1.DeleteDefaultTheme)
	defaultThemeRoute.GET("", auth.Authenticate(), v.h.DefaultThemeHandlerV1.GetByIdDefaultTheme)
	defaultThemeRoute.GET("/all", auth.Authenticate(), v.h.DefaultThemeHandlerV1.GetAllDefaultTheme)

	// Page routes
	pageRoute := route.Group("/page")
	pageRoute.POST("", auth.Authenticate(), v.h.PageHandlerV1.CreatePage)
	pageRoute.PUT("", auth.Authenticate(), v.h.PageHandlerV1.UpdatePage)
	pageRoute.DELETE("", auth.Authenticate(), v.h.PageHandlerV1.DeletePage)
	pageRoute.GET("", auth.Authenticate(), v.h.PageHandlerV1.GetByIdPage)
	pageRoute.GET("/all", auth.Authenticate(), v.h.PageHandlerV1.GetAllPage)
	pageRoute.GET("/admin/all", auth.Authenticate(), v.h.PageHandlerV1.AdminGetAllPage)

	// HeaderFooter routes
	headerFooterRoute := route.Group("/header-footer")
	headerFooterRoute.POST("", auth.Authenticate(), v.h.HeaderFooterHandlerV1.CreateHeaderFooter)
	headerFooterRoute.PUT("", auth.Authenticate(), v.h.HeaderFooterHandlerV1.UpdateHeaderFooter)
	headerFooterRoute.DELETE("", auth.Authenticate(), v.h.HeaderFooterHandlerV1.DeleteHeaderFooter)
	headerFooterRoute.GET("", auth.Authenticate(), v.h.HeaderFooterHandlerV1.GetByIdHeaderFooter)
	headerFooterRoute.GET("/all", auth.Authenticate(), v.h.HeaderFooterHandlerV1.GetAllHeaderFooter)
	headerFooterRoute.GET("/admin/all", auth.Authenticate(), v.h.HeaderFooterHandlerV1.AdminGetAllHeaderFooter)

	// Site routes
	siteRoute := route.Group("/site")
	siteRoute.POST("", auth.Authenticate(), v.h.SiteHandlerV1.CreateSite)
	siteRoute.PUT("", auth.Authenticate(), v.h.SiteHandlerV1.UpdateSite)
	siteRoute.DELETE("", auth.Authenticate(), v.h.SiteHandlerV1.DeleteSite)
	siteRoute.GET("", auth.Authenticate(), v.h.SiteHandlerV1.GetByIdSite)
	siteRoute.GET("/all", auth.Authenticate(), v.h.SiteHandlerV1.GetAllSite)
	siteRoute.GET("/admin/all", auth.Authenticate(), v.h.SiteHandlerV1.AdminGetAllSite)

	// Ticket routes
	ticketRoute := route.Group("/ticket")
	ticketRoute.POST("", auth.Authenticate(), v.h.TicketHandlerV1.CreateTicket)
	ticketRoute.PUT("", auth.Authenticate(), v.h.TicketHandlerV1.ReplayTicket)
	ticketRoute.PUT("/admin", auth.Authenticate(), v.h.TicketHandlerV1.AdminReplayTicket)
	ticketRoute.GET("", auth.Authenticate(), v.h.TicketHandlerV1.GetByIdTicket)
	ticketRoute.GET("/all", auth.Authenticate(), v.h.TicketHandlerV1.GetAllTicket)
	ticketRoute.GET("/admin/all", auth.Authenticate(), v.h.TicketHandlerV1.AdminGetAllTicket)

	// CustomerTicket routes
	customerTicketRoute := route.Group("/customer-ticket")
	customerTicketRoute.POST("", auth.Authenticate(), v.h.CustomerTicketHandlerV1.CreateCustomerTicket)
	customerTicketRoute.PUT("", auth.Authenticate(), v.h.CustomerTicketHandlerV1.ReplayCustomerTicket)
	customerTicketRoute.PUT("/admin", auth.Authenticate(), v.h.CustomerTicketHandlerV1.AdminReplayCustomerTicket)
	customerTicketRoute.GET("", auth.Authenticate(), v.h.CustomerTicketHandlerV1.GetByIdCustomerTicket)
	customerTicketRoute.GET("/all", auth.Authenticate(), v.h.CustomerTicketHandlerV1.GetAllCustomerTicket)
	customerTicketRoute.GET("/admin/all", auth.Authenticate(), v.h.CustomerTicketHandlerV1.AdminGetAllCustomerTicket)

	// UnitPrice routes
	unitPriceRoute := route.Group("/unit-price")
	unitPriceRoute.PUT("", auth.Authenticate(), v.h.UnitPriceHandlerV1.UpdateUnitPrice)
	unitPriceRoute.GET("/calculate", auth.Authenticate(), v.h.UnitPriceHandlerV1.CalculateUnitPrice)
	unitPriceRoute.GET("/all", auth.Authenticate(), v.h.UnitPriceHandlerV1.GetAllUnitPrice)

	// User routes
	userRoute := route.Group("/user")
	userRoute.POST("/register", v.h.UserHandlerV1.RegisterUser)
	userRoute.POST("/login", v.h.UserHandlerV1.LoginUser)
	userRoute.POST("/verify-forget", v.h.UserHandlerV1.RequestVerifyAndForgetUser)
	userRoute.GET("/verify", v.h.UserHandlerV1.VerifyUser)
	userRoute.PUT("/profile", auth.Authenticate(), v.h.UserHandlerV1.UpdateProfileUser)
	userRoute.GET("/profile", auth.Authenticate(), v.h.UserHandlerV1.GetProfileUser)
	userRoute.POST("/charge-credit", auth.Authenticate(), v.h.UserHandlerV1.ChargeCreditRequestUser)
	userRoute.POST("/upgrade-plan", auth.Authenticate(), v.h.UserHandlerV1.UpgradePlanRequestUser)
	userRoute.GET("/admin/all", auth.Authenticate(), v.h.UserHandlerV1.AdminGetAllUser)

	// Customer routes
	customerRoute := route.Group("/customer")
	customerRoute.POST("/register", auth.Authenticate(), v.h.CustomerHandlerV1.RegisterCustomer)
	customerRoute.POST("/login", auth.Authenticate(), v.h.CustomerHandlerV1.LoginCustomer)
	customerRoute.POST("/verify-forget", auth.Authenticate(), v.h.CustomerHandlerV1.RequestVerifyAndForgetCustomer)
	customerRoute.GET("/verify", auth.Authenticate(), v.h.CustomerHandlerV1.VerifyCustomer)
	customerRoute.PUT("/profile", auth.Authenticate(), v.h.CustomerHandlerV1.UpdateProfileCustomer)
	customerRoute.GET("/profile", auth.Authenticate(), v.h.CustomerHandlerV1.GetProfileCustomer)
	customerRoute.GET("/admin/all", auth.Authenticate(), v.h.CustomerHandlerV1.AdminGetAllCustomer)
}
