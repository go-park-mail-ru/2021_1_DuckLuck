package models

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"

type UpdateOrder struct {
	OrderId uint64 `json:"order_id"`
	Status  string `json:"status"`
}

func (u *UpdateOrder) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	u.Status = sanitizer.Sanitize(u.Status)
}
