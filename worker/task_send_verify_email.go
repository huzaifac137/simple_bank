package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/huzaifac137/simple_bank/util"
	"github.com/rs/zerolog/log"
)

const (
	SendVerifyEmailTask = "task:send_verify_email"
)

type SendVerifyEmailPayload struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeSendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload, opts ...asynq.Option) error {
	taskPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(SendVerifyEmailTask, taskPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("task_type", info.Type).Str("task_queue", info.Queue).Bytes("payload", info.Payload).Int("max_retires", info.MaxRetry).Msg("task enqueued successfully")

	return nil

}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	payload := &SendVerifyEmailPayload{}
	if err := json.Unmarshal(task.Payload(), payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// search for user in the database
	user, err := processor.store.GetUser(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Send verification email to the user
	er := util.SendEmail(user.Email)
	if er != nil {
		log.Error().Err(er).Str("user_email", user.Email).Msg("failed to send verification email")
	}
	log.Info().Str("task_type", task.Type()).Bytes("task_payload", task.Payload()).Str("user_email", user.Email).Msg("processed the task")

	return nil

}
