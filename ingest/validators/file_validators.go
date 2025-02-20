package validators

import (
	"thaumaturgy/common/objects/files"
)

func ValidateFileVsFileExtension(filepath string, extension files.KnownFileExtension) error {
	if extension == files.KnownFileExtensionPDF {
		return ValidatePDF(filepath)
	}
	return nil
}

// mime = magic.Magic(mime=True)
// file_mime = mime.from_file(filepath)
// is_text = is_text_file(filepath)
//
// match extension:
//
//	case KnownFileExtension.pdf:
//	    with open(filepath, "rb") as pdf_file:
//	        header = pdf_file.read(4)
//	        if header != b"%PDF":
//	            logger.info(f"File with path {filepath} is a valid PDF")
//	            return False, "not a valid pdf header"
//	    if is_text:
//	        return False, "pdf is text file"
//	    if file_mime != "application/pdf":
//	        logger.error(f"Invalid MIME type for PDF: {file_mime}")
//	        return False, "invalid mime type for pdf"
func ValidatePDF(filepath string) error {
	return nil
}
