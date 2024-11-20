# All this code is mirri's I just copied it into the main repo - nic
import logging
import requests
import psycopg
import json
import time
import os
from typing import List, Union
from datetime import datetime

QUICKWIT_ENDPOINT = "http://quickwit-main:7280"


logger = logging.getLogger(__name__)

db_string = os.environ["DATABASE_CONNECTION_STRING"]


def print_response(r):
    print(
        f"quickwit_api_call:\nstatus:{
            r.status_code}\nresponse:\n{r.text}"
    )


def create_dockets_quickwit_index(index_name: str) -> None:
    request_data = {
        "version": "0.7",
        "index_id": index_name,
        "doc_mapping": {
            "mode": "dynamic",
            "dynamic_mapping": {
                "indexed": True,
                "stored": True,
                "tokenizer": "default",
                "record": "basic",
                "expand_dots": True,
                "fast": True,
            },
            "field_mappings": [
                {"name": "text", "type": "text", "fast": True},
                {
                    "name": "timestamp",
                    "type": "datetime",
                    "input_formats": ["unix_timestamp"],
                    "fast_precision": "seconds",
                    "fast": True,
                },
                {"name": "date_filed", "type": "datetime", "fast": True},
            ],
            "timestamp_field": "timestamp",
        },
        "search_settings": {
            "default_search_fields": ["text", "state", "city", "country"],
        },
        "indexing_settings": {
            "merge_policy": {
                "type": "limit_merge",
                "max_merge_ops": 3,
                "merge_factor": 10,
                "max_merge_factor": 12,
            },
            "resources": {"max_merge_write_throughput": "80mb"},
        },
        "retention": {"period": "10 years", "schedule": "yearly"},
    }
    response = requests.post(
        f"{QUICKWIT_ENDPOINT}/api/v1/indexes",
        headers={"Content-Type": "application/json"},
        json=request_data,
    )
    print_response(response)


def clear_index(index_name: str) -> None:
    r = requests.put(f"{QUICKWIT_ENDPOINT}/api/v1/indexes/{index_name}/clear")
    print_response(r)


def ingest_into_index(index_name: str, data: List[dict]) -> None:
    data_to_post = "\n".join(
        [json.dumps(record, sort_keys=True, default=str) for record in data]
    )
    r = requests.post(
        f"{QUICKWIT_ENDPOINT}/api/v1/{index_name}/ingest",
        headers={"Content-Type": "application/x-ndjson"},
        data=data_to_post,
    )
    print("\n\n\n\n\n==============================")
    print("sent data to quickwit ")  # migrate_docket_to_nypuc()
    print("==============================\n\n\n\n\n")
    print_response(r)


def get_connection():
    return psycopg.connect(db_string)


def resolve_file_schema_for_docket_ingest(
    records: List[dict[str, any]]
) -> List[dict[str, any]]:
    data = []
    for record in records:
        try:
            record["text"] = record.pop("text")
            # record["text"] = ""
            record["source_id"] = str(record.pop("id"))
            record["metadata"] = record.pop("mdata")
            record["timestamp"] = time.time()
            # convert from M/D/Y to datetime rfc3339
            print(f"date filed:\n\n{record['metadata']['date']}\n\n")
            d = record["metadata"]["date"]
            d = d.split("/")
            record["date_filed"] = datetime(int(d[2]), int(d[0]), int(d[1])).strftime(
                "%Y-%m-%dT%H:%M:%SZ"
            )
            data.append(record)
        except Exception as e:
            print("error", e)
            print(record)
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


def reindex_whole_db(index: str, batch_size: int = 10):
    logger.info("reindexing docokets in quickwit")
    # clear_index(index)
    conn = get_connection()
    with conn.cursor(row_factory=psycopg.rows.dict_row) as cur:
        cur.execute(
            """
                    SELECT
                        public.file.*,
                        public.file_metadata.mdata,
                        public.file_text_source.text
                    FROM
                        public.file
                    LEFT JOIN
                        public.file_metadata
                    ON
                        public.file.id = public.file_metadata.id
                    LEFT JOIN
                        public.file_text_source
                    ON
                        public.file.id = public.file_text_source.file_id
                    """
        )
        while True:

            records = cur.fetchmany(batch_size)
            print(f"records:\n\n{records}\n\n")

            if not records:
                break
            data_to_insert = resolve_file_schema_for_docket_ingest(records)
            ingest_into_index(index_name="NY_PUC", data=data_to_insert)

    logger.info("quickwit indexing request complete")


def update_index_given_hash(index_name: str, hash: str, data: dict[str, any]):
    # delete the record with the given hash
    response = requests.post(
        f"{QUICKWIT_ENDPOINT}/api/v1/{index_name}/delete-tasks",
        headers={"Content-Type": "application/json"},
    )

    # insert the new record
    ingest_into_index(index_name, [data])


def migrate_docket_to_nypuc():
    query = "metadata.source:(ny-puc-energyefficiency-filedocs)"
    response = requests.post(
        f"{QUICKWIT_ENDPOINT}/api/v1/dockets/search",
        headers={"Content-Type": "application/json"},
        json={"query": query, "start_offset": 0, "max_hits": 10, "limit": 1000},
    )
    print_response(response)
    pass


if __name__ == "__main__":
    clear_index("NY_PUC")
    print("hitting quickwit api at: ", QUICKWIT_ENDPOINT)
    reindex_whole_db(index="NY_PUC")
