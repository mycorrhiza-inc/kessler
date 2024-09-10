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
from sqlalchemy import select


from common.encounter_schemas import (
    EncounterSchema,
    EventSchema,
)

from common.org_schemas import OrganizationSchema, IndividualSchema


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


class RelationFileAuthoredByIndividual(UUIDAuditBase):
    __tablename__ = "relation_document_authored_by_individual"
    file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
    individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))


class RelationFileAssociatedWithOrganization(UUIDAuditBase):
    __tablename__ = "relation_document_associated_with_organization"
    file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
    organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))


class RelationIndividualsCurrentlyAssociatedOrganization(UUIDAuditBase):
    __tablename__ = "relation_individuals_currently_associated_organization"
    individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
    organization_id: Mapped[UUID] = mapped_column(ForeignKey("organization.id"))


class RelationFilesInEncounter(UUIDAuditBase):
    __tablename__ = "relation_documents_in_encounter"
    file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
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

    # TODO : Figure out better code reuse here:
    async def get_organizations_in_faction(self, faction_id: UUID):
        result = await self.db.execute(
            select(RelationOrganizationsInFaction.organization_id).where(
                RelationOrganizationsInFaction.faction_id == faction_id
            )
        )

        organization_ids = [row[0] for row in result.fetchall()]

        if not organization_ids:
            return []

        # Next, get the names and descriptions of these organizations
        organizations_result = await self.db.execute(
            select(
                OrganizationModel.id,
                OrganizationModel.name,
                OrganizationModel.description,
            ).where(OrganizationModel.id.in_(organization_ids))
        )
        organizations = organizations_result.fetchall()

        def gen_org_schema(org) -> OrganizationSchema:
            return OrganizationSchema(
                id=org.id, name=org.name, description=org.description
            )

        return list(map(gen_org_schema, organizations))

    async def get_individuals_in_faction(self, faction_id: UUID):
        result = await self.db.execute(
            select(RelationIndividualsInFaction.individual_id).where(
                RelationIndividualsInFaction.faction_id == faction_id
            )
        )

        individual_ids = [row[0] for row in result.fetchall()]

        if not individual_ids:
            return []

        individuals_result = await self.db.execute(
            select(
                IndividualModel.id,
                IndividualModel.name,
                IndividualModel.username,
                IndividualModel.chosen_name,
            ).where(IndividualModel.id.in_(individual_ids))
        )
        individuals = individuals_result.fetchall()

        def gen_individual_schema(ind) -> IndividualSchema:
            return IndividualSchema(
                id=ind.id,
                name=ind.name,
                username=ind.username,
                chosen_name=ind.chosen_name,
            )

        return list(map(gen_individual_schema, individuals))

    async def get_factions_in_encounter(self, encounter_id: UUID):
        result = await self.db.execute(
            select(RelationFactionsInEncounter.faction_id).where(
                RelationFactionsInEncounter.encounter_id == encounter_id
            )
        )

        faction_ids = [row[0] for row in result.fetchall()]

        if not faction_ids:
            return []

        factions_result = await self.db.execute(
            select(
                FactionModel.id,
                FactionModel.name,
                FactionModel.description,
            ).where(FactionModel.id.in_(faction_ids))
        )
        factions = factions_result.fetchall()

        def gen_faction_schema(faction) -> FactionSchema:
            return FactionSchema(
                id=faction.id, name=faction.name, description=faction.description
            )

        return list(map(gen_faction_schema, factions))

    async def get_individuals_currently_associated_with_organization(
        self, organization_id: UUID
    ):
        result = await self.db.execute(
            select(
                RelationIndividualsCurrentlyAssociatedOrganization.individual_id
            ).where(
                RelationIndividualsCurrentlyAssociatedOrganization.organization_id
                == organization_id
            )
        )

        individual_ids = [row[0] for row in result.fetchall()]

        if not individual_ids:
            return []

        individuals_result = await self.db.execute(
            select(
                IndividualModel.id,
                IndividualModel.name,
                IndividualModel.username,
                IndividualModel.chosen_name,
            ).where(IndividualModel.id.in_(individual_ids))
        )
        individuals = individuals_result.fetchall()

        def gen_individual_schema(ind) -> IndividualSchema:
            return IndividualSchema(
                id=ind.id,
                name=ind.name,
                username=ind.username,
                chosen_name=ind.chosen_name,
            )

        return list(map(gen_individual_schema, individuals))

    async def get_individuals_associated_with_event(self, event_id: UUID):
        result = await self.db.execute(
            select(RelationIndividualsAssociatedWithEvent.individual_id).where(
                RelationIndividualsAssociatedWithEvent.event_id == event_id
            )
        )

        individual_ids = [row[0] for row in result.fetchall()]

        if not individual_ids:
            return []

        individuals_result = await self.db.execute(
            select(
                IndividualModel.id,
                IndividualModel.name,
                IndividualModel.username,
                IndividualModel.chosen_name,
            ).where(IndividualModel.id.in_(individual_ids))
        )
        individuals = individuals_result.fetchall()

        def gen_individual_schema(ind) -> IndividualSchema:
            return IndividualSchema(
                id=ind.id,
                name=ind.name,
                username=ind.username,
                chosen_name=ind.chosen_name,
            )

        return list(map(gen_individual_schema, individuals))

    async def get_organizations_associated_with_event(self, event_id: UUID):
        result = await self.db.execute(
            select(RelationOrganizationsAssociatedWithEvent.organization_id).where(
                RelationOrganizationsAssociatedWithEvent.event_id == event_id
            )
        )

        organization_ids = [row[0] for row in result.fetchall()]

        if not organization_ids:
            return []

        organizations_result = await self.db.execute(
            select(
                OrganizationModel.id,
                OrganizationModel.name,
                OrganizationModel.description,
            ).where(OrganizationModel.id.in_(organization_ids))
        )
        organizations = organizations_result.fetchall()

        def gen_organization_schema(org) -> OrganizationSchema:
            return OrganizationSchema(
                id=org.id, name=org.name, description=org.description
            )

        return list(map(gen_organization_schema, organizations))


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
#     file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
#     individual_id: Mapped[UUID] = mapped_column(ForeignKey("individual.id"))
#
#
# class DocumentAssociatedWithOrganization(UUIDAuditBase):
#     __tablename__ = "document_associated_with_organization"
#     file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
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
#     file_id: Mapped[UUID] = mapped_column(ForeignKey("file.id"))
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
