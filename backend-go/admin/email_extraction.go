package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/crud"
	"kessler/gen/dbstore"
	"kessler/objects/files"
	"kessler/rag"
	"kessler/routing"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ExtractRelaventEmailsFromOrgUUIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := r.Context()
	// ctx := context.Background()
	org_uuid_str := params["org_uuid"]
	org_uuid, err := uuid.Parse(org_uuid_str)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing org uuid: %v", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	return_obj, err := ExtractRelaventEmailsFromOrgUUID(ctx, rag.DefaultBigLLMModel, org_uuid)
	if err != nil {
		errorstring := fmt.Sprintf("Error extracting emails from org uuid: %v", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(return_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error marshaling return object: %v", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func extractEmails(text string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}
	return matches
}

func extractRelaventEmails(llm rag.LLM, text string, organization_name string) ([]string, error) {
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

func ExtractRelaventEmailsFromFileUUID(ctx context.Context, q dbstore.Queries, llm rag.LLM, file_id uuid.UUID) ([]string, error) {
	file_object, err := crud.CompleteFileSchemaGetFromUUID(ctx, q, file_id)
	if err != nil {
		return []string{}, err
	}
	file_text, err := files.EnglishTextFromCompleteFile(file_object)
	authors := file_object.Authors
	author_list_string := ""
	for _, author := range authors {
		author_list_string += author.AuthorName + "\n"
	}

	emails, err := extractRelaventEmails(llm, file_text, author_list_string)
	if err != nil {
		return emails, err
	}
	return emails, nil
}

type EmailOrgExtraction struct {
	OrganizationName string      `json:"organization_name"`
	OrganizationUUID uuid.UUID   `json:"organization_uuid"`
	FilesCount       int         `json:"files_count"`
	EmailInfo        []EmailInfo `json:"email_info"`
}

type EmailInfo struct {
	Email           string      `json:"email"`
	AssociatedFiles []uuid.UUID `json:"associated_files"`
}

func ExtractRelaventEmailsFromOrgUUID(ctx context.Context, llm rag.LLM, org_id uuid.UUID) (EmailOrgExtraction, error) {
	q := *routing.DBQueriesFromContext(ctx)
	information_map := map[string][]uuid.UUID{}
	org_obj, err := crud.OrgWithFilesGetByID(ctx, &q, org_id)
	if err != nil {
		return EmailOrgExtraction{}, err
	}
	file_obj_list := org_obj.FilesAuthored
	for _, file_obj := range file_obj_list {
		file_id := file_obj.ID
		emails, err := ExtractRelaventEmailsFromFileUUID(ctx, q, llm, file_id)
		if err != nil {
			fmt.Printf("Encountered error getting emails from file with uuid: %v", file_id)
		}
		if err == nil {
			for _, email := range emails {
				results, ok := information_map[email]
				if ok {
					information_map[email] = append(results, file_id)
				} else {
					information_map[email] = []uuid.UUID{file_id}
				}
			}
		}
	}
	email_infos := []EmailInfo{}
	for key, value := range information_map {
		email_info := EmailInfo{Email: key, AssociatedFiles: value}
		email_infos = append(email_infos, email_info)
	}
	return_info := EmailOrgExtraction{
		OrganizationName: org_obj.Name,
		OrganizationUUID: org_obj.ID,
		FilesCount:       len(org_obj.FilesAuthored),
		EmailInfo:        email_infos,
	}
	return return_info, nil
}
