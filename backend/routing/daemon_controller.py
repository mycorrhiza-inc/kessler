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
    FileModel,
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
from sqlalchemy.ext.asyncio import create_async_engine
from litestar.contrib.sqlalchemy.base import UUIDBase

# class UUIDEncoder(json.JSONEncoder):
#     def default(self, obj):
#         if isinstance(obj, UUID):
#             # if the obj is uuid, we simply return the value of uuid
#             return obj.hex
#         return json.JSONEncoder.default(self, obj)








OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")



OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")
OS_HASH_FILEDIR = OS_FILEDIR / Path("raw")
OS_OVERRIDE_FILEDIR = OS_FILEDIR / Path("override")
OS_BACKUP_FILEDIR = OS_FILEDIR / Path("backup")


default_logger = logging.getLogger(__name__)
logging.info("Daemon logging works, and started successfully")


# postgres_connection_string = os.environ["DATABASE_CONNECTION_STRING"]
# if "postgresql://" in postgres_connection_string:
#     postgres_connection_string = postgres_connection_string.replace(
#         "postgresql://", "postgresql+asyncpg://"
#     )
# engine = create_async_engine(
#         "postgresql+asyncpg://scott:tiger@localhost/test",
#         echo=True,
#     )

async def create_global_connection():
    global conn
    conn = utils.sqlalchemy_config.get_engine() 
    await conn.run_sync(UUIDBase.metadata.create_all)


# def jsonify_validate_return(self,):
#     return None
@listener("process_document")
async def process_document(doc_id_str: str, stop_at : str, files_repo : FileRepository) -> None:
    logger = default_logger
    # TODO:: Replace passthrough files repo with actual global repo
    # conn = utils.sqlalchemy_config.get_engine() 
    # await conn.begin().run_sync(UUIDBase.metadata.create_all)
    # # Maybe mirri is right and I should just use sql strings, but code reuse would make this hard
    # # Also this could potentially be a global, but I put it in here to reduce bugs
    # logger.info(f"Executing ackground process doc {doc_id_str} to stage {stop_at}")
    # # conn = utils.sqlalchemy_config.get_engine().begin() 
    # logger.info(str(conn))
    # logger.info(type(conn))
    # files_repo = await provide_files_repo(conn)
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
        files: List[FileModel],
        stop_at : Optional[str] = None,
        regenerate_from : Optional[str] = None
        ) -> None:
        logger = passthrough_request.logger
        if stop_at is None:
            stop_at = "completed"
        if regenerate_from is None:
            regenerate_from = "completed"
        stop_at = DocumentStatus(stop_at)
        regenerate_from= DocumentStatus(regenerate_from)
        for file in files:
            logger.info(f"Validating file {str(file.id)} for processing.")
            file_stage=DocumentStatus(file.stage)
            if docstatus_index(file_stage) > docstatus_index(regenerate_from):
                file_stage = regenerate_from
                file.stage = regenerate_from.value 
                await files_repo.update(file)
                # It was previously this and it was working, but syntax seems to suggest that only the above will work
                # _ = files_repo.update(file)

                logger.info(f"updated fileid {file.id} with new stage {file.stage}")
                await files_repo.session.commit()
            # Dont process the file if it is already processed beyond the stop point.
            if docstatus_index(file_stage) < docstatus_index(stop_at):
                logger.info(f"Sending file {str(file.id)} to be processed in the background.")
                passthrough_request.app.emit("process_document",doc_id_str=str(file.id), stop_at=stop_at.value, files_repo=files_repo)

    @get(path="/daemon/process_file/{file_id:uuid}")
    async def process_file_background(
        self,
        files_repo: FileRepository,
        request : Request,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
        stop_at : Optional[str] = None,
        regenerate_from : Optional[str] = None
        ) -> None:
            obj = await files_repo.get(file_id)
            type_adapter = TypeAdapter(FileSchema)
            validated_obj = type_adapter.validate_python(obj)
            return await self.bulk_process_file_background(files_repo = files_repo, passthrough_request=request, files = [validated_obj], stop_at = stop_at, regenerate_from=regenerate_from)
    @get(path="/daemon/process_all_files")
    async def process_all_background(
        self,
        files_repo: FileRepository,
        request : Request,
        stop_at : Optional[str] = None,
        regenerate_from : Optional[str] = None
        )-> None:
            logger = request.logger
            logger.info("Beginning to process all files.")
            results = await files_repo.list()
            logger.info(f"{len(results)} results")
            # type_adapter = TypeAdapter(list[FileSchema])
            # validated_results = type_adapter.validate_python(results)
            return await self.bulk_process_file_background(files_repo = files_repo, passthrough_request=request, files = results, stop_at = stop_at, regenerate_from=regenerate_from)



        




