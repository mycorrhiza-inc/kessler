import os
import logging
from time import time

from litestar.types import Receive, Scope, Send, Message
from litestar import Request
from litestar.datastructures import MutableScopeHeaders
from litestar.enums import ScopeType
from litestar.middleware import AbstractMiddleware

logger = logging.getLogger(__name__)


class ClerkRequestValidationMiddleware(AbstractMiddleware):
    scopes = {ScopeType.HTTP}
    exclude = []
    exclude_opt_key = "exclude_from_middleware"

    async def __call__(self, scope: Scope, receive: Receive, send: Send) -> None:

        async def send_wrapper(message: "Message") -> None:
            headers = MutableScopeHeaders.from_message(message=message)
            token = headers["Authorization"]
            logger.log(f"token received:\n{token}")

            await send(message)

        await self.app(scope, receive, send_wrapper)
