from uuid import UUID
import uuid
from typing import Annotated, assert_type, Any
import logging

from litestar.params import Parameter
from litestar import Controller, Request
# from util.haystack import query_chroma, get_indexed_by_id

from litestar.handlers.http_handlers.decorators import get, post


from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from typing import List


class SearchQuery(BaseModel):
    query: str


class SearchResult(BaseModel):
    ids: List[str]
    result: List[str] | None = None


class SearchResponse(BaseModel):
    status: str
    message: str | None
    results: List[SearchResult] | None


class IndexFileRequest(BaseModel):
    id: uuid


class IndexFileResponse(BaseModel):
    message: str


class SearchController(Controller):
    """Search Controller"""

    @get(path="/search/{fid:uuid}")
    async def search_collection_by_id(
        self,
        fid: UUID = Parameter(
            title="File ID as hex string", description="File to retieve"
        ),
    ) -> Any:
        # FIXME : Get module working with llamaindex and pgvector
        return "Failure"

    @post(path="/search")
    async def get_file(self, data: SearchQuery) -> SearchResponse:
        # FIXME : Get module working with llamaindex and pgvector
        return "Failure"

