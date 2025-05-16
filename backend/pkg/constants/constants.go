package constants

import (
	"os"
	"path/filepath"
	"strconv"
)

var (
	OPENSCRAPERS_API_URL = func() string {
		if v := os.Getenv("OPENSCRAPERS_API_URL"); v != "" {
			return v
		}
		return "https://openscrapers.kessler.xyz"
	}()
	INTERNAL_KESSLER_API_URL = os.Getenv("INTERNAL_KESSLER_API_URL")
	PUBLIC_KESSLER_API_URL   = os.Getenv("PUBLIC_KESSLER_API_URL")

	DATALAB_API_KEY         = os.Getenv("DATALAB_API_KEY")
	FIREWORKS_EMBEDDING_URL = "https://api.fireworks.ai/inference/v1"
	MARKER_ENDPOINT_URL     = os.Getenv("MARKER_ENDPOINT_URL")

	GROQ_API_KEY      = os.Getenv("GROQ_API_KEY")
	OPENAI_API_KEY    = os.Getenv("OPENAI_API_KEY")
	OCTOAI_API_KEY    = os.Getenv("OCTOAI_API_KEY")
	FIREWORKS_API_KEY = os.Getenv("FIREWORKS_API_KEY")
	DEEPINFRA_API_KEY = os.Getenv("DEEPINFRA_API_KEY")

	MARKER_SERVER_URL       = os.Getenv("MARKER_SERVER_URL")
	MARKER_MAX_POLLS        = getEnvDefaultInt("MARKER_MAX_POLLS", 60)
	MARKER_SECONDS_PER_POLL = getEnvDefaultInt("MARKER_SECONDS_PER_POLL", 10)

	OS_TMPDIR           = filepath.Join(getEnvDefault("TMPDIR", "/tmp/"))
	OS_GPU_COMPUTE_URL  = os.Getenv("GPU_COMPUTE_URL")
	OS_FILEDIR          = "/files/"
	OS_HASH_FILEDIR     = filepath.Join(OS_FILEDIR, "raw")
	OS_OVERRIDE_FILEDIR = filepath.Join(OS_FILEDIR, "override")
	OS_BACKUP_FILEDIR   = filepath.Join(OS_FILEDIR, "backup")

	CLOUD_REGION   = "sfo3"
	S3_ENDPOINT    = "https://sfo3.digitaloceanspaces.com"
	S3_FILE_BUCKET = "kesslerproddocs"

	S3_ACCESS_KEY = os.Getenv("S3_ACCESS_KEY")
	S3_SECRET_KEY = os.Getenv("S3_SECRET_KEY")

	REDIS_HOST = getEnvDefault("REDIS_HOST", "valkey")
	REDIS_PORT = getEnvDefaultInt("REDIS_PORT", 6379)

	REDIS_DOCPROC_PRIORITYQUEUE_KEY = "docproc_queue_priority"
	REDIS_DOCPROC_QUEUE_KEY         = "docproc_queue_background"
	REDIS_DOCPROC_INFORMATION       = "docproc_information"

	REDIS_MAIN_PROCESS_LOOP_ENABLED              = "main_process_daemon_enabled"
	REDIS_MAIN_PROCESS_LOOP_CONFIG               = "main_process_loop_config"
	REDIS_DOCPROC_BACKGROUND_PROCESSING_STOPS_AT = "docproc_background_stop_at"
	REDIS_DOCPROC_CURRENTLY_PROCESSING_DOCS      = "docproc_currently_processing_docs"
)

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
