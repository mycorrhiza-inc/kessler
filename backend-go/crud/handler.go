package crud

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func defineCrudRoutes(router *mux.Router) {
	s := router.PathPrefix("/crud").Subrouter()
	s.HandleFunc("/files/{uuid}", getFileHandler)
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hi there!")
		return
	case http.MethodPost:
		fmt.Fprint(w, "POST request")
		return
	case http.MethodPut:
		fmt.Fprintf(w, "PUT request")
	case http.MethodDelete:
		fmt.Fprintf(w, "DELETE request")
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
	}
}

// func makeGetFileHandler(queries *Queries) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fileIDStr := r.URL.Path[len("/files/"):]
// 		fileID, err := uuid.Parse(fileIDStr)
// 		if err != nil {
// 			http.Error(w, "Invalid file ID", http.StatusBadRequest)
// 			return
// 		}
//
// 		file, err := queries.ReadFile(context.Background(), fileID)
// 		if err != nil {
// 			http.Error(w, "File not found", http.StatusNotFound)
// 			return
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 		if err := json.NewEncoder(w).Encode(file); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }
//
// func makeGetMetadataHandler(queries *Queries) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fileIDStr := r.URL.Path[len("/files/metadata/"):]
// 		fileID, err := uuid.Parse(fileIDStr)
// 		if err != nil {
// 			http.Error(w, "Invalid file ID", http.StatusBadRequest)
// 			return
// 		}
//
// 		file, err := queries.ReadFile(context.Background(), fileID)
// 		if err != nil {
// 			http.Error(w, "File not found", http.StatusNotFound)
// 			return
// 		}
//
// 		var metadata map[string]interface{}
// 		if file.Mdata.Valid {
// 			if err := json.Unmarshal([]byte(file.Mdata.String), &metadata); err != nil {
// 				http.Error(w, "Error parsing metadata", http.StatusInternalServerError)
// 				return
// 			}
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 		if err := json.NewEncoder(w).Encode(metadata); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }
