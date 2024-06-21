import lancedb
from lancedb import DBConnection


def get_lance_connection() -> DBConnection:
    uri = "/tmp/lancedb"
    return lancedb.connect(uri=uri)


def ensure_fts_index():
    lanceconn = get_lance_connection()

    v = lanceconn.open_table("vectors")

    try:
        v.cleanup_old_versions()
        v.create_fts_index("text", replace=True)

    # if this errors then the FTS index is complete
    except Exception as e:
        raise e
        pass
