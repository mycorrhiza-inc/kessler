package files

import "fmt"

func EnglishTextFromCompleteFile(file CompleteFileSchema) (string, error) {
	if len(file.Attachments) == 0 {
		return "", fmt.Errorf("no attachments found")
	}
	attachment := file.Attachments[0]
	return EnglishTextFromAttachment(attachment)
}

func EnglishTextFromAttachment(attachment CompleteAttachmentSchema) (string, error) {
	textList := attachment.Texts
	for _, text := range textList {
		if text.Language == "en" {
			return text.Text, nil
		}
	}
	return "", fmt.Errorf("no text found")
}
