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
	"thaumaturgy/common/objects/files"
	"time"

	"github.com/google/uuid"
)

type DatabaseInteraction int

const (
	Insert DatabaseInteraction = iota
	Update
)

func upsertFullFileToDB(ctx context.Context, obj files.CompleteFileSchema, interact DatabaseInteraction) (*files.CompleteFileSchema, error) {
	if MOCK_DB_CONNECTION {
		return &obj, nil
	}

	originalID := obj.ID
	var url string

	switch interact {
	case Insert:
		url = fmt.Sprintf("%s/v2/public/files/insert", KESSLER_API_URL)
	case Update:
		if obj.ID == uuid.Nil {
			return nil, errors.New("cannot update a file with a null uuid")
		}
		url = fmt.Sprintf("%s/v2/public/files/%s/update", KESSLER_API_URL, obj.ID.String())
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
		"named_docket_id": fileObj.Conversation.DocketID,
		"date_string":     fileObj.Mdata["date"],
		"name":            fileObj.Name,
		"extension":       fileObj.Extension,
		"author_string":   fileObj.Mdata["author"],
	}

	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/v2/admin/file-metadata-match", KESSLER_API_URL)

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

	downloadDir := filepath.Join(OS_TMPDIR, "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", nil, err
	}

	resultPath, err := downloadFile(fileURL, downloadDir)
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
	fileManager := S3FileManager{Logger: defaultLogger}

	if err := validateMetadata(fileObj.Mdata); err != nil {
		return "", nil, err
	}

	hashResult, err := fileManager.saveFileToHash(tmpFilePath)
	if err != nil {
		return "", nil, err
	}

	fileObj.Hash = hashResult.Hash
	os.Remove(tmpFilePath)

	if hashResult.DidExist && disableIngestIfHash {
		return "file already exists", &fileObj, nil
	}

	authorNames, _ := fileObj.Mdata["author"].(string)
	authors, err := splitAuthorField(authorNames)
	if err != nil {
		defaultLogger.Error(fmt.Sprintf("Author splitting failed: %v", err))
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

func splitAuthorField(authorStr string) ([]AuthorInformation, error) {
	if authorStr == "" {
		return nil, nil
	}

	// Simplified version splitting on commas
	authors := strings.Split(authorStr, ",")
	var result []AuthorInformation
	for _, a := range authors {
		result = append(result, AuthorInformation{
			AuthorID:   uuid.Nil,
			AuthorName: strings.TrimSpace(a),
		})
	}
	return result, nil
}

func getListAuthors(authors []AuthorInformation) []string {
	var names []string
	for _, a := range authors {
		names = append(names, a.AuthorName)
	}
	return names
}

// Mock implementations and constants
var (
	MOCK_DB_CONNECTION = false
	KESSLER_API_URL    = "http://localhost:8080"
	OS_TMPDIR          = "/tmp"
)

type hashResult struct {
	Hash     string
	DidExist bool
}

func (s *S3FileManager) saveFileToHash(path string) (*hashResult, error) {
	// Implementation assumed
	return &hashResult{Hash: "mockhash", DidExist: false}, nil
}

func downloadFile(url, dir string) (string, error) {
	// Implementation assumed
	return filepath.Join(dir, "downloaded_file"), nil
}
