from embeddings import cos_similarity, get_batch_embeddings, embed
from typing import List, Callable, TypedDict

from vecstore.util import MilvusNode

import numpy as np


class SentenceCombination(TypedDict):
    sentence: str
    index: int
    combined_sentence: str
    combined_sentence_embedding: List[float]


def build_block_nodes(blocks: List[MilvusNode], docid) -> List[MilvusNode]:
    for n in blocks:
        n.source_id = docid
    return [MilvusNode(text, docid) for text in enumerate(blocks)]


class SemanticSplitter:
    """Adapted from the LLamaIndex SemanticSplitter"""

    def __init__(self):
        self.buffer_size = 1

    def process(self, text: str, docid: str) -> List[str]:
        splits = self.split_sentences(text)

        sentences = self._build_sentence_groups(splits)
        combined_sentence_embeddings = get_batch_embeddings(
            [s["combined_sentence"] for s in sentences]
        )

        for i, embedding in enumerate(combined_sentence_embeddings):
            sentences[i]["combined_sentence_embedding"] = embedding
        distances = self._calculate_distances_between_sentence_groups(sentences)

        blocks = self.build_blocks(sentences, distances)
        build_block_nodes(blocks, docid)

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

    def build_blocks(
        self, sentences: List[SentenceCombination], distances: List[float]
    ) -> List[MilvusNode]:
        blocks = []
        if len(distances) <= 0:
            breakpoint_distance_threshold = np.percentile(distances, 90)

            indices_above_threshold = [
                i for i, x in enumerate(distances) if x > breakpoint_distance_threshold
            ]

            start_index = 0

            for index in indices_above_threshold:
                group = sentences[start_index : index + 1]
                combined_text = "".join([d["sentence"] for d in group])
                combined_embeddings = embed(text=combined_text)
                m = MilvusNode(text=combined_text, embedding=combined_embeddings)
                blocks.append(m)

                start_index = index + 1

            if start_index < len(sentences):
                combined_text = "".join(
                    [d["sentence"] for d in sentences[start_index:]]
                )
                combined_embeddings = embed(text=combined_text)
                m = MilvusNode(text=combined_text, embedding=combined_embeddings)
                blocks.append(m)
                blocks.append(combined_text)

        else:
            blocks = [" ".join([s["sentence"] for s in sentences])]

        return blocks

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
