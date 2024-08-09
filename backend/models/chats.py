from typing import List

from litestar.contrib.sqlalchemy.base import UUIDAuditBase

from sqlalchemy.orm import Mapped


from utils import PydanticBaseModel

from typing import Optional, List, Union

from pydantic import BaseModel

import hashlib


from enum import Enum
from pathlib import Path


def HumanMessage(content: str):
    return {"content": content, "role": "Human"}


def SystemMessage(content: str):
    return {"content": content, "role": "System"}


def AIMessage(content: str):
    return {"ai": content, "role": "System"}


class RAGChat(BaseModel):
    model: Optional[str] = None
    chat_history: List[Dict[str, str]]


class ChatRole(str, Enum):
    user = "user"
    system = "system"
    assistant = "assistant"


class ChatMessage(BaseModel):
    content: str
    role: ChatRole


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


class MessageSchema(PydanticBaseModel):
    id: any
    body: str


class ChatSchema(PydanticBaseModel):
    id: any
    messages: List[MessageSchema]
