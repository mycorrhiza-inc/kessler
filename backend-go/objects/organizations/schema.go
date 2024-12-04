package organizations

import (
	"kessler/objects/files"

	"github.com/google/uuid"
)

type OrganizationSchemaComplete struct {
	ID               uuid.UUID          `json:"id"`
	Name             string             `json:"name"`
	FilesAuthored    []files.FileSchema `json:"files_authored"`
	FilesAuthoredIDs []uuid.UUID        `json:"files_authored_ids"`
}
