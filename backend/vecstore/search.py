from embeddings import embed
from .util import get_milvus_conn
import os
from typing import List

collection_name = os.environ.get("MILVUS_COLLECTION_NAME", "PUC_dockets_prod")


def search(
    query,
    collection_name=collection_name,
    limit=10,
    filter="",
    output_fields: List[str] = ["*"],
) -> List[any]:
    client = get_milvus_conn()
    query_vector = embed(query=query).data[0].embedding
    res = client.search(
        collection_name=collection_name,  # target collection
        data=[query_vector],  # query vector
        limit=limit,  # number of returned entities
        filter=filter,
        # specifies fields to be returned
        output_fields=output_fields,
    )
    return res
