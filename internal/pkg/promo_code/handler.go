package promo_code

import "net/http"

type Handler interface {
	ApplyPromoCodeToOrder(w http.ResponseWriter, r *http.Request)
}