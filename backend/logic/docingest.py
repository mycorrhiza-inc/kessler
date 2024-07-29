from util.niclib import seperate_markdown_string
import yaml
from util.niclib import rand_string, rand_filepath


from typing import Optional, List, Union, Any

import json
import re
import logging
import requests
import subprocess
import urllib
import mimetypes
import os
from pathlib import Path
import shlex
from util.niclib import rand_string, get_blake2
from tempfile import TemporaryFile, NamedTemporaryFile, _TemporaryFileWrapper


from io import BufferedWriter


import shutil
import hashlib
import base64
from util.file_io import S3FileManager


from constants import OS_TMPDIR, OS_FILEDIR


class DocumentIngester:
    def __init__(
        self, logger, savedir=OS_FILEDIR, tmpdir=OS_TMPDIR / Path("kessler/docingest")
    ):
        self.tmpdir = tmpdir
        self.savedir = savedir
        self.logger = logger
        self.rawfile_savedir = savedir / Path("raw")
        # Make sure the different directories always exist
        self.rawfile_savedir.mkdir(exist_ok=True, parents=True)
        self.metadata_backupdir = savedir / Path("metadata")
        # Make sure the different directories always exist
        self.metadata_backupdir.mkdir(exist_ok=True, parents=True)
        self.proctext_backupdir = savedir / Path("processed_text")
        # Make sure the different directories always exist
        self.proctext_backupdir.mkdir(exist_ok=True, parents=True)
        self.file_manager = S3FileManager()

    def url_to_filepath_and_metadata(self, url: str) -> tuple[Path, dict]:
        self.logger.info("collecting filepath metadata")

        try:
            filepath, metadata = self.add_file_from_url_nocall(url)
            self.logger.info("Successfully got file and metadata")
        except Exception as e:
            self.logger.error(f"unable to get file and metadata from url: {url}")
            raise e

        if metadata.get("lang") == None:
            metadata["lang"] = "en"
        # file_metadata = self.file_manager.get_metadata_from_file_obj(filepath, metadata.get("doctype"))
        self.logger.info("Attempted to get metadata from file, adding to main source.")
        # FIXME :
        # metadata.update(file_metadata)

        return (filepath, metadata)

    def get_metada_from_url(self, url: str) -> dict:
        self.logger.info("Getting Metadata from Url")
        response = requests.get(url)

        def get_doctype_from_header_and_url(response, url: str) -> Optional[str]:
            # Guess the file extension from the URL itself
            # This is useful for direct links to files with a clear file extension in the URL
            if url.lower().endswith(
                (
                    ".pdf",
                    ".doc",
                    ".docx",
                    ".xls",
                    ".xlsx",
                    ".ppt",
                    ".pptx",
                    ".md",
                    ".txt",
                    ".epub",
                )
            ):
                return url.split(".")[-1].lower()
            content_type = response.headers.get("Content-Type")
            # If the Content-Type header is not found, return None
            if content_type is None or content_type == None:
                return None
            # Use the mimetypes library to get the corresponding extension
            file_extension = mimetypes.guess_extension(content_type.split(";")[0])
            if file_extension:
                return file_extension.strip(
                    "."
                )  # Remove the . at the beginning of the extension
            else:
                # No matching file extension found, return the MIME type directly
                return content_type.split("/")[-1]

        # Call function to get the document_type
        document_type = get_doctype_from_header_and_url(response, url)
        # Get language from the headers
        language = response.headers.get("Content-Language")
        # Get last modified from the headers
        last_modified = response.headers.get("Last-Modified")

        def url_to_name(url: str) -> str:
            parsed_url = urllib.parse.urlparse(url)
            netloc_path = parsed_url.netloc + parsed_url.path
            return netloc_path.replace("/", "-")

        name = url_to_name(url)
        return {
            "doctype": document_type,
            "language": language,
            "date": last_modified,
            "title": name,
        }

    def get_metadata_from_file(self, path: Path, doctype: str) -> dict:
        if doctype == "md":
            with open(path, "r") as file:
                result = file.read()
                text, metadata = seperate_markdown_string(result)
                # Make sure that the doc doesnt set the source type to something else
                # causing a crash when processing docs.
                metadata["doctype"] = doctype
            return metadata

        return {}

    def rectify_unknown_metadata(self, metadata: dict):
        assert metadata.get("doctype") != None

        def mut_rectify_empty_field(metadata: dict, field: str, defaultval: Any):
            if metadata.get(field) == None:
                metadata[field] = defaultval
            return metadata

        # TODO : Double check and test how mutable values in python work to remove unnecessary assignments.
        metadata = mut_rectify_empty_field(metadata, "title", "unknown")
        metadata = mut_rectify_empty_field(metadata, "author", "unknown")
        metadata = mut_rectify_empty_field(metadata, "language", "en")
        return metadata

    def add_file_from_url_nocall(self, url: str) -> tuple[Any, dict]:
        def rectify_unknown_metadata(metadata: dict):
            assert metadata.get("doctype") != None

            def mut_rectify_empty_field(metadata: dict, field: str, defaultval: Any):
                if metadata.get(field) == None:
                    metadata[field] = defaultval
                return metadata

            # TODO : Double check and test how mutable values in python work to remove unnecessary assignments.
            metadata = mut_rectify_empty_field(metadata, "title", "unknown")
            metadata = mut_rectify_empty_field(metadata, "author", "unknown")
            metadata = mut_rectify_empty_field(metadata, "language", "en")
            return metadata

        metadata = self.get_metada_from_url(url)
        self.logger.info("Got Metadata from Url")
        metadata = rectify_unknown_metadata(metadata)
        self.logger.info(f"Rectified missing metadata, yielding:{metadata}")
        tmpfile = self.download_file_to_file_in_tmpdir(url)
        self.logger.info("Successfully downloaded file from url")
        return (tmpfile, metadata)

    def infer_metadata_from_path(self, filepath: Path) -> dict:
        return_doctype = filepath.suffix
        if return_doctype[0] == ".":
            return_doctype = return_doctype[1:]
        return {"title": filepath.stem, "doctype": filepath.suffix}
