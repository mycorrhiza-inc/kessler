package resultcards

import (
	"time"

	"github.com/google/uuid"
)

// Card data types matching the frontend requirements
type CardData interface {
	GetType() string
}

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
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Timestamp      time.Time            `json:"timestamp"`
	ExtraInfo      string               `json:"extraInfo,omitempty"`
	Index          int                  `json:"index"`
	Type           string               `json:"type"`
	ObjectUUID     uuid.UUID            `json:"object_uuid"`
	AttachmentUUID uuid.UUID            `json:"attachment_uuid"`
	FragmentID     string               `json:"fragment_id"`
	Authors        []DocumentAuthor     `json:"authors"`
	Conversation   DocumentConversation `json:"conversation"`
}

func (d DocumentCardData) GetType() string {
	return "document"
}
