package producttransformer

import (
	"github.com/gin-gonic/gin"
	productdto "go-boilerplate/src/internal/app/dto/product"

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
