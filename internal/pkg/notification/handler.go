package notification

import "net/http"

type Handler interface {
	SubscribeUser(w http.ResponseWriter, r *http.Request)
	UnsubscribeUser(w http.ResponseWriter, r *http.Request)
	GetNotificationPublicKey(w http.ResponseWriter, r *http.Request)
}
