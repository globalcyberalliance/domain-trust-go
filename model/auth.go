// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type (
	Login struct {
		Email    string `json:"email" yaml:"email"`
		Password string `json:"password" yaml:"password"`
	}

	PasswordResetToken struct {
		Created time.Time `json:"created" yaml:"created"`
		Token   string    `json:"token" yaml:"token"`
		UserID  string    `json:"userID" yaml:"userID"`
	}

	PasswordResetTokenFilter struct {
		MetadataFilter

		Token string `json:"token"`
	}
)
