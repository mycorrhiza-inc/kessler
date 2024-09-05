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

from sqlalchemy import ForeignKey
from sqlalchemy import Integer
from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.orm import relationship


class OrganizationSchema(PydanticBaseModel):
    id: UUID
    name: str
    description: Optional[str]
    org_type: str
    parent_org_id: Optional[UUID]
    pseudonames: List[str]  # Names that the organization authors documents under
    current_authors: List[UUID]




class IndividualSchema(PydanticBaseModel):
    id: UUID
    name: str
    current_org: Optional[UUID]


class Faction(PydanticBaseModel):
    name: str
    description: str
    orgs: List[OrganizationSchema]


class EncounterSchema(PydanticBaseModel):
    id: UUID
    name: str
    created_at: datetime
    document_set: List[FileSchema]
    description: str
    factions: List[Faction]


# Actual SQL Models
class FactionModel(UUIDAuditBase):
    __tablename__ = "faction"
    name: Mapped[str]
    description: Mapped[str]


class IndividualModel(UUIDAuditBase):
    __tablename__ = "individual"
    name: Mapped[str]
    username: Mapped[str | None]
    chosen_name: Mapped[str | None]


class OrganizationModel(UUIDAuditBase):
    __tablename__ = "organization"
    name: Mapped[str]
    description: Mapped[str | None]


class EncounterModel(UUIDAuditBase):
    __tablename__ = "encounter"
    name: Mapped[str | None]
    description: Mapped[str | None]


class EventModel(UUIDAuditBase):
    __tablename__ = "event"
    date: Mapped[datetime | None]
    name: Mapped[str | None]
    description: Mapped[str | None]


# ---------
# Testing
# Really stupid way of making tables, if getting errors when initiating the db, comment out all tables beneath this line, then uncomment and run again.
class RelationOrganizationsInFaction(UUIDAuditBase):
    __tablename__ = "relation_organizations_in_faction"
    faction_id: Mapped[UUID] = mapped_column(ForeignKey("faction.id"))
    organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))


class RelationIndividualsInFaction(UUIDAuditBase):
    __tablename__ = "relation_individuals_in_faction"
    faction_id: Mapped[UUID] = mapped_column(ForeignKey("faction.id"))
    individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))


class RelationFactionsInEncounter(UUIDAuditBase):
    __tablename__ = "relation_factions_in_encounter"
    encounter_id: Mapped[UUID] = mapped_column(ForeignKey("encounter.id"))
    faction_id: Mapped[UUID] = mapped_column(ForeignKey("faction.id"))


class RelationDocumentAuthoredByIndividual(UUIDAuditBase):
    __tablename__ = "relation_document_authored_by_individual"
    document_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
    individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))


class RelationDocumentAssociatedWithOrganization(UUIDAuditBase):
    __tablename__ = "relation_document_associated_with_organization"
    document_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
    organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))


class RelationIndividualsCurrentlyAssociatedOrganization(UUIDAuditBase):
    __tablename__ = "relation_individuals_currently_associated_organization"
    individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
    organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))


class RelationDocumentsInEncounter(UUIDAuditBase):
    __tablename__ = "relation_documents_in_encounter"
    document_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
    encounter_id: Mapped[UUID] = mapped_column(ForeignKey("encounter.id"))


class RelationIndividualsAssociatedWithEvent(UUIDAuditBase):
    __tablename__ = "relation_individuals_associated_with_event"
    individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
    event_id: Mapped[UUID] = mapped_column(ForeignKey("event.id"))


class RelationOrganizationsAssociatedWithEvent(UUIDAuditBase):
    __tablename__ = "relation_organizations_associated_with_event"
    organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))
    event_id: Mapped[UUID] = mapped_column(ForeignKey("event.id"))


class RelationFilesAssociatedWithEvent(UUIDAuditBase):
    __tablename__ = "relation_files_associated_with_event"
    file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
    event_id: Mapped[UUID] = mapped_column(ForeignKey("event.id"))


class SQLUtils:
    def __init__(self, async_db_connection: AsyncSession):
        self.db = async_db_connection

    async def get_organizations_in_faction(self, faction_id: UUID):
        
# class OrganizationsInFaction(UUIDAuditBase):
#     __tablename__ = "organizations_in_faction"
#     faction_id: Mapped[UUID] = mapped_column(ForeignKey("faction.id"))
#     organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))
#
#
# class IndividualsInFaction(UUIDAuditBase):
#     __tablename__ = "individuals_in_faction"
#     faction_id: Mapped[UUID] = mapped_column(ForeignKey("faction.id"))
#     individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
#
#
# class FactionsInEncounter(UUIDAuditBase):
#     __tablename__ = "factions_in_encounter"
#     encounter_id: Mapped[UUID] = mapped_column(ForeignKey("encounter.id"))
#     faction_id: Mapped[UUID] = mapped_column(ForeignKey("faction.id"))
#
#
# class DocumentAuthoredByIndividual(UUIDAuditBase):
#     __tablename__ = "document_authored_by_individual"
#     document_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
#     individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
#
#
# class DocumentAssociatedWithOrganization(UUIDAuditBase):
#     __tablename__ = "document_associated_with_organization"
#     document_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
#     organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))
#
#
# class IndividualsCurrentlyAssociatedOrganization(UUIDAuditBase):
#     __tablename__ = "individuals_currently_associated_organization"
#     individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
#     organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))
#
#
# class DocumentsInEncounter(UUIDAuditBase):
#     __tablename__ = "documents_in_encounter"
#     document_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
#     encounter_id: Mapped[UUID] = mapped_column(ForeignKey("encounter.id"))
#
#
# class IndividualsAssociatedWithEvent(UUIDAuditBase):
#     __tablename__ = "individuals_associated_with_event"
#     individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
#     event_id: Mapped[UUID] = mapped_column(ForeignKey("event.id"))
#
#
# class OrganizationsAssociatedWithEvent(UUIDAuditBase):
#     __tablename__ = "organizations_associated_with_event"
#     organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))
#     event_id: Mapped[UUID] = mapped_column(ForeignKey("event.id"))
#
#
# class FilesAssociatedWithEvent(UUIDAuditBase):
#     __tablename__ = "files_associated_with_event"
#     file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
#     event_id: Mapped[UUID] = mapped_column(ForeignKey("event.id"))
