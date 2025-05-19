package validation

import (
	"fmt"
	"kessler/internal/objects/files"
)

func FileHasValidAttachments(file files.CompleteFileSchema) error {
	for _, attachment := range file.Attachments {
		if attachment.Hash.IsZero() {
			return fmt.Errorf("attachment has null hash")
		}
		for _, text := range attachment.Texts {
			if text.Text == "" {
				return fmt.Errorf("attachment text source has no text")
			}
		}
	}
	var nilerr error
	return nilerr
}
