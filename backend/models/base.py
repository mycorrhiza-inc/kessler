from datetime import datetime
from typing import List, Optional
from lancedb.pydantic import LanceModel, Vector

from pydantic import BaseModel

from lance_store.embeddings import func


class ACL(BaseModel):
    # this class is necessary for RLS in lancedb

    # group defaults to none, meaning it is only
    # allowed for the current user. if len(groups) is 1
    # and its only "public" then it must be added by the team
    groups: Optional[List[str] | None] = None
    # here to say which users are allowed to change
    # this resource. Defaults to the user who created it
    admins: List[str]


class Metadata(BaseModel):
    created: datetime
    modified: datetime
    published: Optional[datetime | None]
    parent_resource_id: str
    resource_id: str
    controls: ACL


class BaseLance(LanceModel):
    __tablename__: str
    vector: Vector(func.ndims()) = func.VectorField()
    metadata: Metadata
    # kessler wide resource id
    resource_id: str
