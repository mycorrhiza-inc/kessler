import re
from common.misc_schemas import QueryData

from common.file_schemas import DocumentStatus

import logging
from logic.filelogic import process_fileid_raw
import asyncio
from litestar.contrib.sqlalchemy.base import UUIDBase

from sqlalchemy.ext.asyncio import AsyncSession
import redis
from random import shuffle
from util.redis_utils import (
    clear_file_queue,
    convert_model_to_results_and_push,
    increment_doc_counter,
    pop_from_queue,
)
import traceback

from constants import (
    REDIS_BACKGROUND_DAEMON_TOGGLE,
    REDIS_BACKGROUND_DOCPROC_KEY,
    REDIS_HOST,
    REDIS_PORT,
    REDIS_CURRENTLY_PROCESSING_DOCS,
    REDIS_BACKGROUND_PROCESSING_STOPS_AT,
    REDIS_PRIORITY_DOCPROC_KEY,
)

redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)
default_logger = logging.getLogger(__name__)


async def background_processing_loop() -> None:
    await asyncio.sleep(
        30
    )  # Wait 30 seconds until application has finished loading to start processing background docs
    default_logger.info(
        "Starting the daemon that adds more background documents to process."
    )

    redis_client.set(REDIS_BACKGROUND_DAEMON_TOGGLE, 0)

    async def activity():
        if redis_client.get(REDIS_BACKGROUND_DAEMON_TOGGLE) == 0:
            await asyncio.sleep(300)
            return None
        else:
            await asyncio.sleep(300)

    # Logic to force it to process each loop sequentially
    result = None
    while result is None:
        try:
            result = await activity()
        except Exception as e:
            tb = traceback.format_exc()
            default_logger.error("Encountered error while processing a document")
            default_logger.error(e)
            default_logger.error(tb)
            await asyncio.sleep(2)
            result = None


# Returns a bool depending on if it was actually able to add numdocs
async def add_bulk_background_docs(numdocs: int, stop_at: str = "completed") -> bool:
    logger = default_logger
    engine = utils.sqlalchemy_config.get_engine()
    # Maybe Remove for better perf?
    async with engine.begin() as conn:
        await conn.run_sync(UUIDBase.metadata.create_all)
    session = AsyncSession(engine)
    files_repo = await provide_files_repo(session)
    try:
        data = QueryData(match_stage="unprocessed")
        filters = querydata_to_filters_strict(data)

        file_results = await files_repo.list(*filters)
        shuffle(file_results)
        return_boolean = len(file_results) >= numdocs
        truncated_results = file_results[:numdocs]
        convert_model_to_results_and_push(
            schemas=truncated_results, redis_client=redis_client
        )
    except Exception as e:
        await engine.dispose()
        await session.close()
        raise e
    await session.close()
    await engine.dispose()
    return return_boolean


async def main_processing_loop() -> None:
    await asyncio.sleep(
        10
    )  # Wait 10 seconds until application has finished loading to do anything
    max_concurrent_docs = 30
    redis_client.set(REDIS_CURRENTLY_PROCESSING_DOCS, 0)
    redis_client.set(REDIS_BACKGROUND_PROCESSING_STOPS_AT, "completed")
    # REMOVE FOR PERSIST QUEUES ACROSS RESTARTS:
    #
    clear_file_queue(redis_client=redis_client)
    default_logger.info("Starting the daemon processes docs in the queue.")

    async def activity():
        current_stop_at = redis_client.get(REDIS_BACKGROUND_PROCESSING_STOPS_AT)

        concurrent_docs = int(redis_client.get(REDIS_CURRENTLY_PROCESSING_DOCS))
        if concurrent_docs >= max_concurrent_docs:
            await asyncio.sleep(2)
            return None
        pull_docid = pop_from_queue(redis_client=redis_client)
        if pull_docid is None:
            await asyncio.sleep(2)
            return None
        increment_doc_counter(1, redis_client=redis_client)
        asyncio.create_task(
            process_document(doc_id_str=pull_docid, stop_at=current_stop_at)
        )
        return None

    # Logic to force it to process each loop sequentially
    result = None
    while result is None:
        try:
            result = await activity()
        except Exception as e:
            tb = traceback.format_exc()
            default_logger.error("Encountered error while processing a document")
            default_logger.error(e)
            default_logger.error(tb)
            await asyncio.sleep(2)
            result = None


def initialize_background_loops() -> None:
    asyncio.create_task(main_processing_loop())
    asyncio.create_task(background_processing_loop())


async def process_document(doc_id_str: str, stop_at: str) -> None:
    logger = default_logger
    logger.info(f"Executing background docproc on {doc_id_str} to {stop_at}")
    stop_at = DocumentStatus(stop_at)
    # TODO:: Replace passthrough files repo with actual global repo
    # engine = create_async_engine(
    #     postgres_connection_string,
    engine = utils.sqlalchemy_config.get_engine()
    # Maybe Remove for better perf?
    async with engine.begin() as conn:
        await conn.run_sync(UUIDBase.metadata.create_all)
    session = AsyncSession(engine)
    try:
        files_repo = await provide_files_repo(session)
        await process_fileid_raw(
            doc_id_str, files_repo, logger, stop_at, priority=False
        )
    except Exception as e:
        increment_doc_counter(-1, redis_client=redis_client)
        await engine.dispose()
        await session.close()
        raise e
    increment_doc_counter(-1, redis_client=redis_client)
    await session.close()
    await engine.dispose()
