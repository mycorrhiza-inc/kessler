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
        # TODO: Become more functional
        # ids = map(lambda r: {"uuid": r[0]["entity"]["id"]}, res)

        ids = []
        for r in res:
            # logger.info(r[0]["entity"])
            # logger.info(type(r[0]["entity"]))
            # logger.info(r[0]["entity"].keys())
            # test_dict = json.loads(r[0]["entity"]["_node_content"])
            ids.append(r[0]["entity"])
        if only_uuid:
            # Fix later
            return []

        # TODO: Use an async map for this as soon as python gets that functionality or use an import
        async def get_file(uuid_str: str):
            uuid = UUID(uuid_str)
            logger.info(uuid)
            obj = await files_repo.get(uuid)

            type_adapter = TypeAdapter(FileSchema)

            return type_adapter.validate_python(obj)

        files = []
        for id in ids:
            for field in ["document_id", "doc_id", "ref_doc_id", "id"]:
                try:
                    uuid_str = id[field]
                    file_result = await get_file(uuid_str)
                    files.append(file_result)
                    logger.info("Success on" + field)
                except Exception as e:
                    pass
                    # logger.error(f"Encountered an error while attempting to get file {uuid_str} : {e}")

        return files
