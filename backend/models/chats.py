from contextlib import asynccontextmanager
from typing import AsyncIterator, Annotated, List
import traceback
from uuid import UUID

from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column


from .utils import RepoMixin, sqlalchemy_config, PydanticBaseModel
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession


class MessageModel(UUIDAuditBase):
    __tablename__ = "message"
    sourceModel: Mapped[str]  # Human or ModelName
    Text: Mapped[str]


class ChatModel(UUIDAuditBase):
    pass


class MessageSchema(PydanticBaseModel):
    id: any
    body: str


class ChatSchema(PydanticBaseModel):
    id: any
    messages: List[MessageSchema]
