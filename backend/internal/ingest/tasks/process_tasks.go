package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/ingest/logic"
	"kessler/internal/objects/files"

	"github.com/charmbracelet/log"
	"github.com/hibiken/asynq"
)

func AsynqHandler(mux *asynq.ServeMux) {
	// existing file ingestion
	mux.HandleFunc("ingest:file", HandleIngestNewFileTask)
	// new case ingestion
	mux.HandleFunc(TypeIngestCase, HandleIngestCaseTask)
}

func HandleIngestNewFileTask(ctx context.Context, task *asynq.Task) error {
	var fileObj files.CompleteFileSchema
	if err := json.Unmarshal(task.Payload(), &fileObj); err != nil {
		return err
	}
	file, err := logic.CompleteIngestFileFromAttachmentUrls(ctx, &fileObj)
	if err != nil {
		return err
	}
	log.Info("File added to DB", "file", file)
	return nil
}

// HandleIngestCaseTask processes a case ingestion task
func HandleIngestCaseTask(ctx context.Context, task *asynq.Task) error {
	var caseInfo OpenscrapersCaseInfoPayload
	if err := json.Unmarshal(task.Payload(), &caseInfo); err != nil {
		return fmt.Errorf("failed to unmarshal case payload: %w", err)
	}
	// invoke business logic to persist case and filings
	err := IngestOpenscrapersCase(ctx, &caseInfo)
	if err != nil {
		return fmt.Errorf("error ingesting case: %w", err)
	}
	log.Info("Case ingested successfully", "case_number", caseInfo.CaseNumber)
	return nil
}
