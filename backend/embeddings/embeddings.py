import os
from openai import OpenAI


OCTOAI_API_KEY = os.environ.get("OCTOAI_API_KEY")
OCTOAI_URL = "https://text.octoai.run/v1"


def embed(query: str, model="thenlper/gte-large"):
    client = OpenAI(
        api_key=OCTOAI_API_KEY,
        base_url=OCTOAI_URL
    )
    return client.embeddings.create(
        model=model,
        input=[
            query
        ],
    )
