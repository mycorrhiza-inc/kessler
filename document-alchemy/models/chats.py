from litestar.contrib.sqlalchemy.base import UUIDAuditBase

from sqlalchemy.orm import Mapped


from typing import Optional, List, Union, Dict

from pydantic import BaseModel

import hashlib

from llama_index.core.llms import ChatMessage as LlamaChatMessage

from enum import Enum
from pathlib import Path


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
