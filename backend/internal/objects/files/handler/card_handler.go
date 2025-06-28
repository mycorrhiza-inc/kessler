package handler

import (
	"kessler/internal/objects/files"
	"kessler/internal/search"
)

// Example document card data
//
//	type DocumentAuthor struct {
//		AuthorName      string    `json:"author_name"`
//		IsPerson        bool      `json:"is_person"`
//		IsPrimaryAuthor bool      `json:"is_primary_author"`
//		AuthorID        uuid.UUID `json:"author_id"`
//	}
//
//	type DocumentConversation struct {
//		ConvoName   string    `json:"convo_name"`
//		ConvoNumber string    `json:"convo_number"`
//		ConvoID     uuid.UUID `json:"convo_id"`
//	}
//
//	type DocumentCardData struct {
//		Name           string               `json:"name"`
//		Description    string               `json:"description"`
//		Timestamp      time.Time            `json:"timestamp"`
//		ExtraInfo      string               `json:"extraInfo,omitempty"`
//		Index          int                  `json:"index"`
//		Type           string               `json:"type"`
//		ObjectUUID     uuid.UUID            `json:"object_uuid"`
//		AttachmentUUID uuid.UUID            `json:"attachment_uuid"`
//		FragmentID     string               `json:"fragment_id"`
//		Authors        []DocumentAuthor     `json:"authors"`
//		Conversation   DocumentConversation `json:"conversation"`
//	}
//
// Here is the location where files.CompleteFileSchema) search.DocumentCardData {
// /home/nicole/Documents/mycorrhizae/kessler/backend/internal/objects/files/schema.go
// TODO: Impolement this function DO NOT OVERTHINK EVERYTHING JUST THROW OUT AN IMPLEMENTATION I CAN FULLY IMPROVE
func FileSchemaToDocCard(fileSchema files.CompleteFileSchema) search.DocumentCardData {
}
