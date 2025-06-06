package crud

import (
	"context"
	"fmt"
	"kessler/internal/dbstore"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func TestPostgresConnection(t *testing.T) {
	pgConnString := os.Getenv("DATABASE_CONNECTION_STRING")
	ctx := context.Background()

	// conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		t.Errorf("Unable to connect to database")
	}
	defer conn.Close(ctx)
	queries := dbstore.New(conn)
	files, err := queries.FilesList(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing files: %v\n", err)
		t.Errorf("Error Found")

	}
	truncatedFiles := files[:100]
	log.Info("Successfully listed files:", truncatedFiles)
}
