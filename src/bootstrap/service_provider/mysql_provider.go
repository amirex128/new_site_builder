package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	"github.com/amirex128/new_site_builder/src/config"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func MysqlProvider(cfg *config.Config, logger sflogger.Logger) {
	// Create MySQL configuration directly as a struct
	// Register your database connections with meaningful names and options
	err := sform.RegisterConnection(
		sform.WithLogger(logger),
		sform.WithGlobalOptions(func(db *gorm.DB) {
			db.Debug()
			err := db.AutoMigrate(
				&domain.Address{},
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
					logger.Fatalf("Failed to convert MySQL port to int: %v", err)
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
		logger.Fatalf("Failed to register database connection: %v", err)
	}

}
