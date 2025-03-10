// Create an enum called known file extension with the following file types
package files

import (
	"fmt"
	"strings"
)

type KnownFileExtension string

const (
	KnownFileExtensionPDF     = KnownFileExtension("pdf")
	KnownFileExtensionXLSX    = KnownFileExtension("xlsx")
	KnownFileExtensionDOCX    = KnownFileExtension("docx")
	KnownFileExtensionHTML    = KnownFileExtension("html")
	KnownFileExtensionMD      = KnownFileExtension("md")
	KnownFileExtensionUnknown = KnownFileExtension("unknown")
)

var KnownExtensionsDict = map[string]KnownFileExtension{
	"pdf":  KnownFileExtensionPDF,
	"xlsx": KnownFileExtensionXLSX,
	"docx": KnownFileExtensionDOCX,
	"html": KnownFileExtensionHTML,
	"md":   KnownFileExtensionMD,
}

// GetKnownFileExtension converts a string to a KnownFileExtension type.
// If the extension is not recognized, it returns KnownFileExtensionUnknown and an error.
func FileExtensionFromString(ext string) (KnownFileExtension, error) {
	// Convert to lowercase to make the comparison case-insensitive
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	if extension, exists := KnownExtensionsDict[ext]; exists {
		return extension, nil
	}

	return KnownFileExtensionUnknown, fmt.Errorf("unknown file extension: %s", ext)
}
