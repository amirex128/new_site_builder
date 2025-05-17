package vendtransformer

import (
	vendorrdto "github.com/amirex128/new_site_builder/src/internal/api/dto/vend"
	vendordto "github.com/amirex128/new_site_builder/src/internal/app/dto/vend"
	"github.com/gin-gonic/gin"
	"strconv"
)

type builder struct {
	requestDto *vendorrdto.VendorRequestDto
	params     *vendordto.VendorDto
	ctx        *gin.Context
}

func NewVendorTransformer(ctx *gin.Context) *builder {
	return &builder{
		ctx:    ctx,
		params: &vendordto.VendorDto{},
	}
}

func (b *builder) Build() *vendordto.VendorDto {
	return b.params
}

func (b *builder) SetSuperType() *builder {

	param := b.ctx.Param("superTypeId")

	atoi, err := strconv.Atoi(param)
	if err != nil {
		///
	}

	b.requestDto.SuperTypeId = atoi

	return b
}
