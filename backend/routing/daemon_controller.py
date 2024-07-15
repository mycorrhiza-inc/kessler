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
from litestar.events import listener




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
    FileSchema,
    provide_files_repo,
    DocumentStatus,
    docstatus_index,
    FileRepository
)

from logic.docingest import DocumentIngester
from logic.extractmarkdown import MarkdownExtractor

from typing import List, Optional, Dict


import json

from util.niclib import rand_string

from enum import Enum

import logging
from models import utils
from logic.filelogic import process_fileid_raw
import asyncio

class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)








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


logger = logging.getLogger(__name__)
logging.info("Daemon logging works, and started successfully")



conn = utils.sqlalchemy_config.get_engine().begin() 

# def jsonify_validate_return(self,):
#     return None
@listener("process_document")
async def process_document(doc_id_str: str, stop_at : str) -> None:
    # Maybe mirri is right and I should just use sql strings, but code reuse would make this hard
    # Also this could potentially be a global, but I put it in here to reduce bugs
    files_repo = await FileRepository(session=provide_files_repo(conn))
    stop_at = DocumentStatus(stop_at)
    await process_fileid_raw(doc_id_str,files_repo,logger,stop_at)






class DaemonController(Controller):

    dependencies = {"files_repo": Provide(provide_files_repo)}

    # def jsonify_validate_return(self,):
    #     return None
    #
    async def bulk_process_file_background(
        self,
        files_repo: FileRepository,
        passthrough_request : Request,
        files: List[FileSchema],
        stop_at : Optional[str] = None,
        regenerate_from : Optional[str] = None
        ) -> None:
        if stop_at is None:
            stop_at = "completed"
        if regenerate_from is None:
            regenerate_from = "completed"
        stop_at = DocumentStatus(stop_at)
        regenerate_from= DocumentStatus(regenerate_from)
        for file in files:
            file_stage=DocumentStatus(file.stage)
            if docstatus_index(file_stage) > docstatus_index(regenerate_from):
                file_stage = regenerate_from
                file.state = regenerate_from.value 
                files_repo.update(file)
                await files_repo.session.commit()
            # Dont process the file if it is already processed beyond the stop point.
            if docstatus_index(file_stage) < docstatus_index(stop_at):
                passthrough_request.logger.info(f"Sending file {str(file.id)} to be processed in the background.")
                passthrough_request.app.emit("process_document",doc_id_str=str(file.id), stop_at=stop_at.value)

    @get(path="/daemon/process_fileid/{file_id:uuid}")
    async def process_file_background(
        self,
        files_repo: FileRepository,
        request : Request,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
        stop_at : Optional[str] = None,
        regenerate_from : Optional[str] = None
        ):

            obj = await files_repo.get(file_id)

            type_adapter = TypeAdapter(FileSchema)

            validated_obj = type_adapter.validate_python(obj)
            return self.bulk_process_file_background(files_repo = files_repo, passthrough_request=request, files = [validated_obj], stop_at = stop_at, regenerate_from=regenerate_from)
    @get(path="/daemon/process_all")
    async def process_all_background(
        self,
        files_repo: FileRepository,
        request : Request,
        stop_at : Optional[str] = None,
        regenerate_from : Optional[str] = None
        ):
            results = await files_repo.list()
            logger = request.logger
            logger.info(f"{len(results)} results")
            type_adapter = TypeAdapter(list[FileSchema])
            validated_results = type_adapter.validate_python(results)

            return self.bulk_process_file_background(files_repo = files_repo, passthrough_request=request, files = validated_results, stop_at = stop_at, regenerate_from=regenerate_from)



        




