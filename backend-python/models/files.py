from typing import Annotated, Any, List, Optional

from common.file_schemas import FileSchema, FileSchemaFull
from common.org_schemas import IndividualSchema, OrganizationSchema
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

import asyncio

from common.niclib import Maybe


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


def file_model_to_schema(model: FileModel) -> FileSchema:
    model_dict = model.to_dict()
    if model.mdata is None:
        model.mdata = "{}"
    model_dict["mdata"] = json.loads(model.mdata)
    schema = FileSchema(**model_dict)
    return schema


async def get_partial_file_from_uuid(async_db_connection: AsyncSession, file_id: UUID):
    result = await async_db_connection.execute(
        select(FileModel).where(FileModel.id == file_id)
    )
    db_model = result.scalars().first()
    # YOU CANT TAKE MY MONADS FROM MEEEEEEE (Also the python type system isnt designed for this kind of torture, maybe go is better.)
    return Maybe(file_model_to_schema)(db_model)


async def get_full_file_from_uuid(
    async_db_connection: AsyncSession, file_id: UUID
) -> FileSchemaFull:
    tasks = [
        get_partial_file_from_uuid(async_db_connection, file_id),
        get_texts_from_file_uuid(async_db_connection, file_id),
    ]
    components = await asyncio.gather(*tasks)

    def piece_together_full_fileschema(
        partial_file_schema: FileSchema,
        texts: List[FileTextSchema] = [],
        authors: List[IndividualSchema] = [],
        org: Optional[OrganizationSchema] = None,
    ) -> FileSchemaFull:
        return_dict = partial_file_schema.dict()
        return_dict["texts"] = texts
        return_dict["authors"] = authors
        return_dict["organization"] = org
        return FileSchemaFull(**return_dict)

    return piece_together_full_fileschema(
        partial_file_schema=components[0], texts=components[1]
    )


async def delete_file_from_file_uuid(
    async_db_connection: AsyncSession, file_id: UUID, recursive: bool = False
):
    await async_db_connection.execute(FileModel.delete().where(FileModel.id == file_id))
    if recursive:
        await async_db_connection.execute(
            FileTextSource.delete().where(FileTextSource.file_id == file_id)
        )


async def upsert_file_from_full_schema(
    async_db_connection: AsyncSession,
    file: FileSchemaFull,
    delete_old_recursive: bool = False,
):
    query = select(FileModel).where(FileModel.id == file.id)
    result = await async_db_connection.execute(query)
    existing_file = result.scalars().first()
    file_exists = existing_file is not None

    if delete_old_recursive and file_exists:
        await delete_file_from_file_uuid(async_db_connection, file.id, True)

    async def upsert_file_model(file: FileSchemaFull, file_exists: bool) -> None:
        mdata_str = json.dumps(file.mdata)
        values_dict = {
            "id": file.id,
            "url": file.url,
            "doctype": file.doctype,
            "lang": file.lang,
            "name": file.name,
            "source": file.source,
            "hash": file.hash,
            "mdata": mdata_str,
            "stage": file.stage,
            "summary": file.summary,
            "short_summary": file.short_summary,
        }
        if file_exists:
            update_stmt = (
                FileModel.update().where(FileModel.id == file.id).values(**values_dict)
            )
            await async_db_connection.execute(update_stmt)
        else:
            insert_stmt = FileModel.insert().values(**values_dict)
            await async_db_connection.execute(insert_stmt)

    async def upsert_file_text_source(text_source: FileTextSchema) -> None:
        await async_db_connection.execute(
            select(FileTextSource).where(
                FileTextSource.file_id == text_source.file_id,
                FileTextSource.language == text_source.language,
            )
        )
        file_text = result.scalars().first()
        if file_text is None:
            insert_stmt = (
                FileTextSource.insert()
                .values(
                    id=UUID(),
                    file_id=text_source.file_id,
                    is_original_text=text_source.is_original_text,
                    language=text_source.language,
                    text=text_source.text,
                )
                .on_conflict_do_nothing()
            )
            await async_db_connection.execute(insert_stmt)
        else:
            update_stmt = (
                FileTextSource.update()
                .where(
                    FileTextSource.id == text_source.id,
                )
                .values(
                    text=text_source.text, is_original_text=text_source.is_original_text
                )
            )
            await async_db_connection.execute(update_stmt)

    if file.texts is None:
        file.texts = []
    async_tasks = [upsert_file_from_full_schema(async_db_connection, file)] + [
        upsert_file_text_source(text) for text in file.texts
    ]
    await asyncio.gather(*async_tasks)
    await async_db_connection.commit()


# I wrote these 2 initially, but I think they might be redundant at this point.
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


class FileSchemaWithText(FileSchema):
    id: Annotated[Any, Field(validate_default=True)]
    original_text: str | None = None
    english_text: str | None = None
