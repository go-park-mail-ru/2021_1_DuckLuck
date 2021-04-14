package session

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session Handler

type Handler interface {
	CheckSession(w http.ResponseWriter, r *http.Request)
}
