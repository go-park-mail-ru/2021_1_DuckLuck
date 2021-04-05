package csrf_token

import "net/http"

type Handler interface {
	GetCsrfToken(w http.ResponseWriter, r *http.Request)
}
