package gapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mutex   sync.Mutex
	clients = make(map[string]*Client)
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

func HTTPLoggerRateLimiter(handlerFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {

	go func() {
		time.Sleep(1 * time.Minute)
		mutex.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mutex.Unlock()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Print("HTTPLoggerRateLimiter 1   !!!!!!")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get IP address from request")

		}

		mutex.Lock()
		fmt.Print("HTTPLoggerRateLimiter 2  !!!!!!")
		if _, ok := clients[ip]; !ok {
			clients[ip] = &Client{
				limiter: rate.NewLimiter(1, 6),
			}
		}

		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mutex.Unlock()
			log.Error().Msg("Too many requests, rate limiter!")
			http.Error(w, "Too many requests, please try again later.", http.StatusTooManyRequests)
			return
		}
		mutex.Unlock()

		startTime := time.Now()
		rec := &ResponseRecoder{w, http.StatusOK}
		handlerFunc(rec, r)
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
