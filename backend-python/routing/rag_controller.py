from common.chat_schemas import (
    KeChatMessage,
    cm_to_dict,
    sanitzie_chathistory_llamaindex,
    unvalidate_chat,
    validate_chat,
)
from rag.llamaindex import (
    get_llm_from_model_str,
    create_rag_response_from_query,
    regenerate_vector_database_from_file_table,
    add_document_to_db_from_text,
    generate_chat_completion,
)

import os
from pathlib import Path
from uuid import UUID
import logging

from litestar import Controller, Request


from litestar.params import Parameter

from litestar.handlers.http_handlers.decorators import (
    post,
)


from litestar.di import Provide

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from models.files import (
    FileRepository,
    provide_files_repo,
)


from typing import List, Optional, Union, Any, Dict


from rag.rag_engine import KeRagEngine, convert_search_results_to_frontend_table
from vecstore import search

import json
import asyncio

from constants import lemon_text

from advanced_alchemy.filters import SearchFilter, CollectionFilter


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


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
    id: UUID


class IndexFileResponse(BaseModel):
    message: str


class SimpleChatCompletion(BaseModel):
    model: Optional[str] = None
    chat_history: List[Dict[str, str]]


class RAGQueryResponse(BaseModel):
    model: Optional[str] = None
    prompt: str


class ManualDocument(BaseModel):
    text: str
    metadata: Optional[dict]


class RagController(Controller):
    """Rag Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @post(path="/rag/basic_chat")
    async def basic_chat_no_rag(self, data: SimpleChatCompletion) -> dict:
        model_name = data.model
        validated_chat_history = validate_chat(data.chat_history)
        rag_engine = KeRagEngine(model_name)
        result = await rag_engine.achat(validated_chat_history)
        return {"message": cm_to_dict(result)}

    @post(path="/rag/rag_chat")
    async def rag_chat(
        self, data: SimpleChatCompletion, files_repo: FileRepository
    ) -> dict:
        model_name = data.model
        validated_chat_history = validate_chat(data.chat_history)
        rag_engine = KeRagEngine(model_name)
        result_message, file_schema_citations = await rag_engine.rag_achat(
            validated_chat_history, files_repo
        )
        assert isinstance(result_message, KeChatMessage)
        return {
            "message": cm_to_dict(result_message),
            "citations": file_schema_citations,
        }

    @post(path="/rag/rag_query")
    async def rag_query(
        self, files_repo: FileRepository, data: RAGQueryResponse
    ) -> str:
        model_name = data.model
        # Doesnt Do anything atm
        if model_name is None:
            model_name = "llama-70b"
        # TODO : Add support for custom model stuff.
        query = data.prompt
        response = create_rag_response_from_query(query)
        return response

    @post(path="/dangerous/rag/manaul_add_doc_to_vecdb")
    async def manual_add_doc_vecdb(
        self, files_repo: FileRepository, data: ManualDocument
    ) -> None:
        doc_metadata = data.metadata
        doc_text = data.text
        if doc_text == "":
            doc_text = lemon_text
        add_document_to_db_from_text(doc_text, doc_metadata)

    @post(path="/search")
    async def search(
        self,
        files_repo: FileRepository,
        data: SearchQuery,
        request: Request,
        only_fileobj: bool = False,
        max_results: int = 10,
    ) -> list:
        logger = request.logger
        query = data.query
        # FIXME: Speed up search so its less slow
        if len(query) <= 3:
            return []
        res = search(query=query, output_fields=["source_id", "text"])
        logger.info(res)
        return await convert_search_results_to_frontend_table(
            search_results=res,
            files_repo=files_repo,
            max_results=max_results,
            include_text=True,
        )
