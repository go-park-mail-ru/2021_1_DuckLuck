package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/metrics"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"

	"github.com/lithammer/shortuuid"
)

type statusResponse struct {
	http.ResponseWriter
	status int
}

func (r *statusResponse) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func AccessLog(metric *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requireId := shortuuid.New()
			ctx := context.WithValue(r.Context(), models.RequireIdKey, requireId)
			r = r.WithContext(ctx)
			res := statusResponse{w, 200}

			startTime := time.Now()
			logger.HttpAccessLogStart(r.URL.Path, r.RemoteAddr, r.Method, requireId)
			next.ServeHTTP(&res, r)
			metric.Durations.WithLabelValues(strconv.Itoa(res.status), r.Method,
				r.URL.Path).Observe(time.Since(startTime).Seconds())
			logger.HttpAccessLogEnd(r.URL.Path, r.RemoteAddr, r.Method, requireId, startTime)

			if res.status != http.StatusOK {
				metric.Errors.WithLabelValues(strconv.Itoa(res.status), r.Method, r.URL.Path).Inc()
			} else {
				metric.AccessHits.WithLabelValues(strconv.Itoa(res.status), r.Method, r.URL.Path).Inc()
			}
			metric.TotalHits.Inc()
		})
	}
}
