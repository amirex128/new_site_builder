package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type PlanRepo struct {
	database *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepo {
	return &PlanRepo{
		database: db,
	}
}

func (r *PlanRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Plan], error) {
	var plans []domain.Plan
	var count int64

	query := r.database.Model(&domain.Plan{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&plans)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(plans, paginationRequestDto, count)
}

func (r *PlanRepo) GetByID(id int64) (*domain.Plan, error) {
	var plan *domain.Plan
	result := r.database.First(&plan, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return plan, nil
}

func (r *PlanRepo) Create(plan *domain.Plan) error {
	result := r.database.Create(plan)
	return result.Error
}

func (r *PlanRepo) Update(plan *domain.Plan) error {
	result := r.database.Save(plan)
	return result.Error
}

func (r *PlanRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Plan{}, id)
	return result.Error
}
