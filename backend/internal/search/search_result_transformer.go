package search

import (
	"context"
	"fmt"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"time"

	"go.uber.org/zap"
)

// Update ProcessSearch to use the new transformer
// Updated transformSearchResponse to return card data
func (s *SearchService) transformSearchResponse(ctx context.Context, fuguResponse *fugusdk.SanitizedResponse, query, namespace string, pagination PaginationParams, processTime time.Duration) (*SearchResponse, error) {
	if fuguResponse == nil || len(fuguResponse.Results) == 0 {
		return &SearchResponse{
			Data:        []CardData{},
			Total:       0,
			Page:        pagination.Page,
			PerPage:     pagination.Limit,
			Query:       query,
			Namespace:   namespace,
			ProcessTime: processTime.String(),
		}, nil
	}

	var cards []CardData

	for i, result := range fuguResponse.Results {
		resultType := s.getResultType(result.Facets)

		var card CardData
		var err error

		switch resultType {
		case "conversation":
			card, err = s.HydrateConversation(ctx, result.ID, result.Score, i)
		case "organization":
			card, err = s.HydrateOrganization(ctx, result.ID, result.Score, i)
		default:
			card, err = s.HydrateDocument(ctx, result, i)
			fmt.Printf("card returned::: %v\n\n\n\n", card)
		}

		if err != nil {
			logger.Warn(ctx, "failed to hydrate result",
				zap.String("id", result.ID),
				zap.String("type", resultType),
				zap.Error(err))
			continue // Skip this result
		}

		cards = append(cards, card)
	}

	return &SearchResponse{
		Data:        cards,
		Total:       fuguResponse.Total,
		Page:        pagination.Page,
		PerPage:     pagination.Limit,
		Query:       query,
		Namespace:   namespace,
		ProcessTime: processTime.String(),
	}, nil
}
