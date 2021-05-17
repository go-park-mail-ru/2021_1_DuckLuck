package models

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"

type NotificationCredentials struct {
	Endpoint string `json:"endpoint"`
	Keys NotificationKeys `json:"keys"`
}

func (nc *NotificationCredentials) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	nc.Endpoint = sanitizer.Sanitize(nc.Endpoint)
	nc.Keys.Sanitize()
}

type NotificationKeys struct {
	Auth     string `json:"auth"`
	P256dh   string `json:"p256dh"`
}

func (nk *NotificationKeys) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	nk.Auth = sanitizer.Sanitize(nk.Auth)
	nk.P256dh = sanitizer.Sanitize(nk.P256dh)
}

type NotificationPublicKey struct {
	Key string `json:"key"`
}

type OrderNotification struct {
	Number OrderNumber `json:"order_number"`
	Status string `json:"status"`
}
