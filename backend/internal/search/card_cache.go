
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (s *SearchService) getResultType(facets []string) string {
	for _, facet := range facets {
		if strings.Contains(facet, "/conversation") && !strings.Contains(facet, "/data/conversation") {
			return "conversation"
		}
		if strings.Contains(facet, "/organization") {
			return "organization"
		}
	}
	return "document"
}

func (s *SearchService) getCachedCard(ctx context.Context, key string) (CardData, error) {
	if s.cacheCtrl.Client == nil {
		return nil, fmt.Errorf("cache not available")
	}

	data, err := s.cacheCtrl.Get(key)
	if err != nil {
		return nil, err
	}

	// Try to unmarshal as different card types
	var authorCard AuthorCardData
	if err := json.Unmarshal(data, &authorCard); err == nil && authorCard.Type == "author" {
		return authorCard, nil
	}

	var docketCard DocketCardData
	if err := json.Unmarshal(data, &docketCard); err == nil && docketCard.Type == "docket" {
		return docketCard, nil
	}

	var docCard DocumentCardData
	if err := json.Unmarshal(data, &docCard); err == nil && docCard.Type == "document" {
		return docCard, nil
	}

	return nil, fmt.Errorf("unknown card type in cache")
}
