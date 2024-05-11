from haystack_integrations.document_stores.chroma import ChromaDocumentStore
from haystack import Document

from haystack_integrations.document_stores.chroma import ChromaDocumentStore
from haystack_integrations.components.retrievers.chroma import ChromaEmbeddingRetriever

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
        res = await get_indexed_by_id(fid)
        return res

    @post(path="/search")
    async def get_file(self, data: SearchQuery) -> SearchResponse:
        request = SearchQuery.model_validate(data)
        query = request.query
        results = query_chroma(query=query)
        documents = results["retriever"]["documents"]
        ids = []
        for i in documents:
            ids.append(str(i.id))
        type_adapter = TypeAdapter(SearchResult)
        return type_adapter.validate_python({"ids": ids})

    # @post(path="/search/index/")
    # async def index_file(self, data: IndexFileRequest) -> IndexFileResponse:
    #     query = data.query
