package user

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain/enums"
)

// AdminGetAllUserQuery represents a query to get all users with pagination for admin
type AdminGetAllUserQuery struct {
	common.PaginationRequestDto
}

// CalculatePlanPriceQuery represents a query to calculate plan price
type CalculatePlanPriceQuery struct {
	PlanID *int64 `json:"planId" validate:"required,gt=0" nameFa:"شناسه طرح"`
}

// GetProfileUserQuery represents a query to get user profile
type GetProfileUserQuery struct {
}

// VerifyUserQuery represents a query to verify user
type VerifyUserQuery struct {
	Email *string               `json:"email" form:"email" validate:"required,email" nameFa:"ایمیل"`
	Code  *int                  `json:"code" form:"code" validate:"required" nameFa:"کد"`
	Type  *enums.VerifyTypeEnum `json:"type" validate:"required,enum" nameFa:"نوع"`
}
