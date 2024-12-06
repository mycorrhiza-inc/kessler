package networking

import (
	"net/http"
	"strconv"
)

type BasePaginationNetworkSchema struct {
	Limit  uint `json:"limit"`
	Offset uint `json:"offset"`
}

func PaginationFromUrlParams(r *http.Request) BasePaginationNetworkSchema {
	params := r.URL.Query()
	limit := 30 // default limit
	offset := 0 // default offset

	if limitStr := params.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := params.Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	return BasePaginationNetworkSchema{Limit: uint(limit), Offset: uint(offset)}
}
