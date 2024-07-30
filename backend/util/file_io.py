import boto3

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


from util.niclib import create_markdown_string, seperate_markdown_string
from constants import (
    OS_TMPDIR,
    OS_HASH_FILEDIR,
    OS_BACKUP_FILEDIR,
    CLOUD_REGION,
    S3_SECRET_KEY,
    S3_ACCESS_KEY,
    S3_ENDPOINT
)

default_logger = logging.getLogger(__name__)


class S3FileManager:
    def __init__(self, logger: Optional[Any] = None) -> None:
        if logger is None:
            logger = default_logger
        self.tmpdir = OS_TMPDIR
        self.rawfile_savedir = OS_HASH_FILEDIR
        self.metadata_backupdir = OS_BACKUP_FILEDIR
        self.logger = logger
        self.s3 = boto3.client(
            "s3",
            endpoint_url=S3_ENDPOINT,
            aws_access_key_id=S3_ACCESS_KEY,
            aws_secret_access_key=S3_SECRET_KEY,
            region_name=CLOUD_REGION,
        )


    def save_filepath_to_hash(
        self, filepath: Path, hashpath: Optional[Path] = None
    ) -> tuple[str, Path]:
        if hashpath is None:
            hashpath = self.rawfile_savedir
        filepath.parent.mkdir(exist_ok=True, parents=True)
        self.logger.info("Getting hash")
        b264_hash = self.get_blake2_str(filepath)
        self.logger.info(f"Got hash {b264_hash}")
        saveloc = self.get_default_filepath_from_hash(b264_hash, hashpath)

        self.logger.info(f"Saving file to {saveloc}")
        shutil.copyfile(filepath, saveloc)
        if saveloc.exists():
            self.logger.info(f"Successfully Saved File to: {saveloc}")
        else:
            self.logger.error(f"File could not be saved to : {saveloc}")
        return (b264_hash, saveloc)

    def get_default_filepath_from_hash(
        self, hash: str, hashpath: Optional[Path] = None
    ) -> Path:
        if hashpath is None:
            hashpath = self.rawfile_savedir
        hashpath.parent.mkdir(exist_ok=True, parents=True)
        saveloc = hashpath / Path(hash)
        if saveloc.exists():
            self.logger.info(f"File already at {saveloc}, do not copy any file to it.")
        return saveloc

    def backup_metadata_to_hash(self, metadata: dict, hash: str) -> Path:
        def backup_metadata_to_filepath(metadata: dict, filepath: Path) -> Path:
            with open(filepath, "w+") as ff:
                yaml.dump(metadata, ff)
            return filepath

        savedir = self.metadata_backupdir / Path(str(hash) + ".yaml")
        self.logger.info(f"Backing up metadata to: {savedir}")
        return backup_metadata_to_filepath(metadata, savedir)

    def write_tmpfile_to_path(self, tmp: Any, path: Path):
        path.parent.mkdir(exist_ok=True, parents=True)
        self.logger.info("Seeking to beginning of file")
        # Seek to the beginning of the file
        tmp.seek(0)
        self.logger.info("Attempting to read file contents")
        # Read the file contents
        try:
            file_contents = tmp.read()
            self.logger.info("Attempting to write contents to permanent file")
            # Write the file contents to the desired path
            with open(path, "wb") as dest_file:
                dest_file.write(file_contents)
        except Exception as e:
            self.logger.info(f"The error is: {e}")

    def get_blake2_str(
        self, file_input: Path
    ) -> str:  # TODO: Figure out how df file types work
        self.logger.info("Setting Blake2b as the hash method of choice")
        hasher = hashlib.blake2b
        hash_object = hasher()
        self.logger.info("Created Hash object and initialized hash.")
        if isinstance(file_input, Path):
            f = open(file_input, "rb")
            buf = f.read(65536)
            # self.logger.info(buf)
            while len(buf) > 0:
                hash_object.update(buf)
                buf = f.read(65536)
            return base64.urlsafe_b64encode(hash_object.digest()).decode()
        self.logger.error("Failed to hash file")
        raise Exception("ErrorHashingFile")  # I am really sorry about this

    def backup_processed_text(
        self, text: str, hash: str, metadata: dict, backupdir: Path
    ) -> None:
        savestring = create_markdown_string(
            text, metadata, include_previous_metadata=False
        )
        backuppath = backupdir / Path(hash + ".md")
        # Seems slow to check every time a file is backed up
        backuppath.parent.mkdir(parents=True, exist_ok=True)
        if backuppath.exists():
            backuppath.unlink(missing_ok=True)
        # FIXME: We should probably come up with a better backup protocol then doing everything with hashes
        if backuppath.is_file():
            backuppath.unlink(missing_ok=True)
        with open(backuppath, "w") as text_file:
            text_file.write(savestring)

    def download_file_to_path(self, url: str, savepath: Path) -> Path:
        savepath.parent.mkdir(exist_ok=True, parents=True)
        self.logger.info(f"Downloading file to dir: {savepath}")
        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            with open(savepath, "wb") as f:
                for chunk in r.iter_content(chunk_size=8192):
                    # If you have chunk encoded response uncomment if
                    # and set chunk_size parameter to None.
                    # if chunk:
                    f.write(chunk)
        return savepath

    # TODO : Get types for temporary file
    def download_file_to_tmpfile(self, url: str) -> Any:
        self.logger.info(f"Downloading file to temporary file")
        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            with TemporaryFile("wb") as f:
                for chunk in r.iter_content(chunk_size=8192):
                    # If you have chunk encoded response uncomment if
                    # and set chunk_size parameter to None.
                    # if chunk:
                    f.write(chunk)
                return f

    def check_s3_for_filehash(self, filehash : str) -> Optional[str]:
        

    def download_file_to_file_in_tmpdir(
        self, url: str
    ) -> Any:  # TODO : Get types for temporary file
        savedir = self.tmpdir / Path(rand_string())
        return self.download_file_to_path(url, savedir)

    def create_s3_uri(self, file_name: str, s3_credentials: Optional[Any] = None) -> str:
        return "test"

    def push_file_to_s3(self, filepath: Path, file_upload_name: str) -> str:
        return self.create_s3_uri("test")

    def push_raw_file_to_s3_novalid(self, filepath: Path, hash: str) -> str:
        if not filepath.is_file():
            raise Exception("File does not exist")
        filename = f"raw/{hash}"
        return self.push_file_to_s3(filepath, filename)

    def push_raw_file_to_s3(self, filepath: Path, hash: Optional[str] = None) -> str:
        if not filepath.is_file():
            raise Exception("File does not exist")
        actual_hash = self.get_blake2_str(filepath)
        if hash is not None and actual_hash != hash:
            raise Exception("Hashes did not match, erroring out")
        check_for_exisiting = self.check_s3_for_filehash(actual_hash)
        if check_for_exisiting is None:
            return self.push_raw_file_to_s3_novalid(filepath, actual_hash)
        return check_for_exisiting
