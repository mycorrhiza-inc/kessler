from lance_store.tables import BaseLance, Metadata
from lance_store.embeddings import func


class SummaryMetadata(Metadata):


class Summaries(BaseLance):
    __tablename__ = "summaries"
    text: str = func.SourceField()
