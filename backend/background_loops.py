from models.files import (
    FileSchema,
    FileModel,
    provide_files_repo,
    DocumentStatus,
    docstatus_index,
    FileRepository,
)


from typing import Optional, Dict
import logging
from models import utils
from logic.filelogic import process_fileid_raw
import asyncio
from litestar.contrib.sqlalchemy.base import UUIDBase

from sqlalchemy.ext.asyncio import AsyncSession
import redis

from constants import (
    REDIS_BACKGROUND_DAEMON_TOGGLE,
    REDIS_HOST,
    REDIS_PORT,
    REDIS_PRIORITY_DOCPROC_KEY,
    REDIS_BACKGROUND_DOCPROC_KEY,
    REDIS_CURRENTLY_PROCESSING_DOCS,
)

redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)
default_logger = logging.getLogger(__name__)


async def background_processing_loop() -> None:
    async def activity():
        if redis_client.get(REDIS_BACKGROUND_DAEMON_TOGGLE) == 0:
            await asyncio.sleep(300)
            return None
        

    # Logic to force it to process each loop sequentially
    result = None
    while result is None:
        result = await activity()


async def main_processing_loop() -> None:
    max_concurrent_docs = 30
    redis_client.set(REDIS_CURRENTLY_PROCESSING_DOCS, 0)

    async def activity():
        concurrent_docs = int(redis_client.get(REDIS_CURRENTLY_PROCESSING_DOCS))
        if concurrent_docs >= max_concurrent_docs:
            await asyncio.sleep(2)
            return None
        pull_docid = pop_from_queue()
        if pull_docid is None:
            await asyncio.sleep(2)
            return None
        pull_docinfo =  redis_client.hgetall(pull_docid)
        await process_document(pull_docinfo**)
        return None

        

    # Logic to force it to process each loop sequentially
    result = None
    while result is None:
        result = await activity()

async def initialize_background_loops() -> None:
    asyncio.create_task(main_processing_loop())
    asyncio.create_task(background_processing_loop())

# get type info for this
async def create_async_session() -> Any:
    engine = utils.sqlalchemy_config.get_engine()
    # Maybe Remove for better perf?
    async with engine.begin() as conn:
        await conn.run_sync(UUIDBase.metadata.create_all)
    session = AsyncSession(engine)
    return session


async def create_file_repository() -> FileRepository:
    session = await create_async_session()
    files_repo = await provide_files_repo(session)
    return files_repo


def update_status_in_redis(request_id: int, status: Dict[str, str]) -> None:
    redis_client.hmset(str(request_id), status)


def pop_from_queue() -> Optional[str]:
    # TODO : Clean up code logic
    request_id = redis_client.lpop(REDIS_PRIORITY_DOCPROC_KEY)
    if request_id is None:
        request_id = redis_client.lpop(REDIS_BACKGROUND_DOCPROC_KEY)
    if isinstance(request_id, str):
        return request_id
    default_logger.error(type(request_id))
    raise Exception(
        f"Request id is not string or none and is {type(request_id)} instead."
    )


async def process_document(doc_id_str: str, stop_at: str) -> None:
    logger = default_logger
    logger.info(f"Executing background docproc on {doc_id_str} to {stop_at}")
    stop_at = DocumentStatus(stop_at)
    # TODO:: Replace passthrough files repo with actual global repo
    # engine = create_async_engine(
    #     postgres_connection_string,
    with await create_async_session() as session:
        with await provide_files_repo(session) as files_repo:
            try:
                await process_fileid_raw(
                    doc_id_str, files_repo, logger, stop_at, priority=False
                )
            except Exception as e:
                session.close()
                raise e
        session.close()
