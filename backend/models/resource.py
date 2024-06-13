from typing import Annotated, Any

from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy.orm import Mapped


from .utils import RepoMixin, PydanticBaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from pydantic import Field, field_validator


class ResourceModel(UUIDAuditBase, RepoMixin):
    """unique resource id"""

    __tablename__ = "resource"
    table: Mapped[str]


class ResourceRepository(SQLAlchemyAsyncRepository[ResourceModel]):
    """Resource repository."""

    model_type = ResourceModel


async def provide_files_repo(db_session: AsyncSession) -> ResourceModel:
    """This provides the default Authors repository."""
    return ResourceRepository(session=db_session)


class ResourceSchema(PydanticBaseModel):
    """pydantic schema of the FileModel"""

    id: Annotated[Any, Field(validate_default=True)]
    table: str

    @field_validator("id")
    @classmethod
    def stringify_id(cls, id: any) -> str:
        return str(id)
