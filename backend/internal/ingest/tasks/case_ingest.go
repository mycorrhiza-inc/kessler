package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// IngestOpenscrapersCase processes a case and its associated filings.
// TODO: Implement persistence logic for cases and filings.
func IngestOpenscrapersCase(ctx context.Context, caseInfo *OpenscrapersCaseInfoPayload) error {
	// Example: Log the received case info. Replace with real DB/API calls.
	fmt.Printf("Ingesting case: %s\n", caseInfo.CaseNumber)
	fmt.Printf("Case details: %+v\n", caseInfo)

	// Persist the case (conversation) record.
	// Insert or update conversation based on caseInfo.CaseNumber

	// Iterate over filings and persist each.
	for _, filing := range caseInfo.Filings {
		for _, attachment := range filing.Attachments {
			if attachment.RawAttachment == nil {
			}
		}
		// Create a new filing ingest task.
	}

	return nil
}

func FetchAttachmentDataFromOpenScrapers(attachment AttachmentChildInfo) (RawAttachmentData, error) {
	if attachment.Hash == nil {
		return RawAttachmentData{}, fmt.Errorf("cannot fetch attachment data without hash")
	}

	hashString := attachment.Hash.String()
	url := fmt.Sprintf("https://openscrapers.kessler.xyz/api/raw_attachments/%s/obj", hashString)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RawAttachmentData{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return RawAttachmentData{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RawAttachmentData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result RawAttachmentData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RawAttachmentData{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
