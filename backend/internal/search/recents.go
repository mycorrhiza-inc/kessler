package search

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"kessler/internal/quickwit"
// 	"net/http"

// 	"go.uber.org/zap"
// )

// func GetRecentCaseData(maxHits int, offset int) ([]SearchDataHydrated, error) {
// 	// get all documents with a metadata.date_filed since (x)
// 	log.Info("getting recent case data", zap.Int("maxHits", maxHits), zap.Int("offset", offset))
// 	request := SearchRequest{
// 		Index:       "NY_PUC",
// 		Query:       "",
// 		SortBy:      []string{"date_filed"},
// 		MaxHits:     maxHits,
// 		StartOffset: offset,
// 	}
// 	data, err := SearchQuickwit(request)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Info("received search data", zap.Any("data", data))
// 	if data == nil {
// 		empty := []SearchDataHydrated{}
// 		return empty, nil
// 	}
// 	return data, nil
// }

// func GetCaseDataSince(date string, page int) ([]Hit, error) {
// 	// Go was complaining about unutilized code, assume this is someone in the middle of something, feel free to continue.
// 	// parse the date string
// 	//
// 	// parsedDate, err := convertToUTC(date)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// if the date string is incorrect return a failure
// 	// the failure should be handled on the frontend
// 	maxHits := 40
// 	request := quickwit.QuickwitSearchRequest{
// 		Query:         "",
// 		SnippetFields: "text",
// 		MaxHits:       maxHits,
// 	}

// 	jsonData, err := json.Marshal(request)
// 	// ===== submit request to quickwit =====
// 	if err != nil {
// 		log.Error("error marshalling quickwit request", zap.Error(err))
// 		return nil, err
// 	}

// 	offset := page * maxHits
// 	// get all documents with a metadata.date_filed since (x)
// 	request_url := fmt.Sprintf("%s/api/v1/dockets/search?sort_by=date_filed?max_hits=20?start_offset=%d", quickwit.QuickwitURL, offset)
// 	log.Info("making quickwit request", zap.String("url", request_url))
// 	resp, err := http.Post(
// 		request_url,
// 		"application/json",
// 		bytes.NewBuffer(jsonData),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Info("received quickwit response", zap.Any("response", resp))

// 	defer resp.Body.Close()
// 	cases := []Hit{}

// 	return cases, nil
// }
