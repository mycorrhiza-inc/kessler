package crud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func GetFileWithMeta(config FileHandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(config.dbtx_val)
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			errorstring := fmt.Sprintf("Error parsing file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		file_raw, err := q.GetFileWithMetadata(ctx, parsedUUID)
		if err != nil {

			errorstring := fmt.Sprintf("Error Retriving file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusNotFound)
			return
		}
		var mdata_obj map[string]interface{}
		err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
		if err != nil {
			errorstring := fmt.Sprintf("Error Unmarshalling file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		file := CompleteFileSchema{
			ID:        file_raw.ID,
			Verified:  file_raw.Verified.Bool,
			Extension: file_raw.Extension,
			Lang:      file_raw.Lang,
			Name:      file_raw.Name,
			Hash:      file_raw.Hash,
			IsPrivate: file_raw.Isprivate.Bool,
			Mdata:     mdata_obj,
		}

		response, _ := json.Marshal(file)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func FileSemiCompleteGetFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			errorstring := fmt.Sprintf("Error parsing file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		file, err := SemiCompleteFileGetFromUUID(ctx, q, parsedUUID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, _ := json.Marshal(file)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func SemiCompleteFileGetFromUUID(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) (CompleteFileSchema, error) {
	files_raw, err := q.SemiCompleteFileGet(ctx, uuid)
	if err != nil {
		errorstring := fmt.Sprintf("Error Retriving file %v: %v\n", uuid, err)
		return CompleteFileSchema{}, errors.New(errorstring)
	}
	if len(files_raw) == 0 {
		errorstring := fmt.Sprintf("Error No Files Found for a list of length zero.\n")
		return CompleteFileSchema{}, errors.New(errorstring)
	}
	file_raw := files_raw[0]
	var mdata_obj map[string]interface{}
	err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file metadata %v: %v\n", uuid, err)
		return CompleteFileSchema{}, errors.New(errorstring)
	}
	var extra_obj FileGeneratedExtras
	err = json.Unmarshal(file_raw.ExtraObj, &extra_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file extras %v: %v\n", uuid, err)
		return CompleteFileSchema{}, errors.New(errorstring)
	}
	// Missing info here, it doesnt have the name.
	conv_info := ConversationInformation{ID: file_raw.DocketUuid}
	author_info := make([]AuthorInformation, len(files_raw))
	for i, author_file_raw := range files_raw {
		author_info[i] = AuthorInformation{
			AuthorName:      author_file_raw.OrganizationName.String,
			IsPerson:        author_file_raw.IsPerson.Bool,
			IsPrimaryAuthor: author_file_raw.IsPrimaryAuthor.Bool,
			AuthorID:        author_file_raw.OrganizationID,
		}
	}

	file := CompleteFileSchema{
		ID:           file_raw.ID,
		Verified:     file_raw.Verified.Bool,
		Extension:    file_raw.Extension,
		Lang:         file_raw.Lang,
		Name:         file_raw.Name,
		Hash:         file_raw.Hash,
		Mdata:        mdata_obj,
		Extra:        extra_obj,
		Conversation: conv_info,
		Authors:      author_info,
	}
	return file, nil
}

func FileTextsGetAll(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) ([]FileChildTextSource, error) {
	texts, err := q.ListTextsOfFile(ctx, uuid)
	if err != nil {
		return make([]FileChildTextSource, 0), err
	}
	return_texts := make([]FileChildTextSource, len(texts))
	for i, text := range texts {
		return_texts[i] = FileChildTextSource{
			IsOriginalText: text.IsOriginalText,
			Text:           text.Text,
			Language:       text.Language,
		}
	}
	return return_texts, nil
}

func FileStageGet(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) (DocProcStage, error) {
	stage_str, err := q.StageLogFileGetLatest(ctx, uuid)
	if err != nil {
		return DocProcStage{}, err
	}
	stage := DocProcStage{}
	err = json.Unmarshal(stage_str.Log, &stage)
	if err != nil {
		return stage, err
	}
	return stage, nil
}

func CompleteFileSchemaGetFromUUID(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) (CompleteFileSchema, error) {
	file, err := SemiCompleteFileGetFromUUID(ctx, q, uuid)
	if err != nil {
		return CompleteFileSchema{}, err
	}
	texts, err := FileTextsGetAll(ctx, q, uuid)
	if err != nil {
		return CompleteFileSchema{}, err
	}
	file.DocTexts = texts
	stage, err := FileStageGet(ctx, q, uuid)
	if err != nil {
		return CompleteFileSchema{}, err
	}
	file.Stage = stage
	return file, nil
}

func UnverifedCompleteFileSchemaList(ctx context.Context, q dbstore.Queries, max_responses uint) ([]CompleteFileSchema, error) {
	files, err := q.FilesListUnverified(ctx)
	if err != nil {
		return []CompleteFileSchema{}, err
	}
	unverified_raw_uuids := make([]uuid.UUID, len(files))
	for i, file := range files {
		unverified_raw_uuids[i] = file.ID
	}
	// Shuffle the uuids around to get a random selection while processing
	for i := range unverified_raw_uuids {
		j := rand.Intn(i + 1) // Want to get a range of [0,i] so that there is a possibility of the null swap.
		// Inductive proof this distributes the elements randomly at step k:
		// The element at index k is evenly distributed, since it has a 1/k chance of being the element at k, and a k-1/k chance of selecting from an even distribution of k-1 elements, thus meaning it has an even distribution of 1/k probability of selecting k elements.
		// Same thing for other elements, it has a k-1/k chance of sampling from its EXISTING even distribution of k-1 elements, and a 1/k chance of swapping with the k'th element. Thus it has a even 1/k chance of being selected from every element.
		unverified_raw_uuids[i], unverified_raw_uuids[j] = unverified_raw_uuids[j], unverified_raw_uuids[i]
	}
	if len(unverified_raw_uuids) > int(max_responses) {
		unverified_raw_uuids = unverified_raw_uuids[:max_responses]
	}

	complete_files := []CompleteFileSchema{}
	fileChan := make(chan CompleteFileSchema)
	errChan := make(chan error)
	var wg sync.WaitGroup

	for _, uuid := range unverified_raw_uuids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			complete_file, err := CompleteFileSchemaGetFromUUID(ctx, q, uuid)
			if err != nil {
				fmt.Printf("Error getting file %v: %v\n", uuid, err)
				// errChan <- err
				return
			}
			fileChan <- complete_file
		}()
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(fileChan)
		close(errChan)
	}()

	// Collect results
	for file := range fileChan {
		complete_files = append(complete_files, file)
	}
	return complete_files, nil
}

func ReadFileHandlerFactory(config FileHandlerConfig) http.HandlerFunc {
	private := config.private
	dbtx_val := config.dbtx_val
	return_type := config.return_type

	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		token := r.Header.Get("Authorization")
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		if private {

			isAuthorized, err := checkPrivateFileAuthorization(q, ctx, parsedUUID, token)
			if !isAuthorized {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
			if err != nil {
				fmt.Printf("Ran into the follwing error with authentication $v", err)
			}
		}
		// Since all three of these methods share the same authentication and database connection prerecs
		// switching functionality using an if else, or a cases switch lets code get reused
		// TODO: This is horrible, I need to refactor
		file_params := GetFileParam{
			q, ctx, parsedUUID, private,
		}
		switch return_type {
		case "raw":
			file, err := GetFileObjectRaw(file_params)
			if err != nil {
				error_string := fmt.Sprintf("Error retrieving file object %v", err)
				fmt.Println(error_string)
				http.Error(w, error_string, http.StatusNotFound)
				return
			}
			filehash := file.Hash
			kefiles := NewKeFileManager()
			file_path, err := kefiles.downloadFileFromS3(filehash)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error encountered when getting file with hash %v from s3:%v", filehash, err), http.StatusInternalServerError)
				return
			}
			content, err := os.ReadFile(file_path)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
				return
			}

			mimeType := http.DetectContentType(content)
			// if mimeType == "application/octet-stream" {
			// 	mimeType = "application/pdf" // Default to PDF if mime type can't be determined
			// }

			w.Header().Set("Content-Type", mimeType)
			w.Write(content)
		case "markdown":
			originalLang := r.URL.Query().Get("original_lang") == "true"
			matchLang := r.URL.Query().Get("match_lang")
			// TODO: Add suport for non english text retrieval and original text retrieval
			markdownText, err := GetSpecificFileText(file_params, matchLang, originalLang)
			if err != nil {
				http.Error(w, "Error retrieving texts or no texts found that mach query params", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(markdownText))
		case "object-minimal":
			file, err := GetFileObjectRaw(file_params)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			inflated_schema := CompleteFileSchemaInflateFromPartialSchema(file)

			response, _ := json.Marshal(inflated_schema)

			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		default:
			fmt.Printf("Encountered unreachable code with file return type %v", return_type)
			http.Error(w, "Congradulations for encountering unreachable code about support types!", http.StatusInternalServerError)
		}
	}
}
