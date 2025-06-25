package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"kessler/pkg/util"
	"time"

	"github.com/google/uuid"
	//"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// Hydrate document data
func (s *SearchService) HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int, full_fetch bool) (CardData, error) {
	log := logger.FromContext(ctx)
	// Check cache first
	cacheKey := cache.PrepareKey("search", "document", result.ID)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if doc, ok := cached.(DocumentCardData); ok {
			doc.Index = index // Update index for current search
			return doc, nil
		}
	}

	// Extract metadata
	name := ""
	description := result.Text
	timestamp := time.Now()
	var err error
	extraInfo := ""
	fileID := uuid.Nil
	convoID := uuid.Nil
	authorIDs := []uuid.UUID{}

	if result.Metadata != nil {
		if fileName, ok := result.Metadata["file_name"].(string); ok {
			name = fileName
		}

		if fileIDString, ok := result.Metadata["file_id"].(string); ok {
			fileID, err = uuid.Parse(fileIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse file_id in metadata")
			}
		}

		if convoIDString, ok := result.Metadata["conversation_id"].(string); ok {
			convoID, err = uuid.Parse(convoIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse conversation_id in metadata")
			}
		}

		if authorIDsRaw, ok := result.Metadata["author_ids"].([]string); ok {
			parseUUID := func(val string) (uuid.UUID, error) { return uuid.Parse(val) }
			authorIDs, err = util.MapErrorBubble(authorIDsRaw, parseUUID)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse an author_id in author_ids in metadata")
			}
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

	if len(result.ID) < 36 {
		return DocketCardData{}, fmt.Errorf("Document does not have long enough uuid.")
	}
	parsedAttachmentUUID, err := uuid.Parse(result.ID[:36])
	if err != nil {
		return DocketCardData{}, fmt.Errorf("Could not parse uuid for object")
	}

	// Full fetch: fallback to database for missing metadata
	if full_fetch {
		needFetch := name == "" || description == result.Text
		if needFetch {
			q := dbstore.New(s.db)
			// Fetch file basic info
			if fileRec, err := q.ReadFile(ctx, parsedAttachmentUUID); err == nil {
				if name == "" {
					name = fileRec.Name
				}
				timestamp = fileRec.DatePublished.Time
			} else {
				log.Warn("Failed to read file record", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
			}
			// Fetch metadata record
			if metaRec, err := q.FetchMetadata(ctx, parsedAttachmentUUID); err == nil {
				var m map[string]interface{}
				if err := json.Unmarshal(metaRec.Mdata, &m); err == nil {
					if fn, ok := m["file_name"].(string); ok {
						name = fn
					}
					if desc2, ok := m["description"].(string); ok {
						description = desc2
					}
					if ca, ok := m["created_at"].(string); ok {
						if t2, err := time.Parse(time.RFC3339, ca); err == nil {
							timestamp = t2
						}
					}
					if cn, ok := m["case_number"].(string); ok {
						extraInfo = fmt.Sprintf("Case: %s", cn)
					}
				} else {
					log.Warn("Failed to unmarshal metadata JSON", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
				}
			} else {
				log.Warn("Failed to fetch metadata record", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
			}
		}
	}

	// If no name, error out
	if name == "" {
		return DocumentCardData{}, fmt.Errorf("name is nil")
	}

	// Initialize card
	card := DocumentCardData{
		Name:         name,
		Description:  description,
		Timestamp:    timestamp,
		ExtraInfo:    extraInfo,
		Index:        index,
		Type:         "document",
		ObjectUUID:   fileID,
		FragmentID:   result.ID[36:],
		Authors:      []DocumentAuthor{},
		Conversation: DocumentConversation{},
	}
	log.Info("Successfully Created Initial Card Data", zap.String("file_id", parsedAttachmentUUID.String()))

	// Validate IDs
	if fileID == uuid.Nil {
		return DocumentCardData{}, fmt.Errorf("file_id is nil")
	}
	if convoID == uuid.Nil {
		return DocumentCardData{}, fmt.Errorf("conversation_id is nil")
	}
	if len(authorIDs) == 0 {
		log.Warn("File appears to have no authors", zap.String("raw_id", result.ID), zap.String("file_id", fileID.String()))
	}

	// Prefetch organization IDs and conversation IDs depending on full_fetch
	q := dbstore.New(s.db)
	type orgInfo struct {
		ID      uuid.UUID
		Primary bool
	}
	var orgInfos []orgInfo
	var convUUIDs []uuid.UUID

	if full_fetch {
		// Fetch authorship info from DB
		authorships, err := q.AuthorshipDocumentListOrganizations(ctx, parsedAttachmentUUID)
		if err != nil {
			log.Warn("Failed to list authorships", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
		} else if len(authorships) == 0 {
			log.Warn("No authorships found", zap.String("file_id", parsedAttachmentUUID.String()))
		} else {
			for _, a := range authorships {
				orgInfos = append(orgInfos, orgInfo{ID: a.OrganizationID, Primary: a.IsPrimaryAuthor.Valid && a.IsPrimaryAuthor.Bool})
			}
		}
		// Fetch conversation info from DB
		convInfo, err := q.ConversationIDFetchFromFileID(ctx, parsedAttachmentUUID)
		if err != nil {
			log.Warn("Failed to fetch conversation ID from file ID", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
		} else if len(convInfo) == 0 {
			log.Warn("No conversation info found for file", zap.String("file_id", parsedAttachmentUUID.String()))
		} else {
			convUUIDs = append(convUUIDs, convInfo[0].ConversationUuid)
		}
	} else {
		// Use metadata-provided IDs
		for _, id := range authorIDs {
			orgInfos = append(orgInfos, orgInfo{ID: id, Primary: false})
		}
		if convoID != uuid.Nil {
			convUUIDs = append(convUUIDs, convoID)
		}
	}

	// Lookup organization details
	for _, info := range orgInfos {
		org, err := q.OrganizationRead(ctx, info.ID)
		if err != nil {
			log.Warn("Failed to read organization for authorship", zap.String("org_id", info.ID.String()), zap.Error(err))
			continue
		}
		card.Authors = append(card.Authors, DocumentAuthor{
			AuthorName:      org.Name,
			IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
			IsPrimaryAuthor: info.Primary,
			AuthorID:        org.ID,
		})
	}

	// Lookup conversation details (use first if multiple)
	if len(convUUIDs) > 0 {
		conv, err := q.DocketConversationRead(ctx, convUUIDs[0])
		if err != nil {
			log.Warn("Failed to read conversation details", zap.String("conversation_id", convUUIDs[0].String()), zap.Error(err))
		} else {
			card.Conversation = DocumentConversation{
				ConvoName: conv.Name,
				ConvoID:   conv.ID,
			}
		}
	}

	// Cache the result
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}
