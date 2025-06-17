package database

import (
	"context"
	"kessler/internal/dbstore"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnPool the global connection pool for the server
// Deprecated: Use Init() which returns the pool instead
var ConnPool *pgxpool.Pool

// GetTx returns a new Queries instance
// Deprecated: Use GetQueries with explicit pool parameter
func GetTx() *dbstore.Queries {
	return dbstore.New(ConnPool)
}

// GetQueries returns a new Queries instance with the given pool
func GetQueries(pool dbstore.DBTX) *dbstore.Queries {
	return dbstore.New(pool)
}

// Init initializes the connection pool with the given number of maximum connections
// Returns the pool for explicit dependency injection
func Init(maxConn int32) (*pgxpool.Pool, error) {
	config := pgPoolConfig(maxConn)
	newPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Keep backwards compatibility by setting global
	ConnPool = newPool

	return newPool, nil
}

// InitWithoutGlobal creates a new pool without setting the global variable
// This is the preferred method for new code
func InitWithoutGlobal(maxConn int32) (*pgxpool.Pool, error) {
	config := pgPoolConfig(maxConn)
	return pgxpool.NewWithConfig(context.Background(), config)
}

func pgPoolConfig(maxConns int32) *pgxpool.Config {
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

// Close closes the connection pool
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
