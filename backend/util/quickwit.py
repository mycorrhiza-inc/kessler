import psycopg
import requests
import json
import time
from typing import List, Union


from constants import QUICKWIT_ENDPOINT
from models.util import postgres_connection_string

import logging

logger = logging.get_Logger(__name__)

db_string = postgres_connection_string
def print_response(r):
    logger.info(f"quickwit_api_call:\nstatus:{r.status_code}\nresponse:\n{r.text}")


def create_dockets_quickwit_index() -> None:
    request_data = {
        "version": "0.7",
        "index_id": "dockets",
        "doc_mapping": {
            "mode": "dynamic",
            "dynamic_mapping": {
                "indexed": True,
                "stored": True,
                "tokenizer": "default",
                "record": "basic",
                "expand_dots": True,
                "fast": True
            },
            "field_mappings": [
                {
                    "name": "text",
                    "type": "text",
                    "fast": True
                },
                {
                    "name": "timestamp",
                    "type": "datetime",
                    "input_formats": ["unix_timestamp"],
                    "fast_precision": "seconds",
                    "fast": True
                },
            ],
            "timestamp_field": "timestamp"
        },

        "search_settings": {
            "default_search_fields": [
                "text", "state", "city", "country"],
        },
        "indexing_settings": {
            "merge_policy": {
                "type": "limit_merge",
                "max_merge_ops": 3,
                "merge_factor": 10,
                "max_merge_factor": 12
            },
            "resources": {
                "max_merge_write_throughput": "80mb"
            }
        },
        "retention": {
            "period": "10 years",
            "schedule": "yearly"
        }
    }
    response = requests.post(
        f"{QUICKWIT_ENDPOINT}/api/v1/indexes",
        headers={"Content-Type": "application/json"},
        json=request_data
    )
    print_response(response)


def clear_index(index_name: str) -> None:
    r = requests.put(f"{QUICKWIT_ENDPOINT}/api/v1/indexes/{index_name}/clear")
    print_response(r)


def ingest_into_index(index_name: str, data: List[dict]) -> None:
    data_to_post = '\n'.join(
        [json.dumps(record, sort_keys=True, default=str) for record in data])
    r = requests.post(
        f"{QUICKWIT_ENDPOINT}/api/v1/{index_name}/ingest",
        headers={"Content-Type": "application/x-ndjson"},
        data=data_to_post
    )
    print_response(r)


def get_connection():
    return psycopg.connect(db_string)

def resolve_file_schema_for_docket_ingest(records: List[dict[str, any]]) -> List[dict[str, any]]:
        data = []
        for record in records:
            record['text'] = record.pop('english_text')
            record['source_id'] = str(record.pop('id'))
            record['metadata'] = json.loads(record.pop('mdata'))
            record['timestamp'] = time.time()
            data.append(record)
        return data

def preprocess_id_list_for_sql(ids: List[str]) -> str:
    if len(ids) == 1:
        return f"({ids[0]})"        
    else:
        outstr = " OR ".join(ids)
        return outstr
    
def ingest_files_to_dockets_index(file_ids: Union[str, List[str]]):
    # TODO: valdiate if the file has already been indexed
    conn = get_connection()
    with conn.cursor(row_factory=psycopg.rows.dict_row) as cur:
        ids = preprocess_id_list_for_sql(file_ids)
        cur.execute(f"SELECT * FROM file WHERE ID IN {ids}")
        records = cur.fetchone()

def reindex_whole_db(batch_size: int = 10):
    logger.info("reindexing docokets in quickwit")
    clear_index()
    conn = get_connection()
    with conn.cursor(row_factory=psycopg.rows.dict_row) as cur:
        cur.execute("SELECT * FROM file;")
        while True:
            records = cur.fetchmany(batch_size)
            if not records:
                break
            data_to_insert = resolve_file_schema_for_docket_ingest(records)
            ingest_into_index(data_to_insert)

    logger.info("quickwit indexing request complete")
