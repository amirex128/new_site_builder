package product_review

// CreateProductReviewCommand represents a command to create a new article product_review
type CreateProductReviewCommand struct {
	Rating     *int    `json:"rating" validate:"required,min=1,max=5"`
	Like       *int    `json:"like" validate:"required,min=0"`
	Dislike    *int    `json:"dislike" validate:"required,min=0"`
	Approved   *bool   `json:"approved" validate:"required_bool"`
	ReviewText *string `json:"reviewText" validate:"required_text=1 2000"`
	ProductID  *int64  `json:"productId" validate:"required"`
	SiteID     *int64  `json:"siteId" validate:"required"`
}

// UpdateProductReviewCommand represents a command to update an existing article product_review
type UpdateProductReviewCommand struct {
	ID         *int64  `json:"id" validate:"required"`
	Rating     *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Like       *int    `json:"like,omitempty" validate:"omitempty,min=0"`
	Dislike    *int    `json:"dislike,omitempty" validate:"omitempty,min=0"`
	Approved   *bool   `json:"approved,omitempty" validate:"optional_bool"`
	ReviewText *string `json:"reviewText,omitempty" validate:"optional_text=1 2000"`
	ProductID  *int64  `json:"productId" validate:"required"`
	SiteID     *int64  `json:"siteId,omitempty" validate:"omitempty"`
}

// DeleteProductReviewCommand represents a command to delete a article product_review
type DeleteProductReviewCommand struct {
	ID *int64 `json:"id" validate:"required"`
}
