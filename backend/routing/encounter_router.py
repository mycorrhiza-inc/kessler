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
    description: str
    parent_org_id: Optional[UUID]
    author_names: List[str]  # Names that the organisation authors documents under


class Faction(BaseModel):
    description: str
    position_float: Optional[float] = None
    orgs: List[Organization]


class Encounter(BaseModel):
    id: UUID
    created_at: datetime
    document_set: List[UUID]
    description: str
    factions: List[Faction]


class SeedEncounterData(BaseModel):
    description: Optional[str] = None
    query: Optional[QueryData] = None
    document_uuids: Optional[List[UUID]] = None


class EncounterController(Controller):
    """Encounter Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @post(path="/encounter/seed")
    async def seed(
        self, files_repo: FileRepository, data: SeedEncounterData, request: Request
    ) -> Encounter:
        logger = request.logger
        initial_documents = []
        if data.document_uuids is not None:
            docs = await get_files_from_uuids(files_repo, data.document_uuids)
            initial_documents.append(docs)

        if data.query is not None:
            queries = querydata_to_filters_strict(data.query)
            files = await files_repo.list(queries)
            initial_documents.append(files)
