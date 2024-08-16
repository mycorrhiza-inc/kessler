import traceback

from background_loops import initialize_background_loops
from litestar import Litestar, Router
from litestar.config.cors import CORSConfig
from litestar.repository.filters import LimitOffset
from litestar import MediaType, Request, Response
from litestar.status_codes import HTTP_500_INTERNAL_SERVER_ERROR

from litestar.params import Parameter
from litestar.di import Provide
from litestar.contrib.sqlalchemy.base import UUIDBase

from models import utils
from routing.misc_controller import MiscController
from routing.file_controller import FileController
from routing.rag_controller import RagController

from litestar.plugins.structlog import StructlogPlugin, StructlogConfig
from litestar.logging.config import StructLoggingConfig

from routing.daemon_controller import DaemonController

from util.logging import struct_logging_config


sl_config = StructlogConfig()
sl_config.struct_logging_config = struct_logging_config

struct_log_pluging = StructlogPlugin(sl_config)


async def on_startup() -> None:
    async with utils.sqlalchemy_config.get_engine().begin() as conn:
        # UUIDAuditBase extends UUIDBase so create_all should build both
        await conn.run_sync(UUIDBase.metadata.create_all)
    await initialize_background_loops()


def plain_text_exception_handler(request: Request, exc: Exception) -> Response:
    """Default handler for exceptions subclassed from HTTPException."""
    tb = traceback.format_exc()
    request.logger.error(f"exception: {exc}")
    request.logger.error(f"traceback:\n{tb}")
    status_code = getattr(exc, "status_code", HTTP_500_INTERNAL_SERVER_ERROR)
    details = getattr(exc, "detail", "")

    return Response(
        media_type=MediaType.TEXT,
        content=details,
        status_code=status_code,
    )


async def provide_limit_offset_pagination(
    current_page: int = Parameter(ge=1, query="currentPage", default=1, required=False),
    page_size: int = Parameter(
        query="pageSize",
        ge=1,
        default=10,
        required=False,
    ),
) -> LimitOffset:
    """Add offset/limit pagination.

    Return type consumed by `Repository.apply_limit_offset_pagination()`.

    Parameters
    ----------
    current_page : int
        LIMIT to apply to select.
    page_size : int
        OFFSET to apply to select.
    """
    return LimitOffset(page_size, page_size * (current_page - 1))


cors_config = CORSConfig(allow_origins=["*"])

api_router = Router(
    path="/api",
    route_handlers=[
        FileController,
        RagController,
        DaemonController,
        MiscController,
    ],
)

app = Litestar(
    on_startup=[on_startup],
    plugins=[
        utils.sqlalchemy_plugin,
       struct_log_pluging,
    ],
    route_handlers=[api_router],
    dependencies={
        "limit_offset": Provide(provide_limit_offset_pagination),
    },
    cors_config=cors_config,
    exception_handlers={Exception: plain_text_exception_handler},
)
