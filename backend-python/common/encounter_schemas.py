from pydantic import BaseModel
from uuid import UUID

from typing import List

from datetime import datetime

from common.file_schemas import FileSchema

from common.org_schemas import OrganizationSchema


class Faction(BaseModel):
    name: str
    description: str
    orgs: List[OrganizationSchema]


class EncounterSchema(BaseModel):
    id: UUID
    name: str
    created_at: datetime
    document_set: List[FileSchema]
    description: str
    factions: List[Faction]
