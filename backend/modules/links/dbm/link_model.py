from typing import Optional

from sqlalchemy.orm import Mapped, mapped_column

from litestar.contrib.sqlalchemy.base import UUIDAuditBase, Base as SQLABase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy.ext.asyncio import AsyncSession
from pydantic import validator

from sqlalchemy import ForeignKey


class LinkResourceModel(SQLABase):
    __tablename__ = "LinkResource"
    # used to get Links from resource IDs
    id = mapped_column(ForeignKey("resource.id"), primary_key=True)
    # these are different so we can update link objects without regard to their resource id
    link_id = mapped_column(ForeignKey("link.id"))


class LinkModel(UUIDAuditBase):
    __tablename__ = "Link"
    # path exists if the link has been stored as text
    text_id: Mapped[ForeignKey | None] = mapped_column(
        ForeignKey("text_object.id"))
    url: Mapped[str]
    doctype: Mapped[str] # webpage and the rest
    lang: Mapped[str] # en etc
    title: Mapped[
        str
    ]
    metadata: Mapped[str]
    stage: Mapped[str]  # Either "stage0" "stage1" "stage2" or "stage3"
    summary: Mapped[str]
    short_summary: Mapped[str]

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
