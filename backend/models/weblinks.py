from typing import Optional

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.ext.asyncio import AsyncSession

from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from pydantic import validator


class LinkModel(UUIDAuditBase):
    __tablename__ = "Link"
    # path exists if the link has been stored as text
    text_id: Mapped[ForeignKey | None] = mapped_column(
        ForeignKey("text_object.id"))
    url: Mapped[str]
    doctype: Mapped[str]  # webpage and the rest
    lang: Mapped[str]  # en etc
    title: Mapped[
        str
    ]
    stage: Mapped[str]  # Either "stage0" "stage1" "stage2" or "stage3"
    summary: Mapped[str]
    short_summary: Mapped[str]

    # TODO: implement this
    def bumpLinkResourceTimestamp():
        pass

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class LinkRepository(SQLAlchemyAsyncRepository[LinkModel]):
    """Link repository."""

    model_type = LinkModel


async def provide_Links_repo(db_session: AsyncSession) -> LinkRepository:
    """This provides the default Link repository."""
    return LinkRepository(session=db_session)
