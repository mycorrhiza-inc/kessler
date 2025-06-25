package hydration

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/internal/objects/resultcards"
	"kessler/pkg/logger"
	"time"

	"github.com/google/uuid"
	//"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int, full_fetch bool, s *search.SearchService) (resultcards.CardData, error) {
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
			parse_uuids := func(val string) (uuid.UUID, error) {
				return uuid.Parse(val)
			}
			authorIDs, err = util.MapErrorBubble(authorIDsRaw, parse_uuids)
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
	if full_fetch {
		// Fallback to database if metadata not provided
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

	// If no name, use ID
	if name == "" {
		return DocumentCardData{}, fmt.Errorf("name is nil ")
	}
	card := DocumentCardData{
		Name:         name,
		Description:  description,
		Timestamp:    timestamp,
		ExtraInfo:    extraInfo,
		Index:        index,
		Type:         "document",
		ObjectUUID:   fileID,
		FragmentID:   result.ID[36:],
		Authors:      []DocumentAuthor{}, // Would need to query authorship table
		Conversation: DocumentConversation{},
	}
	log.Info("Successfully Created Initial Card Data", zap.String("file_id", parsedAttachmentUUID.String()))

	if fileID == uuid.Nil {
		return DocumentCardData{}, fmt.Errorf("file_id is nil")
	}

	if convoID == uuid.Nil {
		return DocumentCardData{}, fmt.Errorf("conversation_id")
	}

	if len(authorIDs) == 0 {
		log.Warn("File appears to have no authors", zap.String("raw_id", result.ID), zap.String("file_id", fileID.String()))
	}

	// TODO: SINCE THE CODE WAS MADE WAY FASTER BY PREFETCHING THE ORGANIZATION AND CONVO IDS IN THE SEARCH RESULTS, COULD YOU MOVE THE CODE THAT DOES THAT INITIAL QUERIES INTO THE FULL FETCH BRANCH, BUT STILL LOOK UP STUFF LIKE THE ORGANIZATION NAME AND CONVERSATION NAME IN THE MAIN CODE PATH

	// Try to get authors if this is a file document
	queries := dbstore.New(s.db)
	authorships, err := queries.AuthorshipDocumentListOrganizations(ctx, parsedAttachmentUUID)
	if err != nil {
		log.Warn("Failed to list authorships", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
	} else if len(authorships) == 0 {
		log.Warn("No authorships found", zap.String("file_id", parsedAttachmentUUID.String()))
	} else {
		for _, authorship := range authorships {
			// Get organization details
			org, err := queries.OrganizationRead(ctx, authorship.OrganizationID)
			if err != nil {
				log.Info("Failed to read organization for authorship", zap.String("org_id", authorship.OrganizationID.String()), zap.Error(err))
				continue
			}
			author := DocumentAuthor{
				AuthorName:      org.Name,
				IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
				IsPrimaryAuthor: authorship.IsPrimaryAuthor.Valid && authorship.IsPrimaryAuthor.Bool,
				AuthorID:        org.ID,
			}
			card.Authors = append(card.Authors, author)
		}
	}

	// Try to get conversation info if this is a file document
	// conversation_uuid is stored in public.docket_documents
	conv_info, err := queries.ConversationIDFetchFromFileID(ctx, parsedAttachmentUUID)
	if err != nil {
		log.Warn("Failed to fetch conversation ID from file ID", zap.String("file_id", parsedAttachmentUUID.String()), zap.Error(err))
	} else if len(conv_info) == 0 {
		log.Warn("No conversation info found for file", zap.String("file_id", parsedAttachmentUUID.String()))
	} else {
		// Fetch conversation details
		conv, err := queries.DocketConversationRead(ctx, conv_info[0].ConversationUuid)
		if err != nil {
			log.Info("Failed to read conversation details", zap.String("conversation_id", conv_info[0].ConversationUuid.String()), zap.Error(err))
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
