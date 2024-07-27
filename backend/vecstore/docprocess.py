from rag.SemanticSplitter import SemanticSplitter
from .util import add_nodes


def add_document_to_db(text: str, metadata: dict[str, str], collection_name: str):
    nodes = SemanticSplitter().process(text)
    add_nodes(nodes, collection_name=collection_name, metadata=metadata)
