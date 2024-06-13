from litestar.contrib.sqlalchemy.base import UUIDAuditBase, UUIDBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy import ForeignKey, Column, Table
from sqlalchemy.orm import Mapped, relationship, DeclarativeBase
from sqlalchemy.types import PickleType

from pydantic import BaseModel, Field, field_validator
from contextlib import asynccontextmanager


from typing import List, Annotated, Any

import traceback
import uuid


from utils import sqlalchemy_config

from models.files import FileModel


resource_association_table = Table(
    "resource_association_table",
    DeclarativeBase.metadata,
    Column("left_id", ForeignKey("resource.id"), primary_key=True),
    Column("right_id", ForeignKey("resource.id"), primary_key=True),
)


class ResourceModel(UUIDAuditBase):
    __tablename__ = "resource"
    # how we manage multiple files for one resource
    files: Mapped[List[FileModel]] = relationship("File", back_populates="parent")
    # how we manage multiple links for one resource
    links: Mapped[List[FileModel]] = relationship("File", back_populates="parent")
    # how we manage resource trees
    children: Mapped[List["ResourceModel"]] = relationship(
        secondary=resource_association_table, back_populates="children"
    )
    parents: Mapped[List["ResourceModel"]] = relationship(
        secondary=resource_association_table, back_populates="parents"
    )

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class ResourceSchema(BaseModel):
    id: Annotated[Any, Field(validate_default=True)]
    files: List[any]
    document: List[any]
    children: List[any]
    parents: List[any]

    @field_validator("id")
    @classmethod
    def stringify_id(cls, id: any) -> str:
        return str(id)

    @classmethod
    def update(cls):
        """get the most up to date version of this resource"""


class ResourceRepository(SQLAlchemyAsyncRepository[ResourceModel]):
    """File repository."""

    model_type = ResourceModel
