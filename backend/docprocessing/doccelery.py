from docprocessing.genextras import GenerateExtras
from docprocessing.extractmarkdown import MarkdownExtractor
from docprocessing.docingest import DocumentIngester
from celery import Celery

from typing import Optional, List, Union, Path

app = Celery("tasks", broker="pyamqp://guest@localhost//")

from docprocessing.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor
from docprocessing.genextras import GenerateExtras

from pathlib import Path

import os

OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]


class DocumentProcessing:
    def __init__(
        self,
        gpu_endpoint_url: str = OS_GPU_COMPUTE_URL,
        tmpdir: Path = OS_TMPDIR,
    ):
        self.endpoint_url = gpu_endpoint_url
        self.tmpdir = tmpdir

    @app.task
    def ingest_document(self, input: Union[str, Path]):
        docingest = DocumentIngester(self.tmpdir)
        mdextract = MarkdownExtractor(self.endpoint_url, self.tmpdir)

        return input
