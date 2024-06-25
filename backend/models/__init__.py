from .chunks import Chunks
from .summary import Summaries

from lance_store.connection import get_lance_connection

lancedb_tables = {"text_chunks": Chunks, "doc_summaries": Summaries}


def ensure_lancedb_tables():
    # TODO: lancedb migrations
    # TODO: alter columns to keep the current data
    # in the correct rows
    for table_name, schema in lancedb_tables.items():
        db = get_lance_connection()
        _ = db.create_table(table_name=table_name, schema=schema, overwrite=True)
