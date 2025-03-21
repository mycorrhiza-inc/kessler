package tasks

import (
	"context"
	"encoding/json"
<<<<<<< HEAD:backend/internal/ingest/tasks/process_tasks.go
	"fmt"
	"kessler/internal/ingest/logic"
=======
	"kessler/cmd/ingest/logic"
>>>>>>> main:backend/cmd/ingest/tasks/process_tasks.go
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
