// Code generated from internal models; DO NOT EDIT.
package model

import "time"

const (
	SessionUserKey = "sessionUser"

	UserRoleAdmin  = "admin"
	UserRoleMember = "member"
	UserRoleTrial  = "trial"
)

type (
	User struct {
		Created time.Time `json:"created" yaml:"created"`
		ID      string    `json:"id,omitzero" yaml:"id,omitempty"`

		Email          string `json:"email" yaml:"email"`
		FirstName      string `json:"firstName" yaml:"firstName"`
		LastName       string `json:"lastName" yaml:"lastName"`
		OrganizationID string `json:"organizationID" yaml:"organizationID"`
		Password       string `json:"password,omitempty" yaml:"password,omitempty"`
		Role           string `json:"role" yaml:"role"`
	}

	// UserFilter represents a filter passed to FindUsers().
	UserFilter struct {
		MetadataFilter

		Email          string `json:"email,omitempty" query:"email"`
		FirstName      string `json:"firstName,omitempty" query:"firstName"`
		LastName       string `json:"lastName,omitempty" query:"lastName"`
		OrganizationID string `json:"organizationID,omitempty" yaml:"organizationID" query:"organizationID"`
		Role           string `json:"role,omitempty" query:"role"`
	}

	// UserUpdate represents a set of fields to be updated via UpdateUser().
	UserUpdate struct {
		Email *string `json:"email,omitempty"`

		FirstName      *string `json:"firstName,omitempty"`
		LastName       *string `json:"lastName,omitempty"`
		OrganizationID *string `json:"organizationID,omitempty"`
		Password       *string `json:"password,omitempty"`
		Role           *string `json:"role,omitempty"`
	}

	UserSession struct {
		*User
	}

	UserWithAPIKey struct {
		*User
		*APIKey
	}
)
