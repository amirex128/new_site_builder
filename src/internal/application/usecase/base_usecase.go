package usecase

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/gin-gonic/gin"
)

type BaseUsecase struct {
	Ctx    *gin.Context
	Logger sflogger.Logger
}

func (u *BaseUsecase) SetContext(c *gin.Context) *BaseUsecase {
	u.Ctx = c
	return u
}
