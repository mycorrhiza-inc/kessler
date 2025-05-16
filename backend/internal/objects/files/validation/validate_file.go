package validation

import "kessler/internal/objects/files"

func ValidateFile(file files.CompleteFileSchema) error {
	err := FileHasValidAttachments(file)
	if err != nil {
		return err
	}
	var nillerr error
	return nillerr
}
