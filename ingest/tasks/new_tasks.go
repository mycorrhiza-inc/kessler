// tasks/tasks.go
package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeAddFileScraper      = "task:add_file_scraper"
	TypeProcessExistingFile = "task:process_existing_file"
)

type ProcessFilePayload struct {
	FileID         uuid.UUID         `json:"file_id"`
	DocumentStatus DocumentStatus    `json:"status"`
	Metadata       map[string]string `json:"metadata"`
}

type DocumentStatus struct {
	SkipProcessing bool `json:"skip_processing"`
}

func NewAddFileScraperTask(payload ScraperInfoPayload) (*asynq.Task, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	return asynq.NewTask(TypeAddFileScraper, p), nil
}

func NewProcessFileTask(payload ProcessFilePayload) (*asynq.Task, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	return asynq.NewTask(TypeProcessExistingFile, p, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

type contextKey string

const clientKey = contextKey("asynqClient")

func WithClient(ctx context.Context, client *asynq.Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}

func GetClient(ctx context.Context) *asynq.Client {
	return ctx.Value(clientKey).(*asynq.Client)
}
