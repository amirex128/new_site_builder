package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	"gorm.io/gorm"
	"log"
	"time"
)

func MysqlProvider(logger sflogger.Logger) {
	// Create MySQL configuration directly as a struct
	// Register your database connections with meaningful names and options
	err := sform.RegisterConnection(
		sform.WithLogger(logger),
		sform.WithGlobalOptions(func(db *gorm.DB) {
			db.Debug()
		}),
		sform.WithConnectionDetails("main", &sform.MySQLConfig{
			Username:     "user",
			Password:     "password",
			Host:         "127.0.0.1",
			Port:         3306,
			Database:     "dbname",
			Charset:      "utf8mb4",
			ParseTime:    true,
			Loc:          "Local",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  5 * time.Minute,
		}),
		sform.WithConnectionDetails("read", &sform.MySQLConfig{
			Username:     "user",
			Password:     "password",
			Host:         "127.0.0.1",
			Port:         3306,
			Database:     "dbname",
			Charset:      "utf8mb4",
			ParseTime:    true,
			Loc:          "Local",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  5 * time.Minute,
		}),
	)

	if err != nil {
		log.Fatalf("Failed to register database connection: %v", err)
	}

}
