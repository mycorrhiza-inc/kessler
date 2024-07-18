import os
from pymilvus import MilvusClient

milvus_user = os.environ.get("MILVUS_VEC_USER")
milvus_pass = os.environ.get("MILVUS_VEC_PASS")
milvus_host = os.environ.get("MILVUS_HOST")


def get_milvus_conn(uri: str = milvus_host):
    return MilvusClient(
        uri=uri,
        token=f'{milvus_user}:{milvus_pass}'
    )
