package tasks

import (
	"context"
	"fmt"
	"strings"
	"thaumaturgy/common/objects/conversations"
	"thaumaturgy/common/objects/files"
	"thaumaturgy/common/objects/timestamp"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type KesslerTaskInfo struct {
	TaskID string `json:"task_id"`
	Queue  string `json:"queue"`
	State  string `json:"state"`
	Status string `json:"status"`
}

func GenerateTaskInfoFromInfo(info asynq.TaskInfo) KesslerTaskInfo {
	return KesslerTaskInfo{
		TaskID: info.ID,
		Queue:  info.Queue,
	}
}

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

type CastableIntoScraperInfo interface {
	IntoScraperInfo() (ScraperInfoPayload, error)
}

func (s ScraperInfoPayload) IntoScraperInfo() (ScraperInfoPayload, error) {
	return s, nil
}

func AddScraperTaskCastable(ctx context.Context, castable CastableIntoScraperInfo) (KesslerTaskInfo, error) {
	scraper_info, err := castable.IntoScraperInfo()
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("Error Casting to Scraper Info: %v", err)
	}
	task, err := NewAddFileScraperTask(scraper_info)
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("Error creating task: %v", err)
	}

	// Get client from context
	client := GetClient(ctx)
	info, err := client.Enqueue(task)
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("Error enqueueing task: %v", err)
	}

	kessler_info := GenerateTaskInfoFromInfo(*info)

	return kessler_info, nil
}

type NYPUCDocInfo struct {
	Serial       string `json:"serial"`
	DateFiled    string `json:"date_filed"`
	NYPUCDocType string `json:"nypuc_doctype"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Organization string `json:"organization"`
	ItemNo       string `json:"item_no"`
	FileName     string `json:"file_name"`
	DocketID     string `json:"docket_id"`
}

func (n NYPUCDocInfo) IntoScraperInfo() (ScraperInfoPayload, error) {
	regular_time, err := time.Parse("01/02/2006", n.DateFiled)
	if err != nil {
		return ScraperInfoPayload{}, nil
	}

	return ScraperInfoPayload{
		FileURL:            n.URL,
		Name:               n.Name,
		FileType:           "",
		DocketID:           n.DocketID,
		PublishedDate:      timestamp.KesslerTime(regular_time),
		InternalSourceName: "NYPUC",
		State:              "NY",
		AuthorOrganisation: n.Organization,
		FileClass:          n.NYPUCDocType,
		ItemNumber:         n.ItemNo,
	}, nil
}
