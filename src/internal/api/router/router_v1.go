package router

import (
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/amirex128/new_site_builder/src/internal/api/middleware"
	"github.com/amirex128/new_site_builder/src/internal/contract"
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
	articleRoute.POST("", v.h.ArticleHandlerV1.ArticleCreate, auth.Authenticate())
	articleRoute.PUT("", v.h.ArticleHandlerV1.ArticleUpdate, auth.Authenticate())
	articleRoute.DELETE("", v.h.ArticleHandlerV1.ArticleDelete, auth.Authenticate())
	articleRoute.GET("", v.h.ArticleHandlerV1.ArticleGet, auth.Authenticate())
	articleRoute.GET("/all", v.h.ArticleHandlerV1.ArticleGetAll, auth.Authenticate())
	articleRoute.POST("/filters-sort", v.h.ArticleHandlerV1.ArticleGetByFiltersSort, auth.Authenticate())
	articleRoute.GET("/admin/all", v.h.ArticleHandlerV1.AdminArticleGetAll, auth.Authenticate())

	// Article Category routes
	articleCategoryRoute := route.Group("/article-category")
	articleCategoryRoute.POST("", v.h.ArticleCategoryHandlerV1.CategoryCreate, auth.Authenticate())
	articleCategoryRoute.PUT("", v.h.ArticleCategoryHandlerV1.CategoryUpdate, auth.Authenticate())
	articleCategoryRoute.DELETE("", v.h.ArticleCategoryHandlerV1.CategoryDelete, auth.Authenticate())
	articleCategoryRoute.GET("", v.h.ArticleCategoryHandlerV1.CategoryGet, auth.Authenticate())
	articleCategoryRoute.GET("/all", v.h.ArticleCategoryHandlerV1.CategoryGetAll, auth.Authenticate())
	articleCategoryRoute.GET("/admin/all", v.h.ArticleCategoryHandlerV1.AdminCategoryGetAll, auth.Authenticate())

	// Address routes
	addressRoute := route.Group("/address")
	addressRoute.POST("", v.h.AddressHandlerV1.CreateAddress, auth.Authenticate())
	addressRoute.PUT("", v.h.AddressHandlerV1.UpdateAddress, auth.Authenticate())
	addressRoute.DELETE("", v.h.AddressHandlerV1.DeleteAddress, auth.Authenticate())
	addressRoute.GET("", v.h.AddressHandlerV1.GetByIdAddress, auth.Authenticate())
	addressRoute.GET("/all", v.h.AddressHandlerV1.GetAllAddress, auth.Authenticate())
	addressRoute.GET("/city/all", v.h.AddressHandlerV1.GetAllCity, auth.Authenticate())
	addressRoute.GET("/province/all", v.h.AddressHandlerV1.GetAllProvince, auth.Authenticate())
	addressRoute.GET("/admin/all", v.h.AddressHandlerV1.AdminGetAllAddress, auth.Authenticate())

	// Plan routes
	planRoute := route.Group("/plan")
	planRoute.POST("", v.h.PlanHandlerV1.CreatePlan, auth.Authenticate())
	planRoute.PUT("", v.h.PlanHandlerV1.UpdatePlan, auth.Authenticate())
	planRoute.DELETE("", v.h.PlanHandlerV1.DeletePlan, auth.Authenticate())
	planRoute.GET("", v.h.PlanHandlerV1.GetByIdPlan, auth.Authenticate())
	planRoute.GET("/all", v.h.PlanHandlerV1.GetAllPlan, auth.Authenticate())
	planRoute.GET("/calculate", v.h.PlanHandlerV1.CalculatePlanPrice, auth.Authenticate())

	// Role routes
	roleRoute := route.Group("/role")
	roleRoute.POST("", v.h.RoleHandlerV1.CreateRole, auth.Authenticate())
	roleRoute.PUT("", v.h.RoleHandlerV1.UpdateRole, auth.Authenticate())
	roleRoute.PUT("/customer", v.h.RoleHandlerV1.SetRoleToCustomer, auth.Authenticate())
	roleRoute.PUT("/user", v.h.RoleHandlerV1.SetRoleToUser, auth.Authenticate())
	roleRoute.PUT("/plan", v.h.RoleHandlerV1.SetRoleToPlan, auth.Authenticate())
	roleRoute.GET("/permission/all", v.h.RoleHandlerV1.GetAllPermission, auth.Authenticate())
	roleRoute.GET("/all", v.h.RoleHandlerV1.GetAllRole, auth.Authenticate())
	roleRoute.GET("/permissions", v.h.RoleHandlerV1.GetRolePermissions, auth.Authenticate())

	// FileItem routes
	fileItemRoute := route.Group("/file-item")
	fileItemRoute.POST("", v.h.FileItemHandlerV1.CreateOrDirectoryItem, auth.Authenticate())
	fileItemRoute.PUT("", v.h.FileItemHandlerV1.UpdateFileItem, auth.Authenticate())
	fileItemRoute.DELETE("", v.h.FileItemHandlerV1.DeleteFileItem, auth.Authenticate())
	fileItemRoute.DELETE("/force", v.h.FileItemHandlerV1.ForceDeleteFileItem, auth.Authenticate())
	fileItemRoute.PUT("/restore", v.h.FileItemHandlerV1.RestoreFileItem, auth.Authenticate())
	fileItemRoute.POST("/operation", v.h.FileItemHandlerV1.FileOperation, auth.Authenticate())
	fileItemRoute.GET("/tree", v.h.FileItemHandlerV1.GetTreeDirectory, auth.Authenticate())
	fileItemRoute.GET("/tree/deleted", v.h.FileItemHandlerV1.GetDeletedTreeDirectory, auth.Authenticate())
	fileItemRoute.GET("/download", v.h.FileItemHandlerV1.GetDownloadFileItemById, auth.Authenticate())

	// Basket routes
	basketRoute := route.Group("/basket")
	basketRoute.PUT("", v.h.BasketHandlerV1.UpdateBasket, auth.Authenticate())
	basketRoute.GET("", v.h.BasketHandlerV1.GetBasket, auth.Authenticate())
	basketRoute.GET("/user/all", v.h.BasketHandlerV1.GetAllBasketUser, auth.Authenticate())
	basketRoute.GET("/admin/all", v.h.BasketHandlerV1.AdminGetAllBasketUser, auth.Authenticate())

	// Order routes
	orderRoute := route.Group("/order")
	orderRoute.POST("", v.h.OrderHandlerV1.CreateOrderRequest, auth.Authenticate())
	orderRoute.GET("/customer/all", v.h.OrderHandlerV1.GetAllOrderCustomer, auth.Authenticate())
	orderRoute.GET("/customer/details", v.h.OrderHandlerV1.GetOrderCustomerDetails, auth.Authenticate())
	orderRoute.GET("/user/all", v.h.OrderHandlerV1.GetAllOrderUser, auth.Authenticate())
	orderRoute.GET("/user/details", v.h.OrderHandlerV1.GetOrderUserDetails, auth.Authenticate())
	orderRoute.GET("/admin/all", v.h.OrderHandlerV1.AdminGetAllOrderUser, auth.Authenticate())

	// Payment routes
	paymentRoute := route.Group("/payment")
	paymentRoute.POST("/verify", v.h.PaymentHandlerV1.VerifyPayment, auth.Authenticate())
	paymentRoute.GET("/admin/all", v.h.PaymentHandlerV1.AdminGetAllPayment, auth.Authenticate())

	// Gateway routes
	gatewayRoute := route.Group("/gateway")
	gatewayRoute.POST("", v.h.PaymentHandlerV1.CreateOrUpdateGateway, auth.Authenticate())
	gatewayRoute.GET("", v.h.PaymentHandlerV1.GetByIdGateway, auth.Authenticate())
	gatewayRoute.GET("/admin/all", v.h.PaymentHandlerV1.AdminGetAllGateway, auth.Authenticate())

	// Product Category routes
	productCategoryRoute := route.Group("/product-category")
	productCategoryRoute.POST("", v.h.ProductCategoryHandlerV1.CreateCategory, auth.Authenticate())
	productCategoryRoute.PUT("", v.h.ProductCategoryHandlerV1.UpdateCategory, auth.Authenticate())
	productCategoryRoute.DELETE("", v.h.ProductCategoryHandlerV1.DeleteCategory, auth.Authenticate())
	productCategoryRoute.GET("", v.h.ProductCategoryHandlerV1.GetByIdCategory, auth.Authenticate())
	productCategoryRoute.GET("/all", v.h.ProductCategoryHandlerV1.GetAllCategory, auth.Authenticate())
	productCategoryRoute.GET("/admin/all", v.h.ProductCategoryHandlerV1.AdminGetAllCategory, auth.Authenticate())

	// Product routes
	productRoute := route.Group("/product")
	productRoute.POST("", v.h.ProductHandlerV1.CreateProduct, auth.Authenticate())
	productRoute.PUT("", v.h.ProductHandlerV1.UpdateProduct, auth.Authenticate())
	productRoute.DELETE("", v.h.ProductHandlerV1.DeleteProduct, auth.Authenticate())
	productRoute.GET("", v.h.ProductHandlerV1.GetByIdProduct, auth.Authenticate())
	productRoute.GET("/all", v.h.ProductHandlerV1.GetAllProduct, auth.Authenticate())
	productRoute.POST("/filters-sort", v.h.ProductHandlerV1.GetByFiltersSortProduct, auth.Authenticate())
	productRoute.GET("/admin/all", v.h.ProductHandlerV1.AdminGetAllProduct, auth.Authenticate())

	// Discount routes
	discountRoute := route.Group("/discount")
	discountRoute.POST("", v.h.DiscountHandlerV1.CreateDiscount, auth.Authenticate())
	discountRoute.PUT("", v.h.DiscountHandlerV1.UpdateDiscount, auth.Authenticate())
	discountRoute.DELETE("", v.h.DiscountHandlerV1.DeleteDiscount, auth.Authenticate())
	discountRoute.GET("", v.h.DiscountHandlerV1.GetByIdDiscount, auth.Authenticate())
	discountRoute.GET("/all", v.h.DiscountHandlerV1.GetAllDiscount, auth.Authenticate())
	discountRoute.GET("/admin/all", v.h.DiscountHandlerV1.AdminGetAllDiscount, auth.Authenticate())

	// Product Review routes
	productReviewRoute := route.Group("/product-review")
	productReviewRoute.POST("", v.h.ProductReviewHandlerV1.CreateProductReview, auth.Authenticate())
	productReviewRoute.PUT("", v.h.ProductReviewHandlerV1.UpdateProductReview, auth.Authenticate())
	productReviewRoute.DELETE("", v.h.ProductReviewHandlerV1.DeleteProductReview, auth.Authenticate())
	productReviewRoute.GET("", v.h.ProductReviewHandlerV1.GetByIdProductReview, auth.Authenticate())
	productReviewRoute.GET("/all", v.h.ProductReviewHandlerV1.GetAllProductReview, auth.Authenticate())
	productReviewRoute.GET("/admin/all", v.h.ProductReviewHandlerV1.AdminGetAllProductReview, auth.Authenticate())

	// Website routes
	websiteRoute := route.Group("/website")
	websiteRoute.GET("/page", v.h.WebsiteHandlerV1.GetByDomainPage, auth.Authenticate())
	websiteRoute.GET("/header-footer", v.h.WebsiteHandlerV1.GetByDomainHeaderFooter, auth.Authenticate())
	websiteRoute.GET("/product/search", v.h.WebsiteHandlerV1.ProductSearchList, auth.Authenticate())
	websiteRoute.POST("/article/filters-sort", v.h.WebsiteHandlerV1.GetFiltersSortArticle, auth.Authenticate())
	websiteRoute.POST("/product/filters-sort", v.h.WebsiteHandlerV1.GetFiltersSortProduct, auth.Authenticate())
	websiteRoute.GET("/article/category", v.h.WebsiteHandlerV1.GetArticlesByCategorySlug, auth.Authenticate())
	websiteRoute.GET("/product/category", v.h.WebsiteHandlerV1.GetProductsByCategorySlug, auth.Authenticate())
	websiteRoute.GET("/article", v.h.WebsiteHandlerV1.GetSingleArticleBySlug, auth.Authenticate())
	websiteRoute.GET("/product", v.h.WebsiteHandlerV1.GetSingleProductBySlug, auth.Authenticate())

	// DefaultTheme routes
	defaultThemeRoute := route.Group("/default-theme")
	defaultThemeRoute.POST("", v.h.DefaultThemeHandlerV1.CreateDefaultTheme, auth.Authenticate())
	defaultThemeRoute.PUT("", v.h.DefaultThemeHandlerV1.UpdateDefaultTheme, auth.Authenticate())
	defaultThemeRoute.DELETE("", v.h.DefaultThemeHandlerV1.DeleteDefaultTheme, auth.Authenticate())
	defaultThemeRoute.GET("", v.h.DefaultThemeHandlerV1.GetByIdDefaultTheme, auth.Authenticate())
	defaultThemeRoute.GET("/all", v.h.DefaultThemeHandlerV1.GetAllDefaultTheme, auth.Authenticate())

	// Page routes
	pageRoute := route.Group("/page")
	pageRoute.POST("", v.h.PageHandlerV1.CreatePage, auth.Authenticate())
	pageRoute.PUT("", v.h.PageHandlerV1.UpdatePage, auth.Authenticate())
	pageRoute.DELETE("", v.h.PageHandlerV1.DeletePage, auth.Authenticate())
	pageRoute.GET("", v.h.PageHandlerV1.GetByIdPage, auth.Authenticate())
	pageRoute.GET("/all", v.h.PageHandlerV1.GetAllPage, auth.Authenticate())
	pageRoute.GET("/admin/all", v.h.PageHandlerV1.AdminGetAllPage, auth.Authenticate())

	// HeaderFooter routes
	headerFooterRoute := route.Group("/header-footer")
	headerFooterRoute.POST("", v.h.HeaderFooterHandlerV1.CreateHeaderFooter, auth.Authenticate())
	headerFooterRoute.PUT("", v.h.HeaderFooterHandlerV1.UpdateHeaderFooter, auth.Authenticate())
	headerFooterRoute.DELETE("", v.h.HeaderFooterHandlerV1.DeleteHeaderFooter, auth.Authenticate())
	headerFooterRoute.GET("", v.h.HeaderFooterHandlerV1.GetByIdHeaderFooter, auth.Authenticate())
	headerFooterRoute.GET("/all", v.h.HeaderFooterHandlerV1.GetAllHeaderFooter, auth.Authenticate())
	headerFooterRoute.GET("/admin/all", v.h.HeaderFooterHandlerV1.AdminGetAllHeaderFooter, auth.Authenticate())

	// Site routes
	siteRoute := route.Group("/site")
	siteRoute.POST("", v.h.SiteHandlerV1.CreateSite, auth.Authenticate())
	siteRoute.PUT("", v.h.SiteHandlerV1.UpdateSite, auth.Authenticate())
	siteRoute.DELETE("", v.h.SiteHandlerV1.DeleteSite, auth.Authenticate())
	siteRoute.GET("", v.h.SiteHandlerV1.GetByIdSite, auth.Authenticate())
	siteRoute.GET("/all", v.h.SiteHandlerV1.GetAllSite, auth.Authenticate())
	siteRoute.GET("/admin/all", v.h.SiteHandlerV1.AdminGetAllSite, auth.Authenticate())

	// Ticket routes
	ticketRoute := route.Group("/ticket")
	ticketRoute.POST("", v.h.TicketHandlerV1.CreateTicket, auth.Authenticate())
	ticketRoute.PUT("", v.h.TicketHandlerV1.ReplayTicket, auth.Authenticate())
	ticketRoute.PUT("/admin", v.h.TicketHandlerV1.AdminReplayTicket, auth.Authenticate())
	ticketRoute.GET("", v.h.TicketHandlerV1.GetByIdTicket, auth.Authenticate())
	ticketRoute.GET("/all", v.h.TicketHandlerV1.GetAllTicket, auth.Authenticate())
	ticketRoute.GET("/admin/all", v.h.TicketHandlerV1.AdminGetAllTicket, auth.Authenticate())

	// CustomerTicket routes
	customerTicketRoute := route.Group("/customer-ticket")
	customerTicketRoute.POST("", v.h.CustomerTicketHandlerV1.CreateCustomerTicket, auth.Authenticate())
	customerTicketRoute.PUT("", v.h.CustomerTicketHandlerV1.ReplayCustomerTicket, auth.Authenticate())
	customerTicketRoute.PUT("/admin", v.h.CustomerTicketHandlerV1.AdminReplayCustomerTicket, auth.Authenticate())
	customerTicketRoute.GET("", v.h.CustomerTicketHandlerV1.GetByIdCustomerTicket, auth.Authenticate())
	customerTicketRoute.GET("/all", v.h.CustomerTicketHandlerV1.GetAllCustomerTicket, auth.Authenticate())
	customerTicketRoute.GET("/admin/all", v.h.CustomerTicketHandlerV1.AdminGetAllCustomerTicket, auth.Authenticate())

	// UnitPrice routes
	unitPriceRoute := route.Group("/unit-price")
	unitPriceRoute.PUT("", v.h.UnitPriceHandlerV1.UpdateUnitPrice, auth.Authenticate())
	unitPriceRoute.GET("/calculate", v.h.UnitPriceHandlerV1.CalculateUnitPrice, auth.Authenticate())
	unitPriceRoute.GET("/all", v.h.UnitPriceHandlerV1.GetAllUnitPrice, auth.Authenticate())

	// User routes
	userRoute := route.Group("/user")
	userRoute.POST("/register", v.h.UserHandlerV1.RegisterUser)
	userRoute.POST("/login", v.h.UserHandlerV1.LoginUser)
	userRoute.POST("/verify-forget", v.h.UserHandlerV1.RequestVerifyAndForgetUser)
	userRoute.GET("/verify", v.h.UserHandlerV1.VerifyUser)
	userRoute.PUT("/profile", v.h.UserHandlerV1.UpdateProfileUser, auth.Authenticate())
	userRoute.GET("/profile", v.h.UserHandlerV1.GetProfileUser, auth.Authenticate())
	userRoute.POST("/charge-credit", v.h.UserHandlerV1.ChargeCreditRequestUser, auth.Authenticate())
	userRoute.POST("/upgrade-plan", v.h.UserHandlerV1.UpgradePlanRequestUser, auth.Authenticate())
	userRoute.GET("/admin/all", v.h.UserHandlerV1.AdminGetAllUser, auth.Authenticate())

	// Customer routes
	customerRoute := route.Group("/customer")
	customerRoute.POST("/register", v.h.CustomerHandlerV1.RegisterCustomer, auth.Authenticate())
	customerRoute.POST("/login", v.h.CustomerHandlerV1.LoginCustomer, auth.Authenticate())
	customerRoute.POST("/verify-forget", v.h.CustomerHandlerV1.RequestVerifyAndForgetCustomer, auth.Authenticate())
	customerRoute.GET("/verify", v.h.CustomerHandlerV1.VerifyCustomer, auth.Authenticate())
	customerRoute.PUT("/profile", v.h.CustomerHandlerV1.UpdateProfileCustomer, auth.Authenticate())
	customerRoute.GET("/profile", v.h.CustomerHandlerV1.GetProfileCustomer, auth.Authenticate())
	customerRoute.GET("/admin/all", v.h.CustomerHandlerV1.AdminGetAllCustomer, auth.Authenticate())
}
