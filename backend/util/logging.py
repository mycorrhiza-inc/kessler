# logging configuration
from litestar.logging import LoggingConfig
from litestar.plugins.structlog import StructlogConfig
from litestar.logging.config import StructLoggingConfig

from constants import KESSLER_LOG_DIR

import structlog
import logging

timestamper = structlog.processors.TimeStamper(fmt="iso")

logging_config = LoggingConfig(
    root={"level": logging.getLevelName(logging.INFO), "handlers": ["console"]},
    formatters={
        "standard": {"format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s"},
    },
    handlers={
        "default": {
            "level": "DEBUG",
            "class": "logging.StreamHandler",
        },
        "file": {
            "level": "DEBUG",
            "class": "logging.handlers.WatchedFileHandler",
            "filename": "/logs/kessler.log",
            "formatter": "standard",
        },
    },
    loggers={
        "litestar": {
            "level": "INFO",
            "handlers": ["queue_listener", "file"],
            "propagate": False,
        },
    },
)


struct_logging_config = StructLoggingConfig(
    standard_lib_logging_config=logging_config,
    traceback_line_limit=10,
    pretty_print_tty=True,
)
