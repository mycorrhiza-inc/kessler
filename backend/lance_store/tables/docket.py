from .base import Metadata, BaseLance
from typing import List

from lance_store.embeddings import func


class DocketMetadata(Metadata):
    proceeding_id: str
    industry: List[str]
    supplementary_document: bool
    supplementary_document_type: str
    # possibly have multiple govt entities involved
    government_entity: List[str]
    # allows us to look of ex: "CO" and "Colorado" or "Federal"
    state: List[str]
    paritcipants: List[str]
