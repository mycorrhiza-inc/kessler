import base64
import secrets
import os
import glob

import requests

from datetime import datetime, timezone

import hashlib


from pathlib import Path

import itertools


import math


from typing import Union, Optional, Any


def clean_and_empty_directory(dir: Path):
    files = glob.glob(str(dir / Path("*")))
    for f in files:
        os.remove(f)


# tests


from io import BufferedWriter


def get_hash_str(
    file_input: Union[Path, Any], hasher
) -> str:  # TODO: Figure out how df file types work
    if isinstance(file_input, Path):
        f = open(file_input, "rb")
        buf = f.read(65536)
        while len(buf) > 0:
            hasher().update(buf)
            buf = f.read(65536)
        return base64.urlsafe_b64encode(hasher().digest()).decode()
    if isinstance(file_input, BufferedWriter):
        with tempfile.NamedTemporaryFile() as temp:
            data = payload.get_payload(decode=True)
            temp.write(data)
            return base64.url_safe_b64encode(hasher(data).digest())
    return "ERROR Hashing File"


get_blake2 = lambda filepath: get_hash_str(filepath, hashlib.blake2b)
get_sha256 = lambda filepath: get_hash_str(filepath, hashlib.sha256)


def rand_string() -> str:
    return base64.urlsafe_b64encode(secrets.token_bytes(8)).decode()


def rand_filepath() -> Path:
    return Path(rand_string())


def secs_since_1970() -> int:
    return int(datetime.now(timezone.utc).timestamp())


def download_file(url: str, savedir: Path) -> Path:
    local_filename = savedir  # TODO: Use a temporary directory for downloads or archive it in some other way.
    # NOTE the stream=True parameter below
    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        with open(local_filename, "wb") as f:
            for chunk in r.iter_content(chunk_size=8192):
                # If you have chunk encoded response uncomment if
                # and set chunk_size parameter to None.
                # if chunk:
                f.write(chunk)
    return local_filename


def file2string(path: Path) -> str:
    with open(path, "r") as file:
        data = file.read().rstrip()
    return data


def concatlist(list_of_lists: list) -> list:
    return list(itertools.chain.from_iterable(list_of_lists))
