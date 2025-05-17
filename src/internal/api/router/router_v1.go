package router

import (
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/gin-gonic/gin"
)

type RouterV1 struct {
	h *bootstrap.HandlerManager
}

func (v RouterV1) Routes(route *gin.RouterGroup) {
	route.GET("/search", v.h.ProductHandlerV1.ProductList)
	route.GET("/search", v.h.VendorHandlerV1.VendorList)
}
