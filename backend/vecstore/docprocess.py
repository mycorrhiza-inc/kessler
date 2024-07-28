from rag.SemanticSplitter import SemanticSplitter
from .util import add_nodes, collection_name


def add_document_to_db(
    text: str, metadata: dict[str, str], collection_name: str = collection_name
):
    source_id = metadata.get("source_id")
    nodes = SemanticSplitter().process(text, source_id=str(source_id))
    add_nodes(nodes, collection_name=collection_name, metadata=metadata)


18
