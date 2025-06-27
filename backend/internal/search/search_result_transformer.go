package search

import (
	"context"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"time"

	"go.uber.org/zap"
)

// Update ProcessSearch to use the new transformer
// Updated transformSearchResponse to return card data
func (s *SearchService) transformSearchResponse(ctx context.Context, fuguResponse *fugusdk.SanitizedResponse, query, namespace string, pagination PaginationParams, processTime time.Duration) (*SearchResponse, error) {
	log := logger.FromContext(ctx)
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
		log.Debug("Debugging search result",
			zap.String("result_id", result.ID),
			zap.String("text", result.Text[:100]),
			zap.Any("metadata", result.Metadata),
			zap.Any("facets", result.Facets),
			zap.String("detected_type", resultType),
		)

		var card CardData
		var err error

		switch resultType {
		case "conversation":
			card, err = s.HydrateConversation(ctx, result.ID, result.Score, i)
		case "organization":
			card, err = s.HydrateOrganization(ctx, result.ID, result.Score, i)
		default:
			card, err = s.HydrateDocument(ctx, result, i)
		}

		if err != nil {
			logger.Warn(ctx, "failed to hydrate result",
				zap.String("id", result.ID),
				zap.String("type", resultType),
				zap.Error(err))
			continue // Skip this result
		}

		log.Info("Debugging returned card result",
			zap.Any("card", card),
		)

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
