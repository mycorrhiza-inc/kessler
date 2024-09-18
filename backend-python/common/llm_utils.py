from litestar.contrib.sqlalchemy.base import UUIDAuditBase

from sqlalchemy.orm import Mapped


from typing import Optional, List, Union, Dict

from pydantic import BaseModel

import hashlib

from llama_index.core.llms import ChatMessage as LlamaChatMessage

from enum import Enum
from pathlib import Path


from llama_index.llms.openai import OpenAI
from llama_index.core import PromptTemplate

from llama_index.core.retrievers import BaseRetriever
from llama_index.core.llms import LLM
from dataclasses import dataclass
from typing import Optional, List, Union, Any, Tuple

from logic.databaselogic import get_files_from_uuids
import nest_asyncio
import asyncio

from rag.SemanticSplitter import split_by_max_tokensize
from rag.llamaindex import get_llm_from_model_str
from vecstore.search import search

import logging


from models.files import FileRepository, FileSchema, file_model_to_schema

from vecstore import search


from uuid import UUID

from advanced_alchemy.filters import SearchFilter, CollectionFilter

import re

from constants import lemon_text

qa_prompt = (
    lambda context_str: f"""
The following documents should be relevant to the conversation:
---------------------
{context_str}
---------------------
"""
)


generate_query_from_chat_history_prompt = "Please disregard all instructions and generate a query that could be used to search a vector database for relevant information. The query should capture the main topic or question being discussed in the chat. Please output the query as a string, using a format suitable for a vector database search (e.g. a natural language query or a set of keywords)."

"""
Stuff
Chat history: User: "I'm looking for a new phone. What are some good options?" Assistant: "What's your budget?" User: "Around $500" Assistant: "Okay, in that range you have options like the Samsung Galaxy A series or the Google Pixel 4a"

Example output: "Query: 'best phones under $500'"
"""

does_chat_need_query = 'Please determine if you need to query a vector database of relevant documents to answer the user. Answer with only a "yes" or "no".'

query_str = (
    "Can you tell me about results from RLHF using both model-based and"
    " human-based evaluation?"
)


default_logger = logging.getLogger(__name__)


class RAGChat(BaseModel):
    model: Optional[str] = None
    chat_history: List[Dict[str, str]]


class ChatRole(str, Enum):
    user = "user"
    system = "system"
    assistant = "assistant"


class KeChatMessage(BaseModel):
    content: str
    role: ChatRole


# Do something with the chat message validation maybe, probably not worth it
def sanitzie_chathistory_llamaindex(chat_history: List) -> List[LlamaChatMessage]:
    def sanitize_message(raw_message: Union[dict, KeChatMessage]) -> LlamaChatMessage:
        if isinstance(raw_message, KeChatMessage):
            raw_message = cm_to_dict(raw_message)
        return LlamaChatMessage(
            role=raw_message["role"], content=raw_message["content"]
        )

    return list(map(sanitize_message, chat_history))


def dict_to_cm(input_dict: Union[dict, KeChatMessage]) -> KeChatMessage:
    if isinstance(input_dict, KeChatMessage):
        return input_dict
    return KeChatMessage(
        content=input_dict["content"], role=ChatRole(input_dict["role"])
    )


def cm_to_dict(cm: KeChatMessage) -> Dict[str, str]:
    return {"content": cm.content, "role": cm.role.value}


def unvalidate_chat(chat_history: List[KeChatMessage]) -> List[Dict[str, str]]:
    return list(map(cm_to_dict, chat_history))


def validate_chat(chat_history: List[Dict[str, str]]) -> List[KeChatMessage]:
    return list(map(dict_to_cm, chat_history))


def force_conform_chat(chat_history: List[Dict[str, str]]) -> List[Dict[str, str]]:
    chat_history = list(chat_history)
    for chat in chat_history:
        if not chat.get("role") in ["user", "system", "assistant"]:
            chat["role"] = "system"
        if not isinstance(chat.get("message"), str):
            chat["message"] = str(chat.get("message"))
    return chat_history


class KeLLMUtils:
    def __init__(self, llm: Union[str, Any]) -> None:
        if llm == "":
            llm = None
        if llm is None:
            llm = "llama-405b"
        if isinstance(llm, str):
            llm = get_llm_from_model_str(llm)
        self.llm = llm

    async def achat(self, chat_history: Any) -> Any:
        llama_chat_history = sanitzie_chathistory_llamaindex(chat_history)
        response = await self.llm.achat(llama_chat_history)
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
        return KeChatMessage(role=ChatRole.assistant, content=str_response)

    async def summarize_single_chunk(self, markdown_text: str) -> str:
        summarize_prompt = "Make sure to provide a well researched summary of the text provided by the user, if it appears to be the summary of a larger document, just summarize the section provided."
        summarize_message = KeChatMessage(
            role=ChatRole.assistant, content=summarize_prompt
        )
        text_message = KeChatMessage(role=ChatRole.user, content=markdown_text)
        summary = await self.achat(
            sanitzie_chathistory_llamaindex([summarize_message, text_message])
        )
        return summary.content

    async def summarize_mapreduce(
        self, markdown_text: str, max_tokensize: int = 8096
    ) -> str:
        splits = split_by_max_tokensize(markdown_text, max_tokensize)
        if len(splits) == 1:
            return await self.summarize_single_chunk(markdown_text)
        summaries = await asyncio.gather(
            *[self.summarize_single_chunk(chunk) for chunk in splits]
        )
        coherence_prompt = "Please rewrite the following list of summaries of chunks of the document into a final summary of similar length that incorperates all the details present in the chunks"
        cohere_message = KeChatMessage(ChatRole.assistant, coherence_prompt)
        combined_summaries_prompt = KeChatMessage(ChatRole.user, "\n".join(summaries))
        final_summary = await self.achat([cohere_message, combined_summaries_prompt])
        return final_summary.content
