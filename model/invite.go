// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type (
	Invite struct {
		Created time.Time `json:"created" yaml:"created"`
		ID      string    `json:"id,omitzero" yaml:"id,omitempty"`

		Token              string `json:"token" yaml:"token"`
		UserEmail          string `json:"userEmail" yaml:"userEmail"`
		UserFirstName      string `json:"userFirstName" yaml:"userFirstName"`
		UserLastName       string `json:"userLastName" yaml:"userLastName"`
		UserOrganizationID string `json:"userOrganizationID,omitempty" yaml:"userOrganizationID,omitempty"`
		UserRole           string `json:"userRole,omitempty" yaml:"userRole,omitempty"`
	}

	InviteFilter struct {
		MetadataFilter

		UserEmail          string `json:"userEmail"`
		UserOrganizationID string `json:"userOrganizationID"`
	}
)
