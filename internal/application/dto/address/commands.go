package address

// CreateAddressCommand represents a command to create a new address
type CreateAddressCommand struct {
	Title       *string  `json:"title" nameFa:"عنوان" validate:"required_text=1 100"  errors:"required_text=this field is required,min=this field must be at least 1 character,max=this field must be at most 100 characters,min=this field must be at least 1 character,max=this field must be at most 100 characters"`
	Latitude    *float32 `json:"latitude" nameFa:"عرض جغرافیایی" validate:"required,min=-90,max=90"`
	Longitude   *float32 `json:"longitude" nameFa:"طول جغرافیایی" validate:"required,min=-180,max=180"`
	AddressLine *string  `json:"addressLine" nameFa:"آدرس" validate:"required_text=1 200"`
	PostalCode  *string  `json:"postalCode" nameFa:"کد پستی" validate:"required_text=1 20"`
	CityID      *int64   `json:"cityId" nameFa:"شناسه شهر" validate:"required,gt=0"`
	ProvinceID  *int64   `json:"provinceId" nameFa:"شناسه استان" validate:"required,gt=0"`
}

// UpdateAddressCommand represents a command to update an existing address
type UpdateAddressCommand struct {
	ID          *int64   `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
	Title       *string  `json:"title,omitempty" nameFa:"عنوان" validate:"optional_text=1 100"`
	Latitude    *float32 `json:"latitude,omitempty" nameFa:"عرض جغرافیایی" validate:"omitempty,min=-90,max=90"`
	Longitude   *float32 `json:"longitude,omitempty" nameFa:"طول جغرافیایی" validate:"omitempty,min=-180,max=180"`
	AddressLine *string  `json:"addressLine,omitempty" nameFa:"آدرس" validate:"optional_text=1 500"`
	PostalCode  *string  `json:"postalCode,omitempty" nameFa:"کد پستی" validate:"optional_text=1 20"`
	CityID      *int64   `json:"cityId" nameFa:"شناسه شهر" validate:"required,gt=0"`
	ProvinceID  *int64   `json:"provinceId" nameFa:"شناسه استان" validate:"required,gt=0"`
}

// DeleteAddressCommand represents a command to delete an address
type DeleteAddressCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
}
