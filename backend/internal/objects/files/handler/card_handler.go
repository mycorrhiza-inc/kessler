package handler

import (
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"
	"kessler/internal/search"
	"kessler/pkg/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *FileHandler) FileCardGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "files:FileSemiCompleteGet")
	defer span.End()
	log := logger.FromContext(ctx)

	params := mux.Vars(r)
	fileID := params["uuid"]
	parsedUUID, err := uuid.Parse(fileID)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing file %v: %v", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	q := dbstore.New(h.db)
	file, err := h.SemiCompleteFileGetFromUUID(ctx, q, parsedUUID)
	if err != nil {
		log.Info("encountered error getting file from uuid", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// All identical to the card info so far
	card := FileSchemaToDocCard(file)

	response, _ := json.Marshal(card)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

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

