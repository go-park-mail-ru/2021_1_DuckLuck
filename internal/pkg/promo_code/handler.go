package promo_code

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/promo_code Handler

type Handler interface {
	ApplyPromoCodeToOrder(w http.ResponseWriter, r *http.Request)
}
