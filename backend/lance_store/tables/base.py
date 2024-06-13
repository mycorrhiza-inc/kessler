from datetime import datetime
from typing import List, Optional
from lancedb.pydantic import LanceModel, Vector

from pydantic import BaseModel

from lance_store.embeddings import func


class Metadata(BaseModel):
    created: datetime
    published: Optional[datetime | None]


class ACL(BaseModel):
    # group defaults to none, meaning it is only allowed for the current user
    # other
    groups: Optional[List[str] | None] = None


class BaseLance(LanceModel):
    __tablename__: str
    vector: Vector(func.ndims()) = func.VectorField()
    metadata: Metadata
    # kessler wide resource id
    resource_id: str


class Chunks(BaseLance):
    # every single chunk
    __tablename__ = "text_chunks"
    parent_resource_id: str
    text: str = func.SourceField()
    pos: Optional[tuple[int, int] | None] = None  # row, col


class Summaries(BaseLance):
    __tablename__ = "summaries"
    text: str = func.SourceField()
