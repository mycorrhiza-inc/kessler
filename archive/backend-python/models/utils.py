from typing import TypeVar
import os

from litestar.contrib.sqlalchemy.plugins import (
    AsyncSessionConfig,
    SQLAlchemyAsyncConfig,
    SQLAlchemyInitPlugin,
)

from pydantic import BaseModel as _BaseModel

from sqlalchemy.ext.asyncio import AsyncSession

session_config = AsyncSessionConfig(expire_on_commit=False)

postgres_connection_string = os.environ["DATABASE_CONNECTION_STRING"]
if "postgresql://" in postgres_connection_string:
    postgres_connection_string = postgres_connection_string.replace(
        "postgresql://", "postgresql+asyncpg://"
    )

sqlalchemy_config = SQLAlchemyAsyncConfig(
    connection_string=postgres_connection_string,
    session_config=session_config,
    # extend_existing=True
)


async def provide_async_session(db_session: AsyncSession) -> AsyncSession:
    return db_session


sqlalchemy_plugin = SQLAlchemyInitPlugin(config=sqlalchemy_config)

T = TypeVar("T")


class PydanticBaseModel(_BaseModel):
    """Extend Pydantic's BaseModel to enable ORM mode"""

    model_config = {"from_attributes": True, "arbitrary_types_allowed": True}
