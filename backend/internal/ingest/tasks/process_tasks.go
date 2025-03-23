package tasks

import (
	"context"
	"encoding/json"
	"kessler/internal/ingest/logic"
	"kessler/internal/objects/files"

	"github.com/charmbracelet/log"
	"github.com/hibiken/asynq"
)

func AsynqHandler(mux *asynq.ServeMux) {
	mux.HandleFunc("ingest:file", HandleIngestNewFileTask)
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
