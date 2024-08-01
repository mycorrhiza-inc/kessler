import os
import uuid
from typing import Dict, List, Union
from pymilvus import MilvusClient, FieldSchema, CollectionSchema, DataType

import logging

logger = logging.getLogger(__name__)


collection_name = os.environ.get("MILVUS_COLLECTION_NAME", "PUC_dockets_prod")
milvus_user = os.environ.get("MILVUS_VEC_USER")
milvus_pass = os.environ.get("MILVUS_VEC_PASS")
milvus_host = os.environ.get("MILVUS_HOST")


def get_milvus_conn(uri: str = milvus_host) -> MilvusClient:
    return MilvusClient(uri=uri, token=f"{milvus_user}:{milvus_pass}")


class MilvusRow:
    def __init__(self, text: str, source_id: str, embedding: List[float] = None):
        self.text = text
        self.source_id = source_id
        self.rowtype = "node"
        self.embedding = embedding


class MilvusDoc(MilvusRow):
    def __init__(self, text: str, source_id: str, embedding: List[float]):
        super().__init__(text=text, source_id=source_id, embedding=embedding)
        self.rowtype = "doc"


class MilvusNode(MilvusRow):
    def __init__(self, text: str, source_id: str, embedding: List[float]):
        super().__init__(text=text, source_id=source_id, embedding=embedding)
        self.rowtype = "node"


def drop_collection(collection_name=str):
    conn = get_milvus_conn()
    conn.drop_collection(collection_name=collection_name, timeout=10)
    pass


def check_collction_exists(collection_name=str) -> bool:
    conn = get_milvus_conn()
    conn.list_collections()
    if collection_name in conn.list_collections():
        return True
    return False


def describe_collection_schema(collection_name: str) -> Dict[str, any]:
    conn = get_milvus_conn()
    return conn.describe_collection(
        collection_name=collection_name, timeout=10
    ).to_dict()


def create_doc_node_schema() -> CollectionSchema:
    id_filed = FieldSchema(
        name="id",
        dtype=DataType.INT64,
        is_primary=True,
        auto_id=True,
        description="primary id",
    )
    chunk_id_filed = FieldSchema(
        name="chunk_id",
        dtype=DataType.INT64,
        descriptioin="id of the chunk, should typically be 0 for a docucment",
        defaault_value=0,
    )
    text_field = FieldSchema(
        name="text",
        dtype=DataType.VARCHAR,
        description="text",
        max_length=65535,  # allow the max length of a text field
    )

    source_id_field = FieldSchema(
        name="source_id",
        dtype=DataType.VARCHAR,
        description="the source document id",
        max_length=256,  # allow the max length of a text field
    )
    rowtype_field = FieldSchema(
        name="rowtype",
        dtype=DataType.VARCHAR,
        description="text",
        max_length=256,  # allow the max length of a text field
    )
    embedding_filed = FieldSchema(
        name="embedding",
        dtype=DataType.FLOAT_VECTOR,
        dim=768,
        description="embedding vector",
    )

    # Enable partition key on a field if you need to implement multi-tenancy based on the partition-key field
    # _partition_filed = FieldSchema(
    #     name="tenant",
    #     dtype=DataType.VARCHAR,
    #     max_length=256,  # uuid of the group of allowed users
    #     is_partition_key=True,
    #     defaultv_value="public",
    # )
    schema = CollectionSchema(
        fields=[
            id_filed,
            chunk_id_filed,
            text_field,
            source_id_field,
            rowtype_field,
            embedding_filed,
        ],
        auto_id=True,
        enable_dynamic_field=True,
        description="a collecton of documents and nodes",
    )

    schema.verify()

    return schema


def create_document_collection(collection_name=str, dimension=1024):
    # using the defaults of the octo embedding
    conn = get_milvus_conn()
    schema = create_doc_node_schema()

    index_params = conn.prepare_index_params()
    index_params.add_index(
        field_name="embedding",
        metric_type="COSINE",
        index_type="IVF_FLAT",
        index_name="embedding_vector_index",
        params={"nlist": 128},
    )

    conn.create_collection(
        collection_name=collection_name,
        dimension=dimension,
        index_file_size=1024,
        metric_type="IP",
        timeout=10,
        primary_field_name="id",
        schema=schema,
        index_params=index_params,
    )


def reindex_document_chunks(source_id: Union[List[str], str], collection_name=str):
    if isinstance(source_id, str):
        source_id = [source_id]

    conn = get_milvus_conn()
    for doc in source_id:
        conn.delete(
            collection_name=collection_name,
            filter=f'source_id like "{doc}" AND rowtype not like "doc"',
        )
    # get all rows from a collection as an iterator


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


def add_nodes(
    nodes: Union[MilvusRow, List[MilvusRow]], collection_name: str, metadata: dict
):
    if isinstance(nodes, MilvusRow):
        nodes = [nodes]
    print(f"asdfasdfd: {nodes}")

    conn = get_milvus_conn()
    print(f"collection_name: {collection_name}")
    schema = conn.describe_collection(collection_name)
    print(f"collection schema: {schema}")
    nodes = [node.__dict__ for node in nodes]
    toadd = []
    print(f"there are {len(nodes)} nodes to add")
    # add the same metadata to all nodes
    for node in nodes:
        # add the metadata to the node dict
        node.update(metadata)
        # TODO: set this somewhere else
        node.update({"source_id": str(node["source_id"])})
        node.update({"chunk_id": 0})
        print(f"node: {node}")
        toadd.append(node)

    conn.insert(
        collection_name=collection_name,
        data=nodes,
        timeout=10,
    )
