import os
from pathlib import Path
from uuid import UUID

from litestar import Controller, Request, Response

from litestar.handlers.http_handlers.decorators import (
    post,
)
from litestar.events import listener




from litestar.params import Parameter
from litestar.di import Provide
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
    FileRepository,
)


from typing import List, Optional, Dict



import logging
from models import utils
from logic.filelogic import process_fileid_raw
import asyncio
from litestar.contrib.sqlalchemy.base import UUIDBase

from sqlalchemy.ext.asyncio import AsyncSession

from logic.databaselogic import QueryData, querydata_to_filters_strict , filters_docstatus_processing

from util.gpu_compute_calls import get_total_connections
import random



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


postgres_connection_string = os.environ["DATABASE_CONNECTION_STRING"]
if "postgresql://" in postgres_connection_string:
    postgres_connection_string = postgres_connection_string.replace(
        "postgresql://", "postgresql+asyncpg://"
    )
# engine = create_async_engine(
#         "postgresql+asyncpg://scott:tiger@localhost/test",
#         echo=True,
#     )
is_daemon_running = False

async def create_global_connection():
    global conn
    conn = utils.sqlalchemy_config.get_engine()
    await conn.run_sync(UUIDBase.metadata.create_all)


# def jsonify_validate_return(self,):
#     return None
@listener("process_document")
async def process_document(doc_id_str: str, stop_at: str) -> None:
    logger = default_logger
    logger.info(f"Executing background docproc on {doc_id_str} to {stop_at}")
    stop_at = DocumentStatus(stop_at)
    # TODO:: Replace passthrough files repo with actual global repo
    # engine = create_async_engine(
    #     postgres_connection_string,
    #     echo=True,
    #
    engine = utils.sqlalchemy_config.get_engine()
    logger.info(engine)
    logger.info(type(engine))
    async with engine.begin() as conn:
        await conn.run_sync(UUIDBase.metadata.create_all)
    session = AsyncSession(engine)
    files_repo_2 = await provide_files_repo(session)
    await process_fileid_raw(doc_id_str, files_repo_2, logger, stop_at,priority=False)
    await session.close()


class DaemonController(Controller):
    dependencies = {"files_repo": Provide(provide_files_repo)}

    # def jsonify_validate_return(self,):
    #     return None
    #
    async def bulk_process_file_background(
        self,
        files_repo: FileRepository,
        passthrough_request: Request,
        files: List[FileModel],
        stop_at: DocumentStatus,
        regenerate_from: DocumentStatus ,
        max_documents: Optional[int] = None,
    ) -> None:
        logger = passthrough_request.logger
        if max_documents is None:
            max_documents = -1
        if files is None:
            logger.info("List of files to process was empty")
            return None
        if len(files) == 0:
            logger.info("List of files to process was empty")
            return None
        for file in files:
            if max_documents == 0:
                logger.info("Reached maxiumum file limit, exiting.")
                return None
            logger.info(f"Validating file {str(file.id)} for processing.")
            file_stage = DocumentStatus(file.stage)
            if docstatus_index(file_stage) > docstatus_index(regenerate_from):
                file_stage = regenerate_from
                file.stage = regenerate_from.value
                await files_repo.update(file)
                await files_repo.session.commit()
                logger.info(
                    f"Reverting fileid {
                            file.id} to stage {file.stage}"
                )
            # Dont process the file if it is already processed beyond the stop point.
            if docstatus_index(file_stage) < docstatus_index(stop_at):
                logger.info(
                    f"Sending file {
                        str(file.id)} to be processed in the background."
                )
                # copy_files_repo = copy.deepcopy(files_repo)
                max_documents += -1
                passthrough_request.app.emit(
                    "process_document",
                    doc_id_str=str(file.id),
                    stop_at=stop_at.value,
                )

    @post(path="/daemon/process_file/{file_id:uuid}")
    async def process_file_background(
        self,
        files_repo: FileRepository,
        request: Request,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
        stop_at: Optional[str] = None,
        regenerate_from: Optional[str] = None,
    ) -> None:
        obj = await files_repo.get(file_id)
        if stop_at is None:
            stop_at = "completed"
        if regenerate_from is None:
            regenerate_from = "completed"
        stop_at = DocumentStatus(stop_at)
        regenerate_from = DocumentStatus(regenerate_from)
        return await self.bulk_process_file_background(
            files_repo=files_repo,
            passthrough_request=request,
            files=[obj],
            stop_at=stop_at,
            regenerate_from=regenerate_from,
        )
    async def process_query_background_raw(
        self,
        files_repo: FileRepository,
        passthrough_request: Request,
        data: QueryData,
        stop_at: Optional[str] = None,
        regenerate_from: Optional[str] = None,
        max_documents: Optional[int] = None,
        randomize: bool = False
    ) -> None:
        logger = passthrough_request.logger
        logger.info("Beginning to process all files.")
        if stop_at is None:
            stop_at = "completed"
        if regenerate_from is None:
            regenerate_from = "completed"
        stop_at = DocumentStatus(stop_at)
        regenerate_from = DocumentStatus(regenerate_from)
        filters = querydata_to_filters_strict(data) + filters_docstatus_processing(stop_at=stop_at,regenerate_from=regenerate_from)
        logger.info(filters)

        results = await files_repo.list(*filters)
        # type_adapter = TypeAdapter(list[FileSchema])
        # validated_results = type_adapter.validate_python(results)
        if randomize:
            random.shuffle(results)
        logger.info(f"{len(results)} results")
        return await self.bulk_process_file_background(
            files_repo=files_repo,
            passthrough_request=passthrough_request,
            files=results,
            stop_at=stop_at,
            regenerate_from=regenerate_from,
            max_documents=max_documents,
        )
    @post(path="/daemon/process_all_files")
    async def process_all_background(
        self,
        files_repo: FileRepository,
        request: Request,
        data: QueryData,
        stop_at: Optional[str] = None,
        regenerate_from: Optional[str] = None,
        max_documents: Optional[int] = None,
        randomize: bool = False
    ) -> None:
        return await self.process_query_background_raw( 
            files_repo=files_repo,
            passthrough_request=request,
            data=data,
            stop_at=stop_at,
            regenerate_from=regenerate_from,
            max_documents=max_documents,
            randomize=randomize 
        )

    # # TODO: Refactor so you dont have an open connection all the time.
    # async def background_daemon():
    #     global is_daemon_running
    #     is_daemon_running = True
    #     while True:
    #         try:
    #             await asyncio.sleep(30)
    #             if get_total_connections() < document_threshold:
    #                 await self.process_query_background_raw( 
    #                     files_repo=files_repo,
    #                     passthrough_request=request,
    #                     data=QueryData(),
    #                     stop_at=stop_at,
    #                     max_documents=documents_per_run,
    #                     randomize=True
    #                 )
    #         except Exception as e:
    #             is_daemon_running = False
    #             raise e

    @post(path="/dangerous/daemon/start_background_processing_daemon")
    async def start_background_processing_daemon(
        self,
        files_repo: FileRepository,
        request: Request,
        stop_at : Optional[str],
        documents_per_run : Optional[int],
        document_threshold : Optional[int]
            ) -> str:
        logger = request.logger
        if documents_per_run is None:
            documents_per_run = 10
        if document_threshold is None:
            document_threshold = 10
        global is_daemon_running
        if is_daemon_running:
            return "Daemon is already running, please restart to cease operation."
        while True:
            try:
                await asyncio.sleep(30)
                if get_total_connections() < document_threshold:
                    await self.process_query_background_raw( 
                        files_repo=files_repo,
                        passthrough_request=request,
                        data=QueryData(),
                        stop_at=stop_at,
                        max_documents=documents_per_run,
                        randomize=True
                    )
            except Exception as e:
                is_daemon_running = False
                raise e

        return "Code is in an unreachable state."


