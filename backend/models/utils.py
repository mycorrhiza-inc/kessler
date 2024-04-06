from typing import AsyncIterator
from contextlib import asynccontextmanager

from sqlalchemy.ext.declarative import declared_attr

from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository
from litestar.contrib.sqlalchemy.base import UUIDBase

from litestar.contrib.sqlalchemy.plugins import \
    AsyncSessionConfig, \
    SQLAlchemyAsyncConfig, \
    SQLAlchemyInitPlugin

session_config = AsyncSessionConfig(expire_on_commit=False)

sqlalchemy_config = SQLAlchemyAsyncConfig(
    # connection_string="sqlite+aiosqlite:///instance/kessler.sqlite",
    connection_string="sqlite+aiosqlite:///kessler.sqlite",
    session_config=session_config
)

sqlalchemy_plugin = SQLAlchemyInitPlugin(config=sqlalchemy_config)


class RepoMixin:

    def __init__(cls):
        child_model_type = cls.__class__

        class DerivedRepo(SQLAlchemyAsyncRepository[child_model_type]):
            model_type = child_model_type
        cls.__derived_repository__ = DerivedRepo

    @classmethod
    @asynccontextmanager
    async def repo(cls) -> AsyncIterator['RepoMixin.__derived_repository__']:
        session_factory = sqlalchemy_config.create_session_maker()

        async with session_factory() as db_session:
            try:
                # 
                yield cls.__derived_repository__(session=db_session)
            except Exception as e:
                await db_session.rollback()
            else:
                await db_session.commit()

    @classmethod
    async def create(cls, **kw):
        async with cls.repo() as repo:
            f = cls(**kw)
            obj = await repo.add(f)
            print(obj.__dict__)
            return obj

    @classmethod
    async def add(cls, new_row):
        async with cls.repo() as repo:
            obj = await repo.add(new_row)
            print(obj.__dict__)
            return obj
