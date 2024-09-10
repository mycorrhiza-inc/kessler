from util.file_io import S3FileManager
from vecstore.docprocess import add_document_to_db
import os
from pathlib import Path
from uuid import UUID

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


from models.files import (
    FileModel,
    FileRepository,
    provide_files_repo,
    model_to_schema,
)
from common.file_schemas import FileSchema, DocumentStatus, docstatus_index


from typing import List, Optional, Dict, Annotated, Tuple, Any


import json

from common.niclib import rand_string, paginate_results

from enum import Enum

from sqlalchemy import and_

from logic.databaselogic import QueryData, filter_list_mdata, querydata_to_filters

from constants import (
    OS_TMPDIR,
)


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


# import base64


# import base64


class ThaumaturgyController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    # def jsonify_validate_return(self,):
    #     return None

    @post(path="/thaumaturgy/upsert_file", media_type=MediaType.TEXT)
    async def handle_file_upload(
        self,
        files_repo: FileRepository,
        data: Annotated[UploadFile, Body(media_type=RequestEncodingType.MULTI_PART)],
        request: Request,
        process: bool = True,
        override_hash: bool = False,
    ) -> str:
        logger = request.logger
        logger.info("Process initiated.")
        supplemental_metadata = {"source": "personal"}

        docingest = DocumentIngester(logger)
        input_directory = OS_TMPDIR / Path("formdata_uploads") / Path(rand_string())
        # Ensure the directories exist
        os.makedirs(input_directory, exist_ok=True)
        # Save the PDF to the output directory
        filename = data.filename
        final_filepath = input_directory / Path(filename)
        with open(final_filepath, "wb") as f:
            f.write(data.file.read())
        additional_metadata = docingest.infer_metadata_from_path(final_filepath)
        additional_metadata.update(supplemental_metadata)
        final_metadata = additional_metadata
        if final_metadata.get("lang") is None:
            final_metadata["lang"] = "en"

        file_obj = await add_file_raw(
            final_filepath, final_metadata, process, override_hash, files_repo, logger
        )
        return f"Successfully added document with uuid: {file_obj.id}"
