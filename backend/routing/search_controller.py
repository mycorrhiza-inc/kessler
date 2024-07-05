from uuid import UUID
import uuid
from typing import Annotated, assert_type, Any
import logging
import copy

from litestar.params import Parameter
from litestar import Controller, Request

# from util.haystack import query_chroma, get_indexed_by_id

from litestar.handlers.http_handlers.decorators import get, post


from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel

import lancedb
from lancedb import DBConnection

from lance_store.connection import ensure_fts_index

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

    @post(path="/search/{fid:uuid}")
    async def search_collection_by_id(
        self,
        request: Request,
        data: SearchQuery,
        fid: UUID = Parameter(
            title="File ID as hex string", description="File to retieve"
        ),
    ) -> Any:
        return "failure"

    @post(path="/search")
    async def search(
        self, data: SearchQuery, lanceconn: DBConnection, request: Request
    ) -> SearchResponse:

        v = lanceconn.open_table("vectors")

        def get_text(x):
            return x["text"]

        res = v.search(data.query).to_list()

        # search all dockets for a given item
        def clean_data(data):
            # Make a deep copy of the data to avoid modifying the original data
            cleaned_data = copy.deepcopy(data)

            # Remove the specified fields
            if "vector" in cleaned_data:
                del cleaned_data["vector"]
            if "text" in cleaned_data:
                del cleaned_data["text"]
            if (
                "metadata" in cleaned_data
                and "_node_content" in cleaned_data["metadata"]
            ):
                del cleaned_data["metadata"]["_node_content"]

            return cleaned_data

        # Apply the clean_data function to each item in search results
        filtered_data = [clean_data(item) for item in res]
        return filtered_data
