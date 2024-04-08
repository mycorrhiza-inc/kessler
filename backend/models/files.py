
from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository
from litestar.contrib.sqlalchemy.base import UUIDBase

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from pydantic import BaseModel, ConfigDict, StringConstraints, validator
from contextlib import asynccontextmanager

from utils import RepoMixin

from typing import AsyncIterator, Annotated

import traceback
import uuid
from uuid import UUID


from utils import sqlalchemy_config


class File(BaseModel):
    id: any # TODO: figure out a better type for this UUID :/
    url: str
    path: str
    doctype: str
    lang: str
    name: str
    stage: str   # Either "stage0" "stage1" "stage2" or "stage3"
    summary: str | None
    short_summary: str | None

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class FileModel(UUIDAuditBase, RepoMixin):
    __tablename__ = "file"
    path: Mapped[str]
    doctype: Mapped[str]
    lang: Mapped[str]
    name: Mapped[str]
    stage: Mapped[str]  # Either "stage0" "stage1" "stage2" or "stage3"
    summary: Mapped[str | None]
    short_summary: Mapped[str | None]

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value

    @classmethod
    async def provide_repo(cls, session) -> 'FileRepository':
        return FileRepository(session=session)

    # # define the context manager for each file repo
    @classmethod
    @asynccontextmanager
    async def repo(cls) -> AsyncIterator['FileRepository']:
        session_factory = sqlalchemy_config.create_session_maker()
        async with session_factory() as db_session:
            try:
                yield cls.provide_repo(session=db_session)
            except Exception as e:
                print(traceback.format_exc())
                print("rolling back")
                await db_session.rollback()
            else:
                print("committhing change")
                await db_session.commit()


    @classmethod
    async def updateStage(cls, id, stage):
        async with cls.repo() as repo:
            obj = await cls.find(id)
            obj.stage = stage
            obj = await repo.update(obj)

            return obj


class FileRepository(SQLAlchemyAsyncRepository[FileModel]):
    """File repository."""

    model_type = FileModel
