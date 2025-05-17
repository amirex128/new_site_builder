package producttransformer

import (
	productdto "github.com/amirex128/new_site_builder/src/internal/app/dto/product"
	"github.com/gin-gonic/gin"

	"strconv"
)

type builder struct {
	params *productdto.ProductDto
	ctx    *gin.Context
}

func NewProductTransformer(ctx *gin.Context) *builder {
	return &builder{
		ctx: ctx,
	}
}

func (b *builder) Build() *productdto.ProductDto {
	return b.params
}

func (b *builder) SetSuperType() *builder {

	param := b.ctx.Param("superTypeId")

	atoi, err := strconv.Atoi(param)
	if err != nil {
		///
	}

	b.params.SuperTypeId = atoi

	return b
}
