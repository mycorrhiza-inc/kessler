# from lance_store.tables import BaseLance, Metadata
from lancedb.pydantic import LanceModel
from lance_store.embeddings import func


# class SummaryMetadata(Metadata):


class Summaries(LanceModel):
    __tablename__ = "summaries"
    text: str = func.SourceField()
