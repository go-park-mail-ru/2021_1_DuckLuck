package models

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"

type NotificationCredentials struct {
	UserIdentifier
	Keys NotificationKeys `json:"keys"`
}

func (nc *NotificationCredentials) Sanitize() {
	nc.UserIdentifier.Sanitize()
	nc.Keys.Sanitize()
}

type UserIdentifier struct {
	Endpoint string `json:"endpoint"`
}

func (ui *UserIdentifier) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	ui.Endpoint = sanitizer.Sanitize(ui.Endpoint)
}

type NotificationKeys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
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
	Status string      `json:"status"`
}

type Subscribes struct {
	Credentials map[string]*NotificationKeys `json:"subscribes" valid:"notnull"`
}
