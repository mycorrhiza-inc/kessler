package tasks

import (
	"context"
	"fmt"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/files"
	"kessler/pkg/timestamp"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// FilingInfoPayload is the payload for adding a filing to a case ingestion task.
type FilingInfoPayload struct {
	Filing   FilingChildInfo `json:"filing"`
	CaseInfo CaseInfoMinimal `json:"case_info"`
}

// CastScraperInfoToNewFile converts a FilingInfoPayload into the internal CompleteFileSchema.
func CastScraperInfoToNewFile(info FilingInfoPayload) files.CompleteFileSchema {
	// Attachments
	newAttachments := make([]files.CompleteAttachmentSchema, len(info.Filing.Attachments))
	for i, at := range info.Filing.Attachments {
		md := at.Mdata
		md["url"] = at.URL
		newAttachments[i] = files.CompleteAttachmentSchema{
			Name:      at.Name,
			Lang:      at.Lang,
			Extension: at.Extension,
			Mdata:     md,
		}
	}

	// Merge metadata
	metadata := map[string]any{}
	// start from filing extra metadata
	for k, v := range info.Filing.ExtraMetadata {
		metadata[k] = v
	}
	// incorporate case-level metadata
	for k, v := range info.CaseInfo.ExtraMetadata {
		metadata[k] = v
	}

	// core fields
	fields := map[string]any{
		"case_number": info.CaseInfo.CaseNumber,
		"case_url":    info.CaseInfo.CaseURL,
		"filed_date":  info.Filing.FiledDate,
		"party_name":  info.Filing.PartyName,
		"filing_type": info.Filing.FilingType,
		"description": info.Filing.Description,
	}
	for k, v := range fields {
		if !reflect.ValueOf(v).IsZero() {
			metadata[k] = v
		}
	}

	conv := conversations.ConversationInformation{DocketGovID: info.CaseInfo.CaseNumber}
	return files.CompleteFileSchema{
		ID:           uuid.Nil,
		Name:         info.Filing.Name,
		Conversation: conv,
		Mdata:        metadata,
		Attachments:  newAttachments,
	}
}

// CastableIntoFilingInfo is implemented by types that can be converted to FilingInfoPayload.
type CastableIntoFilingInfo interface {
	IntoScraperInfo() (FilingInfoPayload, error)
}

// IntoScraperInfo allows FilingInfoPayload to satisfy CastableIntoFilingInfo.
func (s FilingInfoPayload) IntoScraperInfo() (FilingInfoPayload, error) {
	return s, nil
}

// AddScraperTaskCastable enqueues a filing ingestion task from any CastableIntoFilingInfo.
func AddScraperTaskCastable(ctx context.Context, castable CastableIntoFilingInfo) (KesslerTaskInfo, error) {
	payload, err := castable.IntoScraperInfo()
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("error casting to FilingInfoPayload: %w", err)
	}
	task, err := NewAddFileScraperTask(payload)
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("error creating add-file task: %w", err)
	}
	return EnqueueTaskFromCtx(ctx, task)
}

// EnqueueTaskFromCtx pushes an Asynq task and returns its metadata.
func EnqueueTaskFromCtx(ctx context.Context, task *asynq.Task) (KesslerTaskInfo, error) {
	client := GetClient(ctx)
	info, err := client.Enqueue(task)
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("error enqueueing task: %w", err)
	}
	return GenerateTaskInfoFromInfo(*info), nil
}

// NYPUCDocInfo maps a NYPUC-specific listing to FilingInfoPayload.
type NYPUCDocInfo struct {
	Serial       string `json:"serial"`
	DateFiled    string `json:"date_filed"`
	NYPUCDocType string `json:"nypuc_doctype"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Organization string `json:"organization"`
	ItemNo       string `json:"item_no"`
	FileName     string `json:"file_name"`
	DocketID     string `json:"docket_id"`
}

// IntoScraperInfo converts NYPUCDocInfo into the new FilingInfoPayload.
func (n NYPUCDocInfo) IntoScraperInfo() (FilingInfoPayload, error) {
	t, err := time.Parse("01/02/2006", n.DateFiled)
	if err != nil {
		return FilingInfoPayload{}, err
	}
	filing := FilingChildInfo{
		Name:        n.FileName,
		FiledDate:   timestamp.RFC3339Time(t),
		PartyName:   n.Organization,
		FilingType:  n.NYPUCDocType,
		Description: n.Name,
		Attachments: []AttachmentChildInfo{{
			Lang:      "en",
			Name:      n.FileName,
			Extension: strings.ToLower(strings.TrimPrefix(strings.ToLower(n.FileName[strings.LastIndex(n.FileName, "."):]), ".")),
			URL:       n.URL,
			Mdata:     map[string]any{},
		}},
		ExtraMetadata: map[string]any{},
	}
	caseInfo := CaseInfoMinimal{CaseNumber: n.DocketID}
	return FilingInfoPayload{Filing: filing, CaseInfo: caseInfo}, nil
}

