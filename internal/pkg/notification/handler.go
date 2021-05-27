package notification

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification Handler

type Handler interface {
	SubscribeUser(w http.ResponseWriter, r *http.Request)
	UnsubscribeUser(w http.ResponseWriter, r *http.Request)
	GetNotificationPublicKey(w http.ResponseWriter, r *http.Request)
}
