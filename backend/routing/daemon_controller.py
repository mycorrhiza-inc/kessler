from typing_extensions import Doc
from lance_store.connection import ensure_fts_index
from rag.llamaindex import add_document_to_db_from_text
import os
from pathlib import Path
from typing import Any
from uuid import UUID
from typing import Annotated

from litestar import Controller, Request, Response

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    MediaType,
)




from sqlalchemy import select


from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


# from models import (
#     FileModel,
#     FileRepository,
#     FileSchema,
#     provide_files_repo,
# )
from models.files import (
    FileModel,
    FileRepository,
    FileSchema,
    FileSchemaWithText,
    provide_files_repo,
    DocumentStatus,
    docstatus_index
)

from logic.docingest import DocumentIngester
from logic.extractmarkdown import MarkdownExtractor

from typing import List, Optional, Dict


import json

from util.niclib import rand_string

from enum import Enum
class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


# TODO : Create test that adds a file once we know what the file DB schema is going to look like


class FileUpdate(BaseModel):
    message: str
    metadata: Dict[str, Any]


class UrlUpload(BaseModel):
    url: str
    metadata: Dict[str, Any] = {}


class UrlUploadList(BaseModel):
    url: List[str]


class FileCreate(BaseModel):
    message: str


class FileUpload(BaseModel):
    message: str


class IndexFileRequest(BaseModel):
    id: UUID






OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")


# import base64


OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")
OS_HASH_FILEDIR = OS_FILEDIR / Path("raw")
OS_OVERRIDE_FILEDIR = OS_FILEDIR / Path("override")
OS_BACKUP_FILEDIR = OS_FILEDIR / Path("backup")


# import base64





class DaemonController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    # def jsonify_validate_return(self,):
    #     return None

    @get(path="/daemon/start_docproc")
    async def get_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> FileSchema:
        obj = await files_repo.get(file_id)

        type_adapter = TypeAdapter(FileSchema)

        return type_adapter.validate_python(obj)

    @get(path="/files/markdown/{file_id:uuid}")
    async def get_markdown(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
        original_lang: bool = False
    ) -> str:
        # Yes I know this is a redundant if, this looks much more readable imo.
        if original_lang == True:
            return "Feature delayed due to only supporting english documents."
        obj = await files_repo.get(file_id)

        type_adapter = TypeAdapter(FileSchemaWithText)

        obj_with_text = type_adapter.validate_python(obj)

        markdown_text = obj_with_text.english_text
        if markdown_text is "":
            markdown_text = "Could not find Document Markdown Text"
        return markdown_text

    @get(path="/files/raw/{file_id:uuid}")
    async def get_raw(
        self,
        files_repo: FileRepository,
        request : Request,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> Response:
        logger = request.logger
        obj = await files_repo.get(file_id)
        if obj is None:
            return Response(content="ID does not exist", status_code=404)

        type_adapter = TypeAdapter(FileSchema)
        obj=type_adapter.validate_python(obj)
        filehash = obj.hash

        file_path = DocumentIngester(logger).get_default_filepath_from_hash(filehash)
        
        if not file_path.is_file():
            return Response(content="File not found", status_code=404)
        
        # Read the file content
        with open(file_path, 'rb') as file:
            file_content = file.read()
        # currently doesnt work unfortunately
        # file_name = obj.name
        # headers = {
        #     "Content-Disposition": f'attachment; filename="{file_name}"'
        # }

        return Response(content=file_content, media_type="application/octet-stream")

        # Return as a result of the get request, the file at file_path. Also make sure to include the correct return type.
