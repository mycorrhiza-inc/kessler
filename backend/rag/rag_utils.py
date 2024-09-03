import re
from typing import Any, Optional
import pytest
from rag.SemanticSplitter import split_by_max_tokensize
from rag.llamaindex import get_llm_from_model_str
from models.chats import KeChatMessage, ChatRole, sanitzie_chathistory_llamaindex

import asyncio


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
        if llm == "":
            llm = None
        if llm is None:
            llm = "llama-70b"
        if isinstance(llm, str):
            llm = get_llm_from_model_str(llm)

        self.llm = llm

    async def achat(self, chat_history: Any) -> Any:
        llama_chat_history = sanitzie_chathistory_llamaindex(chat_history)
        response = await self.llm.achat(llama_chat_history)
        str_response = str(response)

        def remove_prefixes(input_string: str) -> str:
            prefixes = ["assistant: "]
            for prefix in prefixes:
                if input_string.startswith(prefix):
                    input_string = input_string[
                        len(prefix) :
                    ]  # 10 is the length of "assistant: "
            return input_string

        str_response = remove_prefixes(str_response)
        return KeChatMessage(role=ChatRole.assistant, content=str_response)

    async def summarize_single_chunk(self, markdown_text: str) -> str:
        summarize_prompt = "Make sure to provide a well researched summary of the text provided by the user, if it appears to be the summary of a larger document, just summarize the section provided."
        summarize_message = KeChatMessage(
            role=ChatRole.assistant, content=summarize_prompt
        )
        text_message = KeChatMessage(role=ChatRole.user, content=markdown_text)
        summary = await self.achat(
            sanitzie_chathistory_llamaindex([summarize_message, text_message])
        )
        return summary.content

    async def summarize_mapreduce(
        self, markdown_text: str, max_tokensize: int = 8096
    ) -> str:
        splits = split_by_max_tokensize(markdown_text, max_tokensize)
        if len(splits) == 1:
            return await self.summarize_single_chunk(markdown_text)
        summaries = await asyncio.gather(
            *[self.summarize_single_chunk(chunk) for chunk in splits]
        )
        coherence_prompt = "Please rewrite the following list of summaries of chunks of the document into a final summary of similar length that incorperates all the details present in the chunks"
        cohere_message = KeChatMessage(ChatRole.assistant, coherence_prompt)
        combined_summaries_prompt = KeChatMessage(ChatRole.user, "\n".join(summaries))
        final_summary = await self.llm.achat(
            [cohere_message, combined_summaries_prompt]
        )
        return final_summary.content


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
