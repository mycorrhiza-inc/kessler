package organizations

import (
	"kessler/common/objects/files"

	"github.com/google/uuid"
)

type OrganizationSchemaComplete struct {
	ID               uuid.UUID          `json:"id"`
	Name             string             `json:"name"`
	Aliases          []string           `json:"aliases"`
	FilesAuthored    []files.FileSchema `json:"files_authored"`
	FilesAuthoredIDs []uuid.UUID        `json:"files_authored_ids"`
}
type OrganizationQuickwitSchema struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Aliases            []string  `json:"aliases"`
	FilesAuthoredCount int       `json:"files_authored_count"`
}
