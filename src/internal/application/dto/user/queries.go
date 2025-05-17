package user

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllUserQuery represents a query to get all users with pagination for admin
type AdminGetAllUserQuery struct {
	common.PaginationRequestDto
}

// CalculatePlanPriceQuery represents a query to calculate plan price
type CalculatePlanPriceQuery struct {
	PlanID *int64 `json:"planId" validate:"required,gt=0" error:"required=پلن الزامی است|gt=شناسه پلن باید بزرگتر از 0 باشد"`
}

// GetProfileUserQuery represents a query to get user profile
type GetProfileUserQuery struct {
}

// VerifyUserQuery represents a query to verify user
type VerifyUserQuery struct {
	Email *string         `json:"email" validate:"required,email" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد"`
	Code  *int            `json:"code" validate:"required" error:"required=کد تایید الزامی است"`
	Type  *VerifyTypeEnum `json:"type" validate:"required" error:"required=نوع تایید الزامی است"`
}
