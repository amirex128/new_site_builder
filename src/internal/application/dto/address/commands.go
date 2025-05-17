package address

// CreateAddressCommand represents a command to create a new address
type CreateAddressCommand struct {
	CustomerID  *int64   `json:"customerId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه مشتری باید بزرگتر از 0 باشد"`
	UserID      *int64   `json:"userId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه کاربر باید بزرگتر از 0 باشد"`
	Title       *string  `json:"title" validate:"required,max=100" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Latitude    *float32 `json:"latitude" validate:"required,min=-90,max=90" error:"required=عرض جغرافیایی الزامی است|min=عرض جغرافیایی نمی‌تواند کمتر از -90 باشد|max=عرض جغرافیایی نمی‌تواند بیشتر از 90 باشد"`
	Longitude   *float32 `json:"longitude" validate:"required,min=-180,max=180" error:"required=طول جغرافیایی الزامی است|min=طول جغرافیایی نمی‌تواند کمتر از -180 باشد|max=طول جغرافیایی نمی‌تواند بیشتر از 180 باشد"`
	AddressLine *string  `json:"addressLine" validate:"required,max=200" error:"required=آدرس الزامی است|max=آدرس نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	PostalCode  *string  `json:"postalCode" validate:"required,max=20" error:"required=کد پستی الزامی است|max=کد پستی نمی‌تواند بیشتر از 20 کاراکتر باشد"`
	CityID      *int64   `json:"cityId" validate:"required,gt=0" error:"required=شناسه شهر الزامی است|gt=شناسه شهر باید بزرگتر از 0 باشد"`
	ProvinceID  *int64   `json:"provinceId" validate:"required,gt=0" error:"required=شناسه استان الزامی است|gt=شناسه استان باید بزرگتر از 0 باشد"`
}

// UpdateAddressCommand represents a command to update an existing address
type UpdateAddressCommand struct {
	ID          *int64   `json:"id" validate:"required,gt=0" error:"required=شناسه آدرس الزامی است|gt=شناسه آدرس باید بزرگتر از 0 باشد"`
	Title       *string  `json:"title,omitempty" validate:"omitempty,max=100" error:"max=عنوان نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Latitude    *float32 `json:"latitude,omitempty" validate:"omitempty,min=-90,max=90" error:"min=عرض جغرافیایی نمی‌تواند کمتر از -90 باشد|max=عرض جغرافیایی نمی‌تواند بیشتر از 90 باشد"`
	Longitude   *float32 `json:"longitude,omitempty" validate:"omitempty,min=-180,max=180" error:"min=طول جغرافیایی نمی‌تواند کمتر از -180 باشد|max=طول جغرافیایی نمی‌تواند بیشتر از 180 باشد"`
	AddressLine *string  `json:"addressLine,omitempty" validate:"omitempty,max=500" error:"max=آدرس نمی‌تواند بیشتر از 500 کاراکتر باشد"`
	PostalCode  *string  `json:"postalCode,omitempty" validate:"omitempty,pattern=^\\d{5}(-\\d{4})?$" error:"pattern=کد پستی نامعتبر است"`
	CityID      *int64   `json:"cityId" validate:"required,gt=0" error:"required=شناسه شهر الزامی است|gt=شناسه شهر باید بزرگتر از 0 باشد"`
	ProvinceID  *int64   `json:"provinceId" validate:"required,gt=0" error:"required=شناسه استان الزامی است|gt=شناسه استان باید بزرگتر از 0 باشد"`
}

// DeleteAddressCommand represents a command to delete an address
type DeleteAddressCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه آدرس الزامی است|gt=شناسه آدرس باید بزرگتر از 0 باشد"`
}
