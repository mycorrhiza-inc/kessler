from typing import AsyncIterator

from sqlalchemy import ForeignKey
from sqlalchemy import inspect
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from contextlib import asynccontextmanager

from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository
from litestar.contrib.sqlalchemy.base import UUIDBase


from pydantic import validator
import traceback

from sqlalchemy.ext.asyncio import async_sessionmaker, create_async_engine

from .utils import sqlalchemy_config


class FileResourceModel(AuditColumns):
    """
    Used to access and maniuplate files without regard to their resourceid
    """
    __tablename__ = "LinkResource"
    resource_id = mapped_column(ForeignKey("resource.id"))
    # used to get Files from resource IDs
    file_id = mapped_column(ForeignKey("file.id"), primary_key=True)

    # let's make a simple context manager as an example here.


class File(UUIDAuditBase):
    __tablename__ = "file"
    path: Mapped[str]
    doctype: Mapped[str | None]
    lang: Mapped[str | None]
    name: Mapped[
        str | None
    ]  # I dont know if this should be included either in here or as a entry in doc_metadata, expecially since its only ever going to be used by the frontend. However, it might be an important query paramater and it seems somewhat irresponsible to not include.
    stage: Mapped[str | None]  # Either "stage0" "stage1" "stage2" or "stage3"
    summary: Mapped[str | None]
    short_summary: Mapped[str | None]

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value
# TODO: i need to define the repo factory as a mixin for any future work we do
# https://stackoverflow.com/questions/36690588/should-mixins-use-parent-attributes
# this should be extremely easy to

    @classmethod
    @asynccontextmanager
    async def repo(self) -> AsyncIterator['FileRepository']:
        session_factory = sqlalchemy_config.create_session_maker()
        async with session_factory() as db_session:
            try:
                yield FileRepository(session=db_session)
            except Exception as e:
                print(traceback.format_exc())
                print("rolling back")
                await db_session.rollback()
            else:
                print("committhing change")
                await db_session.commit()

    @classmethod
    def printself(cls):
        print(cls.__dict__)
        
    @classmethod
    async def create_self(cls):
        async with cls.repo() as repo:
            obj = await repo.add(cls)
            print(obj.name)
            return obj

    @classmethod
    async def create(cls, **kw):
        async with cls.repo() as repo:
            f = cls(**kw)
            obj = await repo.add(f)
            print(obj.__dict__)
            return obj

    @classmethod
    async def new(cls, f: 'File'):
        async with cls.repo() as repo:
            obj = await repo.add(f)
            print(obj.__dict__)
            return obj



async def newfi():
    """
    to run this test in a python repl
    make the db 
    > import asynio
    > import files
    > asynio.run(files.newfi())
    
    """
    async with sqlalchemy_config.get_engine().begin() as conn:
        # UUIDAuditBase extends UUIDBase so create_all should build both
        print("making sure the db exists")
        await conn.run_sync(UUIDBase.metadata.create_all)

    print("calling create")
    f = File(path="./somewhere")
    f = await File.new(f)
    print(f.__dict__)
    print("calling create")
    await File.create(path="./somewhere")


class FileRepository(SQLAlchemyAsyncRepository[File]):
    """File repository."""

    model_type = File
