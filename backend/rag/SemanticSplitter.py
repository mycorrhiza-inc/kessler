from embeddings import cos_similarity, embed
from typing import List, Callable, TypedDict

from vecstore.util import MilvusNode, MilvusRow

import numpy as np

import logging

logger = logging.getLogger(__name__)


class SentenceCombination(TypedDict):
    sentence: str
    index: int
    combined_sentence: str
    combined_sentence_embedding: List[float]


class SemanticSplitter:
    """Adapted from the LLamaIndex SemanticSplitter"""

    def __init__(self):
        self.buffer_size = 1

    def process(self, text: str, source_id: str) -> List[MilvusRow]:
        splits = self.split_sentences(text)

        sentences = self._build_sentence_groups(splits)
        combined_sentence_embeddings = embed(
            [s["combined_sentence"] for s in sentences]
        )

        for i, embedding in enumerate(combined_sentence_embeddings):
            sentences[i]["combined_sentence_embedding"] = embedding
        distances = self._calculate_distances_between_sentence_groups(sentences)

        chunks = self.build_chunks(sentences, distances)
        blocks = self.build_blocks_from_chunks(chunks, source_id=source_id)

        return blocks

    def _calculate_distances_between_sentence_groups(
        self, sentences: List[SentenceCombination]
    ) -> List[float]:
        distances = []
        for i in range(len(sentences) - 1):
            embedding_current = sentences[i]["combined_sentence_embedding"]
            embedding_next = sentences[i + 1]["combined_sentence_embedding"]

            similarity = cos_similarity(embedding_current, embedding_next)

            distance = 1 - similarity

            distances.append(distance)

        return distances

    def build_blocks_from_chunks(
        self, chunks: List[str], source_id: str
    ) -> List[MilvusRow]:
        blocks = []
        for i, chunk in enumerate(chunks):
            blocks.append(
                MilvusRow(
                    text=chunk,
                    source_id=str(source_id),
                    embedding=embed(chunk)[0],
                )
            )
        return blocks

    def build_chunks(
        self, sentences: List[SentenceCombination], distances: List[float]
    ) -> List[str]:
        chunks = []
        if len(distances) <= 0:
            breakpoint_distance_threshold = np.percentile(distances, 90)

            indices_above_threshold = [
                i for i, x in enumerate(distances) if x > breakpoint_distance_threshold
            ]

            start_index = 0

            # combine sentences into blocks if they are abouve threshold
            for index in indices_above_threshold:
                group = sentences[start_index : index + 1]
                combined_text = "".join([d["sentence"] for d in group])
                chunks.append(combined_text)

                start_index = index + 1

            if start_index < len(sentences):
                combined_text = "".join(
                    [d["sentence"] for d in sentences[start_index:]]
                )
                chunks.append(combined_text)

        else:
            combined_text = " ".join([s["sentence"] for s in sentences])
            chunks.append(combined_text)

        return chunks

    def _build_sentence_groups(
        self, text_splits: List[str]
    ) -> List[SentenceCombination]:
        sentences: List[SentenceCombination] = [
            {
                "sentence": x,
                "index": i,
                "combined_sentence": "",
                "combined_sentence_embedding": [],
            }
            for i, x in enumerate(text_splits)
        ]

        # Group sentences and calculate embeddings for sentence groups
        for i in range(len(sentences)):
            combined_sentence = ""

            for j in range(i - self.buffer_size, i):
                if j >= 0:
                    combined_sentence += sentences[j]["sentence"]

            combined_sentence += sentences[i]["sentence"]

            for j in range(i + 1, i + 1 + self.buffer_size):
                if j < len(sentences):
                    combined_sentence += sentences[j]["sentence"]

            sentences[i]["combined_sentence"] = combined_sentence

        return sentences

    def _build_chunks(self) -> List[str]:
        pass

    def split_sentences(self, text: str) -> List[str]:
        import nltk

        tokenizer = nltk.tokenize.PunktSentenceTokenizer()
        return self.split_by_sentence_tokenizer(text, tokenizer)

    def split_by_sentence_tokenizer(self, text: str, tokenizer) -> List[str]:
        """
        Get the spans and then return the sentences.

        Using the start index of each span
        Instead of using end, use the start of the next span if available
        """
        spans = list(tokenizer.span_tokenize(text))
        sentences = []
        for i, span in enumerate(spans):
            start = span[0]
            if i < len(spans) - 1:
                end = spans[i + 1][0]
            else:
                end = len(text)
            sentences.append(text[start:end])
        return sentences

    def _split_by_tokenizer(self, text: str) -> Callable[[str], List[str]]:
        import nltk

        tokenizer = nltk.tokenize.PunktSentenceTokenizer()

        spans = list(tokenizer.span_tokenize(text))
        sentences = []
        for i, span in enumerate(spans):
            start = span[0]
            if i < len(spans) - 1:
                end = spans[i + 1][0]
            else:
                end = len(text)
            sentences.append(text[start:end])
        return sentences
