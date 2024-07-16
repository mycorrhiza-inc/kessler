import os
import pymilvus

milvus_url = os.environ.get("MILVUS_DB_URL")


def get_milvus_conn(url: str = milvus_url):
    pass
