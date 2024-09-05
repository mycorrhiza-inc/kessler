from pydantic import BaseModel
from uuid import UUID

from typing import List, Optional, Annotated, Any

from datetime import datetime

from .files import FileSchema

from .utils import PydanticBaseModel


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


class OrganizationSchema(PydanticBaseModel):
    id: UUID
    name: str
    description: Optional[str]
    org_type: str
    parent_org_id: Optional[UUID]
    pseudonames: List[str]  # Names that the organisation authors documents under
    current_authors: List[UUID]


class WorkHistory(PydanticBaseModel):
    start_date: datetime
    end_date: Optional[datetime]
    org_id: UUID
    description: str


class AuthorSchema(PydanticBaseModel):
    id: UUID
    name: str
    current_org: Optional[UUID]
    work_history: WorkHistory


class Faction(PydanticBaseModel):
    name: str
    description: str
    position_float: Optional[float] = None
    orgs: List[OrganizationSchema]


class EncounterSchema(PydanticBaseModel):
    id: UUID
    name: str
    created_at: datetime
    document_set: List[FileSchema]
    description: str
    factions: List[Faction]


class OrganizationModel(UUIDAuditBase):
    __tablename__ = "organization"
    name: Mapped[str]
    description: Mapped[str | None]
    parent_org_id: Mapped[UUID | None]
    org_type: Mapped[str | None]
    pseudonames: Mapped[List[str]]
    current_authors: Mapped[List[UUID]]


class AuthorModel(UUIDAuditBase):
    __tablename__ = "author"
    name: Mapped[str]
    current_org: Mapped[Optional[UUID]]
    work_history: Mapped[List["WorkHistoryModel"]]


class WorkHistoryModel(UUIDAuditBase):
    __tablename__ = "work_history"
    start_date: Mapped[datetime]
    end_date: Mapped[Optional[datetime]]
    org_id: Mapped[UUID]
    description: Mapped[str]
    author_id: Mapped[UUID]


class FactionModel(UUIDAuditBase):
    __tablename__ = "faction"
    name: Mapped[str]
    description: Mapped[str]
    position_float: Mapped[Optional[float]]
    orgs: Mapped[List[UUID]]


class EncounterModel(UUIDAuditBase):
    __tablename__ = "encounter"
    name: Mapped[str]
    created_at: Mapped[datetime]
    document_set: Mapped[List[UUID]]
    description: Mapped[str]
    factions: Mapped[List[FactionModel]]
