from sqlalchemy.orm import Mapped
from sqlalchemy.ext.asyncio import AsyncSession

from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from pydantic import validator

from utils import RepoCrudMixin, RepoMixin


class ResourceModel(UUIDAuditBase, RepoMixin, RepoCrudMixin):
    """
    A general Identifier to any given resource
    """

    __tablename__ = "resource"
    metadata: Mapped[str]

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class ResourceRepository(SQLAlchemyAsyncRepository[ResourceModel]):
    """File repository."""

    model_type = ResourceModel


async def provide_resource_repo(db_session: AsyncSession) -> ResourceModel:
    """This provides the default File repository."""
    return ResourceModel(session=db_session)
