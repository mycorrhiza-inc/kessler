package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/ingest/logic"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func AsynqHandler(mux *asynq.ServeMux) {
	// existing file ingestion
	mux.HandleFunc("ingest:file", HandleIngestNewFileTask)
	// new case ingestion
	mux.HandleFunc(TypeIngestCase, HandleIngestCaseTask)
}

func HandleIngestNewFileTask(ctx context.Context, task *asynq.Task) error {
	var fileTaskObj FilingInfoPayload
	if err := json.Unmarshal(task.Payload(), &fileTaskObj); err != nil {
		return err
	}
	completeFileSchema := fileTaskObj.IntoCompleteFile()

	file, err := logic.CompleteIngestFileFromAttachmentUrls(ctx, &completeFileSchema)
	if err != nil {
		return err
	}
	log.Info("File added to DB", zap.String("name", file.Name))
	return nil
}

// HandleIngestCaseTask processes a case ingestion task
func HandleIngestCaseTask(ctx context.Context, task *asynq.Task) error {
	var caseInfo OpenscrapersCaseInfoPayload
	if err := json.Unmarshal(task.Payload(), &caseInfo); err != nil {
		return fmt.Errorf("failed to unmarshal case payload: %w", err)
	}
	// invoke business logic to persist case and filings
	err := IngestOpenscrapersCase(ctx, caseInfo)
	if err != nil {
		return fmt.Errorf("error ingesting case: %w", err)
	}
	log.Info("Case ingested successfully", zap.String("case_number", caseInfo.CaseNumber))
	return nil
}
