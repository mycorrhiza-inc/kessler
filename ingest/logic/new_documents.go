package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"thaumaturgy/common/constants"
	"thaumaturgy/common/objects/authors"
	"thaumaturgy/common/objects/files"
	"thaumaturgy/common/s3utils"
	"thaumaturgy/validators"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type DatabaseInteraction int

const (
	Insert DatabaseInteraction = iota
	Update
)

func upsertFullFileToDB(ctx context.Context, obj files.CompleteFileSchema, interact DatabaseInteraction) (*files.CompleteFileSchema, error) {
	// if constants.MOCK_DB_CONNECTION {
	// 	return &obj, nil
	// }

	originalID := obj.ID
	var url string

	switch interact {
	case Insert:
		url = fmt.Sprintf("%s/v2/public/files/insert", constants.KESSLER_INTERNAL_API_URL)
	case Update:
		if obj.ID == uuid.Nil {
			return nil, errors.New("cannot update a file with a null uuid")
		}
		url = fmt.Sprintf("%s/v2/public/files/%s/update", constants.KESSLER_INTERNAL_API_URL, obj.ID.String())
	default:
		return &obj, nil
	}

	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	newUUID, err := uuid.Parse(response.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID in response: %w", err)
	}

	if newUUID == uuid.Nil {
		return nil, errors.New("received null UUID from server")
	}

	if interact == Insert && newUUID == originalID {
		return nil, errors.New("identical ID returned from server during insert")
	}

	obj.ID = newUUID
	return &obj, nil
}

func checkPgForDuplicateMetadata(ctx context.Context, fileObj files.CompleteFileSchema) (*files.CompleteFileSchema, error) {
	payload := map[string]interface{}{
		"named_docket_id": fileObj.Conversation.DocketGovID,
		"date_string":     fileObj.Mdata["date"],
		"name":            fileObj.Name,
		"extension":       fileObj.Extension,
		"author_string":   fileObj.Mdata["author"],
	}

	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/v2/admin/file-metadata-match", constants.KESSLER_INTERNAL_API_URL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var result files.CompleteFileSchema
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func addURLRaw(ctx context.Context, fileURL string, fileObj files.CompleteFileSchema, disableIngestIfHash bool) (string, *files.CompleteFileSchema, error) {
	existingFile, err := checkPgForDuplicateMetadata(ctx, fileObj)
	if err == nil && existingFile != nil {
		fileObj = *existingFile
		disableIngestIfHash = true
	}

	downloadDir := filepath.Join(constants.OS_TMPDIR, "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", nil, err
	}

	resultPath, err := s3utils.DownloadFile(fileURL, downloadDir)
	if err != nil {
		return "", nil, err
	}

	if fileObj.Extension == "" {
		ext := filepath.Ext(resultPath)
		if ext != "" {
			fileObj.Extension = ext[1:]
		}
	}

	return addFileRaw(ctx, resultPath, fileObj, disableIngestIfHash)
}

func addFileRaw(ctx context.Context, tmpFilePath string, fileObj files.CompleteFileSchema, disableIngestIfHash bool) (string, *files.CompleteFileSchema, error) {
	if err := validateMetadata(fileObj.Mdata); err != nil {
		return "", nil, err
	}
	extension, err := files.FileExtensionFromString(fileObj.Extension)
	if err != nil {
		return "", nil, err
	}
	err = validators.ValidateExtensionFromFilepath(tmpFilePath, extension)
	if err != nil {
		return "", nil, err
	}

	fileManager := s3utils.NewKeFileManager()
	hashResult, err := fileManager.UploadFileToS3(tmpFilePath)
	if err != nil {
		return "", nil, err
	}

	fileObj.Hash = hashResult.String()
	os.Remove(tmpFilePath)

	authorNames, _ := fileObj.Mdata["author"].(string)
	authors, err := splitAuthorField(authorNames)
	if err != nil {
		log.Error("Author splitting failed", "err", err)
	}

	fileObj.Authors = authors
	fileObj.Mdata["authors"] = getListAuthors(authors)

	return "", &fileObj, nil
}

// Helper functions and assumed implementations
func validateMetadata(metadata map[string]interface{}) error {
	if lang, _ := metadata["lang"].(string); lang == "" {
		metadata["lang"] = "en"
	}

	if ext, _ := metadata["extension"].(string); ext != "" && ext[0] == '.' {
		metadata["extension"] = ext[1:]
	}

	// Additional validation logic
	return nil
}

func splitAuthorField(authorStr string) ([]authors.AuthorInformation, error) {
	if authorStr == "" {
		return nil, nil
	}

	// Simplified version splitting on commas
	authors_obj := strings.Split(authorStr, ",")
	var result []authors.AuthorInformation
	for _, a := range authors_obj {
		result = append(result, authors.AuthorInformation{
			AuthorID:   uuid.Nil,
			AuthorName: strings.TrimSpace(a),
		})
	}
	return result, nil
}

func getListAuthors(authors []authors.AuthorInformation) []string {
	var names []string
	for _, a := range authors {
		names = append(names, a.AuthorName)
	}
	return names
}

// Mock implementations and constants

type hashResult struct {
	Hash     string
	DidExist bool
}
