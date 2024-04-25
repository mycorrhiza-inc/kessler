from haystack_integrations.document_stores.chroma import ChromaDocumentStore
from haystack import Document

from haystack_integrations.document_stores.chroma import ChromaDocumentStore
from haystack_integrations.components.retrievers.chroma import ChromaEmbeddingRetriever

from uuid import UUID
from typing import Annotated, assert_type
import logging

from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    patch,
    MediaType,
)

from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body
from litestar.logging import LoggingConfig

from pydantic import TypeAdapter, BaseModel


from models import FileModel, FileRepository, FileSchema, provide_files_repo


from crawler.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor
from docprocessing.genextras import GenerateExtras

from typing import List, Optional, Union

document_store = ChromaDocumentStore()


# for testing purposese
emptyFile = FileModel(
    url="",
    name="",
    doctype="",
    lang="en",
    path="",
    # file=raw_tmpfile,
    doc_metadata={},
    stage="stage0",
    hash = "",
    summary = None,
    short_summary=None,
)

from typing import Any


class VecSearchQuery(BaseModel):
    message: str

class VecSearchResult(BaseModel):
    url : str


class FileController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @get(path="/files/{file_id:uuid}")
    async def get_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
    ) -> FileSchema:
        obj = files_repo.get(file_id)
        return FileSchema.model_validate(obj)

    @get(path="/test")
    async def get_file(self) -> None:
        return None


def AddDocument():
    pass

def GetEmbeddingsFromText():
    pass
    
def queryDocuments():
    pass

