from typing import Annotated, Any, List

from common.file_schemas import FileSchema
from litestar.contrib.sqlalchemy.base import UUIDAuditBase
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy.orm import Mapped


from .utils import PydanticBaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from pydantic import Field, field_validator, TypeAdapter

from uuid import UUID


import json


import logging

from enum import Enum

from sqlalchemy import select


class FileModel(UUIDAuditBase):
    """Database representation of a file"""

    __tablename__ = "file"
    url: Mapped[str | None]
    doctype: Mapped[str | None]
    lang: Mapped[str | None]
    name: Mapped[str | None]
    source: Mapped[str | None]
    hash: Mapped[str | None]
    mdata: Mapped[str | None]
    stage: Mapped[str | None]
    summary: Mapped[str | None]
    short_summary: Mapped[str | None]
    original_text: Mapped[str | None]
    english_text: Mapped[str | None]


# class FileModel(UUIDAuditBase):
#     """Database representation of a file"""
#
#     __tablename__ = "file"
#     url: Mapped[str | None]
#     doctype: Mapped[str | None]
#     name: Mapped[str | None]
#     source: Mapped[str | None]
#     hash: Mapped[str | None]
#     stage: Mapped[str | None]
#     summary: Mapped[str | None]
#     short_summary: Mapped[str | None]


# class FileMetadataSource(UUIDAuditBase):
#     file_id: Mapped[UUID]
#     metadata_key: Mapped[str]
#     value: Mapped[str | None]


class FileTextSource(UUIDAuditBase):
    __tablename__ = "file_text_source"
    file_id: Mapped[UUID]
    is_original_text: Mapped[bool]
    language: Mapped[str]
    text: Mapped[str | None]


class FileTextSchema(PydanticBaseModel):
    file_id: UUID
    is_original_text: bool
    language: str
    text: str


async def get_texts_from_file_uuid(
    async_db_connection: AsyncSession, file_id: UUID
) -> List[FileTextSchema]:
    result = await async_db_connection.execute(
        select(
            FileTextSource.text,
            FileTextSource.language,
            FileTextSource.is_original_text,
        ).where(FileTextSource.file_id == file_id)
    )
    rows = result.fetchall()
    return [
        FileTextSchema(
            file_id=file_id,
            text=row.text,
            language=row.language,
            is_original_text=row.is_original_text,
        )
        for row in rows
    ]


async def get_original_text_from_file_uuid(
    async_db_connection: AsyncSession, file_id: UUID
) -> FileTextSchema | None:
    result = await async_db_connection.execute(
        select(
            FileTextSource.text,
            FileTextSource.language,
            FileTextSource.is_original_text,
        )
        .where(
            FileTextSource.file_id == file_id, FileTextSource.is_original_text == True
        )
        .limit(1)
    )
    row = result.fetchone()
    if row:
        return FileTextSchema(
            file_id=file_id,
            text=row.text,
            language=row.language,
            is_original_text=row.is_original_text,
        )
    return None


async def get_english_text_from_file_uuid(
    async_db_connection: AsyncSession, file_id: UUID
) -> FileTextSchema | None:
    result = await async_db_connection.execute(
        select(
            FileTextSource.text,
            FileTextSource.language,
            FileTextSource.is_original_text,
        )
        .where(FileTextSource.file_id == file_id, FileTextSource.language == "en")
        .limit(1)
    )
    row = result.fetchone()
    if row:
        return FileTextSchema(
            file_id=file_id,
            text=row.text,
            language=row.language,
            is_original_text=row.is_original_text,
        )
    return None


# Write a SQL Migration script that will take the english text of every document, and add it as a entry in the file text source table, with a unique UUID, the UUID of the original file as the file ID, a true value for the original text and "en" as the language?
# INSERT INTO file_text_source (id, file_id, is_original_text, language, text)
# SELECT uuid_generate_v4(), id, true, 'en', english_text
# FROM file
# WHERE english_text IS NOT NULL;

# SQLAlchemyAsyncRepository


class FileRepository(SQLAlchemyAsyncRepository[FileModel]):
    """File repository."""

    model_type = FileModel


default_logger = logging.getLogger(__name__)


async def provide_files_repo(db_session: AsyncSession) -> FileRepository:
    """This provides the default Authors repository."""
    default_logger.info(db_session)
    default_logger.info(type(db_session))
    assert isinstance(db_session, AsyncSession), f"Type is : {type(db_session)}"
    file_repo = FileRepository(session=db_session)
    default_logger.info(file_repo)
    default_logger.info(type(file_repo))
    return file_repo


def model_to_schema(model: FileModel) -> FileSchema:
    metadata_str = model.mdata
    model.mdata = None
    type_adapter = TypeAdapter(FileSchema)
    schema = type_adapter.validate_python(model)
    schema.mdata = json.loads(metadata_str)
    return schema


class FileSchemaWithText(FileSchema):
    id: Annotated[Any, Field(validate_default=True)]
    original_text: str | None = None
    english_text: str | None = None
