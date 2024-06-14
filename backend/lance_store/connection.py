import lancedb
from lancedb import DBConnection


def get_lance_connection() -> DBConnection:
    uri = "/tmp/kessler/lancedb/"
    return lancedb.connect(uri=uri)
