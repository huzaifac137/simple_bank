package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeSendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

const (
	CriticalQueue = "critical"
	DefaultQueue  = "default"
)

func NewRedisTaskDistributor(opts asynq.RedisClientOpt) TaskDistributor {
	return &RedisTaskDistributor{
		client: asynq.NewClient(opts),
	}
}
