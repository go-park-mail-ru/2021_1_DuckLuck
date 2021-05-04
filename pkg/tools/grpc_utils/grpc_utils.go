package grpc_utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/metrics"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"

	"github.com/lithammer/shortuuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	RequireIdKey = "require_key"
)

func AccessInterceptor(metric *metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		md, _ := metadata.FromIncomingContext(ctx)

		statusCode := "200"
		requireId := shortuuid.New()
		logger.GrpcAccessLogStart(info.FullMethod, requireId,
			fmt.Sprintf("%v", req), fmt.Sprintf("%v", md))

		newContext := context.WithValue(ctx, RequireIdKey, requireId)
		reply, err := handler(newContext, req)

		if err != nil {
			statusCode = "500"
			metric.Errors.WithLabelValues(statusCode, info.FullMethod, info.FullMethod).Inc()
		} else {
			metric.AccessHits.WithLabelValues(statusCode, info.FullMethod, info.FullMethod).Inc()
		}

		metric.Durations.WithLabelValues(statusCode, info.FullMethod,
			info.FullMethod).Observe(time.Since(start).Seconds())
		logger.GrpcAccessLogEnd(info.FullMethod, requireId,
			fmt.Sprintf("%v", reply), start)
		metric.TotalHits.Inc()

		return reply, err
	}
}

func MustGetRequireId(ctx context.Context) string {
	requireId, ok := ctx.Value(RequireIdKey).(string)
	if !ok {
		panic("Require id not found")
	}

	return requireId
}
