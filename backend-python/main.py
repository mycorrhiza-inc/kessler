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
from routing.misc_controller import MiscController
from util.logging import logging_config
from routing.file_controller import FileController
from routing.thaumaturgy_controller import ThaumaturgyController

logger = logging.getLogger(__name__)


async def on_startup() -> None:
    async with utils.sqlalchemy_config.get_engine().begin() as conn:
        # UUIDAuditBase extends UUIDBase so create_all should build both
        await conn.run_sync(UUIDBase.metadata.create_all)


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
    path="/api/v1",
    route_handlers=[
        FileController,
        MiscController,
        ThaumaturgyController,
    ],
)

app = Litestar(
    on_startup=[on_startup],
    plugins=[utils.sqlalchemy_plugin],
    route_handlers=[api_router],
    dependencies={
        "limit_offset": Provide(provide_limit_offset_pagination),
    },
    cors_config=cors_config,
    logging_config=logging_config,
    exception_handlers={Exception: plain_text_exception_handler},
)
