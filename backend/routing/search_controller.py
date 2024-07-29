from uuid import UUID
import uuid
from typing import Annotated, assert_type, Any
import logging
import copy

from litestar.params import Parameter
from litestar import Controller, Request
from litestar.di import Provide

# from util.haystack import query_chroma, get_indexed_by_id

from litestar.handlers.http_handlers.decorators import get, post


from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from vecstore import search

from typing import List

# from asyncstdlib import amap

from models.files import (
    FileModel,
    FileRepository,
    FileSchema,
    FileSchemaWithText,
    provide_files_repo,
    DocumentStatus,
    docstatus_index,
    model_to_schema,
)

import json


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

    dependencies = {"files_repo": Provide(provide_files_repo)}

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
        self,
        files_repo: FileRepository,
        data: SearchQuery,
        request: Request,
        only_uuid: bool = False,
    ) -> Any:
        logger = request.logger
        query = data.query
        res = search(query=query)
        return res

