from haystack import Pipeline
from haystack.components.converters import TextFileToDocument
from haystack.components.writers import DocumentWriter

from haystack_integrations.document_stores.chroma import ChromaDocumentStore
from haystack_integrations.components.retrievers.chroma import ChromaQueryTextRetriever

from chromadb.config import Settings
import chromadb
import os
from util.logging import get_logger

local_chroma = os.environ["LOCAL_APP"]
chroma_path = os.environ["CHROMA_PERSIST_PATH"]

logger = get_logger(__name__)

indexing = Pipeline()
indexing.add_component("converter", TextFileToDocument())
indexing.add_component("writer", DocumentWriter(document_store))
indexing.connect("converter", "writer")
indexing.run({"converter": {"sources": chroma_path}})


def getChromaClient():
    if local_chroma:
        chroma_client = chromadb.PersistentClient(path=chroma_path)
        return chroma_client
    else:
        # get the chroma instance from the local docker swarm

        try:
            client = chromadb.HttpClient(
                host=chroma_path,
                port=8000,
                settings=Settings(
                    chroma_client_auth_provider="chromadb.auth.token.TokenAuthClientProvider",
                    chroma_client_auth_credentials="test-token",
                ),
            )

            # this should work with or without authentication - it is a public endpoint
            if client.heartbeat():
                logger.info("connected to chroma")

            # this should work with or without authentication - it is a public endpoint
            chroma_version = client.get_version()
            logger.info(f"chroma version: {chroma_version}")

            # this is a protected endpoint and requires authentication
            try:
                collections = client.list_collections()
                logger.info(f"available collections:\n{collections}")
            except Exception as e:
                logger.fatal(f"Unable to authenticate to the chromadb instance")
                raise e
        except Exception as e:
            logger.fatal(f"Error connecting to remote chroma\n{e}")
            raise e

        return client


# Chroma is used in-memory so we use the same instances in the two pipelines below
document_store = ChromaDocumentStore()


querying = Pipeline()
querying.add_component("retriever", ChromaQueryTextRetriever(document_store))
results = querying.run({"retriever": {"query": "Variable declarations", "top_k": 3}})

for d in results["retriever"]["documents"]:
    print(d.meta, d.score)
