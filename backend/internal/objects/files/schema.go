package files

import (
	"kessler/internal/objects/authors"
	"kessler/internal/objects/conversations"
	"kessler/pkg/hashes"
	"kessler/pkg/timestamp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AttachmentChildTextSource struct {
	IsOriginalText bool   `json:"is_original_text"`
	Text           string `json:"text"`
	Language       string `json:"language"`
}

type FileTextSchema struct {
	FileID         uuid.UUID `json:"file_id"`
	IsOriginalText bool      `json:"is_original_text"`
	Text           string    `json:"text"`
	Language       string    `json:"language"`
}

//	type AttachmentSchema struct {
//		ID   uuid.UUID `json:"id"`
//		Lang string    `json:"lang"`
//		Name string    `json:"name"`
//		Hash string    `json:"hash"`
//	}
type CompleteAttachmentSchema struct {
	ID        uuid.UUID                   `json:"id"`
	FileID    uuid.UUID                   `json:"file_id"`
	Lang      string                      `json:"lang"`
	Name      string                      `json:"name"`
	Hash      hashes.KesslerHash          `json:"hash"`
	URL       string                      `json:"url"`
	Extension string                      `json:"extension"`
	Mdata     map[string]any              `json:"mdata"`
	Texts     []AttachmentChildTextSource `json:"texts"`
}

type FileSchema struct {
	ID            uuid.UUID             `json:"id"`
	Verified      bool                  `json:"verified"`
	Extension     string                `json:"extension"`
	Lang          string                `json:"lang"`
	Name          string                `json:"name"`
	Hash          string                `json:"hash"`
	IsPrivate     bool                  `json:"is_private"`
	DatePublished timestamp.RFC3339Time `json:"date_published"`
}
type FileMetadataSchema map[string]interface{}

type FileGeneratedExtras struct {
	Summary        string  `json:"summary"`
	ShortSummary   string  `json:"short_summary"`
	Purpose        string  `json:"purpose"`
	Impressiveness float64 `json:"impressiveness"`
}

// To heavy to include in a default file schema unless the user specifies they want a smaller version
type CompleteFileSchema struct {
	ID            uuid.UUID                             `json:"id"`
	Verified      bool                                  `json:"verified"`
	Lang          string                                `json:"lang"`
	Name          string                                `json:"name"`
	IsPrivate     bool                                  `json:"is_private"`
	Attachments   []CompleteAttachmentSchema            `json:"attachments"`
	DatePublished timestamp.RFC3339Time                 `json:"date_published"`
	Mdata         FileMetadataSchema                    `json:"mdata"`
	Stage         DocProcStage                          `json:"stage"`
	Extra         FileGeneratedExtras                   `json:"extra"`
	Authors       []authors.AuthorInformation           `json:"authors"`
	Conversation  conversations.ConversationInformation `json:"conversation"`
}

func (input CompleteFileSchema) CompleteFileSchemaPrune() FileSchema {
	return FileSchema{
		ID:            input.ID,
		Verified:      input.Verified,
		Lang:          input.Lang,
		Name:          input.Name,
		IsPrivate:     input.IsPrivate,
		DatePublished: input.DatePublished,
	}
}

func (input FileSchema) CompleteFileSchemaInflateFromPartialSchema() CompleteFileSchema {
	return_schema := CompleteFileSchema{
		ID:            input.ID,
		Verified:      input.Verified,
		Lang:          input.Lang,
		Name:          input.Name,
		IsPrivate:     input.IsPrivate,
		DatePublished: input.DatePublished,
	}
	// TODO: Query Metadata json and also get other stuff
	return return_schema
}

func (updateInfo CompleteFileSchema) ConvertToCreationData() FileCreationDataRaw {
	creationData := FileCreationDataRaw{
		Lang:          updateInfo.Lang,
		Name:          updateInfo.Name,
		Verified:      pgtype.Bool{Bool: updateInfo.Verified, Valid: true},
		DatePublished: updateInfo.DatePublished,
	}
	return creationData
}
