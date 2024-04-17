import chromadb
import os

local_chroma = os.environ["LOCAL_APP"]
chroma_path = os.environ["CHROMA_PERSIST_PATH"]


def getChromaClient():
    if local_chroma:
        chroma_client = chromadb.PersistentClient(path=chroma_path)
        return chroma_client
    else:
        chroma_client = chromadb.HttpClient(host=os.environ["CHROMA_URL"], port=os.environ["CHROMA_PORT"])
        return chroma_client
