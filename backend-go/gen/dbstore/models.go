// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbstore

import (
	"github.com/jackc/pgx/v5/pgtype"
)

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
	Url          pgtype.Text
	Doctype      pgtype.Text
	Lang         pgtype.Text
	Name         pgtype.Text
	Source       pgtype.Text
	Hash         pgtype.Text
	Mdata        pgtype.Text
	Stage        pgtype.Text
	Summary      pgtype.Text
	ShortSummary pgtype.Text
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
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

type Individual struct {
	Name       string
	Username   pgtype.Text
	ChosenName pgtype.Text
	ID         pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
}

type Organization struct {
	Name        string
	Description pgtype.Text
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type RelationDocumentsEncounter struct {
	DocumentID  pgtype.UUID
	EncounterID pgtype.UUID
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type RelationDocumentsIndividualsAuthor struct {
	DocumentID   pgtype.UUID
	IndividualID pgtype.UUID
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type RelationDocumentsOrganization struct {
	DocumentID     pgtype.UUID
	OrganizationID pgtype.UUID
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
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

type RelationIndividualsFaction struct {
	FactionID    pgtype.UUID
	IndividualID pgtype.UUID
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type RelationIndividualsOrganization struct {
	IndividualID   pgtype.UUID
	OrganizationID pgtype.UUID
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
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

type User struct {
	ID        pgtype.UUID
	Username  pgtype.Text
	StripeID  pgtype.Text
	Email     string
	CreatedAt pgtype.Timestamp
}

type UserfilesAcl struct {
	ID          pgtype.UUID
	UsergroupID pgtype.UUID
	OwnerID     pgtype.UUID
	CreatedAt   pgtype.Timestamp
}

type UserfilesPrivateAccessControl struct {
	OperatorID    pgtype.UUID
	OperatorTable string
	ObjectID      pgtype.UUID
	ObjectTable   string
	ID            pgtype.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

type UserfilesPrivateFile struct {
	Url          pgtype.Text
	Doctype      pgtype.Text
	Lang         pgtype.Text
	Name         pgtype.Text
	Source       pgtype.Text
	Hash         pgtype.Text
	Mdata        pgtype.Text
	Stage        pgtype.Text
	Summary      pgtype.Text
	ShortSummary pgtype.Text
	UsergroupID  pgtype.UUID
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type UserfilesPrivateFileTextSource struct {
	FileID         pgtype.UUID
	IsOriginalText bool
	Language       string
	Text           pgtype.Text
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type Usergroup struct {
	ID        pgtype.UUID
	Name      string
	CreatedAt pgtype.Timestamp
}
