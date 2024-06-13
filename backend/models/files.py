from typing import Annotated, Any

from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship
from uuid import UUID


from .utils import RepoMixin, PydanticBaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from pydantic import Field, field_validator

from .resources import ResourceModel


class FileModel(UUIDAuditBase, RepoMixin):
    """Database representation of a file"""
    __tablename__ = "file"
    url: Mapped[str | None] = None
    doctype: Mapped[str]
    lang: Mapped[str]
    name: Mapped[str | None]
    source: Mapped[str]
    hash: Mapped[str]
    stage: Mapped[str]
    summary: Mapped[str | None]
    short_summary: Mapped[str | None]
    original_text: Mapped[str | None]
    english_text: Mapped[str | None]

    resource_id = mapped_column(UUID, ForeignKey("resource.id"))
    resource = relationship(ResourceModel)


class FileRepository(SQLAlchemyAsyncRepository[FileModel]):
    """File repository."""

    model_type = FileModel


async def provide_files_repo(db_session: AsyncSession) -> FileRepository:
    """This provides the default Authors repository."""
    return FileRepository(session=db_session)


class FileSchema(PydanticBaseModel):
    """pydantic schema of the FileModel"""

    id: Annotated[Any, Field(validate_default=True)]
    url: str | None = None
    path: str | None = None
    doctype: str
    lang: str
    name: str | None = None
    source: str
    stage: str
    metadata_str: str
    summary: str | None = None
    short_summary: str | None = None

    @field_validator("id")
    @classmethod
    def stringify_id(cls, id: any) -> str:
        return str(id)


class FileSchemaWithText(FileSchema):
    original_text: str | None = None
    english_text: str | None = None
