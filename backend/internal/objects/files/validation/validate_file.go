package validation

import (
	"fmt"
	"kessler/internal/objects/files"
)

func ValidateFile(file files.CompleteFileSchema) error {
	err := FileHasValidAttachments(file)
	if err != nil {
		return err
	}
	if len(file.Attachments) == 0 {
		return fmt.Errorf("file must have at least 1 attachment")
	}
	if file.Name == "" {
		return fmt.Errorf("file must have a nonempty name")
	}
	var nillerr error
	return nillerr
}
