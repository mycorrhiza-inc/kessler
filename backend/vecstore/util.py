import os
from pymilvus import MilvusClient

milvus_user = os.environ.get("MILVUS_VEC_USER")
milvus_pass = os.environ.get("MILVUS_VEC_PASS")
milvus_host = os.environ.get("MILVUS_HOST")


def get_milvus_conn(uri: str = milvus_host):
    return MilvusClient(uri=uri, token=f"{milvus_user}:{milvus_pass}")


def drop_collection(collection_name=str):
    conn = get_milvus_conn()
    conn.drop_collection(collection_name=collection_name, timeout=10)
    pass


def create_document_collection(collection_name=str, dimension=1024):
    # using the defaults of the octo embedding
    conn = get_milvus_conn()
    conn.create_collection(
        collection_name=collection_name,
        dimension=dimension,
        index_file_size=1024,
        metric_type="IP",
        timeout=10,
    )


def reindex_collection(collection_name=str):
    conn = get_milvus_conn()

    # get all rows from a collection as an iterator
    rows = conn.query(
        collection_name=collection_name,
        filter=None,
        expr=None,
        output_fields=["*"],
        timeout=10,
    )
