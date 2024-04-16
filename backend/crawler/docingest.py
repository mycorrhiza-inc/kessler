from util.niclib import rand_string, rand_filepath


from typing import Optional, List, Union, Any

import json
import re


# from habanero import Crossref


import logging

# Note: Refactoring imports.py

import requests

from typing import Optional, List, Union


# from langchain.vectorstores import FAISS

import subprocess
import shutil
import urllib
import mimetypes
import os
import pickle


from pathlib import Path
import shlex


from util.niclib import rand_string


from tempfile import TemporaryFile

OS_TMPDIR = Path(os.environ["TMPDIR"])


class DocumentIngester:
    def __init__(self,logger, tmpdir=OS_TMPDIR / Path("kessler/docingest")):
        self.tmpdir = tmpdir
        self.logger=logger
        # self.crossref = Crossref()

    def url_to_file_and_metadata(self, url: str) -> tuple[Path, dict]:
        self.logger.warn("File function running")
        parsed_url = urllib.parse.urlparse(url)
        domain = (
            parsed_url.netloc.split(".")[-2] + "." + parsed_url.netloc.split(".")[-1]
        )
        self.logger.info("Domain Successfully processed")
        # TODO:  youtube and arxiv document adding, or refactor and use in a crawler
        # if domain in ["youtube.com", "youtu.be"]:
        #     filepath, metadata = self.get_file_from_ytdlp(url)
        # elif domain in ["arxiv.org"]:
        #     filepath, metadata = self.get_file_from_arxiv(url)
        
        self.logger.info("Proceeding to Download Regular file")
        filepath, metadata = self.add_file_from_url_nocall(url)
        self.logger.info("Successfully got file and metadata")
        if metadata.get("lang") == None:
            metadata["lang"] = "en"
        return (filepath, metadata)

    def get_file_from_arxiv(self, url: str) -> tuple[Path, dict]:
        def extract_arxiv_id(url: str) -> Optional[str]:
            # Regular expression to match arXiv ID patterns in the URL
            arxiv_regex = re.compile(
                r"arxiv\.org/(?:abs|html|pdf|e-print)/(\d+\.\d+v?\d*)(?:\.pdf)?",
                re.IGNORECASE,
            )

            # Search for matches using the regex
            match = arxiv_regex.search(url)

            # If a match is found, return the ID in the required format
            if match:
                arxiv_id = match.group(1)
                return arxiv_id
            return None

        arxiv_id = extract_arxiv_id(url)
        if arxiv_id == None:
            logging.warn("Failed to find arxiv id, falling back to HTML,")
            return self.add_file_from_url_nocall(url)
        htmlurl = f"https://arxiv.org/html/{arxiv_id}v1"

        htmlresponse = requests.get(htmlurl)

        # TODO : Generalize this function into a general metadata searcher, this fails if the doi is not found for example.
        try:
            # metadata = self.crossref.works(ids=f"10.48550/arXiv.{arxiv_id}")
            assert 0==1 # TODO : Fix crossref metadata lookup
        except:
            self.logger.warning(
                "Not able to lookup metadata based on doi, defaulting to extracting html metadata from arxiv."
            )
            metadata = self.get_metada_from_url(f"https://arxiv.org/abs/{arxiv_id}")

        if "HTML is not available for the source" in htmlresponse.text:
            metadata["doctype"] = "pdf"
            pdfdir = download_file(
                f"https://arxiv.org/pdf/{arxiv_id}.pdf", self.tmpdir / rand_filepath()
            )
            return (pdfdir, metadata)
        htmlpath = self.tmpdir / Path(rand_string())
        with open(htmlpath, "w") as file:
            file.write(htmlresponse.text)
        metadata["doctype"] = "html"
        return (htmlpath, metadata)

    def get_file_from_ytdlp(self, url: str) -> Optional[tuple[Path, dict]]:
        filename = rand_string()
        ytdlp_path = self.tmpdir / Path(filename)
        video_path = self.tmpdir / Path(filename + ".mkv")
        json_path = self.tmpdir / Path(filename + ".info.json")
        json_filepath = self.tmpdir / Path(rand_string())
        command = f"yt-dlp --remux-video mkv --write-info-json -o {shlex.quote(str(ytdlp_path))} {shlex.quote(url)}"
        self.logger.info(f"Calling youtube dlp with call: {command}")
        try:
            result = subprocess.run(
                command,
                shell=True,
                check=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            with open(json_path, "r", encoding="utf-8") as info_file:
                info_data = json.load(info_file)
                metadata = {
                    "title": info_data.get("title"),
                    "author": info_data.get("uploader"),
                    "language": info_data.get("language"),
                    "url": url,
                    "doctype": "mkv",
                }
                self.logger.info((video_path, metadata))
                return (video_path, metadata)
        except subprocess.CalledProcessError as e:
            self.logger.critical(f"An error occurred when using yt-dlp: {e.stderr}")
            return None

    def get_metada_from_url(self, url: str) -> dict:
        self.logger.info("Getting Metadata from Url")

        response = requests.get(url)
        
        def get_doctype_from_header_and_url(response, url: str) -> Optional[str]:
            # Guess the file extension from the URL itself
            # This is useful for direct links to files with a clear file extension in the URL
            if url.lower().endswith(
                (".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx")
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
    def download_file_to_path(self,url: str, savepath: Path) -> Path:
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

    def download_file_to_tmpfile(self,url: str): # TODO : Get types for temporary file
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
    def rectify_unknown_metadata(self,metadata : dict):
        assert metadata.get("doctype") != None
        def mut_rectify_empty_field(metadata : dict, field : str, defaultval : Any):
            if metadata.get(field)==None:
                metadata[field]=defaultval
            return metadata
        # TODO : Double check and test how mutable values in python work to remove unnecessary assignments.
        metadata = mut_rectify_empty_field(metadata,"title","unknown")
        metadata = mut_rectify_empty_field(metadata,"author","unknown")
        metadata = mut_rectify_empty_field(metadata,"language","en")
        return metadata


    def add_file_from_url_nocall(self, url: str) -> tuple[Any, dict]:
        metadata = self.get_metada_from_url(url)
        self.logger.info("Got Metadata from Url")
        metadata = self.rectify_unknown_metadata(metadata)
        self.logger.info(f"Rectified missing metadata, yielding:{metadata}")
        tmpfile = self.download_file_to_tmpfile(url)
        self.logger.info("Successfully downloaded file from url")
        return (tmpfile, metadata)

    def infer_metadata_from_path(self, filepath: Path) -> dict:
        return {"title": filepath.stem, "doctype": filepath.suffix}
