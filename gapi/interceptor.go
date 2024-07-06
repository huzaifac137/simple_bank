package gapi

import (
	"context"
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
