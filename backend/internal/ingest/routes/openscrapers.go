package routes

import (
	"kessler/internal/ingest/tasks"
	"kessler/pkg/timestamp"
	"net/http"
	"time"
)

// @Summary	Add OpenScraper Ingest Task
// @Description	Creates a new OpenScraper-specific ingestion task
// @Tags		tasks
// @Accept	json
// @Produce	json
// @Param		body	body	OpenScraperFiling	true	"OpenScraper filing information"
// @Success	200	{object}	tasks.KesslerTaskInfo
// @Failure	400	{string}	string	"Error decoding request body"
// @Failure	500	{string}	string	"Error adding task"
// @Router	/add-task/ingest/openscraper [post]
func HandleOpenScraperIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[OpenScraperFiling](w, r)
}

type OpenScraperFiling struct {
	CaseNumber    string                  `json:"case_number"`
	FiledDate     string                  `json:"filed_date"`
	PartyName     string                  `json:"party_name"`
	FilingType    string                  `json:"filing_type"`
	Description   string                  `json:"description"`
	Attachments   []OpenScraperAttachment `json:"attachments"`
	ExtraMetadata map[string]interface{}  `json:"extra_metadata,omitempty"`
}

type OpenScraperAttachment struct {
	Name          string                 `json:"name"`
	URL           string                 `json:"url"`
	DocumentType  string                 `json:"document_type,omitempty"`
	ExtraMetadata map[string]interface{} `json:"extra_metadata,omitempty"`
}

// IntoScraperInfo converts OpenScraperFiling into the new FilingInfoPayload.
func (o OpenScraperFiling) IntoScraperInfo() (tasks.FilingInfoPayload, error) {
	filedTime, err := time.Parse("2006-01-02", o.FiledDate) // Adjust date format as needed
	if err != nil {
		return tasks.FilingInfoPayload{}, err
	}
	// Build attachments
	attachments := make([]tasks.AttachmentChildInfo, len(o.Attachments))
	for i, attach := range o.Attachments {
		mdata := make(map[string]any)
		// standard
		if attach.DocumentType != "" {
			mdata["document_type"] = attach.DocumentType
		}
		for k, v := range attach.ExtraMetadata {
			mdata[k] = v
		}
		attachments[i] = tasks.AttachmentChildInfo{
			Lang:  "en",
			Name:  attach.Name,
			URL:   attach.URL,
			Mdata: mdata,
		}
	}
	// Build FilingChildInfo
	filing := tasks.FilingChildInfo{
		Name:          o.PartyName,
		FiledDate:     timestamp.RFC3339Time(filedTime),
		PartyName:     o.PartyName,
		FilingType:    o.FilingType,
		Description:   o.Description,
		Attachments:   attachments,
		ExtraMetadata: o.ExtraMetadata,
	}
	// Minimal case info
	caseInfo := tasks.CaseInfoMinimal{CaseNumber: o.CaseNumber}
	return tasks.FilingInfoPayload{Filing: filing, CaseInfo: caseInfo}, nil
}
