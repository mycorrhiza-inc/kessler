package admin

import (
	"context"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/quickwit"
	"net/http"

	"github.com/google/uuid"
)

func HandleQuickwitIngestFromPostgresFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := QuickwitIngestFromPostgresMain(dbtx_val, ctx)
		if err != nil {
			errorstring := fmt.Sprintf("Error ingesting from postgres: %v", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Sucessfully ingested from postgres"))
	}
}

func QuickwitIngestFromPostgresMain(dbtx_val dbstore.DBTX, ctx context.Context) error {
	q := *dbstore.New(dbtx_val)
	indexName := "idk what the right index name is"
	err := quickwit.ClearIndex(indexName)
	if err != nil {
		return err
	}

	files, err := q.FilesList(ctx)
	if err != nil {
		return err
	}

	ids := make([]uuid.UUID, len(files))
	for i, file := range files {
		ids[i] = file.ID
	}
	complete_files, err := CompleteFileSchemasFromUUIDs(ctx, q, ids)
	if err != nil {
		return err
	}
	fmt.Printf("Got %d complete files\n", len(complete_files))
	for _, complete_file := range complete_files {
		err = quickwit.IndexFile(complete_file)
		if err != nil {
			return err
		}
	}

	return nil
}
