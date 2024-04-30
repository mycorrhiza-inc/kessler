from haystack_integrations.document_stores.chroma import ChromaDocumentStore
from haystack_integrations.components.retrievers.chroma import ChromaEmbeddingRetriever
from haystack_integrations.components.retrievers.chroma import ChromaQueryTextRetriever
from haystack.components.converters import MarkdownToDocument
from haystack.components.embedders import OpenAIDocumentEmbedder
from haystack.components.retrievers.in_memory import InMemoryEmbeddingRetriever
from haystack.components.writers import DocumentWriter
from haystack.utils import Secret, deserialize_secrets_inplace
from haystack import Document, Pipeline, PredefinedPipeline

from sqlalchemy import select
from models.files import FileModel

from typing import Optional, List
from pathlib import Path
from uuid import UUID
from logging import getLogger
import os
from models.files import provide_files_repo, FileModel, FileRepository

from models.utils import sqlalchemy_config

logger = getLogger(__name__)


octo_ef = OpenAIDocumentEmbedder(api_key=Secret.from_env_var(
    "OCTO_API_KEY"), model="thenlper/gte-large", api_base_url="https://text.octoai.run/v1/embeddings")

chroma_path = os.environ["CHROMA_PERSIST_PATH"]
chroma_store = ChromaDocumentStore(
    persist_path=chroma_path)
chroma_embedding_retriever = ChromaEmbeddingRetriever(
    document_store=chroma_store)
chroma_text_retriever = ChromaQueryTextRetriever(document_store=chroma_store)

chroma_pipeline = Pipeline()
chroma_pipeline.add_component("OpenAIDocumentEmbedder", octo_ef)
chroma_pipeline.add_component("writer", DocumentWriter(chroma_store))


async def indexDocByID(fid: UUID):
    # find file
    logger.info(f'INDEXER: indexing document with fid: {fid}')
    session_factory = sqlalchemy_config.create_session_maker()
    async with session_factory() as db_session:
        try:
            file_repo = FileRepository(session=db_session)
        except Exception as e:
            logger.error("unable to get file model repo", e)
        logger.info("created file repo")
        f = await file_repo.get(fid)
        logger.info(f'INDEXER: found file')
        # get
        docs = [
            Document(
                id=str(f.id),
                content=f.english_text,
                meta=f.doc_metadata
            )
        ]

        try:
            f.stage = "indexing"
            await file_repo.update(f, auto_commit=True)
            indexing = Pipeline()
            indexing.add_component("writer", DocumentWriter(chroma_store))
            logger.info("indexing document")
            indexing.run({"writer": {"documents": docs}})
            logger.info("completed indexing ")
        except Exception as e:
            logger.critical(f'Failed to index document with id {fid}', e)
            raise e


async def get_indexed_by_id(fid: UUID):
    searching = Pipeline()
    querying = Pipeline()
    querying.add_component("retriever", ChromaQueryTextRetriever(chroma_store))
    results = querying.run({
        "retriever": {
            "query": f'id: str(fid)', 
            "top_k": 3
        }
    })
    return results


async def indexDocByHash(hash: str):
    file_repo = FileModel.repo()
    stmt = select(FileModel).where(FileModel.hash == hash)
    f = file_repo.get(statement=stmt)


def query_chroma(query: str, top_k: int = 5):
    querying = Pipeline()
    querying.add_component("retriever", ChromaQueryTextRetriever(chroma_store))
    results = querying.run({
        "retriever": {
            "query": query, 
            "top_k": top_k
        }
    })
    return results
