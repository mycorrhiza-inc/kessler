from pydantic import BaseModel
from typing import Optional
from uuid import UUID

from uuid import UUID

from litestar import Controller, Request, Response

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    MediaType,
)


from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body


from common.file_schemas import (
    FileSchema,
    DocumentStatus,
    FileTextSchema,
    docstatus_index,
)


from typing import List, Optional, Dict, Annotated, Tuple, Any


import json

from common.niclib import rand_string, paginate_results

from enum import Enum


from constants import (
    OS_TMPDIR,
)

from common.misc_schemas import QueryData

# REDIS_DOCPROC_PRIORITYQUEUE_KEY = "docproc_queue_priority"
# REDIS_DOCPROC_QUEUE_KEY = "docproc_queue_background"
#
# REDIS_DOCPROC_INFORMATION = "docproc_information"
#
# REDIS_DOCPROC_BACKGROUND_DAEMON_TOGGLE = "docproc_background_daemon"
# REDIS_DOCPROC_BACKGROUND_PROCESSING_STOPS_AT = "docproc_background_stop_at"
# REDIS_DOCPROC_CURRENTLY_PROCESSING_DOCS = "docproc_currently_processing_docs"
from constants import (
    REDIS_DOCPROC_PRIORITYQUEUE_KEY,
    REDIS_DOCPROC_QUEUE_KEY,
    REDIS_DOCPROC_INFORMATION,
    REDIS_DOCPROC_BACKGROUND_DAEMON_TOGGLE,
    REDIS_DOCPROC_BACKGROUND_PROCESSING_STOPS_AT,
    REDIS_DOCPROC_CURRENTLY_PROCESSING_DOCS,
    REDIS_HOST,
    REDIS_PORT,
)

from background_loops import clear_file_queue

import redis

redis_client = redis.Redis(REDIS_HOST, port=REDIS_PORT)


def push_to_queue(request: str, priority: bool):
    if priority:
        pushkey = REDIS_DOCPROC_QUEUE_KEY
    else:
        pushkey = REDIS_DOCPROC_PRIORITYQUEUE_KEY

    redis_client.rpush(pushkey, request)


class DaemonState(BaseModel):
    enable_background_processing: Optional[bool] = None
    stop_at_background_docprocessing: Optional[str] = None
    clear_queue: Optional[bool] = None


class DaemonController(Controller):
    @post(path="/dangerous/daemon/control_background_processing_daemon")
    async def control_background_processing_daemon(self, data: DaemonState) -> str:
        daemon_toggle = data.enable_background_processing
        stop_at = data.stop_at_background_docprocessing
        clear_queue = data.clear_queue
        if daemon_toggle is not None:
            redis_client.set(REDIS_DOCPROC_BACKGROUND_DAEMON_TOGGLE, int(daemon_toggle))
        if stop_at is not None:
            target = DocumentStatus(stop_at).value
            redis_client.set(REDIS_DOCPROC_BACKGROUND_PROCESSING_STOPS_AT, target)
        if clear_queue is not None:
            if clear_queue:
                clear_file_queue()
        return "Success!"

    # async def process_force_downgrade_raw(
    #     self,
    #     files_repo: FileRepository,
    #     data: QueryData,
    #     logger: Optional[Any],
    #     regenerate_from: Optional[str] = None,
    # ) -> None:
    #     if logger is None:
    #         logger = default_logger
    #     logger.info("Beginning to process all files.")
    #     if regenerate_from is None:
    #         regenerate_from = "completed"
    #     regenerate_from = DocumentStatus(regenerate_from)
    #     filters = querydata_to_filters_strict(data)
    #     logger.info(filters)
    #     results = await files_repo.list(*filters)
    #
    #     for file in results:
    #         file_stage = DocumentStatus(file.stage)
    #         if docstatus_index(file_stage) > docstatus_index(regenerate_from):
    #             file_stage = regenerate_from
    #             file.stage = regenerate_from.value
    #             await files_repo.update(file)
    #             logger.info(
    #                 f"Reverting fileid {
    #                         file.id} to stage {file.stage}"
    #             )
    #
    #     await files_repo.session.commit()
    #
    # @post(path="/daemon/force_downgrade")
    # async def force_downgrade(
    #     self,
    #     files_repo: FileRepository,
    #     request: Request,
    #     data: QueryData,
    #     regenerate_from: Optional[str] = None,
    # ) -> None:
    #     return await self.process_force_downgrade_raw(
    #         files_repo=files_repo,
    #         data=data,
    #         logger=request.logger,
    #         regenerate_from=regenerate_from,
    #     )
    #
    # @post(path="/daemon/process_file/{file_id:uuid}")
    # async def process_file_background(
    #     self,
    #     files_repo: FileRepository,
    #     request: Request,
    #     file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    #     stop_at: Optional[str] = None,
    # ) -> None:
    #     obj = await files_repo.get(file_id)
    #     if stop_at is None:
    #         stop_at = "completed"
    #     stop_at = DocumentStatus(stop_at)
    #     await bulk_process_file_background(
    #         files_repo=files_repo,
    #         logger=request.logger,
    #         files=[obj],
    #         stop_at=stop_at,
    #     )
    #
    # @post(path="/daemon/process_all_files")
    # async def process_all_background(
    #     self,
    #     files_repo: FileRepository,
    #     request: Request,
    #     data: QueryData,
    #     stop_at: Optional[str] = None,
    #     regenerate_from: Optional[str] = None,
    #     max_documents: Optional[int] = None,
    #     randomize: bool = False,
    # ) -> None:
    #     return await self.process_query_background_raw(
    #         files_repo=files_repo,
    #         data=data,
    #         stop_at=stop_at,
    #         regenerate_from=regenerate_from,
    #         max_documents=max_documents,
    #         randomize=randomize,
    #         logger=request.logger,
    #     )
    #
    # async def process_query_background_raw(
    #     self,
    #     files_repo: FileRepository,
    #     data: QueryData,
    #     stop_at: Optional[str] = None,
    #     regenerate_from: Optional[str] = None,
    #     max_documents: Optional[int] = None,
    #     randomize: bool = False,
    #     logger: Any = default_logger,
    # ) -> None:
    #     logger.info("Beginning to process all files.")
    #     if stop_at is None:
    #         stop_at = "completed"
    #     stop_at = DocumentStatus(stop_at)
    #     filters = querydata_to_filters_strict(data) + filters_docstatus_processing(
    #         stop_at=stop_at
    #     )
    #     logger.info(filters)
    #
    #     results = await files_repo.list(*filters)
    #     # type_adapter = TypeAdapter(list[FileSchema])
    #     # validated_results = model_to_schema(results)
    #     if randomize:
    #         random.shuffle(results)
    #     logger.info(f"{len(results)} results")
    #     await bulk_process_file_background(
    #         files_repo=files_repo,
    #         files=results,
    #         stop_at=stop_at,
    #         max_documents=max_documents,
    #         logger=logger,
    #     )
