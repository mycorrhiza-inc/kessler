from llama_index.llms.openai import OpenAI
from llama_index.core import PromptTemplate

from llama_index.core.retrievers import BaseRetriever
from llama_index.core.llms import LLM
from dataclasses import dataclass
from typing import Optional, List

import nest_asyncio
import asyncio

from models.chats import KeChatMessage
nest_asyncio.apply()

qa_prompt = PromptTemplate(
    """\
Context information is below.
---------------------
{context_str}
---------------------
Given the context information and not prior knowledge, answer the query.
Query: {query_str}
Answer: \
"""
)

query_str = (
    "Can you tell me about results from RLHF using both model-based and"
    " human-based evaluation?"
)

async def generate_query_from_chat_history(chat_history : List[KeChatMessage]) -> str:
    



def generate_response(retrieved_nodes, query_str, qa_prompt, llm):
    context_str = "\n\n".join([r.get_content() for r in retrieved_nodes])
    fmt_qa_prompt = qa_prompt.format(context_str=context_str, query_str=query_str)
    response = llm.complete(fmt_qa_prompt)
    return str(response), fmt_qa_prompt
