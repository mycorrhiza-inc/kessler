import chromadb
import os

local_chroma = os.environ["LOCAL_CHROMA"]
chroma_path = os.environ["CHROMA_PATH"]


def getChromaClient():
    if local_chroma:
        chroma_client = chromadb.PersistentClient(path=chroma_path)
        return chroma_client
    else:
        chroma_client = chromadb.HttpClient(host='localhost', port=8000)
        return chroma_client
