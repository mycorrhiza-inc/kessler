package routes

import (
	"kessler/internal/ingest/tasks"
	"kessler/pkg/timestamp"
	"net/http"
	"time"
)

// @Summary		Add OpenScraper Ingest Task
// @Description	Creates a new OpenScraper-specific ingestion task
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			body	body		OpenScraperFiling	true	"OpenScraper filing information"
// @Success		200		{object}	tasks.KesslerTaskInfo
// @Failure		400		{string}	string	"Error decoding request body"
// @Failure		500		{string}	string	"Error adding task"
// @Router			/add-task/ingest/openscraper [post]
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

func (o OpenScraperFiling) IntoScraperInfo() (tasks.ScraperInfoPayload, error) {
	filedTime, err := time.Parse("2006-01-02", o.FiledDate) // Adjust date format as needed
	if err != nil {
		return tasks.ScraperInfoPayload{}, err
	}

	attachments := make([]tasks.AttachmentChildPayload, len(o.Attachments))
	for i, attach := range o.Attachments {
		mdata := make(map[string]interface{})

		// Add standard fields
		if attach.DocumentType != "" {
			mdata["document_type"] = attach.DocumentType
		}

		// Merge with extra metadata
		for k, v := range attach.ExtraMetadata {
			mdata[k] = v
		}

		attachments[i] = tasks.AttachmentChildPayload{
			Lang:  "en",
			Name:  attach.Name,
			URL:   attach.URL,
			Mdata: mdata,
		}
	}

	return tasks.ScraperInfoPayload{
		Attachments:        attachments,
		DocketID:           o.CaseNumber,
		PublishedDate:      timestamp.RFC3339Time(filedTime),
		InternalSourceName: "OpenScraper",
		AuthorOrganisation: o.PartyName,
		FileClass:          o.FilingType,
		// Add additional fields as required by your ScraperInfoPayload
	}, nil
}
