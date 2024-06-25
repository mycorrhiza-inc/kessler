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
        # It should be possible to replace this entire try catch block w/ just this line. Doing this for now to avoid any uninentional changes in behavior if the table already exists.

    # if this errors then the FTS index is complete
    # Nic: Quick question, is this code irrelavent since all it does is raise the exception that would have been raised if the code was outside of a try block.
    # except Exception as e:
    #    raise e
