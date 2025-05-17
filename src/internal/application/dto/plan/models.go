package plan

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
)

// PlanDto represents a plan data transfer object
type PlanDto struct {
	ID             int64                 `json:"id"`
	Name           string                `json:"name"`
	Description    string                `json:"description,omitempty"`
	Price          int64                 `json:"price"`
	DiscountType   user.DiscountTypeEnum `json:"discountType,omitempty"`
	Discount       int64                 `json:"discount,omitempty"`
	Duration       int                   `json:"duration"`
	Feature        string                `json:"feature,omitempty"`
	SmsCredits     int                   `json:"smsCredits,omitempty"`
	EmailCredits   int                   `json:"emailCredits,omitempty"`
	StorageCredits int                   `json:"storageCredits,omitempty"`
	AiCredits      int                   `json:"aiCredits,omitempty"`
	AiImageCredits int                   `json:"aiImageCredits,omitempty"`
	Roles          []string              `json:"roles,omitempty"`
	NewPrice       int64                 `json:"newPrice,omitempty"`
	NewDiscount    int64                 `json:"newDiscount,omitempty"`
}
