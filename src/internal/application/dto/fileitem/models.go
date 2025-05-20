package fileitem

// MediaUrlRequest represents a request for media URL with ordering information
type MediaUrlRequest struct {
	ID    *int64 `json:"id" nameFa:"شناسه" validate:"required"`
	Order *int   `json:"order" nameFa:"ترتیب" validate:"required"`
}
