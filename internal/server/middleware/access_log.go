package middleware

import (
	"context"
	"net/http"
	"regexp"
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
			metric.ActualConnections.Inc()

			requireId := shortuuid.New()
			ctx := context.WithValue(r.Context(), models.RequireIdKey, requireId)
			r = r.WithContext(ctx)
			res := statusResponse{w, http.StatusOK}

			logger.HttpAccessLogStart(r.URL.Path, r.RemoteAddr, r.Method, requireId)
			startTime := time.Now()

			re := regexp.MustCompile(".+/?[a-z]+")
			url := re.FindStringSubmatch(r.URL.Path)[0]

			next.ServeHTTP(&res, r)
			metric.Durations.WithLabelValues(strconv.Itoa(res.status), r.Method,
				url).Observe(time.Since(startTime).Seconds())
			logger.HttpAccessLogEnd(url, r.RemoteAddr, r.Method, requireId, startTime)

			if res.status != http.StatusOK {
				metric.Errors.WithLabelValues(strconv.Itoa(res.status), r.Method, url).Inc()
			} else {
				metric.AccessHits.WithLabelValues(strconv.Itoa(res.status), r.Method, url).Inc()
			}

			metric.TotalHits.Inc()
			metric.ActualConnections.Desc()
		})
	}
}
