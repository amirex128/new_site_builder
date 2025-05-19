package router

import (
	"context"
	"fmt"
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"

	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"

	"github.com/gin-gonic/gin"
)

func InitServer(handlers *bootstrap.HandlerManager, container contract.IContainer, logger sflogger.Logger, cnf contract.IConfig) {

	RegisterRoutes(handlers, container)

	err := sfrouting.StartServer(fmt.Sprintf(":%s", cnf.GetString("app_port")))
	if err != nil {
		logger.Errorf("Failed to start server: %s", err.Error())
		return
	}
}

// RegisterRoutes registers all routes
func RegisterRoutes(h *bootstrap.HandlerManager, container contract.IContainer) {
	// Register user routes in a group
	sfrouting.RegisterRouterGroup("/api/v1", &RouterV1{h: h, container: container})
}

func runServer(r *gin.Engine, logger sflogger.Logger) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		// service connections
		logger.InfoWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Startup, "Started", nil)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FatalWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Status.Error, err.Error(), nil)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.InfoWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Shutdown, "Shutdown Server ...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.InfoWithCategory(sflogger.Category.System.Shutdown, sflogger.SubCategory.Status.Error, err.Error(), nil)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.InfoWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Status.Timeout, "timeout of 5 seconds.", nil)
	}

	logger.InfoWithCategory(sflogger.Category.System.Shutdown, sflogger.SubCategory.Status.Success, "Server exiting", nil)
}
