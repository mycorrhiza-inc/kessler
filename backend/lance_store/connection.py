import lancedb
from lancedb import DBConnection


def get_lance_connection() -> DBConnection:
    uri = "/tmp/lancedb"
    return lancedb.connect(uri=uri)


def ensure_fts_index():
    pass

    lanceconn = get_lance_connection()

    try:
        v = lanceconn.open_table("vectors")
        v.cleanup_old_versions()
        v.create_fts_index("text", replace=True)
    except FileNotFoundError as e:
        print(f"Encountered FileNotFoundError: {e}")
