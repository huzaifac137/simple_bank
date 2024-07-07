package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GRPCUnaryInterceptor() *grpc.ServerOption {
	interceptorOpt := grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		startTime := time.Now()

		log.Info().Msg("Unary GRPC Request Detected")
		result, err := handler(ctx, req)
		timeTaken := time.Since(startTime)

		if err != nil {
			errStatus, ok := status.FromError(err)
			if ok {
				log.Error().Err(err).Str("method", info.FullMethod).Int("status_code", int(errStatus.Code())).Str("status_message", (errStatus.Code().String())).Str("duration", timeTaken.String()+"ms").Msg(err.Error())

			}
			return nil, err
		}
		resStatus, ok := status.FromError(err)
		if ok {

			log.Info().Str("method", info.FullMethod).Int("status_code", int(resStatus.Code())).Str("status_message", resStatus.Code().String()).Dur("duration", timeTaken).Msg("Unary GRPC Request Completed")
		}

		return result, nil
	})

	return &interceptorOpt
}

type ResponseRecoder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *ResponseRecoder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecoder{w, http.StatusOK}
		handler.ServeHTTP(rec, r)
		timeTaken := time.Since(startTime)

		logger := log.Info()

		if rec.StatusCode != http.StatusOK {
			logger = log.Error()
			logger.Str("method", r.Method).Str("url", r.URL.String()).Int("status_code", rec.StatusCode).Str("status_text", http.StatusText(rec.StatusCode)).Dur("duration", timeTaken).Msg("HTTP Request Failed")
		} else {
			logger.Str("method", r.Method).Str("url", r.URL.String()).Int("status_code", rec.StatusCode).Str("status_text", http.StatusText(rec.StatusCode)).Dur("duration", timeTaken).Msg("HTTP Request Completed")
		}
	})
}
