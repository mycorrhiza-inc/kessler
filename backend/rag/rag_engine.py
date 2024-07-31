from llama_index.llms.openai import OpenAI
from llama_index.core import PromptTemplate

from llama_index.core.retrievers import BaseRetriever
from llama_index.core.llms import LLM
from dataclasses import dataclass
from typing import Optional, List

import nest_asyncio
import asyncio

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


async def acombine_results(
    texts,
    query_str,
    qa_prompt,
    llm,
    cur_prompt_list,
    num_children=10,
):
    fmt_prompts = []
    for idx in range(0, len(texts), num_children):
        text_batch = texts[idx : idx + num_children]
        context_str = "\n\n".join([t for t in text_batch])
        fmt_qa_prompt = qa_prompt.format(context_str=context_str, query_str=query_str)
        fmt_prompts.append(fmt_qa_prompt)
        cur_prompt_list.append(fmt_qa_prompt)

    tasks = [llm.acomplete(p) for p in fmt_prompts]
    combined_responses = await asyncio.gather(*tasks)
    new_texts = [str(r) for r in combined_responses]

    if len(new_texts) == 1:
        return new_texts[0]
    else:
        return await acombine_results(
            new_texts, query_str, qa_prompt, llm, num_children=num_children
        )


async def agenerate_response_hs(
    retrieved_nodes, query_str, qa_prompt, llm, num_children=10
):
    """Generate a response using hierarchical summarization strategy.

    Combine num_children nodes hierarchically until we get one root node.

    """
    fmt_prompts = []
    node_responses = []
    for node in retrieved_nodes:
        context_str = node.get_content()
        fmt_qa_prompt = qa_prompt.format(context_str=context_str, query_str=query_str)
        fmt_prompts.append(fmt_qa_prompt)

    tasks = [llm.acomplete(p) for p in fmt_prompts]
    node_responses = await asyncio.gather(*tasks)

    response_txt = acombine_results(
        [str(r) for r in node_responses],
        query_str,
        qa_prompt,
        llm,
        fmt_prompts,
        num_children=num_children,
    )

    return response_txt, fmt_prompts


@dataclass
class Response:
    response: str
    source_nodes: Optional[List] = None

    def __str__(self):
        return self.response


class MyQueryEngine:
    """My query engine.

    Uses the tree summarize response synthesis module by default.

    """

    def __init__(
        self,
        retriever: BaseRetriever,
        qa_prompt: PromptTemplate,
        llm: LLM,
        num_children=10,
    ) -> None:
        self._retriever = retriever
        self._qa_prompt = qa_prompt
        self._llm = llm
        self._num_children = num_children

    def query(self, query_str: str):
        retrieved_nodes = self._retriever.retrieve(query_str)
        response_txt, _ = generate_response_hs(
            retrieved_nodes,
            query_str,
            self._qa_prompt,
            self._llm,
            num_children=self._num_children,
        )
        response = Response(response_txt, source_nodes=retrieved_nodes)
        return response

    async def aquery(self, query_str: str):
        retrieved_nodes = await self._retriever.aretrieve(query_str)
        response_txt, _ = await agenerate_response_hs(
            retrieved_nodes,
            query_str,
            self._qa_prompt,
            self._llm,
            num_children=self._num_children,
        )
        response = Response(response_txt, source_nodes=retrieved_nodes)
        return response
