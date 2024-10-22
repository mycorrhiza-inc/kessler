from pydantic import BaseModel
from uuid import UUID

from typing import Optional



class OrganizationSchema(BaseModel):
    id: UUID
    name: str
    description: Optional[str]


class IndividualSchema(BaseModel):
    id: UUID
    name: str
    current_org: Optional[UUID]
