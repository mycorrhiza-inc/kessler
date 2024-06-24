import requests
import json

from typing import Union, List

import numpy as np


from lancedb.embeddings.registry import EmbeddingFunctionRegistry
from lancedb.embeddings.base import TextEmbeddingFunction

registry = EmbeddingFunctionRegistry.get_instance()


@registry.register("uttu-embeddings")
class UttuEmbedFunc(TextEmbeddingFunction):
    def ndims(self) -> int:
        # we are using Alibaba-NLP/gte-large-en-v1.5
        return 1024

    def compute_query_embeddings(self, query: str, *args, **kwargs):
        return self.compute_source_embeddings(query, *args, **kwargs)

    def compute_source_embeddings(self, query: str, *args, **kwargs):
        # this is where do we do the request
        data = {
            "input": [query],
            # If we do have a 3090 we can run the much larger embedding model based of 7b transformer architectures that are currently at the top of MTEB leaderboard.
            "model": "Alibaba-NLP/gte-large-en-v1.5",
        }

        resp = requests.post("http://uttu-fedora:7997/embeddings", json=data)

        print(resp.text)

        body = json.loads(resp.text)
        data = body["data"][0]
        return data["embedding"]

    def generate_embeddings(
        self, texts: Union[List[str], np.ndarray], *args, **kwargs
    ) -> List[np.array]:
        data = {
            "input": texts,
            "model": "Alibaba-NLP/gte-large-en-v1.5",
        }

        resp = requests.post("http://uttu-fedora:7997/embeddings", json=data)

        print(resp.text)

        body = json.loads(resp.text)
        data = body["data"]
        return data.embedding


func = UttuEmbedFunc().create()
