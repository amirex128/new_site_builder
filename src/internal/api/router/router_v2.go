package router

import (
	"github.com/gin-gonic/gin"
	"go-boilerplate/src/bootstrap"
)

type RouterV2 struct {
	h *bootstrap.HandlerManager
}

func (v RouterV2) Routes(route *gin.RouterGroup) {
	//route.GET("/search", v.h.SampleHandlerV2.Simple)
}
