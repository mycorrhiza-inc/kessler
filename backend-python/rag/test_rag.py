import re
from typing import Any, Optional
import pytest
from rag.SemanticSplitter import split_by_max_tokensize
from rag.llamaindex import get_llm_from_model_str
from common.llm_utils import KeChatMessage, ChatRole, sanitzie_chathistory_llamaindex

import asyncio


@pytest.mark.parametrize(
    "markdown_text, expected_output",
    [
        (
            """
            Here is a [link](http://example.com) and another [link](http://example.org).
            
            | Header1 | Header2 |
            |---------|---------|
            | Row1Col1| Row1Col2|
            | Row2Col1| Row2Col2|
            
            Another [link](http://example.net).
            """,
            """
            Here is a  and another .
            
            
            
            Another .
            """,
        )
    ],
)
class LLMUtils:
    def __init__(self, llm: Optional[Any]):
        if llm == "":
            llm = None
        if llm is None:
            llm = "llama-70b"
        if isinstance(llm, str):
            llm = get_llm_from_model_str(llm)

        self.llm = llm


# Tests
def test_strip_links_and_tables():
    markdown_text = """
        Here is a [link](http://example.com) and another [link](http://example.org).
        
        | Header1 | Header2 |
        |---------|---------|
        | Row1Col1| Row1Col2|
        | Row2Col1| Row2Col2|
        
        Another [link](http://example.net).
        """
    expected_output = """
        Here is a  and another .
        
        
        
        Another .
        """
    assert (
        strip_links_and_tables(markdown_text) == expected_output
    ), f"Got {strip_links_and_tables(markdown_text)}, expected {expected_output}"


if __name__ == "__main__":
    pytest.main()
