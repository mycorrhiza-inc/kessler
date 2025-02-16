package files

import (
	"thaumaturgy/objects/authors"
	"thaumaturgy/objects/conversations"
	"thaumaturgy/objects/timestamp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FileChildTextSource struct {
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

func (child_source FileChildTextSource) ChildTextSouceToRealTextSource(id uuid.UUID) FileTextSchema {
	return FileTextSchema{
		FileID:         id,
		IsOriginalText: child_source.IsOriginalText,
		Text:           child_source.Text,
		Language:       child_source.Language,
	}
}

type FileSchema struct {
	ID            uuid.UUID             `json:"id"`
	Verified      bool                  `json:"verified"`
	Extension     string                `json:"extension"`
	Lang          string                `json:"lang"`
	Name          string                `json:"name"`
	Hash          string                `json:"hash"`
	IsPrivate     bool                  `json:"is_private"`
	DatePublished timestamp.KesslerTime `json:"date_published"`
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
	Extension     string                                `json:"extension"`
	Lang          string                                `json:"lang"`
	Name          string                                `json:"name"`
	Hash          string                                `json:"hash"`
	IsPrivate     bool                                  `json:"is_private"`
	DatePublished timestamp.KesslerTime                 `json:"date_published"`
	Mdata         FileMetadataSchema                    `json:"mdata"`
	Stage         DocProcStage                          `json:"stage"`
	Extra         FileGeneratedExtras                   `json:"extra"`
	Authors       []authors.AuthorInformation           `json:"authors"`
	Conversation  conversations.ConversationInformation `json:"conversation"`
	DocTexts      []FileChildTextSource                 `json:"doc_texts"`
}

func (input CompleteFileSchema) CompleteFileSchemaPrune() FileSchema {
	return FileSchema{
		ID:            input.ID,
		Verified:      input.Verified,
		Extension:     input.Extension,
		Lang:          input.Lang,
		Name:          input.Name,
		Hash:          input.Hash,
		IsPrivate:     input.IsPrivate,
		DatePublished: input.DatePublished,
	}
}

func (input FileSchema) CompleteFileSchemaInflateFromPartialSchema() CompleteFileSchema {
	return_schema := CompleteFileSchema{
		ID:            input.ID,
		Verified:      input.Verified,
		Extension:     input.Extension,
		Lang:          input.Lang,
		Name:          input.Name,
		Hash:          input.Hash,
		IsPrivate:     input.IsPrivate,
		DatePublished: input.DatePublished,
	}
	// TODO: Query Metadata json and also get other stuff
	return return_schema
}

func (updateInfo CompleteFileSchema) ConvertToCreationData() FileCreationDataRaw {
	creationData := FileCreationDataRaw{
		Extension:     updateInfo.Extension,
		Lang:          updateInfo.Lang,
		Name:          updateInfo.Name,
		Hash:          updateInfo.Hash,
		Verified:      pgtype.Bool{Bool: updateInfo.Verified, Valid: true},
		DatePublished: updateInfo.DatePublished,
	}
	return creationData
}
