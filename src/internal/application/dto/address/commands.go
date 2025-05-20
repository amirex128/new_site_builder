package address

// CreateAddressCommand represents a command to create a new address
type CreateAddressCommand struct {
	CustomerID  *int64   `json:"customerId,omitempty" validate:"omitempty,gt=0"`
	UserID      *int64   `json:"userId,omitempty" validate:"omitempty,gt=0"`
	Title       *string  `json:"title" validate:"required_text=1 100"`
	Latitude    *float32 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   *float32 `json:"longitude" validate:"required,min=-180,max=180"`
	AddressLine *string  `json:"addressLine" validate:"required_text=1 200"`
	PostalCode  *string  `json:"postalCode" validate:"required_text=1 20"`
	CityID      *int64   `json:"cityId" validate:"required,gt=0"`
	ProvinceID  *int64   `json:"provinceId" validate:"required,gt=0"`
}

// UpdateAddressCommand represents a command to update an existing address
type UpdateAddressCommand struct {
	ID          *int64   `json:"id" validate:"required,gt=0"`
	Title       *string  `json:"title,omitempty" validate:"optional_text=1 100"`
	Latitude    *float32 `json:"latitude,omitempty" validate:"omitempty,min=-90,max=90"`
	Longitude   *float32 `json:"longitude,omitempty" validate:"omitempty,min=-180,max=180"`
	AddressLine *string  `json:"addressLine,omitempty" validate:"optional_text=1 500"`
	PostalCode  *string  `json:"postalCode,omitempty" validate:"optional_text=1 20"`
	CityID      *int64   `json:"cityId" validate:"required,gt=0"`
	ProvinceID  *int64   `json:"provinceId" validate:"required,gt=0"`
}

// DeleteAddressCommand represents a command to delete an address
type DeleteAddressCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
