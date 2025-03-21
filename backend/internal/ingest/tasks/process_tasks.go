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
	mux.HandleFunc("ingest:file", HandleIngestNewFileTask)
}

func HandleIngestNewFileTask(ctx context.Context, task *asynq.Task) error {
	var fileObj files.CompleteFileSchema
	if err := json.Unmarshal(task.Payload(), &fileObj); err != nil {
		return err
	}
	file_URL, ok := fileObj.Mdata["url"].(string)
	if !ok {
		return fmt.Errorf("url field missing or not a string in file metadata")
	}
	file, err := logic.AddURLRaw(ctx, file_URL, &fileObj)
	if err != nil {
		return err
	}
	log.Info("File added to DB", "file", file)
	return nil
}
