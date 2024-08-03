import os
from pathlib import Path


DATALAB_API_KEY = os.environ["DATALAB_API_KEY"]
FIREWORKS_EMBEDDING_URL = "https://api.fireworks.ai/inference/v1"
MARKER_ENDPOINT_URL = os.environ["MARKER_ENDPOINT_URL"]

GROQ_API_KEY = os.environ["GROQ_API_KEY"]
OPENAI_API_KEY = os.environ["OPENAI_API_KEY"]
OCTOAI_API_KEY = os.environ["OCTOAI_API_KEY"]
FIREWORKS_API_KEY = os.environ["FIREWORKS_API_KEY"]


OS_TMPDIR = Path(os.getenv("TMPDIR", "/tmp/"))
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")
OS_HASH_FILEDIR = OS_FILEDIR / Path("raw")
OS_OVERRIDE_FILEDIR = OS_FILEDIR / Path("override")
OS_BACKUP_FILEDIR = OS_FILEDIR / Path("backup")


CLOUD_REGION = "sfo3"
S3_ENDPOINT = "https://sfo3.digitaloceanspaces.com"
S3_FILE_BUCKET = "kesslerproddocs"

S3_ACCESS_KEY = os.environ["S3_ACCESS_KEY"]
S3_SECRET_KEY = os.environ["S3_SECRET_KEY"]

REDIS_HOST = os.getenv("REDIS_HOST", "valkey")
REDIS_PORT = int(os.getenv("REDIS_PORT", 6379))


REDIS_PRIORITY_DOCPROC_KEY = "docproc_queue_priority"
REDIS_BACKGROUND_DOCPROC_KEY = "docproc_queue_background"

REDIS_DOCPROC_INFORMATION = "docproc_information"

REDIS_BACKGROUND_DAEMON_ON = "background_daemon"
