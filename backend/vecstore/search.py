from embeddings import embed
from .util import get_milvus_conn
import os
from typing import List

from .util import collection_name

import logging


logger = logging.getLogger(__name__)


def search(
    query,
    collection_name=collection_name,
    limit=10,
    filter="",
    output_fields: List[str] = ["source_id"],
) -> List[any]:
    print(f"searching for '{query}'")
    client = get_milvus_conn()
    query_vector = embed(query=query)[0]
    res = client.search(
        collection_name=collection_name,  # target collection
        data=[query_vector],  # query vector
        limit=limit,  # number of returned entities
        filter=filter,
        # specifies fields to be returned
        output_fields=output_fields,
    )
    print(f"res: {res}")
    return res
