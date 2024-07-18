from rag.llamaindex import (
    create_rag_response_from_query,
    regenerate_vector_database_from_file_table,
    add_document_to_db_from_text,
    generate_chat_completion,
    sanitzie_chathistory_llamaindex,
)
from llama_index.llms.openai import OpenAI
from llama_index.llms.groq import Groq
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


from models.files import FileModel, FileRepository, FileSchema, provide_files_repo


from crawler.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor

from typing import List, Optional, Union, Any, Dict

from util.niclib import get_blake2

import json


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


class SimpleChatCompletion(BaseModel):
    model: Optional[str] = None
    chat_history: List[Dict[str, str]]


class RAGChat(BaseModel):
    model: Optional[str] = None
    chat_history: List[Dict[str, str]]


class RAGQueryResponse(BaseModel):
    model: Optional[str] = None
    prompt: str


class ManualDocument(BaseModel):
    text: str
    metadata: Optional[dict]


OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")


GROQ_API_KEY = os.environ["GROQ_API_KEY"]
OPENAI_API_KEY = os.environ["OPENAI_API_KEY"]


def validate_chat(chat_history: List[Dict[str, str]]) -> bool:
    if not isinstance(chat_history, list):
        return False
    found_problem = False
    for chat in chat_history:
        if not isinstance(chat, dict):
            found_problem = True
        if not chat.get("role") in ["user", "system", "assistant"]:
            found_problem = True
        if not isinstance(chat.get("message"), str):
            found_problem = True
    return not found_problem


def force_conform_chat(chat_history: List[Dict[str, str]]) -> List[Dict[str, str]]:
    chat_history = list(chat_history)
    for chat in chat_history:
        if not chat.get("role") in ["user", "system", "assistant"]:
            chat["role"] = "system"
        if not isinstance(chat.get("message"), str):
            chat["message"] = str(chat.get("message"))
    return chat_history


class RagController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @post(path="/rag/basic_chat")
    async def basic_chat_no_rag(
        self, files_repo: FileRepository, data: SimpleChatCompletion
    ) -> dict:
        model_name = data.model
        if model_name == "":
            model_name = None
        if model_name is None:
            model_name = "llama3-70b-8192"
        chat_history = data.chat_history
        chat_history = force_conform_chat(chat_history)
        assert validate_chat(chat_history), chat_history
        llama_chat_history = sanitzie_chathistory_llamaindex(chat_history)
        if model_name in ["llama3-70b-8192"]:
            groq_llm = Groq(
                model=model_name, request_timeout=60.0, api_key=GROQ_API_KEY
            )
            response = groq_llm.chat(llama_chat_history)
        if model_name in ["gpt-4o"]:
            openai_llm = OpenAI(
                model=model_name, request_timeout=60.0, api_key=OPENAI_API_KEY
            )
            response = openai_llm.chat(llama_chat_history)
        str_response = str(response)

        def remove_prefixes(input_string: str) -> str:
            prefixes = ["assistant: "]
            for prefix in prefixes:
                if input_string.startswith(prefix):
                    input_string = input_string[
                        len(prefix):
                    ]  # 10 is the length of "assistant: "
            return input_string

        str_response = remove_prefixes(str_response)
        return {"role": "assistant", "content": str_response}

    @post(path="/rag/rag_chat")
    async def rag_chat(
        self, files_repo: FileRepository, data: SimpleChatCompletion
    ) -> dict:
        chat_history = data.chat_history
        ai_message_response = generate_chat_completion(chat_history)
        return ai_message_response

    @post(path="/rag/rag_query")
    async def rag_query(
        self, files_repo: FileRepository, data: RAGQueryResponse
    ) -> str:
        model_name = data.model
        if model_name is None:
            model_name = "llama3-70b-8192"
        # TODO : Add support for custom model stuff.
        query = data.prompt
        response = create_rag_response_from_query(query)
        return response

    @post(path="/rag/manaul_add_doc_to_vecdb")
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
        regenerate_vector_database_from_file_table()
        return ""
