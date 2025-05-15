package logic

import (
	"context"
	"errors"
	"fmt"
	"kessler/internal/ingest/validators"
	"kessler/internal/objects/files"
	"kessler/pkg/logger"
	"kessler/pkg/s3utils"
	"log/slog"
	"os"

	"go.uber.org/zap"
	// Assume these are implemented in other packages
)

var OS_HASH_FILEDIR = os.Getenv("OS_HASH_FILEDIR")

func ProcessFile(ctx context.Context, complete_file files.CompleteFileSchema) error {
	log := logger.GetLogger("process_file")
	file, err := ProcessFileRaw(ctx, &complete_file, files.DocStatusCompleted)
	if err != nil {
		log.Warn("Encountered error processing file", zap.String("name", complete_file.Name), zap.Error(err))
	}
	file, err = upsertFullFileToDB(ctx, complete_file, DatabaseInteractionInsert)
	if err != nil {
		log.Error("Could not upload file to database", zap.String("name", file.Name), zap.Error(err))
		return fmt.Errorf("Could not upload file to db")
	}
	return nil
}

func ProcessFileRaw(ctx context.Context, obj *files.CompleteFileSchema, stopAt files.DocProcStatus) (files.CompleteFileSchema, error) {
	if obj == nil {
		return files.CompleteFileSchema{}, nil
	}
	logger := slog.Default()

	// if obj.Lang == "" {
	// 	return *obj, errors.New("language is required")
	// }

	currentStage := obj.Stage.DocProcStatus

	texts := make(map[string]string)

	for i := 0; i < 1000; i++ {
		if currentStage.Index() >= stopAt.Index() {
			obj.Stage = files.DocProcStage{
				PGStage:       files.PGStageCompleted,
				IsCompleted:   true,
				DocProcStatus: currentStage,
			}
			return *obj, nil
		}

		var nextStage files.DocProcStatus
		var err error

		switch currentStage {
		case files.DocStatusUnprocessed:
			nextStage, err = processStageHandleExtension(ctx, obj)
		case files.DocStatusBeginProcessing:
			nextStage, err = processGenerateRawText(ctx, obj, texts)
		case files.DocStatusRawTextCompleted:
			nextStage, err = processTranslateRawText(ctx, obj, texts)
		case files.DocStatusTextCompleted:
			nextStage, err = createLLMExtras(ctx, obj)
		case files.DocStatusSummarizationCompleted:
			nextStage, err = processEmbeddings(ctx, obj)
		case files.DocStatusEmbeddingsCompleted:
			nextStage = files.DocStatusCompleted
		default:
			err = fmt.Errorf("invalid processing stage: %v", currentStage)
		}

		if err != nil {
			logger.Error("processing error", "error", err, "stage", currentStage)
			obj.Stage = files.DocProcStage{
				PGStage:            files.PGStageErrored,
				IsErrored:          true,
				IsCompleted:        true,
				ProcessingErrorMsg: fmt.Sprintf("Encountered Processing Error: %v", err),
				IngestErrorMsg:     obj.Stage.IngestErrorMsg,
				DocProcStatus:      currentStage,
			}
			return *obj, fmt.Errorf("processing error at stage %d: %w", currentStage, err)
		}
		currentStage = nextStage
	}

	return files.CompleteFileSchema{}, errors.New("exceeded maximum processing iterations")
}

func processStageHandleExtension(ctx context.Context, obj *files.CompleteFileSchema) (files.DocProcStatus, error) {
	for index, attachment := range obj.Attachments {
		validExtension, err := files.FileExtensionFromString(attachment.Extension)
		if err != nil {
			return files.DocStatusUnprocessed, fmt.Errorf("invalid file extension: %w", err)
		}
		s3_obj := s3utils.NewKeFileManager()

		err = validators.ValidateExtensionFromHash(*s3_obj, attachment.Hash, validExtension)
		if err != nil {
			obj.Stage.SkipProcessing = true
			return files.DocStatusUnprocessed, fmt.Errorf("file validation failed: %v", err)
		}

		obj.Attachments[index].Extension = string(validExtension)
		if validExtension == files.KnownFileExtensionXLSX {
			return files.DocStatusUnprocessed, nil
		}

	}
	return files.DocStatusBeginProcessing, nil
}

func processGenerateRawText(ctx context.Context, obj *files.CompleteFileSchema, texts map[string]string) (files.DocProcStatus, error) {
	for _, attachment := range obj.Attachments {
		doesnt_have_text := len(attachment.Texts) == 0
		if doesnt_have_text {
			return files.DocStatusBeginProcessing, fmt.Errorf("PDF transcription via crimson not implemented yet")
		}

	}
	// Only english text is supported for now so we jump straight past translation
	return files.DocStatusTextCompleted, nil
}

func processTranslateRawText(ctx context.Context, obj *files.CompleteFileSchema, texts map[string]string) (files.DocProcStatus, error) {
	return files.DocStatusRawTextCompleted, errors.New("not implemented")
}

func createLLMExtras(ctx context.Context, obj *files.CompleteFileSchema) (files.DocProcStatus, error) {
	return files.DocStatusTextCompleted, errors.New("not implemented")
}

func processEmbeddings(ctx context.Context, obj *files.CompleteFileSchema) (files.DocProcStatus, error) {
	return files.DocStatusSummarizationCompleted, errors.New("not implemented")
}
