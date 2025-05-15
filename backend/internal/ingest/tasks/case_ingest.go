package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/internal/ingest/logic"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/files"
	"kessler/pkg/constants"
	"net/http"
	"reflect"
	"time"
)

// IngestOpenscrapersCase processes a case and its associated filings.
// TODO: Implement persistence logic for cases and filings.
func IngestOpenscrapersCase(ctx context.Context, caseInfo *OpenscrapersCaseInfoPayload) error {
	// Example: Log the received case info. Replace with real DB/API calls.
	fmt.Printf("Ingesting case: %s\n", caseInfo.CaseNumber)
	fmt.Printf("Case details: %+v\n", caseInfo)

	minimal_case_info := caseInfo.IntoCaseInfoMinimal()
	err := IngestCaseSpecificData(minimal_case_info)
	if err != nil {
		return err
	}

	// Iterate over filings and persist each.
	for _, filing := range caseInfo.Filings {
		for _, attachment := range filing.Attachments {
			if reflect.ValueOf(attachment.RawAttachment.Hash).IsZero() {
				raw_att, err := FetchAttachmentDataFromOpenScrapers(attachment)
				if err != nil {
					return err
				}
				attachment.RawAttachment = raw_att
			}
			inclusive_filing_info := FilingInfoPayload{
				Filing:   filing,
				CaseInfo: minimal_case_info,
			}
			complete_filing := inclusive_filing_info.IntoCompleteFile()
			_, err := logic.ProcessFileRaw(ctx, &complete_filing, files.DocStatusCompleted)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func FetchAttachmentDataFromOpenScrapers(attachment AttachmentChildInfo) (RawAttachmentData, error) {
	if reflect.ValueOf(attachment.Hash).IsZero() {
		return RawAttachmentData{}, fmt.Errorf("cannot fetch attachment data without hash")
	}

	hashString := attachment.Hash.String()
	fetch_obj_url := fmt.Sprintf("%s/api/raw_attachments/%s/obj", constants.OPENSCRAPERS_API_URL, hashString)

	req, err := http.NewRequest("GET", fetch_obj_url, nil)
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

	fetch_file_url := fmt.Sprintf("%s/api/raw_attachments/%s/raw", constants.OPENSCRAPERS_API_URL, hashString)
	result.GetAttachmentUrl = fetch_file_url

	return result, nil
}

func IngestCaseSpecificData(caseInfoMinimal CaseInfoMinimal) error {
	verify_url := fmt.Sprintf("%s/v2/public/conversations/verify", constants.INTERNAL_KESSLER_API_URL)

	// Transform CaseInfoMinimal to ConversationInformation
	ke_convo_data := conversations.ConversationInformation{
		DocketGovID:    caseInfoMinimal.CaseNumber,
		Name:           caseInfoMinimal.CaseName,
		Description:    caseInfoMinimal.Description,
		MatterType:     caseInfoMinimal.CaseType,
		IndustryType:   caseInfoMinimal.Industry,
		Metadata:       "{}",     // Default empty JSON
		Extra:          "{}",     // Default empty JSON
		DocumentsCount: 0,        // Initial value
		State:          "active", // Default state
	}

	payload, err := json.Marshal(ke_convo_data)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation data: %w", err)
	}

	req, err := http.NewRequest("POST", verify_url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("verification failed: status %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
