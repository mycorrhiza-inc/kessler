package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kessler/internal/constants"
	"kessler/internal/hashes"
	"kessler/internal/ingest/external"
	"kessler/internal/ingest/validators"
	"kessler/internal/objects/files"
	"kessler/internal/s3utils"
	"log/slog"
	"os"

	"github.com/google/uuid"
	// Assume these are implemented in other packages
)

var OS_HASH_FILEDIR = os.Getenv("OS_HASH_FILEDIR")

type DocumentStatus int

const (
	DocumentStatusUnprocessed DocumentStatus = iota
	DocumentStatusStage1
	DocumentStatusStage2
	DocumentStatusStage3
	DocumentStatusCompleted
	DocumentStatusEmbeddingsCompleted
	DocumentStatusSummarizationCompleted
)

func (d DocumentStatus) Index() int {
	return int(d)
}

type PGStage int

const (
	PGStageUnprocessed PGStage = iota
	PGStageCompleted
	PGStageErrored
)

type DocProcStage struct {
	PGStage            PGStage
	ProcessingErrorMsg string
	IngestErrorMsg     string
	DatabaseErrorMsg   string
	DocProcStage       DocumentStatus
	IsErrored          bool
	IsCompleted        bool
	SkipProcessing     bool
}

type FileTextSchema struct {
	IsOriginalText bool
	Language       string
	Text           string
}

type CompleteFileSchema struct {
	ID        uuid.UUID
	Hash      string
	Lang      string
	Extension string
	DocTexts  []FileTextSchema
	Stage     DocProcStage
	Extra     json.RawMessage
}

type MarkdownExtractor struct {
	logger   *slog.Logger
	tmpDir   string
	priority bool
}

func ProcessFileRaw(ctx context.Context, obj *CompleteFileSchema, stopAt *DocumentStatus, priority bool) (*CompleteFileSchema, error) {
	logger := slog.Default()

	if obj.Hash == "" {
		return nil, errors.New("hash is required")
	}

	if obj.Lang == "" {
		return nil, errors.New("language is required")
	}

	currentStage := obj.Stage.DocProcStage
	defaultStop := DocumentStatusCompleted
	if stopAt == nil {
		stopAt = &defaultStop
	}

	mdextract := &MarkdownExtractor{
		logger:   logger,
		tmpDir:   constants.OS_TMPDIR,
		priority: priority,
	}

	texts := make(map[string]string)

	for i := 0; i < 1000; i++ {
		if currentStage.Index() >= stopAt.Index() {
			obj.Stage = DocProcStage{
				PGStage:      PGStageCompleted,
				IsCompleted:  true,
				DocProcStage: currentStage,
			}
			return obj, nil
		}

		var nextStage DocumentStatus
		var err error

		switch currentStage {
		case DocumentStatusUnprocessed:
			nextStage, err = processStageHandleExtension(ctx, obj)
		case DocumentStatusStage1:
			nextStage, err = processStageOne(ctx, obj, mdextract, texts)
		case DocumentStatusStage2:
			nextStage, err = processStageTwo(ctx, obj, mdextract, texts)
		case DocumentStatusStage3:
			nextStage, err = createLLMExtras(ctx, obj)
		case DocumentStatusSummarizationCompleted:
			nextStage, err = processEmbeddings(ctx, obj)
		case DocumentStatusEmbeddingsCompleted:
			nextStage = DocumentStatusCompleted
		default:
			err = fmt.Errorf("invalid processing stage: %v", currentStage)
		}

		if err != nil {
			logger.Error("processing error", "error", err, "stage", currentStage)
			obj.Stage = DocProcStage{
				PGStage:            PGStageErrored,
				IsErrored:          true,
				IsCompleted:        true,
				ProcessingErrorMsg: fmt.Sprintf("Encountered Processing Error: %v", err),
				IngestErrorMsg:     obj.Stage.IngestErrorMsg,
				DocProcStage:       currentStage,
			}
			return obj, fmt.Errorf("processing error at stage %d: %w", currentStage, err)
		}
		currentStage = nextStage
	}

	return nil, errors.New("exceeded maximum processing iterations")
}

func processStageHandleExtension(ctx context.Context, obj *CompleteFileSchema) (DocumentStatus, error) {
	validExtension, err := files.FileExtensionFromString(obj.Extension)
	if err != nil {
		return DocumentStatusUnprocessed, fmt.Errorf("invalid file extension: %w", err)
	}
	s3_obj := s3utils.NewKeFileManager()
	obj_hash, err := hashes.HashFromString(obj.Hash)
	if err != nil {
		return DocumentStatusUnprocessed, fmt.Errorf("invalid hash: %w", err)
	}

	err = validators.ValidateExtensionFromHash(*s3_obj, obj_hash, validExtension)
	if err != nil {
		obj.Stage.SkipProcessing = true
		return DocumentStatusUnprocessed, fmt.Errorf("file validation failed: %v", err)
	}

	obj.Extension = string(validExtension)
	if validExtension == files.KnownFileExtensionXLSX {
		return DocumentStatusCompleted, nil
	}
	return DocumentStatusStage1, nil
}

func processStageOne(ctx context.Context, obj *CompleteFileSchema, extractor *MarkdownExtractor, texts map[string]string) (DocumentStatus, error) {
	validated_hash, err := hashes.HashFromString(obj.Hash)
	if err != nil {
		return DocumentStatusStage1, fmt.Errorf("invalid hash: %w", err)
	}
	processedText, err := external.TranscribePDFFromHash(validated_hash)
	if err != nil {
		// if _, ok := err.(*common.FormatError); ok {
		// 	obj.Stage.SkipProcessing = true
		// }
		return DocumentStatusStage1, err
	}

	obj.DocTexts = append(obj.DocTexts, FileTextSchema{
		IsOriginalText: true,
		Language:       obj.Lang,
		Text:           processedText,
	})

	if obj.Lang == "en" {
		texts["englishText"] = processedText
		return DocumentStatusStage3, nil
	}

	texts["originalText"] = processedText
	return DocumentStatusStage2, nil
}

func processStageTwo(ctx context.Context, obj *CompleteFileSchema, extractor *MarkdownExtractor, texts map[string]string) (DocumentStatus, error) {
	return DocumentStatusStage2, errors.New("not implemented")
	// if obj.Lang == "en" {
	// 	return DocumentStatusStage3, errors.New("invalid state: document already in english")
	// }
	//
	// translated, err := extractor.TranslateText(texts["originalText"], obj.Lang)
	// if err != nil {
	// 	return DocumentStatusStage2, fmt.Errorf("translation failed: %w", err)
	// }
	//
	// obj.DocTexts = append(obj.DocTexts, FileTextSchema{
	// 	IsOriginalText: false,
	// 	Language:       "en",
	// 	Text:           translated,
	// })
	// texts["englishText"] = translated
	// return DocumentStatusStage3, nil
}

func createLLMExtras(ctx context.Context, obj *CompleteFileSchema) (DocumentStatus, error) {
	return DocumentStatusStage3, errors.New("not implemented")
	// extras, err := common.GenerateExtras(obj)
	//
	//	if err != nil {
	//		return DocumentStatusStage3, fmt.Errorf("LLM extras generation failed: %w", err)
	//	}
	//
	// obj.Extra = extras
	// return DocumentStatusSummarizationCompleted, nil
}

func processEmbeddings(ctx context.Context, obj *CompleteFileSchema) (DocumentStatus, error) {
	return DocumentStatusSummarizationCompleted, errors.New("not implemented")
	//	if err := common.InsertEmbeddings(obj); err != nil {
	//		return DocumentStatusSummarizationCompleted, fmt.Errorf("embeddings failed: %w", err)
	//	}
	//
	// return DocumentStatusEmbeddingsCompleted, nil
}
