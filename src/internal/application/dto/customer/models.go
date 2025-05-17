package customer

// CustomerDto represents a customer data transfer object
type CustomerDto struct {
	ID           int64    `json:"id"`
	FirstName    string   `json:"firstName,omitempty"`
	LastName     string   `json:"lastName,omitempty"`
	Email        string   `json:"email,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	NationalCode string   `json:"nationalCode,omitempty"`
	SiteID       int64    `json:"siteId"`
	AvatarID     string   `json:"avatarId,omitempty"`
	AddressIDs   []int64  `json:"addressIds,omitempty"`
	Roles        []string `json:"roles,omitempty"`
}
