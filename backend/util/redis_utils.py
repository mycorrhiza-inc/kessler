from constants import (
    REDIS_BACKGROUND_DOCPROC_KEY,
    REDIS_CURRENTLY_PROCESSING_DOCS,
    REDIS_HOST,
    REDIS_PORT,
    REDIS_PRIORITY_DOCPROC_KEY,
)
from models.files import DocumentStatus, FileModel, FileRepository, docstatus_index
from typing import List, Tuple, Any, Union, Optional, Dict
import redis
import logging
import sys

# TODO : Mabye asycnify all the redis calls

default_redis_client = redis.Redis(
    host=REDIS_HOST, port=REDIS_PORT, decode_responses=True
)
default_logger = logging.getLogger(__name__)


def pop_from_queue(redis_client: Optional[Any] = None) -> Optional[str]:
    if redis_client is None:
        redis_client = default_redis_client
    # TODO : Clean up code logic
    request_id = redis_client.lpop(REDIS_PRIORITY_DOCPROC_KEY)
    if request_id is None:
        request_id = redis_client.lpop(REDIS_BACKGROUND_DOCPROC_KEY)
    if isinstance(request_id, str) or request_id is None:
        return request_id
    default_logger.error(type(request_id))
    raise Exception(
        f"Request id is not string or none and is {type(request_id)} instead."
    )


def update_status_in_redis(request_id: int, status: Dict[str, str]) -> None:
    redis_client.hmset(str(request_id), status)


def increment_doc_counter(
    increment: int,
    redis_client: Optional[Any] = None,
) -> None:
    if redis_client is None:
        redis_client = default_redis_client
    counter = redis_client.get(REDIS_CURRENTLY_PROCESSING_DOCS)
    default_logger.info(counter)
    redis_client.set(REDIS_CURRENTLY_PROCESSING_DOCS, int(counter) + increment)


def convert_model_to_results_and_push(
    schemas: Union[FileModel, List[FileModel]],
    redis_client: Optional[Any] = None,
) -> None:
    def convert_model_to_results(schemas: List[FileModel]) -> list:
        return_list = []
        for schema in schemas:
            str_id = str(schema.id)
            # Order doesnt matter since the list is shuffled anyway
            return_list.append(str_id)
        return return_list

    if redis_client is None:
        redis_client = default_redis_client
    if isinstance(schemas, FileModel):
        schemas = [schemas]
    id_list = convert_model_to_results(schemas)
    redis_client.rpush(REDIS_BACKGROUND_DOCPROC_KEY, *id_list)


def clear_file_queue(
    redis_client: Optional[Any] = None,
) -> None:
    if redis_client is None:
        redis_client = default_redis_client
    redis_client.ltrim(REDIS_PRIORITY_DOCPROC_KEY, 0, 0)
    redis_client.ltrim(REDIS_BACKGROUND_DOCPROC_KEY, 0, 0)


async def bulk_process_file_background(
    files_repo: FileRepository,
    files: List[FileModel],
    stop_at: DocumentStatus,
    max_documents: Optional[int] = None,
    logger: Optional[Any] = None,
    redis_client: Optional[Any] = None,
) -> bool:
    if redis_client is None:
        redis_client = default_redis_client
    if logger is None:
        logger = default_logger
    if max_documents is None:
        max_documents = 1000  # Default to prevent server from crashing by accidentially not including a value
    if files is None or len(files) == 0:
        logger.info("List of files to process was empty")
        return max_documents == 0

    def sanitize(x: Any) -> list:
        if x is None:
            return []
        return list(x)

    currently_processing_docs = sanitize(
        redis_client.lrange(REDIS_BACKGROUND_DOCPROC_KEY, 0, -1)
    ) + sanitize(redis_client.lrange(REDIS_PRIORITY_DOCPROC_KEY, 0, -1))

    def should_process(file: FileModel) -> bool:
        if not docstatus_index(file.stage) < docstatus_index(stop_at):
            return False
        # Set up a toggle for this at some point in time
        return file.id not in currently_processing_docs

    await files_repo.session.commit()
    files_to_convert = list(filter(should_process, files))[:max_documents]

    convert_model_to_results_and_push(schemas=files_to_convert)
    return len(files_to_convert) == max_documents
