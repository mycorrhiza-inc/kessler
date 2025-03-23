package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kessler/internal/constants"
	"kessler/pkg/hashes"
	"kessler/pkg/s3utils"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func getIntEnv(key string) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		log.Fatal("Environment variable not set", "key", key)
	}
	var val int
	fmt.Sscanf(valStr, "%d", &val)
	return val
}

func PollMarkerEndpointForResponse(requestCheckURL string, maxPolls int, pollWait int) (string, error) {
	for polls := 0; polls < maxPolls; polls++ {
		time.Sleep(time.Duration(pollWait) * time.Second)
		resp, err := http.Get(requestCheckURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var pollData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&pollData)
		if err != nil {
			return "", err
		}

		status, ok := pollData["status"].(string)
		if !ok {
			return "", fmt.Errorf("status not found in response")
		}

		switch status {
		case "complete":
			markdown, ok := pollData["markdown"].(string)
			if !ok {
				return "", fmt.Errorf("markdown not found in response")
			}
			if markdown == "" {
				return "", fmt.Errorf("got empty string from markdown server")
			}
			log.Info("Processed document after polls", "polls", polls, "text", markdown[:50])
			return markdown, nil
		case "error":
			errorMsg, ok := pollData["error"].(string)
			if !ok {
				return "", fmt.Errorf("error not found in response")
			}
			log.Error("Pdf server encountered an error", "polls", polls, "error", errorMsg)
			return "", fmt.Errorf("pdf server encountered an error after polls: %s", errorMsg)
		default:
			if status != "processing" {
				return "", fmt.Errorf("pdf processing failed. status was unrecognized %s after polls %d", status, polls)
			}
		}
	}

	return "", fmt.Errorf("polling for marker API result timed out")
}

func TranscribePdfS3URI(s3URI string, externalProcess bool, priority bool) (string, error) {
	baseURL := constants.MARKER_ENDPOINT_URL
	var queryStr string
	if priority {
		queryStr = "?priority=true"
	} else {
		queryStr = "?priority=false"
	}
	markerURLEndpoint := baseURL + "/api/v1/marker/direct_s3_url_upload" + queryStr

	data := map[string]string{"s3_url": s3URI}
	log.Info("Sending request to marker server", "data", data)
	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	log.Info("Sending request to marker server", "body", string(body))

	resp, err := http.Post(
		markerURLEndpoint,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	requestCheckURLLeaf, ok := response["request_check_url_leaf"].(string)
	if !ok {
		return "", fmt.Errorf("request_check_url_leaf not found in response")
	}
	requestCheckURL := baseURL + requestCheckURLLeaf
	log.Info("Got response from marker server, polling to see when file is finished processing", "requestCheckURL", requestCheckURL)
	return PollMarkerEndpointForResponse(requestCheckURL, constants.MARKER_MAX_POLLS, constants.MARKER_SECONDS_PER_POLL)
}

func TranscribePDFFromHash(hash hashes.KesslerHash) (string, error) {
	s3_client := s3utils.NewKeFileManager()

	s3_uri, err := s3_client.GetURIFromHash(hash)
	if err != nil {
		return "", err
	}
	return TranscribePdfS3URI(s3_uri, true, true)
}
