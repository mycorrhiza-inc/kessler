package tasks

import (
	"strings"
	"thaumaturgy/common/objects/conversations"
	"thaumaturgy/common/objects/files"
	"thaumaturgy/common/objects/timestamp"

	"github.com/google/uuid"
)

type ScraperInfoPayload struct {
	FileURL               string                `json:"file_url"`
	Text                  string                `json:"text"`
	FileType              string                `json:"file_type"`
	DocketID              string                `json:"docket_id"`
	PublishedDate         timestamp.KesslerTime `json:"published_date"`
	Name                  string                `json:"name"`
	InternalSourceName    string                `json:"internal_source_name"`
	State                 string                `json:"state"`
	AuthorIndividual      string                `json:"author_individual"`
	AuthorIndividualEmail string                `json:"author_individual_email"`
	AuthorOrganisation    string                `json:"author_organisation"`
	FileClass             string                `json:"file_class"`
	Lang                  string                `json:"lang"`
	ItemNumber            string                `json:"item_number"`
}

func CastScraperInfoToNewFile(info ScraperInfoPayload) files.CompleteFileSchema {
	metadata := map[string]any{
		"url":                 strings.TrimSpace(info.FileURL),
		"docket_id":           strings.TrimSpace(info.DocketID),
		"extension":           strings.TrimSpace(info.FileType),
		"lang":                strings.TrimSpace(info.Lang),
		"title":               strings.TrimSpace(info.Name),
		"source":              strings.TrimSpace(info.InternalSourceName),
		"date":                info.PublishedDate,
		"file_class":          strings.TrimSpace(info.FileClass),
		"author_organisation": strings.TrimSpace(info.AuthorOrganisation),
		"author":              strings.TrimSpace(info.AuthorIndividual),
		"author_email":        strings.TrimSpace(info.AuthorIndividualEmail),
		"item_number":         strings.TrimSpace(info.ItemNumber),
	}
	docket_info := conversations.ConversationInformation{
		DocketGovID: strings.TrimSpace(info.DocketID),
	}
	return files.CompleteFileSchema{
		ID:           uuid.Nil,
		Name:         strings.TrimSpace(info.Name),
		Conversation: docket_info,
		Mdata:        metadata,
	}
}
