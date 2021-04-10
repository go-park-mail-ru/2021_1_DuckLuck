package models

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"
)

// All order information
// This models saved in database
type Order struct {
	Recipient OrderRecipient `json:"recipient" valid:"notnull"`
	Address   OrderAddress   `json:"address" valid:"notnull"`
}

func (o *Order) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	o.Recipient.FirstName = sanitizer.Sanitize(o.Recipient.FirstName)
	o.Recipient.LastName = sanitizer.Sanitize(o.Recipient.LastName)
	o.Recipient.Email = sanitizer.Sanitize(o.Recipient.Email)
	o.Address.Address = sanitizer.Sanitize(o.Address.Address)
}

// All order information
// This model preview info for user
type PreviewOrder struct {
	Products  []*PreviewCartArticle `json:"products" valid:"notnull"`
	Recipient OrderRecipient        `json:"recipient" valid:"notnull"`
	Price     TotalPrice            `json:"price" valid:"notnull"`
	Address   OrderAddress          `json:"address" valid:"notnull"`
}

// Order address for delivery
type OrderAddress struct {
	Address string `json:"address" valid:"utfletter, stringlength(3|30)"`
}

// Info about order recipient
type OrderRecipient struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(3|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(3|30)"`
	Email     string `json:"email" valid:"email"`
}
