package admin

import (
	"kessler/rag"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func extractEmails(text string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}
	return matches
}

func extractRelaventEmails(llm rag.LLMModel, text string, organization_name string) ([]string, error) {
	emails := extractEmails(text)
	prompt := `Given the following list of email addresses and an organization name, identify which emails likely belong to members of the organization.
Consider email patterns, domains, and professional email formats.

Organization: ` + organization_name + `
Email addresses: ` + strings.Join(emails, ", ") + `

Return only the email addresses that are likely associated with the organization, one per line.`

	response, err := rag.SimpleInstruct(llm, prompt)
	if err != nil {
		return emails, err
	}
	relavent_emails := []string{}

	// Split response into lines and clean each line
	for _, line := range strings.Split(response, "\n") {
		email := strings.TrimSpace(line)
		if email != "" && slices.Contains(emails, email) {
			relavent_emails = append(relavent_emails, email)
		}
	}
	return relavent_emails, nil
}

func extractRelaventEmailsFromFileUUID(file_id uuid.UUID) []string {
}
