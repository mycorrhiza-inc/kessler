import os
from typing import Union, List

from typing import List
from openai import OpenAI

from numpy.linalg import norm
from numpy import dot

from constants import FIREWORKS_API_KEY, FIREWORKS_EMBEDDING_URL


def embed(
    query: Union[str, List[str]], model="nomic-ai/nomic-embed-text-v1.5"
) -> List[float]:
    if isinstance(query, str):
        query = [query]

    client = OpenAI(api_key=FIREWORKS_API_KEY, base_url=FIREWORKS_EMBEDDING_URL)

    try:
        response = client.embeddings.create(model=model, input=query)
        embeddings = [r.embedding for r in response.data]
        return embeddings

    except Exception as e:
        print(e)
        return []


def cos_similarity(
    a: List[float],
    b: List[float],
) -> float:
    cos_sim = dot(a, b) / (norm(a) * norm(b))
    return cos_sim
