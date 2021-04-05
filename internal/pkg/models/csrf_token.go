package models

import uuid "github.com/satori/go.uuid"

const (
	CsrfTokenHeaderName = "X-CSRF-TOKEN"
	ExpireCsrfToken     = 900
)

type CsrfToken struct {
	Value string `json:"token"`
}

func NewCsrfToken() *CsrfToken {
	newValue := uuid.NewV4()
	return &CsrfToken{
		Value: newValue.String(),
	}
}
