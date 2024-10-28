from io import BufferedWriter
import base64
import secrets
import os
import glob
import requests
import hashlib
import itertools
import math
import yaml
import re

from datetime import datetime, timezone

from pathlib import Path

from typing import Union, Optional, Any, Tuple
from typing import Callable

import logging

import asyncio

import tokenizers
import math

default_logger = logging.getLogger(__name__)


def paginate_results(
    results: list, num_results: Optional[int], page: Optional[int]
) -> Tuple[list, int]:
    if num_results is None:
        return (results, 1)
    if page is None or page < 1:
        page = 1
    rectify = lambda x: max(0, min(x, len(results)))
    avalible_pages = math.ceil(len(results) / num_results)
    default_logger.info(f"Length of Results: {len(results)}")
    if page > avalible_pages:
        page = avalible_pages
    start_page_index = rectify((page - 1) * num_results)
    end_page_index = rectify((page) * num_results)
    default_logger.info(f"Start Pagination index: {start_page_index}")
    default_logger.info(f"End Pagination index: {end_page_index}")
    return (results[start_page_index:end_page_index], avalible_pages)


def Maybe(func: Callable) -> Callable:
    return lambda x: (None if x is None else func(x))


def fizbuzz(maxiters: int) -> str:
    loops = math.ceil(maxiters / 15)
    return_str = ""
    n = 0
    for i in range(0, loops):
        return_str = (
            return_str
            + f"FizzBuzz\n{n+1}\n{n+2}\nFizz\n{n+4}\nBuzz\nFizz\n{n+7}\n{n+8}\nFizz\nBuzz\n{n+11}\nFizz\n{n+13}\n{n+14}\n"
        )
        n += 15
    return return_str


def clean_and_empty_directory(dir: Path):
    files = glob.glob(str(dir / Path("*")))
    for f in files:
        os.remove(f)


def create_markdown_string(
    text: str, metadata: dict, include_previous_metadata: bool = True
) -> str:
    text_without_metadata, extra_metadata = seperate_markdown_string(text)
    if include_previous_metadata and (extra_metadata != {}):
        metadata = extra_metadata.update(metadata)
    yaml_metadata = yaml.safe_dump(metadata, default_flow_style=False)

    # Construct the markdown string with front matter
    markdown_with_frontmatter = f"---\n{yaml_metadata}\n---\n{text_without_metadata}"
    return markdown_with_frontmatter


def seperate_markdown_string(mdstr_with_metadata: str) -> Tuple[str, dict]:
    # Define regex pattern for front matter
    frontmatter_pattern = re.compile(r"^---\s*\n(.*?)\s*\n---\s*\n", re.DOTALL)

    match = frontmatter_pattern.match(mdstr_with_metadata)

    if match:
        # Extract metadata
        frontmatter = match.group(1)
        # Remove front matter from markdown text to get main body
        main_body = mdstr_with_metadata[match.end() :]
        #
        try:
            # Parse the YAML content
            yaml_dict = yaml.safe_load(frontmatter)
            return (main_body, yaml_dict)
        except yaml.YAMLError as e:
            print(f"Error parsing YAML: {e}")
            return (main_body, {})
    else:
        # If no front matter, return markdown text and empty dictionary
        return (mdstr_with_metadata, {})


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
    # TODO: Use a temporary directory for downloads or archive it in some other way.
    local_filename = savedir
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


# Yes I should use a library for this but I really just need this to work for now
async def amap(function: Any, iterable: list) -> list:
    return_list = []
    for iter in iterable:
        return_list.append(await function(iter))
    return return_list


async def amap_fast(function: Any, iterable: list) -> list:
    tasks = list(map(function, iterable))
    return await asyncio.gather(*tasks)


tokenizer = tokenizers.Tokenizer.from_pretrained("bert-base-uncased")


def token_split(string: str, max_length: int, overlap: int = 0) -> list:
    tokenlist = tokenizer.encode(string).ids
    num_chunks = math.ceil(len(tokenlist) / max_length)
    chunk_size = math.ceil(
        len(tokenlist) / num_chunks
    )  # Takes in an integer $n$ and then outputs the nth token for use in a map call.

    def make_index(token_id: int) -> str:
        begin_index, end_index = (
            chunk_size * token_id,
            chunk_size * (token_id + 1) + overlap,
        )
        tokens = tokenlist[begin_index:end_index]
        return_string = tokenizer.decode(tokens)
        return return_string

    chunk_ids = range(0, num_chunks - 1)
    return list(map(make_index, chunk_ids))


# class MarkerServer(BaseModel):
#     url: str
#     connections: int = 0
#
#
# def create_server_list(urls: List[str]) -> List[MarkerServer]:
#     return_servers = []
#     for url in urls:
#         return_servers.append(MarkerServer(url=url, connections=0))
#     return return_servers
#
#
# global_marker_servers = create_server_list(global_marker_server_urls)
#
#
# def get_total_connections() -> int:
#     def total_connections_list(marker_servers: List[MarkerServer]) -> int:
#         total = 0
#         for marker_server in marker_servers:
#             total = total + marker_server.connections
#         return total
#
#     global global_marker_servers
#     return total_connections_list(global_marker_servers)
#
#
# def get_least_connection() -> MarkerServer:
#     def least_connection_list(marker_servers: List[MarkerServer]) -> MarkerServer:
#         min_conns = 999999
#         min_conn_server = None
#         for marker_server in marker_servers:
#             if marker_server.connections < min_conns:
#                 min_conn_server = marker_server
#                 min_conns = marker_server.connections
#         if min_conn_server is None:
#             raise Exception("Marker Server Not Found in List")
#         return min_conn_server
#
#     global global_marker_servers
#     return least_connection_list(global_marker_servers)
#
#
