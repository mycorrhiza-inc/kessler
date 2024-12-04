package files

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	// "github.com/jackc/pgx/v5/pgtype"
)

type FileSchema struct {
	ID        uuid.UUID `json:"id"`
	Verified  bool      `json:"verified"`
	Extension string    `json:"extension"`
	Lang      string    `json:"lang"`
	Name      string    `json:"name"`
	Hash      string    `json:"hash"`
	IsPrivate bool      `json:"is_private"`
}
type FileMetadataSchema map[string]interface{}

type CompleteFileSchema struct {
	ID           uuid.UUID               `json:"id"`
	Verified     bool                    `json:"verified"`
	Extension    string                  `json:"extension"`
	Lang         string                  `json:"lang"`
	Name         string                  `json:"name"`
	Hash         string                  `json:"hash"`
	IsPrivate    bool                    `json:"is_private"`
	Mdata        FileMetadataSchema      `json:"mdata"`
	Stage        DocProcStage            `json:"stage"`
	Extra        FileGeneratedExtras     `json:"extra"`
	Authors      []AuthorInformation     `json:"authors"`
	Conversation ConversationInformation `json:"conversation"`
	// To heavy to include in a default file schema unless the user specifies they want it
	DocTexts []FileChildTextSource `json:"doc_texts"`
}

func (input CompleteFileSchema) CompleteFileSchemaPrune() FileSchema {
	return FileSchema{
		ID:        input.ID,
		Verified:  input.Verified,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
}

func (input FileSchema) CompleteFileSchemaInflateFromPartialSchema() CompleteFileSchema {
	return_schema := CompleteFileSchema{
		ID:        input.ID,
		Verified:  input.Verified,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
	// TODO: Query Metadata json and also get other stuff
	return return_schema
}

func (updateInfo CompleteFileSchema) ConvertToCreationData() FileCreationDataRaw {
	creationData := FileCreationDataRaw{
		Extension: updateInfo.Extension,
		Lang:      updateInfo.Lang,
		Name:      updateInfo.Name,
		Hash:      updateInfo.Hash,
		Verified:  pgtype.Bool{Bool: updateInfo.Verified, Valid: true},
	}
	return creationData
}
