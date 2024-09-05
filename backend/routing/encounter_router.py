from uuid import UUID
import logging

from litestar import Controller, Request


from litestar.params import Parameter

from litestar.handlers.http_handlers.decorators import (
    post,
)


from litestar.di import Provide

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from models.files import (
    FileModel,
    FileRepository,
    FileSchema,
    model_to_schema,
    provide_files_repo,
)


from typing import List, Optional, Union, Any, Dict


from rag.rag_engine import KeRagEngine, convert_search_results_to_frontend_table
from vecstore import search

import json
import asyncio

from constants import lemon_text


from datetime import date, datetime

from logic.databaselogic import (
    QueryData,
    get_files_from_uuids,
    querydata_to_filters_strict,
)


class Organization(BaseModel):
    id: UUID
    name : str
    description: str
    parent_org_id: Optional[UUID]
    author_names: List[str]  # Names that the organisation authors documents under


class Faction(BaseModel):
    name: str
    description: str
    position_float: Optional[float] = None
    orgs: List[Organization]


class EncounterSchema(BaseModel):
    id: UUID
    name : str
    created_at: datetime
    document_set: List[FileSchema]
    description: str
    factions: List[Faction]


class SeedEncounterData(BaseModel):
    name : str
    description: Optional[str] = None
    query: Optional[QueryData] = None
    document_uuids: Optional[List[UUID]] = None


# TODO:
# - Generate a database of organizations, and associate each with an individual author name. Add a couple quick llm prompts to see if the author is a person or an org.
# - Add support for a internet rag agent as a tool that can be called by our larger rag agents. This is a good idea to in house over time, but for the next 2 months, I think its a good idea to try to use an API from perplexity or something. Specifically because getting rag agents to work with generic internet data is really hard, and is also something we are trying to avoid due the inevitable loss in quality. But I think its a very good idea, since there are lots of questions that our agents will need to answer, that will require a simple google search.


class EncounterController(Controller):
    """Encounter Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @post(path="/encounter/seed")
    async def seed(
        self, files_repo: FileRepository, data: SeedEncounterData, request: Request
    ) -> EncounterSchema:
        logger = request.logger
        initial_documents = []
        if data.document_uuids is not None:
            docs = await get_files_from_uuids(files_repo, data.document_uuids)
            initial_documents.append(docs)

        if data.query is not None:
            queries = querydata_to_filters_strict(data.query)
            files = await files_repo.list(queries)
            initial_documents.append(files)
        description = data.description
        if description is None:
            description = ""
        id = UUID()
        initial_encounter = EncounterSchema(
            description=description,
            document_set=initial_documents,
            created_at=datetime.now(),
            factions=[],
            id=id,
        )
        await self.refine_initial_seed(files_repo, initial_encounter)
        return initial_encounter

    async def refine_initial_seed(
        self, files_repo: FileRepository, encounter: EncounterSchema
    ) -> EncounterSchema:
        return await self.refine_seed(files_repo, encounter)

    async def refine_seed(
        self, files_repo: FileRepository, encounter: EncounterSchema
    ) -> EncounterSchema:
        async def generate_encounter_description():
            # Take the document summaries, and the previous description if applicable and generate a new description for the encounter.
            encounter.description = lemon_text
        async def search_for_more_documents():
            # Generate search queries for more documents and add them to the document list.

        async def generate_organisation_list():
            # Generate a list of orgs, include them all in 1 faction for now, the special "unknown" faction.


        if encounter.description == "":
            await generate_encounter_description()

        return encounter