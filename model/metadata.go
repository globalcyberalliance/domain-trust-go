// Code generated from internal models; DO NOT EDIT.
package model

const (
	DefaultMetadataLimit = 100
	MaxMetadataLimit     = 10000
)

type (
	// MetadataFilter represents a filter.
	MetadataFilter struct {
		Limit int `json:"limit,omitempty" query:"limit"`
	}
)
