package productinitializer

import (
	productdto "go-boilerplate/src/internal/app/dto/product"

	"go-boilerplate/src/internal/contract/repository"
)

type builder struct {
	params     *productdto.ProductDto
	simpleRepo repository.ISampleRepository
}

func NewProductInitializer(simpleRepo repository.ISampleRepository, dto *productdto.ProductDto) *builder {
	return &builder{
		params:     dto,
		simpleRepo: simpleRepo,
	}
}

func (b *builder) Build() *productdto.ProductDto {
	return b.params
}

func (b *builder) SetDelivery() *builder {

	return b
}
func (b *builder) SetPickUp() *builder {

	return b
}
func (b *builder) SetCPC() *builder {

	return b
}
