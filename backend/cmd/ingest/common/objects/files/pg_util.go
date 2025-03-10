package files

import (
	"thaumaturgy/common/objects/timestamp"

	"github.com/jackc/pgx/v5/pgtype"
)

type FileCreationDataRaw struct {
	Extension     string
	Lang          string
	Name          string
	Hash          string
	IsPrivate     pgtype.Bool
	Verified      pgtype.Bool
	DatePublished timestamp.KesslerTime
}
