from embeddings import embed, batch_embed, cos_similarity
from typing import List, Callable


class SentenceCombination(TypedDict):
    sentence: str
    index: int
    combined_sentence: str
    combined_sentence_embedding: List[float]


class MilvusNode:
    def __init__(self, text: str, index: int):
        self.text = text
        self.index = index


class SemanticSplitter:
    """Adapted from the LLamaIndex SemanticSplitter"""

    def __init__(self, text: str):
        self.text = text
        self.sentences = []

    def process(self) -> List[str]:
        splits = self.split_sentences(self.text)

        sentences = self._build_sentence_groups(splits)

        combined_sentence_embeddings = batch_embed(
            [s["combined_sentence"] for s in sentences]
        )

        for i, embedding in enumerate(combined_sentence_embeddings):
            sentences[i]["combined_sentence_embedding"] = embedding

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

    def _split_sentences(self) -> List[str]:
        pass

    def calculate_similarity(self, sentence1: str, sentence2: str) -> float:
        a = embed(sentence1)
        b = embed(sentence2)

    def split_sentences(self, text: str) -> Callable[[str], List[str]]:
        import nltk

        tokenizer = nltk.tokenize.PunktSentenceTokenizer()
        return lambda text: self.split_by_sentence_tokenizer_ld(text, tokenizer)

    def split_by_sentence_tokenizer_ld(text: str, tokenizer) -> List[str]:
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
