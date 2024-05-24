from hashlib import blake2b
import os
from pathlib import Path
from typing import Any
from uuid import UUID
from typing import Annotated, assert_type
import logging

from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    patch,
    MediaType,
)


from sqlalchemy import select
from sqlalchemy.exc import IntegrityError, NoResultFound
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body
from litestar.logging import LoggingConfig

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from models import FileModel, FileRepository, FileSchema, provide_files_repo


from crawler.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor
from docprocessing.genextras import GenerateExtras
from util.haystack import indexDocByID

from typing import List, Optional, Union, Any, Dict


from util.niclib import get_blake2

from util.haystack import indexDocByID, get_indexed_by_id

import json


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)

class SimpleChatCompletion(BaseModel):
    model: Optional[str]
    chat_history: List[Dict[str, str]]

class RAGChat(BaseModel):
    model: Optional[str]
    chat_history: List[Dict[str, str]]
class RAGQueryResponse(BaseModel):
    model: Optional[str]
    prompt : str
OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")




from llama_index.llms.groq import Groq





GROQ_API_KEY = os.environ["GROQ_API_KEY"]
groq_llm = Groq(
    model="llama3-70b-8192", request_timeout=360.0, api_key=GROQ_API_KEY
)


from rag.llamaindex import create_rag_response_from_query, regenerate_vector_database_from_file_table

def validate_chat(chat_history : List[Dict[str, str]]) -> bool:
    if not isinstance(chat_history, list):
       return False
    found_problem = False
    for chat in chat_history:
        if not isinstance(chat, dict):
            found_problem = True
        if not chat.get("role") in ["user","system","assistant"]:
            found_problem = True 
        if not isinstance(chat.get("message"),str):
            found_problem = True 
    return not found_problem

def force_conform_chat(chat_history : List[Dict[str, str]]) -> List[Dict[str, str]]:
    chat_history = list(chat_history)
    for chat in chat_history:
        if not chat.get("role") in ["user","system","assistant"]:
            chat["role"] = "system" 
        if not isinstance(chat.get("message"),str):
            chat["message"] = str(chat.get("message"))
    return chat_history

class RagController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @post(path="/rag/simple_chat_completion")
    async def simple_chat_completion(
        self,
        files_repo: FileRepository,
        data : SimpleChatCompletion
    ) -> str:
        model_name = data.model
        if model_name is None:
            model_name = "llama3-70b-8192" 
        groq_llm = Groq(
            model=model_name, request_timeout=360.0, api_key=GROQ_API_KEY
        )
        chat_history = data.chat_history
        chat_history = force_conform_chat(chat_history)
        assert validate_chat(chat_history), chat_history
        response = groq_llm.chat(chat_history)
        return response


    @post(path="/rag/rag_chat")
    async def rag_chat(
        self,
        files_repo: FileRepository,
        data : SimpleChatCompletion
    ) -> str:
        model_name = data.model
        if model_name is None:
            model_name = "llama3-70b-8192" 
        groq_llm = Groq(
            model=model_name, request_timeout=360.0, api_key=GROQ_API_KEY
        )
        chat_history = data.chat_history
        chat_history = force_conform_chat(chat_history)
        assert validate_chat(chat_history), chat_history
        response = groq_llm.chat(chat_history)
        return response

    @post(path="/rag/rag_query")
    async def rag_query(
        self,
        files_repo: FileRepository,
        data : RAGQueryResponse
    ) -> str:
        model_name = data.model
        if model_name is None:
            model_name = "llama3-70b-8192" 
        # TODO : Add support for custom model stuff.
        query = data.prompt
        response = create_rag_response_from_query(query)
        return response

    @post(path="/dangerous/regenerate_vector_database")
    async def regen_vecdb(
        self,
        files_repo: FileRepository,
    ) -> str:
        regenerate_vector_database_from_file_table()
        return response
