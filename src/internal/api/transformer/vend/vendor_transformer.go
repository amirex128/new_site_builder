package vendtransformer

import (
	"github.com/gin-gonic/gin"
	vendorrdto "go-boilerplate/src/internal/api/dto/vend"
	vendordto "go-boilerplate/src/internal/app/dto/vend"
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
