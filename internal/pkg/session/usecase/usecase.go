package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/pkg/errors"

	"net/http"
	"time"
)


type UseCase struct {
	SessionRepo session.Repository
}

func (h *UseCase) Check(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie(models.SessionCookieName)
	if err != nil {
		return nil, errors.Wrap(err, "don't find cookie = session_id")
	}

	sess, err := h.SessionRepo.GetByValue(sessionCookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "user unauthorized")
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, models.SessionContextKey, sess)
	return sess, nil
}

func (h *UseCase) Create(w http.ResponseWriter, r *http.Request, userId uint64) (*models.Session, error) {
	sess := models.NewSession(userId)
	err := h.SessionRepo.Add(sess)
	if err != nil {
		return nil, errors.Wrap(err, "the session was not added")
	}

	cookie := &http.Cookie {
		Name: models.SessionCookieName,
		Value: sess.Value,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path: "/",
	}

	http.SetCookie(w, cookie)
	ctx := r.Context()
	ctx = context.WithValue(ctx, models.SessionContextKey, sess)
	return sess, nil
}

func (h *UseCase) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || sess == nil {
		return errors.New("no cookie in context")
	}

	err := h.SessionRepo.DestroyByValue(sess.Value)
	if err != nil {
		return errors.Wrap(err, "no destroy cookie")
	}

	cookie := http.Cookie{
		Name:    models.SessionCookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}