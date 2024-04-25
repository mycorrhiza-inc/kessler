from uuid import UUID
from typing import Annotated, assert_type
import logging

from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
)

from litestar.params import Body
from litestar.logging import LoggingConfig

from pydantic import BaseModel


from typing import List


class SearchQuery(BaseModel):
    query: str


class SearchResult(BaseModel):
    id: any
    chunks: List[str]


class SearchResponse(BaseModel):
    status: str
    message: str | None
    results: List[SearchResult] | None


class SearchController(Controller):
    """Search Controller"""

    @get(path="/search")
    async def get_file(
        self,
        data: dict[str, str]
    ) -> SearchResponse:
        request = SearchQuery.model_validate(data)
        query = request.query
		