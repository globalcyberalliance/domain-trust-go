// Code generated from internal models; DO NOT EDIT.
package model

import "time"

const (
	OrganizationDefaultUserQuota = 3

	OrganizationRatingHighConfidence = "high-confidence" // A high-confidence provider, such as a police officer or judge.
	OrganizationRatingMedConfidence  = "med-confidence"  // A medium-confidence provider with high-scale SOC, such as Verizon, IBM, or BT.
	OrganizationRatingLowConfidence  = "low-confidence"  // A low-confidence provider, such as SpamHaus or DMARC@scale.
	OrganizationRatingPredictive     = "predictive"      // A predictive intelligence provider, such as before.ai or Opora.
	OrganizationRatingTrial          = "trial"

	OrganizationRoleICANN     = "icann"
	OrganizationRoleOther     = "other"
	OrganizationRoleRegistrar = "registrar"
	OrganizationRoleRegistry  = "registry"
	OrganizationRoleReseller  = "reseller"

	OrganizationStatusActive      = "active"
	OrganizationStatusDeactivated = "deactivated"
)

type (
	Organization struct {
		Created time.Time `json:"created,omitzero" yaml:"created,omitempty"`
		ID      string    `json:"id,omitempty" yaml:"id,omitempty"`

		Name      string `json:"name" yaml:"name"`
		Rating    string `json:"rating" yaml:"rating"`
		Role      string `json:"role" yaml:"role"`
		Status    string `json:"status,omitempty" yaml:"status,omitempty"`
		UserQuota int8   `json:"userQuota,omitempty" yaml:"userQuota,omitempty"`
	}

	OrganizationFilter struct {
		MetadataFilter

		Name   string `json:"name" query:"name"`
		Rating string `json:"rating" query:"rating"`
		Role   string `json:"role" query:"role"`
		Status string `json:"status" query:"status"`
	}

	OrganizationUpdate struct {
		Name      *string `json:"name,omitempty"`
		Rating    *string `json:"rating,omitempty"`
		Role      *string `json:"role,omitempty"`
		Status    *string `json:"status,omitempty"`
		UserQuota *int8   `json:"userQuota,omitempty"`
	}
)
