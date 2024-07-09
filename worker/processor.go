package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
)

var (
	queues = map[string]int{
		CriticalQueue: 10,
		DefaultQueue:  5,
	}
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{
		Queues: queues,
	})
	return &RedisTaskProcessor{server: server, store: store}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(SendVerifyEmailTask, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)

}
