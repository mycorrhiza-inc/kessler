// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbstore

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type StageState string

const (
	StageStatePending    StageState = "pending"
	StageStateProcessing StageState = "processing"
	StageStateCompleted  StageState = "completed"
	StageStateErrored    StageState = "errored"
)

func (e *StageState) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StageState(s)
	case string:
		*e = StageState(s)
	default:
		return fmt.Errorf("unsupported scan type for StageState: %T", src)
	}
	return nil
}

type NullStageState struct {
	StageState StageState
	Valid      bool // Valid is true if StageState is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStageState) Scan(value interface{}) error {
	if value == nil {
		ns.StageState, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StageState.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStageState) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StageState), nil
}

type DocketConversation struct {
	ID          uuid.UUID
	DocketID    string
	State       string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	DeletedAt   pgtype.Timestamp
	Name        string
	Description string
}

type DocketDocument struct {
	DocketID  uuid.UUID
	FileID    uuid.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type File struct {
	ID        uuid.UUID
	Lang      string
	Name      string
	Extension string
	Isprivate pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Hash      string
	Verified  pgtype.Bool
}

type FileExtra struct {
	ID        uuid.UUID
	Isprivate pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ExtraObj  []byte
}

type FileMetadatum struct {
	ID        uuid.UUID
	Isprivate pgtype.Bool
	Mdata     []byte
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type FileTextSource struct {
	FileID         uuid.UUID
	IsOriginalText bool
	Language       string
	Text           string
	ID             uuid.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type Organization struct {
	Name        string
	Description string
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	IsPerson    pgtype.Bool
}

type OrganizationAlias struct {
	OrganizationAlias string
	OrganizationID    uuid.UUID
	ID                uuid.UUID
	CreatedAt         pgtype.Timestamptz
	UpdatedAt         pgtype.Timestamptz
}

type PrivateAccessControl struct {
	OperatorID    uuid.UUID
	OperatorTable string
	ObjectID      uuid.UUID
	ObjectTable   string
	ID            uuid.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type RelationDocumentsEncounter struct {
	DocumentID  uuid.UUID
	EncounterID uuid.UUID
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type RelationDocumentsOrganizationsAuthorship struct {
	DocumentID      uuid.UUID
	OrganizationID  uuid.UUID
	ID              uuid.UUID
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	IsPrimaryAuthor pgtype.Bool
}

type RelationUsersUsergroup struct {
	UserID      uuid.UUID
	UsergroupID uuid.UUID
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamp
}

type StageLog struct {
	ID        uuid.UUID
	Status    NullStageState
	Log       []byte
	CreatedAt pgtype.Timestamptz
	FileID    uuid.UUID
}

type User struct {
	ID        uuid.UUID
	Username  string
	StripeID  string
	Email     string
	CreatedAt pgtype.Timestamp
}

type UserfilesThaumaturgyApiKey struct {
	KeyName       pgtype.Text
	KeyBlake3Hash string
	ID            uuid.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type Usergroup struct {
	ID        uuid.UUID
	Name      string
	CreatedAt pgtype.Timestamp
}
