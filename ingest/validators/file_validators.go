package validators

import (
	"fmt"
	"net/http"
	"os"
	"thaumaturgy/common/objects/files"
)

func ValidateExtensionFromFilepath(filepath string, extension files.KnownFileExtension) error {
	if extension == files.KnownFileExtensionPDF {
		return ValidatePDF(filepath)
	}
	return nil
}

// ValidatePDF validates whether the given file is a valid PDF file.
// It performs checks on the file's header, MIME type, and whether it's a text file.
func ValidatePDF(filepath string) error {
	// Open the file for reading
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed after processing

	// Read the first 4 bytes to check the PDF header
	header := make([]byte, 512)
	_, err = file.Read(header)
	if err != nil {
		return err
	}
	if string(header[:4]) != "%PDF" {
		err := fmt.Errorf("File %s does not have a valid PDF header", filepath)
		return err
	}

	mimeType := http.DetectContentType(header)
	// Check if the MIME type indicates a text file
	if mimeType == "text/plain" {
		err := fmt.Errorf("File %s is a text file", filepath)
		return err
	}
	if mimeType != "application/pdf" {
		err := fmt.Errorf("Invalid MIME type for PDF: %s", mimeType)
		return err
	}

	return nil
}
