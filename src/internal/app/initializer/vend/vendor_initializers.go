package vendorinitializer

import (
	vendordto "go-boilerplate/src/internal/app/dto/vend"
	"go-boilerplate/src/internal/contract/repository"
)

type builder struct {
	simpleRepo repository.ISampleRepository
	params     *vendordto.VendorDto
}

func NewVendorInitializer(simpleRepo repository.ISampleRepository, dto *vendordto.VendorDto) *builder {
	return &builder{
		simpleRepo: simpleRepo,
		params:     dto,
	}
}

func (b *builder) Build() *vendordto.VendorDto {
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
