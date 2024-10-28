package crud

import (
	"fmt"

	"github.com/google/uuid"
)

type PGStage string

const (
	StageStatePending    PGStage = "pending"
	StageStateProcessing PGStage = "processing"
	StageStateCompleted  PGStage = "completed"
	StageStateErrored    PGStage = "errored"
)

type DocProcStatus string

type DocProcStage struct {
	PGStage         PGStage       `json:"pg_stage"`
	DocProcStatus   DocProcStatus `json:"docproc_stage"`
	IsErrored       bool          `json:"is_errored"`
	IsCompleted     bool          `json:"is_completed"`
	ErrorMsg        string        `json:"error_msg"`
	ErrorStacktrace string        `json:"error_stacktrace"`
}
type CompleteFileSchema struct {
	ID        uuid.UUID
	Extension string
	Lang      string
	Name      string
	Hash      string
	IsPrivate bool
	MdataStr  string
	Texts     []FileTextSchema
	Stage     DocProcStage
}

func CompleteFileSchemaPrune(input CompleteFileSchema) FileSchema {
	return FileSchema{
		ID:        input.ID,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
}

func FileSchemaToComplete(input FileSchema) (CompleteFileSchema, error) {
	return_schema := CompleteFileSchema{
		ID:        input.ID,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
	// TODO: Query Metadata json and also get other stuff
	return return_schema, fmt.Errorf("not implemented")
}
