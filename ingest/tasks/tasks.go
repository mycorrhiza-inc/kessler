// tasks/tasks.go
package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeAddFileScraper      = "task:add_file_scraper"
	TypeProcessExistingFile = "task:process_existing_file"
)

type ScraperInfoPayload struct {
	FileURL            string    `json:"file_url"`
	FileType           string    `json:"file_type"`
	DocketID           string    `json:"docket_id"`
	PublishedDate      time.Time `json:"published_date"`
	Name               string    `json:"name"`
	InternalSourceName string    `json:"internal_source_name"`
	FileClass          string    `json:"file_class"`
}

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

func HandleAddFileScraperTask(ctx context.Context, t *asynq.Task) error {
	var p ScraperInfoPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	log.Printf("Processing file from scraper: %s", p.FileURL)

	// Generate UUID for the file
	fileID := uuid.New()

	// Create process file task
	processTask, err := NewProcessFileTask(ProcessFilePayload{
		FileID:   fileID,
		Metadata: make(map[string]string),
	})
	if err != nil {
		return fmt.Errorf("failed to create process task: %v", err)
	}

	// Get client from context
	client := GetClient(ctx)
	if _, err := client.Enqueue(processTask); err != nil {
		return fmt.Errorf("failed to enqueue process task: %v", err)
	}

	return nil
}

func HandleProcessFileTask(ctx context.Context, t *asynq.Task) error {
	var p ProcessFilePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	if p.DocumentStatus.SkipProcessing {
		log.Printf("Skipping processing for file: %s", p.FileID)
		return nil
	}

	log.Printf("Processing file: %s", p.FileID)
	// Add your actual processing logic here

	return nil
}

type contextKey string

const clientKey = contextKey("asynqClient")

func WithClient(ctx context.Context, client *asynq.Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}

func GetClient(ctx context.Context) *asynq.Client {
	return ctx.Value(clientKey).(*asynq.Client)
}
