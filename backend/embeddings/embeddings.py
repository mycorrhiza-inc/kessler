import os
from typing import List
from openai import OpenAI

from numpy.linalg import norm
from numpy import dot


OCTOAI_API_KEY = os.environ.get("OCTOAI_API_KEY")
OCTOAI_URL = "https://text.octoai.run/v1"


def embed(query: str, model="thenlper/gte-large") -> List[float]:
    client = OpenAI(api_key=OCTOAI_API_KEY, base_url=OCTOAI_URL)
    try:
        response = client.embeddings.create(
            model=model,
            input=[query],
        )
        return response.data[0].embedding
    except Exception as e:
        print(e)
        return []


def get_batch_embeddings(
    queries: List[str], model="thenlper/gte-large"
) -> List[List[float]]:
    client = OpenAI(api_key=OCTOAI_API_KEY, base_url=OCTOAI_URL)

    responses = client.embeddings.create(
        model=model,
        input=queries,
    )
    embeddings = []
    for r in responses.data:
        embeddings.append(r.embedding)

    return embeddings


def cos_similarity(
    a: List[float],
    b: List[float],
) -> float:
    cos_sim = dot(a, b) / (norm(a) * norm(b))
    return cos_sim
