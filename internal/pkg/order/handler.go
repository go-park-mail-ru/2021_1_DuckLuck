package order

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order Handler

type Handler interface {
	GetOrderFromCart(w http.ResponseWriter, r *http.Request)
	AddCompletedOrder(w http.ResponseWriter, r *http.Request)
}
