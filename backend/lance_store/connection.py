import lancedb


def get_lance_connection():
    uri = "/tmp/kessler/lancedb/"
    return lancedb.connect(uri=uri)
