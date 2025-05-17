package product

import "time"

// CouponCommand represents a coupon for a product
type CouponCommand struct {
	Quantity   *int              `json:"quantity" validate:"required" error:"required=تعداد کوپن الزامی است"`
	Type       *DiscountTypeEnum `json:"type" validate:"required" error:"required=نوع تخفیف الزامی است"`
	Value      *int64            `json:"value" validate:"required" error:"required=مقدار تخفیف الزامی است"`
	ExpiryDate *time.Time        `json:"expiryDate" validate:"required,gtfield=time.Now" error:"required=تاریخ انقضا الزامی است|gtfield=تاریخ انقضا باید در آینده باشد"`
}

// ProductAttributeCommand represents a product attribute
type ProductAttributeCommand struct {
	ID    *int64                    `json:"id,omitempty" validate:"omitempty" error:""`
	Type  *ProductAttributeTypeEnum `json:"type" validate:"required" error:"required=نوع ویژگی الزامی است"`
	Name  *string                   `json:"name" validate:"required,max=100" error:"required=نام ویژگی الزامی است|max=نام ویژگی نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Value *string                   `json:"value" validate:"required,max=500" error:"required=مقدار ویژگی الزامی است|max=مقدار ویژگی نمی‌تواند بیشتر از 500 کاراکتر باشد"`
}

// ProductVariantCommand represents a product variant
type ProductVariantCommand struct {
	ID    *int64  `json:"id,omitempty" validate:"omitempty" error:""`
	Name  *string `json:"name" validate:"required,max=200" error:"required=نام تنوع محصول الزامی است|max=نام تنوع محصول نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Price *int64  `json:"price" validate:"required" error:"required=قیمت الزامی است"`
	Stock *int    `json:"stock" validate:"required" error:"required=موجودی الزامی است"`
}

// OrderBasketItem represents an item in an order basket for price calculation
type OrderBasketItem struct {
	BasketItemID       *int64 `json:"basketItemId" validate:"required" error:"required=شناسه آیتم سبد خرید الزامی است"`
	Quantity           *int   `json:"quantity" validate:"required,max=1000" error:"required=تعداد الزامی است|max=تعداد نمی‌تواند بیشتر از 1000 باشد"`
	ProductID          *int64 `json:"productId" validate:"required" error:"required=محصول الزامی است"`
	ProductVariationID *int64 `json:"productVariationId" validate:"required" error:"required=تنوع محصول الزامی است"`
}
