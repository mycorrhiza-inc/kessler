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

type KesslerTaskInfo struct {
	TaskID string `json:"task_id"`
	Queue  string `json:"queue"`
	State  string `json:"state"`
	Status string `json:"status"`
}

func GenerateTaskInfoFromInfo(info asynq.TaskInfo) KesslerTaskInfo {
	return KesslerTaskInfo{
		TaskID: info.ID,
		Queue:  info.Queue,
	}
}

// I am trying to migrate some legacy code right now, typically this endpoing used to handle data from multiple different sources, but now all of those different data sources are handled and brought into a common format elsewhere. Namely these are accessible using api endpoints like:
//
//
// class GenericCase(BaseModel, extra=Extra.allow):
//     """Model representing case data.
//
//     Attributes:
//         case_number (str): The unique case number.
//         case_type (Optional[str]): The type of the case (e.g., civil, regulatory).
//         description (Optional[str]): A detailed description of the case.
//         industry (Optional[str]): The industry related to the case.
//         petitioner (Optional[str]): The name of the petitioner in the case.
//         hearing_officer (Optional[str]): The hearing officer for the case.
//         opened_date (Optional[date]): The date the case was opened.
//         closed_date (Optional[date]): The date the case was closed.
//         filings (Optional[list[Filing]]): A list of filings associated with the case.
//     """
//
//     case_number: str
//     case_name: str = ""
//     case_url: str = ""
//     case_type: Optional[str] = None
//     description: Optional[str] = None
//     industry: Optional[str] = None
//     petitioner: Optional[str] = None
//     hearing_officer: Optional[str] = None
//     opened_date: Optional[RFC3339Time] = None
//     closed_date: Optional[RFC3339Time] = None
//     filings: Optional[list[GenericFiling]] = None
//     extra_metadata: Dict[str, Any] = {}
//     indexed_at: RFC3339Time = rfc_time_now()
//
//     def model_post_init(self, __context: Any) -> None:
//         if self.model_extra:
//             self.extra_metadata.update(self.model_extra)
// and 
//
// class GenericFiling(BaseModel, extra=Extra.allow):
//     """Model representing filing data within a case.
//
//     Attributes:
//         filed_date (date): The date the filing was made.
//         party_name (str): The name of the party submitting the filing.
//         filing_type (str): The type of filing (e.g., brief, testimony).
//         description (str): A description of the filing.
//         attachments (Optional[list[Attachment]]): A list of associateda ttachments.
//     """
//
//     name: str = ""
//     filed_date: RFC3339Time
//     party_name: str
//     filing_type: str
//     description: str
//     attachments: List[GenericAttachment] = []
//     extra_metadata: Dict[str, Any] = {}
//
// Could you refactor this project to use these 2 data types as the primary interface. Currently its using a type of ScraperInfoPayload in /ingest/tasks/schemas.go
// Also KesslerTime is the same as an RFC3339Time

    def model_post_init(self, __context: Any) -> None:
        if self.model_extra:
            self.extra_metadata.update(self.model_extra)

type AttachmentChildPayload struct {
	Lang      string         `json:"lang"`
	Name      string         `json:"name"`
	Extension string         `json:"extension"`
	URL       string         `json:"url"`
	Mdata     map[string]any `json:"mdata"`
}

type ScraperInfoPayload struct {
	Text                  string                   `json:"text"`
	FileType              string                   `json:"file_type"`
	DocketID              string                   `json:"docket_id"`
	PublishedDate         timestamp.RFC3339Time    `json:"published_date" example:"2024-02-27T12:34:56Z"`
	Name                  string                   `json:"name"`
	InternalSourceName    string                   `json:"internal_source_name"`
	State                 string                   `json:"state"`
	AuthorIndividual      string                   `json:"author_individual"`
	AuthorIndividualEmail string                   `json:"author_individual_email"`
	AuthorOrganisation    string                   `json:"author_organisation"`
	FileClass             string                   `json:"file_class"`
	Lang                  string                   `json:"lang"`
	ItemNumber            string                   `json:"item_number"`
	ExtraMetadata         map[string]any           `json:"extra_metadata"`
	Attachments           []AttachmentChildPayload `json:"attachments"`
}

func CastScraperInfoToNewFile(info ScraperInfoPayload) files.CompleteFileSchema {
	new_attachments := make([]files.CompleteAttachmentSchema, len(info.Attachments))
	for i, attachment := range info.Attachments {
		metadata := attachment.Mdata
		metadata["url"] = attachment.URL
		new_attachments[i] = files.CompleteAttachmentSchema{
			Name:      attachment.Name,
			Lang:      attachment.Lang,
			Extension: attachment.Extension,
			Mdata:     metadata,
		}
	}
	metadata := info.ExtraMetadata
	metadata_fields := map[string]any{
		// "url":                 strings.TrimSpace(info.FileURL),
		"docket_id":           strings.TrimSpace(info.DocketID),
		"extension":           strings.TrimSpace(info.FileType),
		"lang":                strings.TrimSpace(info.Lang),
		"title":               strings.TrimSpace(info.Name),
		"source":              strings.TrimSpace(info.InternalSourceName),
		"date":                info.PublishedDate,
		"file_class":          strings.TrimSpace(info.FileClass),
		"author_organisation": strings.TrimSpace(info.AuthorOrganisation),
		"author":              strings.TrimSpace(info.AuthorIndividual),
		"author_email":        strings.TrimSpace(info.AuthorIndividualEmail),
		"item_number":         strings.TrimSpace(info.ItemNumber),
	}
	for key, value := range metadata_fields {
		if !(reflect.ValueOf(value).IsZero()) {
			metadata[key] = value
		}
	}

	docket_info := conversations.ConversationInformation{
		DocketGovID: strings.TrimSpace(info.DocketID),
	}
	return files.CompleteFileSchema{
		ID:           uuid.Nil,
		Name:         strings.TrimSpace(info.Name),
		Conversation: docket_info,
		Mdata:        metadata,
	}
}

type CastableIntoScraperInfo interface {
	IntoScraperInfo() (ScraperInfoPayload, error)
}

func (s ScraperInfoPayload) IntoScraperInfo() (ScraperInfoPayload, error) {
	return s, nil
}

func AddScraperTaskCastable(ctx context.Context, castable CastableIntoScraperInfo) (KesslerTaskInfo, error) {
	scraper_info, err := castable.IntoScraperInfo()
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("Error Casting to Scraper Info: %v", err)
	}
	task, err := NewAddFileScraperTask(scraper_info)
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("Error creating task: %v", err)
	}
	return EnqueueTaskFromCtx(ctx, task)
}

func EnqueueTaskFromCtx(ctx context.Context, task *asynq.Task) (KesslerTaskInfo, error) {
	client := GetClient(ctx)
	info, err := client.Enqueue(task)
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("Error enqueueing task: %v", err)
	}

	kessler_info := GenerateTaskInfoFromInfo(*info)
	return kessler_info, nil
}

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

func (n NYPUCDocInfo) IntoScraperInfo() (ScraperInfoPayload, error) {
	regular_time, err := time.Parse("01/02/2006", n.DateFiled)
	if err != nil {
		return ScraperInfoPayload{}, err
	}

	attachments := []AttachmentChildPayload{
		{
			Lang:  "en",
			Name:  n.FileName,
			URL:   n.URL,
			Mdata: map[string]any{},
		},
	}

	return ScraperInfoPayload{
		// FileURL:            n.URL,
		Attachments:        attachments,
		DocketID:           n.DocketID,
		PublishedDate:      timestamp.RFC3339Time(regular_time),
		InternalSourceName: "NYPUC",
		State:              "NY",
		AuthorOrganisation: n.Organization,
		FileClass:          n.NYPUCDocType,
		ItemNumber:         n.ItemNo,
	}, nil
}
