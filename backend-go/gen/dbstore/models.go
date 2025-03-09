// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

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

type Attachment struct {
	ID        uuid.UUID
	FileID    uuid.UUID
	Lang      string
	Name      string
	Extension string
	Hash      string
	Mdata     []byte
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type AttachmentExtra struct {
	ID        uuid.UUID
	ExtraObj  []byte
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type AttachmentTextSource struct {
	ID             uuid.UUID
	AttachmentID   uuid.UUID
	IsOriginalText bool
	Language       string
	Text           string
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type DocketConversation struct {
	ID            uuid.UUID
	DocketGovID   string
	State         string
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
	Name          string
	Description   string
	MatterType    string
	IndustryType  string
	Metadata      string
	Extra         string
	DatePublished pgtype.Timestamptz
}

type DocketDocument struct {
	ConversationUuid uuid.UUID
	FileID           uuid.UUID
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
}

type Encounter struct {
	Name        pgtype.Text
	Description pgtype.Text
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type Event struct {
	Date        pgtype.Timestamptz
	Name        pgtype.Text
	Description pgtype.Text
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type Faction struct {
	Name        string
	Description string
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type File struct {
	ID            uuid.UUID
	Lang          string
	Name          string
	Extension     string
	Isprivate     pgtype.Bool
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
	Hash          string
	Verified      pgtype.Bool
	DatePublished pgtype.Timestamptz
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

type Filter struct {
	ID          uuid.UUID
	Name        string
	State       string
	FilterType  string
	Description pgtype.Text
	IsActive    pgtype.Bool
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type FilterDatasetMapping struct {
	ID        uuid.UUID
	FilterID  uuid.UUID
	DatasetID uuid.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Job struct {
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	JobPriority int32
	JobName     string
	JobStatus   string
	JobType     string
	JobData     []byte
}

type JobsLog struct {
	ID        uuid.UUID
	JobID     uuid.UUID
	CreatedAt pgtype.Timestamp
	Status    string
	Message   pgtype.Text
}

type JuristictionInformation struct {
	ID             uuid.UUID
	Country        pgtype.Text
	State          pgtype.Text
	Municipality   pgtype.Text
	Agency         pgtype.Text
	ProceedingName pgtype.Text
	Extra          []byte
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type MultiselectValue struct {
	ID        uuid.UUID
	FilterID  uuid.UUID
	Value     string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
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

type RelationFactionsEncounter struct {
	EncounterID uuid.UUID
	FactionID   uuid.UUID
	ID          uuid.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type RelationFilesEvent struct {
	FileID    uuid.UUID
	EventID   uuid.UUID
	ID        uuid.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type RelationIndividualsEvent struct {
	IndividualID uuid.UUID
	EventID      uuid.UUID
	ID           uuid.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type RelationOrganizationsEvent struct {
	OrganizationID uuid.UUID
	EventID        uuid.UUID
	ID             uuid.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type RelationOrganizationsFaction struct {
	FactionID      uuid.UUID
	OrganizationID uuid.UUID
	ID             uuid.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
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

type Testmat struct {
	ID               uuid.UUID
	Name             string
	Extension        string
	Lang             string
	Verified         pgtype.Bool
	Hash             string
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	DatePublished    pgtype.Timestamptz
	Mdata            []byte
	ExtraObj         []byte
	ConversationUuid pgtype.UUID
	DocketGovID      pgtype.Text
	FileText         string
	Organizations    []byte
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
	ID            pgtype.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type Usergroup struct {
	ID        uuid.UUID
	Name      string
	CreatedAt pgtype.Timestamp
}
