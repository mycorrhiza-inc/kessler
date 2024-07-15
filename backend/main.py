import logging
import traceback

from litestar import Litestar, Router
from litestar.config.cors import CORSConfig
from litestar.repository.filters import LimitOffset
from litestar import MediaType, Request, Response
from litestar.status_codes import HTTP_500_INTERNAL_SERVER_ERROR

from litestar.params import Parameter
from litestar.di import Provide
from litestar.contrib.sqlalchemy.base import UUIDBase

from models import utils
from routing.test_controller import TestController
from util.logging import logging_config
from routing.file_controller import FileController
from routing.search_controller import SearchController
from routing.rag_controller import RagController

from lance_store.connection import get_lance_connection

from litestar.events import listener
from lance_store.connection import ensure_fts_index
from rag.llamaindex import initialize_db_table
import threading


from routing.daemon_controller import DaemonController, process_document

logger = logging.getLogger(__name__)


added_docs = -1


def full_fts_reindex() -> None:
    threading.Timer(10.0, full_fts_reindex).start()
    global added_docs
    if added_docs > 0 or added_docs == -1:
        ensure_fts_index()
        added_docs = 0
        logger.info("detected new doc, successfully reindexed FTS")
        return


@listener("increment_processed_docs")
def increment_processed_docs(num: int) -> None:
    if num > 0:
        global added_docs
        added_docs += num


async def on_startup() -> None:
    logger.debug("running startup")
    logger.debug("ensuring tables exist")
    try:
        initialize_db_table()
    except Exception as e:
        logger.error(f"catastrophic failure setting up lancedb: {e}")
        # What if we didnt?
        # raise e

    async with utils.sqlalchemy_config.get_engine().begin() as conn:
        # UUIDAuditBase extends UUIDBase so create_all should build both
        await conn.run_sync(UUIDBase.metadata.create_all)


def plain_text_exception_handler(request: Request, exc: Exception) -> Response:
    """Default handler for exceptions subclassed from HTTPException."""
    tb = traceback.format_exc()
    request.logger.warn(f"exception: {exc}")
    request.logger.warn(f"traceback:\n{tb}")
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
    route_handlers=[FileController, SearchController, RagController, TestController, DaemonController],
)

app = Litestar(
    on_startup=[on_startup],
    plugins=[utils.sqlalchemy_plugin],
    route_handlers=[api_router],
    dependencies={
        "limit_offset": Provide(provide_limit_offset_pagination),
        # do a lance connections for each request
        "lanceconn": Provide(get_lance_connection, sync_to_thread=True),
    },
    cors_config=cors_config,
    logging_config=logging_config,
    exception_handlers={Exception: plain_text_exception_handler},
    listeners=[increment_processed_docs, process_document],
)
