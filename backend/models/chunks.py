from lance_store.tables import BaseLance
from lancedb.pydantic import LanceModel, Vector
from lance_store.embeddings import func

from datetime import datetime

from typing import Optional


class Chunks(LanceModel):
    # every single chunk
    vector: Vector(func.ndims()) = func.VectorField()
    text: str = func.SourceField()
    created: datetime
    modified: datetime
    published: Optional[datetime | None]
    # since its a UUID the chance of collision is minimal
    parent_id: str
    chunk_id: str
    # row, col of the parent
    location: Optional[tuple[int, int] | None] = None
