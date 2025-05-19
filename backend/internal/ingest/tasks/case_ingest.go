package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/internal/ingest/logic"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/files/validation"
	"kessler/pkg/constants"
	"kessler/pkg/logger"
	"net/http"
	"reflect"
	"time"

	"go.uber.org/zap"
)

var log = logger.GetLogger("tasks")

// IngestOpenscrapersCase processes a case and its associated filings.
// TODO: Implement persistence logic for cases and filings.
func IngestOpenscrapersCase(ctx context.Context, caseInfo OpenscrapersCaseInfoPayload) error {
	ctx = context.Background()
	// Example: Log the received case info. Replace with real DB/API calls.
	log.Info("Ingesting case: %s\n", zap.String("case number", caseInfo.CaseNumber))
	log.Info("Case details: %+v\n", zap.Int("filings length", len(caseInfo.Filings)))

	minimal_case_info := caseInfo.IntoCaseInfoMinimal()
	if caseInfo.CaseName == "" {
		caseInfo.CaseName = caseInfo.Description
	}
	err := IngestCaseSpecificData(minimal_case_info)
	if err != nil {
		return err
	}

	// Iterate over filings and persist each.
	for filing_index, filing := range caseInfo.Filings {
		log.Info("Processing nth filing in case", zap.Int("filing_index", filing_index))
		log.Info("Nth Filing has this many attachments", zap.Int("filing_index", filing_index), zap.Int("number_of_attachments", len(filing.Attachments)))
		for attachment_index, attachment := range filing.Attachments {
			log.Info("Processing nth attachment in nth filing in case", zap.Int("filing_index", filing_index), zap.Int("attachment_index", attachment_index))
			if reflect.ValueOf(attachment.RawAttachment.Hash).IsZero() {
				raw_att, err := FetchAttachmentDataFromOpenScrapers(attachment)
				if err != nil {
					log.Error("Encountered error getting attachment data from openscrapers", zap.Error(err))
					// return fmt.Errorf("couldnt get attachment info from openscrapers: %s", err)
				}
				if err == nil {
					caseInfo.Filings[filing_index].Attachments[attachment_index].RawAttachment = raw_att
				}
			}
		}
		inclusive_filing_info := FilingInfoPayload{
			Filing:   filing,
			CaseInfo: minimal_case_info,
		}
		complete_filing := inclusive_filing_info.IntoCompleteFile()
		err := validation.ValidateFile(complete_filing)
		if err != nil {
			log.Error("file was not properly formatted", zap.Error(err))
			return err
		}
		log.Info("Successfully completed conversion into complete file", zap.String("name", complete_filing.Name))
		err = logic.ProcessFile(ctx, complete_filing)
		log.Warn("Made it past the line??")
		if err != nil {
			log.Error("Encountered error processing file", zap.Error(err), zap.String("name", complete_filing.Name))
		}
		logger.Log.Info("Successfully ingested file", zap.String("name", complete_filing.Name))
	}

	return nil
}

// "request failed: Get \"https://openscrapers.kessler.xyz/api/raw_attachments/q24nB9T-EtQ4UAakxSqwVnUl4VNsDZ1FnpgD516x6k8=/obj\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)"
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

		resp_body_bytes := []byte{}
		resp.Body.Read(resp_body_bytes)
		resp_body_string := string(resp_body_bytes)
		log.Error("Encountered error fetching file from openscrapers", zap.Int("code", resp.StatusCode), zap.String("body", resp_body_string))
		return RawAttachmentData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result RawAttachmentData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RawAttachmentData{}, fmt.Errorf("failed to decode response: %w", err)
	}

	fetch_file_url := fmt.Sprintf("%s/api/raw_attachments/%s/raw", constants.OPENSCRAPERS_API_URL, hashString)
	result.GetAttachmentUrl = fetch_file_url
	var nilerr error

	return result, nilerr
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
