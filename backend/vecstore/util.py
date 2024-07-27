import os
from typing import Dict, List, Union
from pymilvus import MilvusClient, FieldSchema, CollectionSchema, DataType

milvus_user = os.environ.get("MILVUS_VEC_USER")
milvus_pass = os.environ.get("MILVUS_VEC_PASS")
milvus_host = os.environ.get("MILVUS_HOST")


def get_milvus_conn(uri: str = milvus_host) -> MilvusClient:
    return MilvusClient(uri=uri, token=f"{milvus_user}:{milvus_pass}")


class MilvusRow:
    def __init__(self, text: str, source_id: str = None, embedding: List[float] = None):
        self.text = text
        self.source_id = source_id
        self.doc_type = "node"
        self.embedding = embedding


class MilvusDoc(MilvusRow):
    def __init__(self, text: str, k_uuid: str):
        super().__init__(text=text, k_uuid=k_uuid)
        self.doc_type = "doc"


class MilvusNode(MilvusRow):
    def __init__(self, text: str, k_uuid: str):
        super().__init__(text=text, k_uuid=k_uuid)
        self.doc_type = "node"


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


def describe_collection_schema(collection_name: str) -> Dict[any]:
    conn = get_milvus_conn()
    return conn.describe_collection(
        collection_name=collection_name, timeout=10
    ).to_dict()


def create_doc_node_schema() -> CollectionSchema:
    id_field = FieldSchema(
        name="id",
        dtype=DataType.INT64,
        is_primary=True,
        auto_id=True,
        description="primary id",
    )
    chunk_id_field = FieldSchema(
        name="chunk_id",
        dtype=DataType.INT64,
        descriptioin="id of the chunk, should typically be 0 for a docucment",
        defaault_value=0,
    )
    text_filed = FieldSchema(
        name="text",
        dtype=DataType.VARCHAR,
        description="text",
        max_length=65535,  # allow the max length of a text field
    )
    source_id_filed = FieldSchema(
        name="source_id",
        dtype=DataType.VARCHAR,
        description="the source document id",
        max_length=256,  # allow the max length of a text field
    )
    rowtype_filed = FieldSchema(
        name="rowtype",
        dtype=DataType.VARCHAR,
        description="texte",
        max_length=256,  # allow the max length of a text field
    )
    embedding_field = FieldSchema(
        name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1024, description="vector"
    )

    # Enable partition key on a field if you need to implement multi-tenancy based on the partition-key field
    partition_field = FieldSchema(
        name="tenant",
        dtype=DataType.VARCHAR,
        max_length=256,  # uuid of the group of allowed users
        is_partition_key=True,
        defaultv_value="public",
    )

    schema = CollectionSchema(
        fields=[
            id_field,
            rowtype_filed,
            text_filed,
            source_id_filed,
            chunk_id_field,
            embedding_field,
            partition_field,
        ],
        auto_id=True,
        enable_dynamic_field=True,
        description="a collecton of documents and nodes",
    )
    schema.validate()

    return schema


def create_document_collection(collection_name=str, dimension=1024):
    # using the defaults of the octo embedding
    conn = get_milvus_conn()
    schema = create_doc_node_schema()
    conn.create_collection(
        collection_name=collection_name,
        dimension=dimension,
        index_file_size=1024,
        metric_type="IP",
        timeout=10,
        primary_field_name="id",
        schema=schema,
        index_params=None,  # Used for index specific pareams
    )


def reindex_document_chunks(docids: List[dict], collection_name=str):
    conn = get_milvus_conn()
    for doc in docids:
        conn.delete(
            collection_name=collection_name,
            filter=f'source_id like "{doc}" AND doc_type not like "doc"',
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

    conn = get_milvus_conn()
    nodes = [node.__dict__ for node in nodes]
    # add the same metadata to all nodes
    for i, node in enumerate(node):
        # add the metadata to the node dict
        nodes[i].update(metadata)

    conn.insert(
        collection_name=collection_name,
        data=nodes,
        timeout=10,
    )
