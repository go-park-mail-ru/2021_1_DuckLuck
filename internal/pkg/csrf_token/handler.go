package csrf_token

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token Handler

type Handler interface {
	GetCsrfToken(w http.ResponseWriter, r *http.Request)
}
