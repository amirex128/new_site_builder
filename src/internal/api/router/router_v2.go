package router

import (
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/gin-gonic/gin"
)

type RouterV2 struct {
	h *bootstrap.HandlerManager
}

func (v RouterV2) Routes(route *gin.RouterGroup) {
	//route.GET("/search", v.h.SampleHandlerV2.Simple)
}
