package site

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateSiteCommand represents a command to create a new site
type CreateSiteCommand struct {
	Domain     *string               `json:"domain" validate:"required,domain" nameFa:"آدرس سایت"`
	DomainType *enums.DomainTypeEnum `json:"domainType" validate:"required,enum" nameFa:"نوع آدرس"`
	Name       *string               `json:"name" validate:"required_text=1 200" nameFa:"نام سایت"`
	Status     *enums.StatusEnum     `json:"status" validate:"required,enum" nameFa:"وضعیت"`
	SiteType   *enums.SiteTypeEnum   `json:"siteType" validate:"required,enum" nameFa:"نوع سایت"`
}

// UpdateSiteCommand represents a command to update an existing site
type UpdateSiteCommand struct {
	ID         *int64                `json:"id" validate:"required,gt=0" nameFa:"شناسه سایت"`
	Domain     *string               `json:"domain,omitempty" validate:"domain_optional" nameFa:"آدرس سایت"`
	DomainType *enums.DomainTypeEnum `json:"domainType,omitempty" validate:"enum_optional" nameFa:"نوع آدرس"`
	Name       *string               `json:"name,omitempty" validate:"optional_text=1 200" nameFa:"نام سایت"`
	Status     *enums.StatusEnum     `json:"status,omitempty" validate:"enum_optional" nameFa:"وضعیت"`
	SiteType   *enums.SiteTypeEnum   `json:"siteType,omitempty" validate:"enum_optional" nameFa:"نوع سایت"`
}

// DeleteSiteCommand represents a command to delete a site
type DeleteSiteCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" nameFa:"شناسه سایت"`
}
