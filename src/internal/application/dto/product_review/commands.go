package product_review

// CreateProductReviewCommand represents a command to create a new product product_review
type CreateProductReviewCommand struct {
	Rating     *int    `json:"rating" validate:"required,min=1,max=5" error:"required=امتیاز الزامی است|min=امتیاز باید حداقل 1 باشد|max=امتیاز نمی‌تواند بیشتر از 5 باشد"`
	Like       *int    `json:"like" validate:"required,min=0" error:"required=تعداد لایک الزامی است|min=تعداد لایک نمی‌تواند کمتر از 0 باشد"`
	Dislike    *int    `json:"dislike" validate:"required,min=0" error:"required=تعداد دیسلایک الزامی است|min=تعداد دیسلایک نمی‌تواند کمتر از 0 باشد"`
	Approved   *bool   `json:"approved" validate:"required" error:"required=وضعیت تایید نظر الزامی است"`
	ReviewText *string `json:"reviewText" validate:"required,max=2000" error:"required=متن نظر الزامی است|max=متن نظر نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	ProductID  *int64  `json:"productId" validate:"required" error:"required=محصول الزامی است"`
	SiteID     *int64  `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// UpdateProductReviewCommand represents a command to update an existing product product_review
type UpdateProductReviewCommand struct {
	ID         *int64  `json:"id" validate:"required" error:"required=نظر الزامی است"`
	Rating     *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5" error:"min=امتیاز باید بین 1 تا 5 باشد|max=امتیاز باید بین 1 تا 5 باشد"`
	Like       *int    `json:"like,omitempty" validate:"omitempty,min=0" error:"min=تعداد لایک نمی‌تواند کمتر از 0 باشد"`
	Dislike    *int    `json:"dislike,omitempty" validate:"omitempty,min=0" error:"min=تعداد دیسلایک نمی‌تواند کمتر از 0 باشد"`
	Approved   *bool   `json:"approved,omitempty" validate:"omitempty" error:""`
	ReviewText *string `json:"reviewText,omitempty" validate:"omitempty,max=2000" error:"max=متن نظر نباید بیشتر از 2000 کاراکتر باشد"`
	ProductID  *int64  `json:"productId" validate:"required" error:"required=محصول الزامی است"`
	SiteID     *int64  `json:"siteId,omitempty" validate:"omitempty" error:""`
}

// DeleteProductReviewCommand represents a command to delete a product product_review
type DeleteProductReviewCommand struct {
	ID *int64 `json:"id" validate:"required" error:"required=نظر محصول الزامی است"`
}
