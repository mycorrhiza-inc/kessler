from llama_index.llms.openai import OpenAI
from llama_index.core import PromptTemplate

from llama_index.core.retrievers import BaseRetriever
from llama_index.core.llms import LLM
from dataclasses import dataclass
from typing import Optional, List, Union, Any

import nest_asyncio
import asyncio

from models.chats import ChatRole, KeChatMessage, sanitzie_chathistory_llamaindex
from rag.llamaindex import get_llm_from_model_str
from vecstore.search import search

import logging


from models.files import FileRepository, model_to_schema

from vecstore import search


from uuid import UUID

from advanced_alchemy.filters import SearchFilter, CollectionFilter


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
    uuid_filter = CollectionFilter(field_name="id", values=uuid_list)
    file_models = await files_repo.list(uuid_filter)
    file_results = list(map(model_to_schema, file_models))
    if include_text:
        for index in range(len(file_results)):
            file_results[index].display_text = text_list[index]

    return file_results


class KeRagEngine:
    def __init__(self, llm: Union[str, Any]) -> None:
        if llm == "":
            llm = None
        if llm is None:
            llm = "llama-405b"
        if isinstance(llm, str):
            llm = get_llm_from_model_str(llm)
        self.llm = llm

    async def achat_basic(self, chat_history: List[KeChatMessage]) -> KeChatMessage:
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

        return check_yes_no(
            (await self.achat_basic(chat_history + [check_message])).content
        )

    async def rag_achat(
        self, chat_history: List[KeChatMessage], logger: Optional[logging.Logger] = None
    ) -> KeChatMessage:
        if logger is None:
            logger = default_logger
        if not await self.does_chat_need_query(chat_history):
            return await self.achat_basic(chat_history)

        async def generate_query_from_chat_history(
            chat_history: List[KeChatMessage],
        ) -> str:
            querygen_addendum = KeChatMessage(
                role=ChatRole.system, content=generate_query_from_chat_history_prompt
            )
            completion = await self.achat_basic(chat_history + [querygen_addendum])
            return completion

        query = await generate_query_from_chat_history(chat_history)
        res = search(query=query)
        logger.info(res)
        final_result = "rag functionality not implemented yet"
        chat = KeChatMessage(role=ChatRole.assistant, content=final_result)
        return chat
