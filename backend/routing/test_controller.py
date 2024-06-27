import os
from pathlib import Path

from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    MediaType,
)


import json


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


# TODO : Create test that adds a file once we know what the file DB schema is going to look like


# litestar only


OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")


# import base64


OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")
OS_HASH_FILEDIR = OS_FILEDIR / Path("raw")
OS_OVERRIDE_FILEDIR = OS_FILEDIR / Path("override")
OS_BACKUP_FILEDIR = OS_FILEDIR / Path("backup")


# import base64


class FileController(Controller):
    """File Controller"""

    @get(path="/test/run_suite")
    async def run_test_suite(self) -> str:
        return "Test complete"
