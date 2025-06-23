package search

import (
	"context"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthorCardData struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	ExtraInfo   string    `json:"extraInfo,omitempty"`
	Index       int       `json:"index"`
	Type        string    `json:"type"`
	ObjectUUID  uuid.UUID `json:"object_uuid"`
}

func (a AuthorCardData) GetType() string {
	return "author"
}

type DocketCardData struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Index       int       `json:"index"`
	Type        string    `json:"type"`
	ObjectUUID  uuid.UUID `json:"object_uuid"`
}

func (d DocketCardData) GetType() string {
	return "docket"
}

type DocumentAuthor struct {
	AuthorName      string    `json:"author_name"`
	IsPerson        bool      `json:"is_person"`
	IsPrimaryAuthor bool      `json:"is_primary_author"`
	AuthorID        uuid.UUID `json:"author_id"`
}

type DocumentConversation struct {
	ConvoName string    `json:"convo_name"`
	ConvoID   uuid.UUID `json:"convo_id"`
}

type DocumentCardData struct {
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Timestamp    time.Time            `json:"timestamp"`
	ExtraInfo    string               `json:"extraInfo,omitempty"`
	Index        int                  `json:"index"`
	Type         string               `json:"type"`
	ObjectUUID   uuid.UUID            `json:"object_uuid"`
	FragmentID   string               `json:"fragment_id"`
	Authors      []DocumentAuthor     `json:"authors"`
	Conversation DocumentConversation `json:"conversation"`
}

func (d DocumentCardData) GetType() string {
	return "document"
}

// BuildAuthorCard fetches an organization by ID and constructs AuthorCardData.
func BuildAuthorCard(ctx context.Context, db dbstore.DBTX, id string, index int) (AuthorCardData, error) {
	orgID, err := uuid.Parse(id)
	if err != nil {
		return AuthorCardData{}, fmt.Errorf("invalid organization ID: %w", err)
	}
	queries := dbstore.New(db)
	org, err := queries.OrganizationRead(ctx, orgID)
	if err != nil {
		return AuthorCardData{}, fmt.Errorf("failed to read organization: %w", err)
	}
	extraInfo := ""
	if org.IsPerson.Valid && org.IsPerson.Bool {
		extraInfo = "Individual contributor"
	} else {
		extraInfo = "Organization"
	}
	return AuthorCardData{
		Name:        org.Name,
		Description: org.Description,
		Timestamp:   org.CreatedAt.Time,
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "author",
		ObjectUUID:  org.ID,
	}, nil
}

// BuildDocketCard fetches a conversation by ID and constructs DocketCardData.
func BuildDocketCard(ctx context.Context, db dbstore.DBTX, id string, index int) (DocketCardData, error) {
	convID, err := uuid.Parse(id)
	if err != nil {
		return DocketCardData{}, fmt.Errorf("invalid conversation ID: %w", err)
	}
	queries := dbstore.New(db)
	conv, err := queries.DocketConversationRead(ctx, convID)
	if err != nil {
		return DocketCardData{}, fmt.Errorf("failed to read conversation: %w", err)
	}
	return DocketCardData{
		Name:        conv.Name,
		Description: conv.Description,
		Timestamp:   conv.CreatedAt.Time,
		Index:       index,
		Type:        "docket",
		ObjectUUID:  conv.ID,
	}, nil
}

// BuildDocumentCard fetches minimal metadata and, if fetchDetails is true, authors and conversation info.
func BuildDocumentCard(ctx context.Context, db dbstore.DBTX, result fugusdk.FuguSearchResult, index int, fetchDetails bool) (DocumentCardData, error) {
	log := logger.FromContext(ctx)
	// Extract metadata
	name := ""
	description := result.Text
	timestamp := time.Now()
	extraInfo := ""

	if result.Metadata != nil {
		if fileName, ok := result.Metadata["file_name"].(string); ok {
			name = fileName
		}
		if desc, ok := result.Metadata["description"].(string); ok {
			description = desc
		}
		if createdAt, ok := result.Metadata["created_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				timestamp = t
			}
		}
		if caseNumber, ok := result.Metadata["case_number"].(string); ok {
			extraInfo = fmt.Sprintf("Case: %s", caseNumber)
		}
	}
	if name == "" {
		name = fmt.Sprintf("Document %s", result.ID)
	}
	if len(result.ID) < 36 {
		return DocumentCardData{}, fmt.Errorf("document ID too short to parse UUID")
	}
	parsedUUID, err := uuid.Parse(result.ID[:36])
	if err != nil {
		return DocumentCardData{}, fmt.Errorf("could not parse UUID: %w", err)
	}
	card := DocumentCardData{
		Name:         name,
		Description:  description,
		Timestamp:    timestamp,
		ExtraInfo:    extraInfo,
		Index:        index,
		Type:         "document",
		ObjectUUID:   parsedUUID,
		FragmentID:   result.ID[36:],
		Authors:      []DocumentAuthor{},
		Conversation: DocumentConversation{},
	}
	log.Info("Created DocumentCardData", zap.String("file_id", parsedUUID.String()))

	if fetchDetails {
		queries := dbstore.New(db)
		// Fetch authors
		auths, err := queries.AuthorshipDocumentListOrganizations(ctx, parsedUUID)
		if err != nil {
			log.Warn("Failed to list authorships", zap.String("file_id", parsedUUID.String()), zap.Error(err))
		} else {
			for _, a := range auths {
				org, err := queries.OrganizationRead(ctx, a.OrganizationID)
				if err != nil {
					log.Warn("Failed to read org for authorship", zap.String("org_id", a.OrganizationID.String()), zap.Error(err))
					continue
				}
				card.Authors = append(card.Authors, DocumentAuthor{
					AuthorName:      org.Name,
					IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
					IsPrimaryAuthor: a.IsPrimaryAuthor.Valid && a.IsPrimaryAuthor.Bool,
					AuthorID:        org.ID,
				})
			}
		}
		// Fetch conversation info
		convInfo, err := queries.ConversationIDFetchFromFileID(ctx, parsedUUID)
		if err != nil {
			log.Warn("Failed to fetch conversation ID", zap.String("file_id", parsedUUID.String()), zap.Error(err))
		} else if len(convInfo) > 0 {
			conv, err := queries.DocketConversationRead(ctx, convInfo[0].ConversationUuid)
			if err != nil {
				log.Warn("Failed to read conversation details", zap.String("conv_id", convInfo[0].ConversationUuid.String()), zap.Error(err))
			} else {
				card.Conversation = DocumentConversation{
					ConvoName: conv.Name,
					ConvoID:   conv.ID,
				}
			}
		}
	}
	return card, nil
}
