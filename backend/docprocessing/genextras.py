import re
from util.niclib import rand_string, rand_filepath
import logging

# Note: Refactoring imports.py

import requests

from typing import Optional, List, Union


# from langchain.vectorstores import FAISS

import subprocess
import os


from pathlib import Path

from util.gpu_compute_calls import GPUComputeEndpoint

import re


from util.llm_prompts import LLM, token_split


class GenerateExtras:
    def __init__(self, logger, endpoint_url: str, tmpdir: Path):
        self.tmpdir = tmpdir
        self.endpoint_url = endpoint_url
        self.logger = logger
        self.small_llm = GPUComputeEndpoint(self.endpoint_url).llm_from_model_name(
            "small"
        )
        self.large_llm = GPUComputeEndpoint(self.endpoint_url).llm_from_model_name(
            "large"
        )
        # TODO : Add database connection.

    # Gets the links from the text of a document id:
    def extract_markdown_links(self, markdown_document_text: str) -> list[str]:
        # This regex pattern is designed to match typical markdown link structures
        markdown_url_pattern = r"!?\[.*?\]\((.*?)\)"
        # Find all non-overlapping matches in the markdown text
        urls = re.findall(markdown_url_pattern, markdown_document_text)
        return urls

    def summarize_document_text(
        self, document_text: str, max_chunk_size: int = 5000
    ) -> str:
        summary = self.small_llm.summarize_document_text(document_text, max_chunk_size)
        return summary

    def llm_postprocess_audio(self, raw_text: str, chunk_size=1000) -> str:
        final_text = self.small_llm.llm_postprocess_audio(raw_text, chunk_size)
        return final_text

    def gen_short_sum_from_long_sum(self, long_sum_text: str) -> str:
        short_summary = self.small_llm.gen_short_sum_from_long_sum(long_sum_text)
        return short_summary
