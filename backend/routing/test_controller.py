import os
from pathlib import Path

from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    MediaType,
)


import json


class TestController(Controller):
    """File Controller"""

    @get(path="/test/run_suite")
    async def run_test_suite(self) -> str:
        return "Test complete"
