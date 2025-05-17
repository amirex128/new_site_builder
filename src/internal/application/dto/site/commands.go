package site

// CreateSiteCommand represents a command to create a new site
type CreateSiteCommand struct {
	Domain     *string         `json:"domain" validate:"required,domain" error:"required=دامنه الزامی است|domain=دامنه نامعتبر است"`
	DomainType *DomainTypeEnum `json:"domainType" validate:"required" error:"required=نوع دامنه الزامی است"`
	Name       *string         `json:"name" validate:"required,max=200" error:"required=نام سایت الزامی است|max=نام سایت نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Status     *StatusEnum     `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	SiteType   *SiteTypeEnum   `json:"siteType" validate:"required" error:"required=نوع سایت الزامی است"`
}

// UpdateSiteCommand represents a command to update an existing site
type UpdateSiteCommand struct {
	ID         *int64          `json:"id" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Domain     *string         `json:"domain,omitempty" validate:"omitempty,domain" error:"domain=دامنه نامعتبر است"`
	DomainType *DomainTypeEnum `json:"domainType,omitempty" validate:"omitempty" error:""`
	Name       *string         `json:"name,omitempty" validate:"omitempty,max=200" error:"max=نام سایت نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Status     *StatusEnum     `json:"status,omitempty" validate:"omitempty" error:""`
	SiteType   *SiteTypeEnum   `json:"siteType,omitempty" validate:"omitempty" error:""`
}

// DeleteSiteCommand represents a command to delete a site
type DeleteSiteCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
}
