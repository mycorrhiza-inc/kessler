from pydantic import BaseModel


from common.niclib import token_split

from llama_index.core.llms import ChatMessage as LlamaChatMessage

from enum import Enum
from pathlib import Path


from typing import Optional, List, Union, Any, Tuple, Dict

import asyncio

from rag.SemanticSplitter import split_by_max_tokensize

import logging

from llama_index.llms.groq import Groq
from llama_index.llms.openai import OpenAI
from llama_index.llms.fireworks import Fireworks


from constants import OPENAI_API_KEY, OCTOAI_API_KEY, GROQ_API_KEY, FIREWORKS_API_KEY


def get_llm_from_model_str(model_name: Optional[str]):
    if model_name is None:
        model_name = "llama-405b"
    if model_name in ["llama-8b", "llama-3.1-8b-instant"]:
        actual_name = "llama-3.1-8b-instant"
        return Groq(model=actual_name, request_timeout=60.0, api_key=GROQ_API_KEY)
    if model_name in [
        "llama-70b",
        "llama3-70b-8192",
        "llama-3.1-70b-versatile",
    ]:
        actual_name = "llama-3.1-70b-versatile"
        return Groq(model=actual_name, request_timeout=60.0, api_key=GROQ_API_KEY)
    if model_name in ["llama-405b", "llama-3.1-405b-reasoning"]:
        actual_name = "accounts/fireworks/models/llama-v3p1-405b-instruct"
        return Fireworks(model=actual_name, api_key=FIREWORKS_API_KEY)
    if model_name in ["gpt-4o"]:
        return OpenAI(model=model_name, request_timeout=60.0, api_key=OPENAI_API_KEY)
    else:
        raise Exception("Model String Invalid or Not Supported")


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
            role=ChatRole.system, content=summarize_prompt
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
        cohere_message = KeChatMessage(role=ChatRole.system, content=coherence_prompt)
        combined_summaries_prompt = KeChatMessage(
            role=ChatRole.user, content="\n".join(summaries)
        )
        final_summary = await self.achat([cohere_message, combined_summaries_prompt])
        return final_summary.content

    async def simple_instruct(self, content: str, instruct: str) -> str:
        history = [
            KeChatMessage(content=instruct, role=ChatRole.system),
            KeChatMessage(content=content, role=ChatRole.user),
        ]
        completion = await self.achat(history)
        return completion.content

    async def mapreduce_llm_instruction_across_string(
        self, content: str, chunk_size: int, instruction: str, join_str: str
    ) -> str:
        # Replace with semantic splitter
        split = token_split(content, chunk_size)

        async def clean_chunk(chunk: str) -> str:

            history = [
                KeChatMessage(content=instruction, role=ChatRole.system),
                KeChatMessage(content=chunk, role=ChatRole.user),
            ]
            completion = await self.llm.achat(history)
            return completion.content

        tasks = [clean_chunk(chunk) for chunk in split]
        results = await asyncio.gather(*tasks)
        return join_str.join(results)
