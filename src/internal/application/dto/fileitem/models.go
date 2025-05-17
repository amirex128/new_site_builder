package fileitem

// MediaUrlRequest represents a request for media URL with ordering information
type MediaUrlRequest struct {
	ID    *int64 `json:"id" validate:"required"`
	Order *int   `json:"order" validate:"required"`
}
