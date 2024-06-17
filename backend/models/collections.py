from contextlib import asynccontextmanager
from typing import AsyncIterator, Annotated, List
import traceback
from uuid import UUID

from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column


from .utils import sqlalchemy_config, PydanticBaseModel
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession


class CollectionModel(UUIDAuditBase):
    __tablename__ = "document_collection"
    # used to check chroma named collections
    name: Mapped[str]


class CollectionSchema(PydanticBaseModel):
    id: any
    name: str
