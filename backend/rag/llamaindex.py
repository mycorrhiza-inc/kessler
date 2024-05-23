import logging
import sys
import os

import asyncio
from sqlalchemy import select
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from typing import List, Dict, Tuple, Optional

# Import the FileModel from file.py
from models.files import FileModel

from llama_index.core import SimpleDirectoryReader, StorageContext
from llama_index.core import VectorStoreIndex
from llama_index.core import Settings
from llama_index.vector_stores.postgres import PGVectorStore
import textwrap
import openai
from llama_index.llms.groq import Groq


from sqlalchemy import make_url
from llama_index.core.response_synthesizers import CompactAndRefine
from llama_index.core.retrievers import QueryFusionRetriever
from llama_index.core.query_engine import RetrieverQueryEngine

import os


import psycopg2


logger = logging.getLogger()


logger.setLevel(logging.DEBUG)

handler = logging.StreamHandler(sys.stderr)


# Uncomment to see debug logs
logging.basicConfig(stream=sys.stdout, level=logging.DEBUG)
logging.getLogger().addHandler(logging.StreamHandler(stream=sys.stdout))


GROQ_API_KEY = os.environ["GROQ_API_KEY"]
Settings.llm = Groq(
    model="llama3-70b-8192", request_timeout=360.0, api_key=GROQ_API_KEY
)

openai.api_key = os.environ["OPENAI_API_KEY"]
# TODO : Change embedding model to use not openai.
# Settings.embed_model = OllamaEmbedding(model_name="nomic-embed-text")

connection_string_unknown = os.environ["DATABASE_CONNECTION_STRING"]

# FIXME : Go ahead and try to figure out how to get this to work asynchronously assuming that is an important thing to do.
if "postgresql+asyncpg://" in connection_string_unknown:
    sync_postgres_connection_string = connection_string_unknown.replace(
        "postgresql+asyncpg://", "postgresql://"
    )
else:
    sync_postgres_connection_string = connection_string_unknown

async_postgres_connection_string = sync_postgres_connection_string.replace("postgresql://","postgresql+asyncpg://")


db_name = "postgres"
vec_table_name = "demo_vectordb"
file_table_name = "file"


conn = psycopg2.connect(sync_postgres_connection_string)
conn.autocommit = True




url = make_url(sync_postgres_connection_string)
hybrid_vector_store = PGVectorStore.from_params(
    database=db_name,
    host=url.host,
    password=url.password,
    port=url.port,
    user=url.username,
    table_name=vec_table_name,
    embed_dim=1536,  # openai embedding dimension
    hybrid_search=True,
    text_search_config="english",
)
storage_context = StorageContext.from_defaults(
    vector_store=hybrid_vector_store
)

initial_documents = []
hybrid_index = VectorStoreIndex.from_documents(
    initial_documents, storage_context=storage_context
)

vector_retriever = hybrid_index.as_retriever(
    vector_store_query_mode="default",
    similarity_top_k=5,
)
text_retriever = hybrid_index.as_retriever(
    vector_store_query_mode="sparse",
    similarity_top_k=5,  # interchangeable with sparse_top_k in this context
)
retriever = QueryFusionRetriever(
    [vector_retriever, text_retriever],
    similarity_top_k=5,
    num_queries=1,  # set this to 1 to disable query generation
    mode="relative_score",
    use_async=False,
)

response_synthesizer = CompactAndRefine()
query_engine = RetrieverQueryEngine(
    retriever=retriever,
    response_synthesizer=response_synthesizer,
)
from llama_index.core import Document




async def add_document_to_db_from_hash(hash_str: str) -> None:
    async def query_file_table_for_hash(hash: str) -> Tuple[any, any]:
        # Create an async engine and session
        engine = create_async_engine(async_postgres_connection_string, echo=True)
        async_session_maker = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

        async with async_session_maker() as session:
            async with FileModel.repo() as repo:
                # Create a query to select the first row matching the given id
                stmt = select(FileModel).where(FileModel.hash == hash)
                result = await session.execute(stmt)
                file_row = result.scalars().first()

                if file_row:
                    english_text = file_row.english_text
                    document_metadata = file_row.doc_metadata
                else:
                    english_text = None
                    document_metadata = None

                return (english_text, document_metadata)
    return_tuple = await query_file_table_for_hash(hash_str)
    english_text = return_tuple[0]
    doc_metadata = return_tuple[1]
    if not english_text is None:
        assert isinstance(english_text,str)
        assert isinstance(doc_metadata,dict)
        # TODO : Add support for metadata filtering 
        additional_document = Document(text = english_text, metadata = doc_metadata)
        additional_document.doc_id = str(hash) 
        # FIXME : Make sure the UUID matches the other function, and dryify this entire fucking mess.
        hybrid_index.insert(additional_document)
    return None


async def add_all_documents_to_db() -> None:
    async def query_file_table_for_all_rows() -> List[Tuple[str, dict]]:
        # Create an async engine and session
        engine = create_async_engine(async_postgres_connection_string, echo=True)
        async_session_maker = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

        async with async_session_maker() as session:
            async with FileModel.repo() as repo:
                # Create a query to select all rows
                stmt = select(FileModel)
                result = await session.execute(stmt)
                file_rows = result.scalars().all()

                documents = []
                for file_row in file_rows:
                    if file_row.english_text is not None and isinstance(file_row.doc_metadata, dict):
                        documents.append((file_row.english_text, file_row.doc_metadata))
                
                return documents
    
    documents = await query_file_table_for_all_rows()
    
    for english_text, doc_metadata in documents:
        if not english_text is None: 
            additional_document = Document(text=english_text, metadata=doc_metadata)
            additional_document.doc_id = str(doc_metadata.get('hash', ''))  # Setting a value for `doc_id`, customize as needed
            # FIXME: Ensure the UUID generation strategy matches the rest of your system.
            hybrid_index.insert(additional_document)
    return None

async def regenerate_vector_database_from_file_table() -> None:
    await add_all_documents_to_db()
    return None




def create_rag_response_from_query( query : str ):
    return str(query_engine.query(query))


