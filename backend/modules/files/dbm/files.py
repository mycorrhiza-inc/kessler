from sqlalchemy.orm import Mapped

from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy.ext.asyncio import AsyncSession
from pydantic import validator


from pathlib import Path


class FileModel(UUIDAuditBase):
    __tablename__ = "file"
    hash: Mapped[
        str
    ]  # Blake2. For the file database this should absolutely be the primary key,
    path: Mapped[str]  # No type for os.pathlib type Path
    doctype: Mapped[str]
    lang: Mapped[str]
    name: Mapped[
        str
    ]  # I dont know if this should be included either in here or as a entry in doc_metadata, expecially since its only ever going to be used by the frontend. However, it might be an important query paramater and it seems somewhat irresponsible to not include.
    doc_metadata: Mapped[str]
    stage: Mapped[str]  # Either "stage0" "stage1" "stage2" or "stage3"
    summary: Mapped[str]
    short_summary: Mapped[str]

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class FileRepository(SQLAlchemyAsyncRepository[FileModel]):
    """File repository."""

    model_type = FileModel


async def provide_files_repo(db_session: AsyncSession) -> FileRepository:
    """This provides the default File repository."""
    return FileRepository(session=db_session)
