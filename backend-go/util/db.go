package util

import (
	"context"
	"kessler/gen/dbstore"
	"charmbracelet/log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PgPoolConfig(maxConns int32) *pgxpool.Config {
	const defaultMaxConns = int32(30)
	if maxConns == 0 {
		maxConns = defaultMaxConns
	}
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	DATABASE_URL := os.Getenv("DATABASE_CONNECTION_STRING")

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = maxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout
	// Removed to clean up logging in golang
	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		// log.Info("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		// log.Info("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		// log.Info("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func CreateDBContextWithTimeout(timeout time.Duration, maxConns int) context.Context {
	// Create a context with the specified timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Initialize the connection pool.
	connPool, err := pgxpool.NewWithConfig(context.Background(), PgPoolConfig(int32(maxConns)))
	if err != nil {
		log.Fatal("Failed to create database connection pool: ", err)
	}

	// Attach the connection pool to the context.
	ctx = context.WithValue(ctx, "db", connPool)

	// Close the pool and release context resources when the context is done.
	go func() {
		<-ctx.Done()     // Wait until the context is canceled or times out.
		connPool.Close() // Close the connection pool.
		cancel()         // Release the context's resources.
	}()

	return ctx
}

func DBTXFromContext(ctx context.Context) dbstore.DBTX {
	pool := ctx.Value("db").(*pgxpool.Pool)
	// connection, err := pool.Acquire(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	return pool
}

func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
	dbtx := DBTXFromContext(ctx)
	q := dbstore.New(dbtx)
	return q
}

func DBQueriesFromRequest(r *http.Request) *dbstore.Queries {
	ctx := r.Context()
	q := DBQueriesFromContext(ctx)
	return q
}
