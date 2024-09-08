import logging
import traceback

from litestar import Litestar, Router
from litestar.config.cors import CORSConfig
from litestar import MediaType, Request, Response
from litestar.status_codes import HTTP_500_INTERNAL_SERVER_ERROR

from litestar.di import Provide

from models import utils
from util.logging import logging_config

from routiing.daemon_controller import DaemonController

logger = logging.getLogger(__name__)


async def on_startup() -> None:
    pass


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


cors_config = CORSConfig(allow_origins=["*"])

api_router = Router(
    path="/thaumaturgy/api/v1",
    route_handlers=[],
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
