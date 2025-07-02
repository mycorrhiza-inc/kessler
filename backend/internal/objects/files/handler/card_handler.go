package handler

import (
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"
	"kessler/internal/search"
	"kessler/pkg/hashes"
	"kessler/pkg/logger"
	"kessler/pkg/util"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type FillingAttachmentInfo struct {
	AttachmentUUID      uuid.UUID          `json:"attachment_uuid"`
	AttachmentHash      hashes.KesslerHash `json:"attachment_hash"`
	AttachmentName      string             `json:"attachment_name"`
	AttachmentExtension string             `json:"attachment_extension"`
}

type PageInfo struct {
	CardInfo       search.DocumentCardData `json:"card_info"`
	AttachemntInfo []FillingAttachmentInfo `json:"attachments"`
}

func (h *FileHandler) FilePageInfoGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "files:FileSemiCompleteGet")
	defer span.End()
	log := logger.FromContext(ctx)

	params := mux.Vars(r)
	fileIDRaw := params["uuid"]
	fileUUID, err := uuid.Parse(fileIDRaw)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing file %v: %v", fileIDRaw, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	q := dbstore.New(h.db)
	file, err := h.SemiCompleteFileGetFromUUID(ctx, q, fileUUID)
	if err != nil {
		log.Info("encountered error getting file from uuid", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// All identical to the card info so far
	card := FileSchemaToDocCard(file)

	attachments, err := q.AttachmentListByFileId(ctx, fileUUID)
	if err != nil {
		log.Error("Encountered error fetching attachments", zap.Error(err))
	}
	extractAttachments := func(attach dbstore.Attachment) (FillingAttachmentInfo, error) {
		parsedHash, err := hashes.HashFromString(attach.Hash)
		if err != nil {
			return FillingAttachmentInfo{}, err
		}
		extension := strings.TrimSpace(attach.Extension)
		if extension == "" {
			extension = "pdf"
		}
		return FillingAttachmentInfo{
			AttachmentUUID:      attach.ID,
			AttachmentName:      attach.Name,
			AttachmentHash:      parsedHash,
			AttachmentExtension: extension,
		}, nil
	}
	attachInfos, err := util.MapErrorBubble(attachments, extractAttachments)

	info := PageInfo{
		CardInfo:       card,
		AttachemntInfo: attachInfos,
	}

	response, _ := json.Marshal(info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

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
		FileUUID:       fileSchema.ID,
		AttachmentUUID: attachID,
		FragmentID:     "",
		Authors:        docAuthors,
		Conversation:   docConv,
	}
}
