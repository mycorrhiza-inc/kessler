from uuid import UUID
from pydantic import BaseModel

from pydantic import Field, field_validator, TypeAdapter

from typing import Annotated, Any, List


from enum import Enum

from common.org_schemas import OrganizationSchema, IndividualSchema


class FileTextSchema(BaseModel):
    file_id: UUID
    is_original_text: bool
    language: str
    text: str


class FileSchema(BaseModel):
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


class FileSchemaFull(BaseModel):
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
    texts: List[FileTextSchema] | None = []
    authors: List[IndividualSchema] | None = []
    organization: OrganizationSchema | None = None


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
