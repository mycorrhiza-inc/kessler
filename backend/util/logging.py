# logging configuration
from litestar.logging import LoggingConfig
from litestar.plugins.structlog import StructlogConfig
import logging

logging_config = LoggingConfig(
    root={"level": logging.getLevelName(logging.INFO), "handlers": ["console"]},
    formatters={
        "standard": {"format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s"}
    },
)

structlog_config = StructlogConfig(middleware_logging_config=logging_config)
