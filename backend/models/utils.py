from typing import AsyncIterator, Generic, TypeVar
from contextlib import asynccontextmanager
import traceback

from sqlalchemy.ext.declarative import declared_attr

from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository
from litestar.contrib.sqlalchemy.base import UUIDBase

from litestar.contrib.sqlalchemy.plugins import (
    AsyncSessionConfig,
    SQLAlchemyAsyncConfig,
    SQLAlchemyInitPlugin,
)

session_config = AsyncSessionConfig(expire_on_commit=False)

sqlalchemy_config = SQLAlchemyAsyncConfig(
    connection_string="sqlite+aiosqlite:///instance/kessler.sqlite",
    # connection_string="sqlite+aiosqlite:///kessler.sqlite",
    session_config=session_config,
)

sqlalchemy_plugin = SQLAlchemyInitPlugin(config=sqlalchemy_config)

T = TypeVar("T")


class RepoMixin:
    """
    Implements common class level functionality for all objects.
    Wraps the repository
    """

    # TODO: generalize this
    # @classmethod
    # async def repo_factory(cls, RepoClass: Generic[T]) -> AsyncIterator[T]:
    #     session_factory = sqlalchemy_config.create_session_maker()

    #     async with session_factory() as db_session:
    #         try:
    #             yield RepoClass(session=db_session)
    #         except Exception as e:
    #             print(traceback.format_exc())
    #             print("rolling back")
    #             await db_session.rollback()
    #         else:
    #             print("committhing change")
    #             await db_session.commit()

    @classmethod
    async def new(cls, **kw):
        async with cls.repo() as repo:
            new = cls(**kw)
            obj = await repo.add(new)
            print(f"created {obj.__dict__}")
            return obj

    @classmethod
    async def create(cls, newObj):
        """Passing an object to create will add it"""
        async with cls.repo() as repo:
            obj = await repo.add(newObj)
            print(f"added {obj.__dict__}")
            return obj

    @classmethod
    async def remove(cls, id):
        async with cls.repo() as repo:
            obj = await repo.delete(id)
            print(f"deleted {obj.__dict__}")
            return obj

    @classmethod
    async def find(cls, id):
        async with cls.repo() as repo:
            obj = await repo.get(id)
            return obj
