package serviceprovider

import (
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	"github.com/amirex128/new_site_builder/src/config"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

func MysqlProvider(cfg *config.Config, logger sflogger.Logger) {
	// Create MySQL configuration directly as a struct
	// Register your database connections with meaningful names and options
	err := sform.RegisterConnection(
		sform.WithLogger(logger),
		sform.WithRetryOptions(&sform.RetryOptions{
			MaxRetries:     100,
			InitialBackoff: time.Second,
			MaxBackoff:     15 * time.Second,
			BackoffFactor:  1.5,
		}),
		sform.WithGlobalOptions(func(db *gorm.DB) {
			db.Debug()
			err := db.AutoMigrate(
				&domain.Address{},
				&domain.AddressCustomer{},
				&domain.AddressUser{},
				&domain.Article{},
				&domain.Basket{},
				&domain.BasketItem{},
				&domain.ArticleCategory{},
				&domain.City{},
				&domain.Comment{},
				&domain.Coupon{},
				&domain.Credit{},
				&domain.Customer{},
				&domain.CustomerComment{},
				&domain.CustomerRole{},
				&domain.CustomerTicket{},
				&domain.CustomerTicketMedia{},
				&domain.DefaultTheme{},
				&domain.Discount{},
				&domain.FileItem{},
				&domain.Gateway{},
				&domain.HeaderFooter{},
				&domain.Media{},
				&domain.Order{},
				&domain.OrderItem{},
				&domain.Page{},
				&domain.PageArticleUsage{},
				&domain.PageHeaderFooterUsage{},
				&domain.PageProductUsage{},
				&domain.ParbadPayment{},
				&domain.ParbadTransaction{},
				&domain.Payment{},
				&domain.Permission{},
				&domain.PermissionRole{},
				&domain.Plan{},
				&domain.Product{},
				&domain.ProductAttribute{},
				&domain.ProductCategory{},
				&domain.ProductReview{},
				&domain.ProductVariant{},
				&domain.Province{},
				&domain.ReturnItem{},
				&domain.Role{},
				&domain.RolePlan{},
				&domain.RoleUser{},
				&domain.Setting{},
				&domain.Site{},
				&domain.Storage{},
				&domain.Ticket{},
				&domain.TicketMedia{},
				&domain.UnitPrice{},
				&domain.User{},
			)
			if err != nil {
				logger.Errorf("Error migrating database: %v", err)
			}
		}),
		sform.WithConnectionDetails("main", &sform.MySQLConfig{
			Username: cfg.MysqlUsername,
			Password: cfg.MysqlPassword,
			Host:     cfg.MysqlHost,
			Port: func() int {
				port, err := strconv.Atoi(cfg.MysqlPort)
				if err != nil {
					logger.Errorf("Failed to convert MySQL port to int: %v", err)
				}
				return port
			}(),
			Database:     cfg.MysqlDatabase,
			Charset:      "utf8mb4",
			ParseTime:    true,
			Loc:          "Local",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  5 * time.Minute,
		}),
	)

	if err != nil {
		logger.Errorf("Failed to register database connection: %v", err)
	}

}
