package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetRecentCaseData(page int) ([]SearchData, error) {
	// get all documents with a metadata.date_filed since (x)
	log.Printf("gettings ssssssflkjadflhdsfuhlifadlhf")
	request := SearchRequest{
		Index:       "NY_PUC",
		Query:       "",
		SortBy:      []string{"date_filed"},
		MaxHits:     20,
		StartOffset: page,
	}
	data, err := SearchQuickwit(request)
	if err != nil {
		return nil, err
	}
	log.Printf("data: \n%v", data)
	return data, nil
}

func GetCaseDataSince(date string, page int) ([]Hit, error) {
	// Go was complaining about unutilized code, assume this is someone in the middle of something, feel free to continue.
	// parse the date string
	//
	// parsedDate, err := convertToUTC(date)

	// if err != nil {
	// 	return nil, err
	// }

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
	log.Printf("request_url: \n%s", request_url)
	resp, err := http.Post(
		request_url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	log.Printf("response: \n%v", resp)

	defer resp.Body.Close()
	cases := []Hit{}

	return cases, nil
}
