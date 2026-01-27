// Code generated from internal models; DO NOT EDIT.
package model

type Option struct {
	ID    string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type OptionFilter struct {
}

type OptionUpdate struct {
	Value *string `json:"value"`
}
