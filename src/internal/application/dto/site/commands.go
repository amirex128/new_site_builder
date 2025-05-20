package site

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateSiteCommand represents a command to create a new site
type CreateSiteCommand struct {
	Domain     *string               `json:"domain" validate:"required,domain"`
	DomainType *enums.DomainTypeEnum `json:"domainType" validate:"required,enum"`
	Name       *string               `json:"name" validate:"required_text=1 200"`
	Status     *enums.StatusEnum     `json:"status" validate:"required,enum"`
	SiteType   *enums.SiteTypeEnum   `json:"siteType" validate:"required,enum"`
}

// UpdateSiteCommand represents a command to update an existing site
type UpdateSiteCommand struct {
	ID         *int64                `json:"id" validate:"required,gt=0"`
	Domain     *string               `json:"domain,omitempty" validate:"domain_optional"`
	DomainType *enums.DomainTypeEnum `json:"domainType,omitempty" validate:"enum_optional"`
	Name       *string               `json:"name,omitempty" validate:"optional_text=1 200"`
	Status     *enums.StatusEnum     `json:"status,omitempty" validate:"enum_optional"`
	SiteType   *enums.SiteTypeEnum   `json:"siteType,omitempty" validate:"enum_optional"`
}

// DeleteSiteCommand represents a command to delete a site
type DeleteSiteCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
