from typing import Annotated, Any, List

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


class FileMetadataSource(UUIDAuditBase):
    file_id: Mapped[UUID]
    metadata_key: Mapped[str]
    value: Mapped[str | None]


class FileTextSource(UUIDAuditBase):
    __tablename__ = "file_text_source"
    file_id: Mapped[UUID]
    is_original_text: Mapped[bool]
    language: Mapped[str]
    text: Mapped[str | None]


# Write a SQL Migration script that will take the english text of every document, and add it as a entry in the file text source table, with a unique UUID, the UUID of the original file as the file ID, a true value for the original text and "en" as the language?
# INSERT INTO file_text_source (id, file_id, is_original_text, language, text)
# SELECT uuid_generate_v4(), id, true, 'en', english_text
# FROM file
# WHERE english_text IS NOT NULL;


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


class FileSchema(PydanticBaseModel):
    """pydantic schema of the FileModel"""

    id: Annotated[Any, Field(validate_default=True)]
    url: str | None = None
    hash: str | None = None
    doctype: str | None = None
    lang: str | None = None
    name: str | None = None
    source: str | None = None
    stage: str | None = None
    short_summary: str | None = None
    summary: str | None = None
    organization_id: UUID | None = None
    mdata: dict | None = None
    display_text: str | None = None

    # Good idea to do this for dict based mdata, instead wrote a custom function for it
    @field_validator("id")
    @classmethod
    def stringify_id(cls, id: any) -> str:
        return str(id)


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


class DocumentStatus(str, Enum):
    unprocessed = "unprocessed"
    completed = "completed"
    encounters_analyzed = "encounters_analyzed"
    organization_assigned = "organization_assigned"
    summarization_completed = "summarization_completed"
    embeddings_completed = "embeddings_completed"
    stage3 = "stage3"
    stage2 = "stage2"
    stage1 = "stage1"


# I am deeply sorry for not reading the python documentation ahead of time and storing the stage of processed strings instead of ints, hopefully this can atone for my mistakes


# This should probably be a method on documentstatus, but I dont want to fuck around with it for now
def docstatus_index(docstatus: DocumentStatus) -> int:
    match docstatus:
        case DocumentStatus.unprocessed:
            return 0
        case DocumentStatus.stage1:
            return 1
        case DocumentStatus.stage2:
            return 2
        case DocumentStatus.stage3:
            return 3
        case DocumentStatus.embeddings_completed:
            return 4
        case DocumentStatus.summarization_completed:
            return 5
        case DocumentStatus.organization_assigned:
            return 6
        case DocumentStatus.encounters_analyzed:
            return 7
        case DocumentStatus.completed:
            return 1000
