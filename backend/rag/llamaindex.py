import logging
import sys
import os



import asyncio
from sqlalchemy import select
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from typing import List, Dict, Tuple, Optional

# Import the FileModel from file.py
from file import FileModel, sqlalchemy_config


logger = logging.getLogger()


logger.setLevel(logging.DEBUG)

handler = logging.StreamHandler(sys.stderr)
""""
# Postgres Vector Store
In this notebook we are going to show how to use [Postgresql](https://www.postgresql.org) and  [pgvector](https://github.com/pgvector/pgvector)  to perform vector searches in LlamaIndex

If you're opening this Notebook on colab, you will probably need to install LlamaIndex 🦙.
"""

"""Running the following cell will install Postgres with PGVector in Colab."""


# import logging
# import sys

# Uncomment to see debug logs
logging.basicConfig(stream=sys.stdout, level=logging.DEBUG)
logging.getLogger().addHandler(logging.StreamHandler(stream=sys.stdout))

from llama_index.core import SimpleDirectoryReader, StorageContext
from llama_index.core import VectorStoreIndex
from llama_index.core import Settings
from llama_index.vector_stores.postgres import PGVectorStore
import textwrap
import openai
from llama_index.llms.groq import Groq

from llama_index.readers.database import DatabaseReader


GROQ_API_KEY = os.environ["GROQ_API_KEY"]
Settings.llm = Groq(
    model="llama3-70b-8192", request_timeout=360.0, api_key=GROQ_API_KEY
)

openai.api_key = os.environ["OPENAI_API_KEY"]
# TODO : Change embedding model to use not openai.
# Settings.embed_model = OllamaEmbedding(model_name="nomic-embed-text")
"""### Setup OpenAI
The first step is to configure the openai key. It will be used to created embeddings for the documents loaded into the index
"""

import os

# os.environ["OPENAI_API_KEY"] = "<your key>"

"""Download Data"""

# !mkdir -p 'data/paul_graham/'
# !wget 'https://raw.githubusercontent.com/run-llama/llama_index/main/docs/docs/examples/data/paul_graham/paul_graham_essay.txt' -O 'data/paul_graham/paul_graham_essay.txt'

"""### Loading documents
Load the documents stored in the `data/paul_graham/` using the SimpleDirectoryReader
"""

documents = SimpleDirectoryReader(datadir).load_data()
logger.info("Document ID:", documents[0].doc_id)

"""### Create the Database
Using an existing postgres running at localhost, create the database we'll be using.
"""

import psycopg2

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
vec_table_name = "vector_db"
file_table_name = "file"


conn = psycopg2.connect(connection_string)
conn.autocommit = True


"""### Hybrid Search

To enable hybrid search, you need to:
1. pass in `hybrid_search=True` when constructing the `PGVectorStore` (and optionally configure `text_search_config` with the desired language)
2. pass in `vector_store_query_mode="hybrid"` when constructing the query engine (this config is passed to the retriever under the hood). You can also optionally set the `sparse_top_k` to configure how many results we should obtain from sparse text search (default is using the same value as `similarity_top_k`).
"""

from sqlalchemy import make_url

url = make_url(connection_string)


reader = DatabaseReader(
    dbname=db_name,
    host=url.host,
    password=url.password,
    port=url.port,
    user=url.username,
)


from llama_index.core import Document

async def add_document_to_db_from_uuid(uuid_str: str) -> None:
    async def query_file_table_for_id(id: str) -> Tuple[any, any]:
        # Create an async engine and session
        engine = create_async_engine(async_postgres_connection_string, echo=True)
        async_session_maker = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

        async with async_session_maker() as session:
            async with FileModel.repo() as repo:
                # Create a query to select the first row matching the given id
                stmt = select(FileModel).where(FileModel.id == id)
                result = await session.execute(stmt)
                file_row = result.scalars().first()

                if file_row:
                    english_text = file_row.english_text
                    document_metadata = file_row.doc_metadata
                else:
                    english_text = None
                    document_metadata = None

                return (english_text, document_metadata)
    return_tuple = await query_file_table_for_id(id)
    english_text = return_tuple[0]
    doc_metadata = return_tuple[1]
    assert isinstance(english_text,str)
    assert isinstance(doc_metadata,dict)
    # TODO : Add support for metadata filtering 
    additional_document = Document(text = english_text, metadata = doc_metadata)
    additional_document.doc_id = str(uuid_str)
    hybrid_index.insert(additional_document)
    return None



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
    assert isinstance(english_text,str)
    assert isinstance(doc_metadata,dict)
    # TODO : Add support for metadata filtering 
    additional_document = Document(text = english_text, metadata = doc_metadata)
    additional_document.doc_id = str(hash) 
    # FIXME : Make sure the UUID matches the other function, and dryify this entire fucking mess.
    hybrid_index.insert(additional_document)
    return None

hybrid_response = hybrid_query_engine.query(
    "Who does Paul Graham think of with the word schtick"
)


"""#### Improving hybrid search with QueryFusionRetriever

Since the scores for text search and vector search are calculated differently, the nodes that were found only by text search will have a much lower score.

You can often improve hybrid search performance by using `QueryFusionRetriever`, which makes better use of the mutual information to rank the nodes.
"""

from llama_index.core.response_synthesizers import CompactAndRefine
from llama_index.core.retrievers import QueryFusionRetriever
from llama_index.core.query_engine import RetrieverQueryEngine

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

# response = query_engine.query(
#     "Who does Paul Graham think of with the word schtick, and why?"
# )
def create_rag_response_from_query( query : str ):
    return str(query_engine.query(query))


