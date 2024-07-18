from embeddings import embed
from .util import get_milvus_conn
import os

collection_name = os.environ.get("MILVUS_COLLECTION_NAME", "PUC_dockets_prod")


def search(query, collection_name=collection_name, limit=10, filter=""):
    client = get_milvus_conn()
    query_vector = embed(query=query).data[0].embedding
    res = client.search(
        collection_name=collection_name,  # target collection
        data=[query_vector],  # query vector
        limit=limit,  # number of returned entities
        filter=filter,
        # specifies fields to be returned
        output_fields=['*']
    )
    return res
