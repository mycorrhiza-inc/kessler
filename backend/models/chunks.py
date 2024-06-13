from lance_store.tables import BaseLance
from lance_store.embeddings import func
from typing import Optional


class Chunks(BaseLance):
    # every single chunk
    __tablename__ = "text_chunks"
    text: str = func.SourceField()
    location: Optional[tuple[int, int] | None] = None  # row, col
