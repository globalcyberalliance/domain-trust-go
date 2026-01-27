// Code generated from internal models; DO NOT EDIT.
package model

const (
	AbsoluteMaxLimit     = 999999999
	DefaultMetadataLimit = 100
	MaxMetadataLimit     = 10000
)

// MetadataFilter represents a filter.
type MetadataFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`
}
