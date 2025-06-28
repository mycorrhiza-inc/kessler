package handler

import (
	"time"

	"github.com/google/uuid"
	"kessler/internal/objects/files"
	"kessler/internal/search"
)

// FileSchemaToDocCard converts a CompleteFileSchema to a search.DocumentCardData for displaying document cards.
// This is a basic implementation and can be customized further.
func FileSchemaToDocCard(fileSchema files.CompleteFileSchema) search.DocumentCardData {
	// Map authors
	docAuthors := make([]search.DocumentAuthor, len(fileSchema.Authors))
	for i, a := range fileSchema.Authors {
		docAuthors[i] = search.DocumentAuthor{
			AuthorName:      a.AuthorName,
			IsPerson:        a.IsPerson,
			IsPrimaryAuthor: a.IsPrimaryAuthor,
			AuthorID:        a.AuthorID,
		}
	}

	// Map conversation
	docConv := search.DocumentConversation{
		ConvoName:   fileSchema.Conversation.Name,
		ConvoNumber: fileSchema.Conversation.DocketGovID,
		ConvoID:     fileSchema.Conversation.ID,
	}

	// Select first attachment if available
	var attachID uuid.UUID
	if len(fileSchema.Attachments) > 0 {
		attachID = fileSchema.Attachments[0].ID
	}

	// Build the document card data
	return search.DocumentCardData{
		Name:           fileSchema.Name,
		Description:    fileSchema.Extra.Summary,
		Timestamp:      time.Time(fileSchema.DatePublished),
		ExtraInfo:      fileSchema.Extra.ShortSummary,
		Index:          0,
		Type:           "document",
		ObjectUUID:     fileSchema.ID,
		AttachmentUUID: attachID,
		FragmentID:     "",
		Authors:        docAuthors,
		Conversation:   docConv,
	}
}