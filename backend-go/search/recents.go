package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CaseFile struct {
	DocketId string `json:"docket_id"`
	FileId   string `json:"file_id"`
}

func GetRecentCaseData(page int) ([]CaseFile, error) {
	request := QuickwitSearchRequest{
		Query:         "",
		SnippetFields: "text",
		MaxHits:       20,
	}

	jsonData, err := json.Marshal(request)

	// ===== submit request to quickwit =====
	log.Printf("jsondata: \n%s", jsonData)
	if err != nil {
		log.Printf("Error Marshalling quickwit request: %s", err)
		return nil, err
	}

	offset := page * 20
	// get all documents with a metadata.date_filed since (x)
	request_url := fmt.Sprintf("%s/api/v1/dockets/search?sort_by=date_filed?max_hits=20?start_offset=%d", quickwitURL, offset)
	resp, err := http.Post(
		request_url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	cases := []CaseFile{}

	return cases, nil
}

func convertToUTC(dateStr string) (string, error) {
	// Define common date formats to try
	formats := []string{
		time.RFC3339,           // 2024-10-23T14:35:22Z
		time.RFC1123,           // Mon, 02 Jan 2006 15:04:05 MST
		"2006-01-02 15:04:05",  // 2024-10-23 14:35:22
		"02/01/2006 15:04:05",  // 23/10/2024 14:35:22
		"02 Jan 2006 15:04:05", // 23 Oct 2024 14:35:22
		"2006-01-02",           // 2024-10-23
		"02/01/2006",           // 23/10/2024
	}

	var parsedTime time.Time
	var err error

	// Attempt to parse the date string with each format
	for _, layout := range formats {
		parsedTime, err = time.Parse(layout, dateStr)
		if err == nil {
			// Successfully parsed, convert to UTC
			return parsedTime.UTC().Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("could not parse date: %s", dateStr)
}

func GetCaseDataSince(date string, page int) ([]CaseFile, error) {
	// parse the date string
	//
	parsedDate, err := convertToUTC(date)

	if err != nil {
		return nil, err
	}

	// if the date string is incorrect return a failure
	// the failure should be handled on the frontend
	request := QuickwitSearchRequest{
		Query:         "",
		SnippetFields: "text",
		MaxHits:       20,
	}

	jsonData, err := json.Marshal(request)

	// ===== submit request to quickwit =====
	log.Printf("jsondata: \n%s", jsonData)
	if err != nil {
		log.Printf("Error Marshalling quickwit request: %s", err)
		return nil, err
	}

	offset := page * 20
	// get all documents with a metadata.date_filed since (x)
	request_url := fmt.Sprintf("%s/api/v1/dockets/search?sort_by=date_filed?max_hits=20?start_offset=%d", quickwitURL, offset)
	resp, err := http.Post(
		request_url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	cases := []CaseFile{}

	return cases, nil
}
