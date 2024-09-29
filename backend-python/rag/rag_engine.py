from llama_index.llms.openai import OpenAI
from llama_index.core import PromptTemplate

from llama_index.core.retrievers import BaseRetriever
from llama_index.core.llms import LLM
from dataclasses import dataclass
from typing import Optional, List, Union, Any, Tuple

from logic.databaselogic import get_files_from_uuids
import nest_asyncio
import asyncio

from common.llm_utils import (
    ChatRole,
    KeChatMessage,
    sanitzie_chathistory_llamaindex,
    KeLLMUtils,
)
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
from pydantic import BaseModel

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


def strip_links_and_tables(markdown_text):
    # Remove markdown links
    no_links = re.sub(r"\[.*?\]\(.*?\)", "", markdown_text)
    # Remove markdown tables
    no_tables = re.sub(r"\|.*?\|", "", no_links)
    return no_tables


async def convert_search_results_to_frontend_table_old(
    search_results: List[Any],
    files_repo: FileRepository,
    max_results: int = 10,
    include_text: bool = True,
):
    logger = default_logger
    res = search_results[0]
    res = res[:max_results]
    uuid_list = []
    text_list = []
    # TODO: Refactor for less checks and ugliness
    for result in res:
        logger.info(result)
        logger.info(result["entity"])
        uuid = UUID(result["entity"]["source_id"])
        uuid_list.append(uuid)
        if include_text:
            text_list.append(result["entity"]["text"])
            # text_list.append(lemon_text)
    file_models = await get_files_from_uuids(files_repo, uuid_list)
    file_results = list(map(file_model_to_schema, file_models))
    if include_text:
        for index in range(len(file_results)):
            file_results[index].display_text = text_list[index]

    return file_results


async def convert_search_results_to_frontend_table(
    search_results: List[Any],
    files_repo: FileRepository,
    max_results: int = 10,
    include_text: bool = True,
):
    logger = default_logger
    res = search_results[0]
    res = res[:max_results]
    uuid_list = []
    text_list = []
    # TODO: Refactor for less checks and ugliness
    for result in res:
        logger.info(result)
        logger.info(result["entity"])
        uuid = UUID(result["entity"]["source_id"])
        uuid_list.append(uuid)
        if include_text:
            text_list.append(result["entity"]["text"])
            # text_list.append(lemon_text)
    file_models = await get_files_from_uuids(files_repo, uuid_list)
    search_data_results = []
    for index, file_model in enumerate(file_models):
        search_data = SearchData(
            name=file_model.name if file_model.name else "",
            text=text_list[index] if include_text else "",
            docID="",  # Not a validated field in the mdata so far, discuss at some point.
            sourceID=str(file_model.id),
        )
        search_data_results.append(search_data)

    return search_data_results


class KeRagEngine(KeLLMUtils):
    def __init__(self, llm: Union[str, Any]) -> None:
        super().__init__(llm)

    async def does_chat_need_query(self, chat_history: List[KeChatMessage]) -> bool:
        does_chat_need_query = 'Please determine if you need to query a vector database of relevant documents to answer the user. Answer with only a "yes" or "no".'
        check_message = KeChatMessage(
            role=ChatRole.assistant, content=does_chat_need_query
        )

        def check_yes_no(test_str: str) -> bool:
            test_str = test_str.lower()
            if test_str.startswith("yes"):
                return True
            if test_str.startswith("no"):
                return False
            raise ValueError("Expected yes or no got: " + test_str)

        return check_yes_no((await self.achat(chat_history + [check_message])).content)

    async def rag_achat(
        self,
        chat_history: List[KeChatMessage],
        files_repo: FileRepository,
        logger: Optional[logging.Logger] = None,
    ) -> Tuple[KeChatMessage, List[FileSchema]]:
        if logger is None:
            logger = default_logger
        if not await self.does_chat_need_query(chat_history):
            return await self.achat(chat_history)

        async def generate_query_from_chat_history(
            chat_history: List[KeChatMessage],
        ) -> str:
            querygen_addendum = KeChatMessage(
                role=ChatRole.system, content=generate_query_from_chat_history_prompt
            )
            completion = await self.achat(chat_history + [querygen_addendum])
            return completion

        def generate_context_msg_from_search_results(
            search_results: List[Any], max_results: int = 3
        ) -> KeChatMessage:
            logger = default_logger
            res = search_results[0]
            res = res[:max_results]
            return_prompt = "Here is a list of documents that might be relevant to the following chat:"
            # TODO: Refactor for less checks and ugliness
            for result in res:
                uuid_str = result["entity"]["source_id"]
                text = result["entity"]["text"]
                return_prompt += f"\n\n{uuid_str}:\n{text}"
            return KeChatMessage(role=ChatRole.assistant, content=return_prompt)

        query = await generate_query_from_chat_history(chat_history)
        res = search(query=query, output_fields=["source_id", "text"])
        logger.info(res)
        context_msg = generate_context_msg_from_search_results(res)
        # TODO: Get these 2 async func calls to happen simultaneously
        final_message = await self.achat([context_msg] + chat_history)
        return_schemas = await convert_search_results_to_frontend_table(res, files_repo)

        return (final_message, return_schemas)
