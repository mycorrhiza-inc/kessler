package handler

import (
	"context"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"

	"github.com/google/uuid"
)

func DeduplicateFileAttachments(ctx context.Context, q *dbstore.Queries, file *files.CompleteFileSchema) (*files.CompleteFileSchema, error) {
	attachments := file.Attachments
	for index, attachment := range attachments {
		new_attachment, err := DeduplicateSingularAttachment(ctx, q, &attachment)
		if err != nil {
			return file, fmt.Errorf("deduplication error: %w", err)
		}
		attachments[index] = *new_attachment
	}
	file.Attachments = attachments
	return file, nil
}

func DeduplicateSingularAttachment(ctx context.Context, q *dbstore.Queries, attachment *files.CompleteAttachmentSchema) (*files.CompleteAttachmentSchema, error) {
	if attachment.ID == uuid.Nil {
		results, err := q.AttachmentListByHash(ctx, attachment.Hash.String())
		if err != nil {
			return &files.CompleteAttachmentSchema{}, fmt.Errorf("database error: %w", err)
		}
		if len(results) > 0 {
			// Idk if this is a good idea or not, but fuck it ship it
			attachment.ID = results[0].ID
		}
		return attachment, nil
	}
	return attachment, nil
}
