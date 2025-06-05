package serviceprovider

import (
	"fmt"
	"github.com/amirex128/new_site_builder/config"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	"gorm.io/gorm"
)

func MysqlProvider(cfg *config.Config, logger sflogger.Logger) {
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
				&domain2.Configuration{},
				&domain2.Address{},
				&domain2.AddressCustomer{},
				&domain2.AddressUser{},
				&domain2.Article{},
				&domain2.Basket{},
				&domain2.BasketItem{},
				&domain2.ArticleCategory{},
				&domain2.City{},
				&domain2.Comment{},
				&domain2.Coupon{},
				&domain2.Credit{},
				&domain2.Customer{},
				&domain2.CustomerComment{},
				&domain2.CustomerRole{},
				&domain2.CustomerTicket{},
				&domain2.CustomerTicketMedia{},
				&domain2.DefaultTheme{},
				&domain2.Discount{},
				&domain2.FileItem{},
				&domain2.Gateway{},
				&domain2.HeaderFooter{},
				&domain2.Media{},
				&domain2.Order{},
				&domain2.OrderItem{},
				&domain2.Page{},
				&domain2.PageArticleUsage{},
				&domain2.PageHeaderFooterUsage{},
				&domain2.PageProductUsage{},
				&domain2.ParbadPayment{},
				&domain2.ParbadTransaction{},
				&domain2.Payment{},
				&domain2.Permission{},
				&domain2.PermissionRole{},
				&domain2.Plan{},
				&domain2.Product{},
				&domain2.ProductAttribute{},
				&domain2.ProductCategory{},
				&domain2.ProductReview{},
				&domain2.ProductVariant{},
				&domain2.Province{},
				&domain2.ReturnItem{},
				&domain2.Role{},
				&domain2.RolePlan{},
				&domain2.RoleUser{},
				&domain2.Setting{},
				&domain2.Site{},
				&domain2.Storage{},
				&domain2.Ticket{},
				&domain2.TicketMedia{},
				&domain2.UnitPrice{},
				&domain2.User{},
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
		logger.ErrorWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, fmt.Sprintf("Failed to register database connection: %v", err), nil)
	}
	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded Mysql", nil)

}
