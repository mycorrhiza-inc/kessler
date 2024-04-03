from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy.ext.asyncio import AsyncSession
from pydantic import validator


class ResourceModel(UUIDAuditBase):
    __tablename__ = "resource"

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
