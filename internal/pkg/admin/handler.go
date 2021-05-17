package admin

import "net/http"

type Handler interface {
	ChangeOrderStatus(w http.ResponseWriter, r *http.Request)
}
