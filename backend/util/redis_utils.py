from constants import (
    REDIS_BACKGROUND_DOCPROC_KEY,
    REDIS_CURRENTLY_PROCESSING_DOCS,
    REDIS_HOST,
    REDIS_PORT,
)
from models.files import FileModel
from typing import List, Tuple, Any, Union, Optional
import redis

# TODO : Mabye asycnify all the redis calls

default_redis_client = redis.Redis(
    host=REDIS_HOST, port=REDIS_PORT, decode_responses=True
)


def increment_doc_counter(
    increment: int,
    redis_client: Optional[Any] = None,
) -> None:
    if redis_client is None:
        redis_client = default_redis_client
    counter = redis_client.get(REDIS_CURRENTLY_PROCESSING_DOCS)
    redis_client.set(REDIS_CURRENTLY_PROCESSING_DOCS, counter + increment)


def convert_model_to_results(
    schemas: List[FileModel], stop_at: str
) -> Tuple[dict, list]:
    return_dict = {}
    return_list = []
    for schema in schemas:
        str_id = str(schema.id)
        return_dict[str_id] = {"doc_id_str": str_id, "stop_at": stop_at}
        # Order doesnt matter since the list is shuffled anyway
        return_list.append(str_id)
    return (return_dict, return_list)


def convert_model_to_results_and_push(
    schemas: Union[FileModel, List[FileModel]],
    stop_at: str,
    redis_client: Optional[Any] = None,
) -> None:
    if redis_client is None:
        redis_client = default_redis_client
    if isinstance(schemas, FileModel):
        schemas = [schemas]
    data_dictionary, id_list = convert_model_to_results(schemas, stop_at)
    redis_client.mset(data_dictionary)
    redis_client.rpush(REDIS_BACKGROUND_DOCPROC_KEY, id_list)
