from rag.llamaindex import (
    get_llm_from_model_str,
    create_rag_response_from_query,
    regenerate_vector_database_from_file_table,
    add_document_to_db_from_text,
    generate_chat_completion,
    sanitzie_chathistory_llamaindex,
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
    FileModel,
    FileRepository,
    FileSchema,
    model_to_schema,
    provide_files_repo,
)


from typing import List, Optional, Union, Any, Dict


from vecstore import search

import json


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


from constants import (
    OS_TMPDIR,
    OS_GPU_COMPUTE_URL,
    OS_FILEDIR,
    OS_HASH_FILEDIR,
    OS_OVERRIDE_FILEDIR,
    OS_BACKUP_FILEDIR,
)


class RagController(Controller):
    """Rag Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @post(path="/rag/basic_chat")
    async def basic_chat_no_rag(self, data: SimpleChatCompletion) -> dict:
        model_name = data.model
        if model_name == "":
            model_name = None
        if model_name is None:
            model_name = "llama-405b"
        chat_history = data.chat_history
        chat_history = force_conform_chat(chat_history)
        assert validate_chat(chat_history), chat_history
        llama_chat_history = sanitzie_chathistory_llamaindex(chat_history)
        chosen_llm = get_llm_from_model_str(model_name)
        response = await chosen_llm.achat(llama_chat_history)
        str_response = str(response)

        def remove_prefixes(input_string: str) -> str:
            prefixes = ["assistant: "]
            for prefix in prefixes:
                if input_string.startswith(prefix):
                    input_string = input_string[
                        len(prefix) :
                    ]  # 10 is the length of "assistant: "
            return input_string

        str_response = remove_prefixes(str_response)
        return {"role": "assistant", "content": str_response}

    @post(path="/rag/rag_chat")
    async def rag_chat(self, data: SimpleChatCompletion) -> dict:
        chat_history = data.chat_history
        ai_message_response = generate_chat_completion(chat_history)
        return ai_message_response

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
            # Congrats for finding the portal easter egg!
            doc_text = "This is the cannonical example document with the following advice for dealing with adversity in life: All right, I've been thinking, when life gives you lemons, don't make lemonade! Make life take the lemons back! Get mad! I don't want your damn lemons! What am I supposed to do with these? Demand to see life's manager! Make life rue the day it thought it could give Cave Johnson lemons! Do you know who I am? I'm the man whose gonna burn your house down - with the lemons!"
        add_document_to_db_from_text(doc_text, doc_metadata)

    @post(path="/dangerous/rag/regenerate_vector_database")
    async def regen_vecdb(
        self,
        files_repo: FileRepository,
    ) -> str:
        await regenerate_vector_database_from_file_table()
        return ""

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
        only_fileobj: bool = False,
    ) -> Any:
        logger = request.logger
        query = data.query
        res = search(query=query)
        res = res[0]
        for result in res:
            logger.info(result["entity"])
            uuid = UUID((result["entity"]["source_id"]))
            logger.info(f"Asking PG for data on file: {uuid}")
            schema = model_to_schema(await files_repo.get(uuid))
            result["file"] = schema
        if only_fileobj:
            return list(map(lambda r: r["file"], res))
        return res

        # return list(map(create_rag_response_from_query, res))

    # @post(path="/search/{fid:uuid}")
    # async def search_collection_by_id(
    #     self,
    #     request: Request,
    #     data: SearchQuery,
    #     fid: UUID = Parameter(
    #         title="File ID as hex string", description="File to retieve"
    #     ),
    # ) -> Any:
    #     return "failure"
