from contextlib import asynccontextmanager
from typing import AsyncIterator, Annotated
import traceback

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.ext.asyncio import AsyncSession

from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from pydantic import validator

from models.utils import sqlalchemy_config

from haystack.components.converters import HTMLToDocument


converter = HTMLToDocument()
results = converter.run(sources=["path/to/sample.html"])
documents = results["documents"]
print(documents[0].content)


# 'This is a text from the HTML file.'
class WeblinkModel(UUIDAuditBase):
    __tablename__ = "Link"
    # path exists if the link has been stored as text
    text_id: Mapped[ForeignKey | None] = mapped_column(ForeignKey("text_object.id"))
    url: Mapped[str]
    doctype: Mapped[str]  # webpage and the rest
    lang: Mapped[str]  # en etc
    title: Mapped[str]
    stage: Mapped[str]  # Either "stage0" "stage1" "stage2" or "stage3"
    summary: Mapped[str]
    short_summary: Mapped[str]

    @classmethod
    @asynccontextmanager
    async def repo(cls) -> AsyncIterator["WeblinkRepository"]:
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

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class WeblinkRepository(SQLAlchemyAsyncRepository[WeblinkModel]):
    """Link repository."""

    model_type = WeblinkModel


async def provide_weblinks_repo(db_session: AsyncSession) -> WeblinkRepository:
    """This provides the default Link repository."""
    return WeblinkRepository(session=db_session)
