// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbstore

import (
	"database/sql/driver"
	"fmt"

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
	ID        pgtype.UUID
	DocketID  string
	State     string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type DocketDocument struct {
	DocketID  pgtype.UUID
	FileID    pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Encounter struct {
	Name        pgtype.Text
	Description pgtype.Text
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type Event struct {
	Date        pgtype.Timestamptz
	Name        pgtype.Text
	Description pgtype.Text
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type Faction struct {
	Name        string
	Description string
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type File struct {
	ID        pgtype.UUID
	Lang      pgtype.Text
	Name      pgtype.Text
	Extension pgtype.Text
	Isprivate pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Hash      pgtype.Text
}

type FileExtra struct {
	ID        pgtype.UUID
	Isprivate pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ExtraObj  []byte
}

type FileMetadatum struct {
	ID        pgtype.UUID
	Isprivate pgtype.Bool
	Mdata     []byte
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type FileTextSource struct {
	FileID         pgtype.UUID
	IsOriginalText bool
	Language       string
	Text           pgtype.Text
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type JuristictionInformation struct {
	ID             pgtype.UUID
	Country        pgtype.Text
	State          pgtype.Text
	Municipality   pgtype.Text
	Agency         pgtype.Text
	ProceedingName pgtype.Text
	Extra          []byte
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type Organization struct {
	Name        string
	Description pgtype.Text
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	IsPerson    pgtype.Bool
}

type PrivateAccessControl struct {
	OperatorID    pgtype.UUID
	OperatorTable string
	ObjectID      pgtype.UUID
	ObjectTable   string
	ID            pgtype.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type RelationDocumentsEncounter struct {
	DocumentID  pgtype.UUID
	EncounterID pgtype.UUID
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type RelationDocumentsOrganizationsAuthorship struct {
	DocumentID      pgtype.UUID
	OrganizationID  pgtype.UUID
	ID              pgtype.UUID
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	IsPrimaryAuthor pgtype.Bool
}

type RelationFactionsEncounter struct {
	EncounterID pgtype.UUID
	FactionID   pgtype.UUID
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type RelationFilesEvent struct {
	FileID    pgtype.UUID
	EventID   pgtype.UUID
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type RelationIndividualsEvent struct {
	IndividualID pgtype.UUID
	EventID      pgtype.UUID
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type RelationOrganizationsEvent struct {
	OrganizationID pgtype.UUID
	EventID        pgtype.UUID
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type RelationOrganizationsFaction struct {
	FactionID      pgtype.UUID
	OrganizationID pgtype.UUID
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type RelationUsersUsergroup struct {
	UserID      pgtype.UUID
	UsergroupID pgtype.UUID
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamp
}

type StageLog struct {
	ID        pgtype.UUID
	Status    NullStageState
	Log       []byte
	CreatedAt pgtype.Timestamptz
	FileID    pgtype.UUID
}

type User struct {
	ID        pgtype.UUID
	Username  pgtype.Text
	StripeID  pgtype.Text
	Email     string
	CreatedAt pgtype.Timestamp
}

type UserfilesThaumaturgyApiKey struct {
	KeyName       pgtype.Text
	KeyBlake3Hash string
	ID            pgtype.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type Usergroup struct {
	ID        pgtype.UUID
	Name      string
	CreatedAt pgtype.Timestamp
}
