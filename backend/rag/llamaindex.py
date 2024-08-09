from llama_index.embeddings.octoai import OctoAIEmbedding
import logging
import sys
import os

from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from typing import List, Tuple, Optional

# Import the FileModel from file.py
from models.files import FileModel

from llama_index.core import StorageContext
from llama_index.core import VectorStoreIndex
from llama_index.core import Settings

from llama_index.core.node_parser import SentenceWindowNodeParser

# from llama_index.vector_stores.postgres import PGVectorStore
# import psycopg2
import openai
from llama_index.llms.groq import Groq
from llama_index.llms.openai import OpenAI
from llama_index.llms.octoai import OctoAI
from llama_index.llms.fireworks import Fireworks


from sqlalchemy import make_url
from llama_index.core.response_synthesizers import CompactAndRefine
from llama_index.core.retrievers import QueryFusionRetriever
from llama_index.core.query_engine import RetrieverQueryEngine

# import MilvusVectorStore
from llama_index.vector_stores.milvus import MilvusVectorStore


from llama_index.core import Document


from constants import OPENAI_API_KEY, OCTOAI_API_KEY, GROQ_API_KEY, FIREWORKS_API_KEY

logger = logging.getLogger()


logger.setLevel(logging.DEBUG)

handler = logging.StreamHandler(sys.stderr)


def get_llm_from_model_str(model_name: Optional[str]):
    if model_name is None:
        model_name = "llama-405b"
    if model_name in ["llama-8b", "llama-3.1-8b-instant"]:
        actual_name = "llama-3.1-8b-instant"
        return Groq(model=actual_name, request_timeout=60.0, api_key=GROQ_API_KEY)
    if model_name in [
        "llama-70b",
        "llama3-70b-8192",
        "llama-3.1-70b-versatile",
    ]:
        actual_name = "llama-3.1-70b-versatile"
        return Groq(model=actual_name, request_timeout=60.0, api_key=GROQ_API_KEY)
    if model_name in ["llama-405b", "llama-3.1-405b-reasoning"]:
        actual_name = "accounts/fireworks/models/llama-v3p1-405b-instruct"
        return Fireworks(model=actual_name, api_key=FIREWORKS_API_KEY)
    if model_name in ["gpt-4o"]:
        return OpenAI(model=model_name, request_timeout=60.0, api_key=OPENAI_API_KEY)
    else:
        raise Exception("Model String Invalid or Not Supported")


# Uncomment to see debug logs
logging.basicConfig(stream=sys.stdout, level=logging.DEBUG)
logging.getLogger().addHandler(logging.StreamHandler(stream=sys.stdout))


openai.api_key = OPENAI_API_KEY


# TODO : Change embedding model to use not openai.
# Settings.embed_model = OllamaEmbedding(model_name="nomic-embed-text")
Settings.embed_model = OctoAIEmbedding(api_key=OCTOAI_API_KEY)
Settings.llm = get_llm_from_model_str("gpt-4o")
# Settings.llm = get_llm_from_model_str("llama-70b")
connection_string_unknown = os.environ["DATABASE_CONNECTION_STRING"]

# FIXME : Go ahead and try to figure out how to get this to work asynchronously assuming that is an important thing to do.
if "postgresql+asyncpg://" in connection_string_unknown:
    sync_postgres_connection_string = connection_string_unknown.replace(
        "postgresql+asyncpg://", "postgresql://"
    )
else:
    sync_postgres_connection_string = connection_string_unknown

async_postgres_connection_string = sync_postgres_connection_string.replace(
    "postgresql://", "postgresql+asyncpg://"
)


db_name = "postgres"
vec_table_name = "demo_vectordb"
file_table_name = "file"


url = make_url(sync_postgres_connection_string)


def get_hybrid_vector_store() -> MilvusVectorStore:
    milvus_user = os.environ.get("MILVUS_VEC_USER")
    milvus_pass = os.environ.get("MILVUS_VEC_PASS")
    milvus_host = os.environ.get("MILVUS_HOST")

    from vecstore.search import collection_name

    hybrid_vector_store = MilvusVectorStore(
        uri=milvus_host,
        dim=1024,
        overwrite=False,
        # TODO: Change collection name for prod
        collection_name=collection_name,
        token=f"{milvus_user}:{milvus_pass}",
    )
    return hybrid_vector_store


def get_storage_context() -> StorageContext:
    hybrid_vector_store = get_hybrid_vector_store()
    storage_context = StorageContext.from_defaults(vector_store=hybrid_vector_store)
    return storage_context


def create_hybrid_index() -> VectorStoreIndex:
    hybrid_vector_store = get_hybrid_vector_store()

    hybrid_index = VectorStoreIndex.from_vector_store(vector_store=hybrid_vector_store)

    return hybrid_index


def create_retriever(hybrid_index: VectorStoreIndex) -> QueryFusionRetriever:
    vector_retriever = hybrid_index.as_retriever(
        vector_store_query_mode="default",
        similarity_top_k=5,
    )
    text_retriever = hybrid_index.as_retriever(
        vector_store_query_mode="sparse",
        similarity_top_k=5,
    )
    retriever = QueryFusionRetriever(
        [vector_retriever, text_retriever],
        similarity_top_k=5,
        num_queries=1,
        mode="relative_score",
        use_async=False,
    )

    return retriever


def create_query_engine(retriever: QueryFusionRetriever) -> RetrieverQueryEngine:
    response_synthesizer = CompactAndRefine()
    query_engine = RetrieverQueryEngine(
        retriever=retriever,
        response_synthesizer=response_synthesizer,
    )

    return query_engine


hybrid_index = create_hybrid_index()
retriever = create_retriever(hybrid_index)
query_engine = create_query_engine(retriever)


async def get_document_list_from_file_table() -> list:
    async def query_file_table_for_all_rows() -> List[Tuple[str, dict]]:
        # Create an async engine and session
        engine = create_async_engine(async_postgres_connection_string, echo=True)
        async_session_maker = sessionmaker(
            engine, expire_on_commit=False, class_=AsyncSession
        )

        async with async_session_maker() as session:
            # Create a query to select all rows
            stmt = session.select(FileModel)
            result = await session.execute(stmt)
            file_rows = result.scalars().all()

            documents = []
            for file_row in file_rows:
                if file_row.english_text is not None and isinstance(
                    file_row.doc_metadata, dict
                ):
                    documents.append((file_row.english_text, file_row.doc_metadata))

            return documents

    documents = await query_file_table_for_all_rows()
    document_list = []
    for english_text, doc_metadata in documents:
        if english_text is not None:
            additional_document = Document(text=english_text, metadata=doc_metadata)
            additional_document.doc_id = str(doc_metadata.get("hash"))
            document_list.append(document_list)
    return document_list


def add_document_to_db(doc: Document) -> None:
    # split the document into sentences
    parser = SentenceWindowNodeParser(include_metadata=True)
    nodes = parser.get_nodes_from_documents([doc])
    hybrid_index.insert_nodes(nodes)


def vecstore_metadata_error():
    logger.error("No metadata found for document. Skipping.")
    raise ValueError("No metadata found for document. Must include a source_id.")


# WARN: DANGEROUS


def add_document_to_db_from_text(text: str, metadata: Optional[dict] = None) -> None:
    if metadata is None or metadata.get("source_id") is None:
        logger.error("No metadata found for document. Skipping.")
        raise ValueError("No metadata found for document. Must include a source_id.")

    try:
        document = Document(text=str(text), metadata=metadata)
        add_document_to_db(document)
    except Exception as e:
        logger.error(f"Encountered error while adding document: {e}")
        logger.error("Trying again with no metadata")
        document = Document(text=str(text))
        add_document_to_db(document)
    return None


async def add_document_to_db_from_hash(hash_str: str) -> None:
    async def query_file_table_for_hash(hash: str) -> Tuple[any, any]:
        # Create an async engine and session
        # TODO: Refactor this to use the same session or common interface
        engine = create_async_engine(async_postgres_connection_string, echo=True)
        async_session_maker = sessionmaker(
            engine, expire_on_commit=False, class_=AsyncSession
        )

        async with async_session_maker() as session:
            # Create a query to select the first row matching the given id
            stmt = session.select(FileModel).where(FileModel.hash == hash)
            result = await session.execute(stmt)
            file_row = result.scalars().first()

            if file_row:
                english_text = file_row.english_text
                document_metadata = file_row.doc_metadata
                document_metadata["source_id"] = str(file_row.id)
            else:
                vecstore_metadata_error()

            return (english_text, document_metadata)

    return_tuple = await query_file_table_for_hash(hash_str)
    english_text = return_tuple[0]
    doc_metadata = return_tuple[1]
    if english_text is not None:
        assert isinstance(english_text, str)
        assert isinstance(doc_metadata, dict)
        # TODO : Add support for metadata filtering
        additional_document = Document(text=english_text, metadata=doc_metadata)
        additional_document.doc_id = str(hash)
        # FIXME : Make sure the UUID matches the other function, and dryify this entire fucking mess.
        add_document_to_db(additional_document)
    else:
        assert False, "English text not present for document."
    return None


async def regenerate_vector_database_from_file_table() -> None:
    document_list = await get_document_list_from_file_table()
    storage_context = get_storage_context()
    global hybrid_index
    hybrid_index = VectorStoreIndex.from_documents(
        document_list, storage_context=storage_context
    )


# Chat engine for rag
#
def create_rag_response_from_query(query: str):
    return str(query_engine.query(query))


def generate_chat_completion(chat_history: List[dict]) -> dict:
    llama_index_chat_history = sanitzie_chathistory_llamaindex(chat_history)
    chat_engine = hybrid_index.as_chat_engine(
        chat_mode="react", verbose=True, chat_history=llama_index_chat_history
    )
    response = chat_engine.chat("")
    response_str = str(response)
    chat_engine.reset()
    return {"role": "assistant", "content": response_str}
