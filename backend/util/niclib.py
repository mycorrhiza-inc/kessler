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

    match = frontmatter_pattern.match(markdown_text)

    if match:
        # Extract metadata
        frontmatter = match.group(1)
        # Remove front matter from markdown text to get main body
        main_body = markdown_text[match.end() :]
        #
        try:
            # Parse the YAML content
            yaml_dict = yaml.safe_load(frontmatter)
            return (main_body, yaml_dict)
        except yaml.YAMLError as e:
            print(f"Error parsing YAML: {e}")
            return (main_body, {})
        metadata = yaml.safe_load(frontmatter)
        return (main_body, metadata)
    else:
        # If no front matter, return markdown text and empty dictionary
        return (markdown_text, {})


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
