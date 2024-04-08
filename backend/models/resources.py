from litestar.contrib.sqlalchemy.base import UUIDAuditBase, UUIDBase, AuditColumns
from litestar.contrib.sqlalchemy.repository import SQLAlchemyAsyncRepository

from sqlalchemy import ForeignKey, Column, Table
from sqlalchemy.orm import Mapped, relationship, DeclarativeBase
from sqlalchemy.types import PickleType

from pydantic import BaseModel, ConfigDict, StringConstraints, validator
from contextlib import asynccontextmanager

from utils import RepoMixin

from typing import AsyncIterator, List, Optional

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


class ResourceModel(UUIDAuditBase, RepoMixin):
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

    @classmethod
    async def provide_repo(cls, session) -> "ResourceRepository":
        return ResourceRepository(session=session)

    # # define the context manager for each file repo
    @classmethod
    @asynccontextmanager
    async def repo(cls) -> AsyncIterator["ResourceRepository"]:
        session_factory = sqlalchemy_config.create_session_maker()
        async with session_factory() as db_session:
            try:
                yield cls.provide_repo(session=db_session)
            except Exception as e:
                print(traceback.format_exc())
                print("rolling back")
                await db_session.rollback()
            else:
                print("committhing change")
                await db_session.commit()


class Resource(BaseModel):
    id: uuid  # TODO: figure out a better type for this UUID :/
    files: List[any]
    links: List[any]
    children: List[any]
    parents: List[any]

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value

    def update(cls):
        """get the most up to date version of this resource"""


class ResourceRepository(SQLAlchemyAsyncRepository[ResourceModel]):
    """File repository."""

    model_type = ResourceModel
