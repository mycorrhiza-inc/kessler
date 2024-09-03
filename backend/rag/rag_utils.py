import re
from typing import Any, Optional
import pytest
from rag.llamaindex import get_llm_from_model_str
from models.chats import KeChatMessage, ChatRole, sanitzie_chathistory_llamaindex


def strip_links_and_tables(markdown_text):
    # Remove markdown links
    no_links = re.sub(r"\[.*?\]\(.*?\)", "", markdown_text)
    # Remove markdown tables
    no_tables = re.sub(r"\|.*?\|", "", no_links)
    return no_tables


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
        if llm is None:
            llm = "llama-70b"
        if isinstance(llm, str):
            llm = get_llm_from_model_str(llm)

        self.llm = llm

    async def summarize_single_chunk(self, markdown_text: str) -> str:
        summarize_prompt = f"Make sure to provide a well researched summary of the "
        summarize_message = KeChatMessage(
            role=ChatRole.assistant, content=summarize_prompt
        )
        text_message = KeChatMessage(role=ChatRole.user, content=markdown_text)
        summary = await self.llm.achat(
            sanitzie_chathistory_llamaindex([summarize_message, text_message])
        )
        return summary


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
