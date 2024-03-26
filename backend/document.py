from sqlalchemy import UUID
from sqlalchemy.orm import Mapped

from litestar import Controller

from litestar.pagination import OffsetPagination

from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository
from litestar.handlers.http_handlers.decorators import get, post, delete, patch
from litestar.params import Parameter
from litestar.repository.filters import LimitOffset
from sqlalchemy.ext.asyncio import AsyncSession, AsyncEngine
from pydantic import TypeAdapter

from db import BaseModel

class DocumentModel(UUIDAuditBase):
    __tablename__ = "document"
    path: Mapped[str]
    name: Mapped[str]