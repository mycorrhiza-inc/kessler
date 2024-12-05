package files

import "fmt"

func EnglishTextFromCompleteFile(file CompleteFileSchema) (string, error) {
	textList := file.DocTexts
	for _, text := range textList {
		if text.Language == "en" {
			return text.Text, nil
		}
	}
	return "", fmt.Errorf("no text found")
}
